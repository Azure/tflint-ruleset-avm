package rules_test

import (
	"os"
	"testing"

	"github.com/Azure/tflint-ruleset-avm/attrvalue"
	"github.com/Azure/tflint-ruleset-avm/rules"
	"github.com/prashantv/gostub"
	"github.com/spf13/afero"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func TestModtmProviderVersionRule(t *testing.T) {
	cases := []struct {
		desc     string
		config   string
		rule     tflint.Rule
		expected helper.Issues
	}{
		{
			desc:   "no terraform block no issue, it's not our business",
			config: ``,
			rule: rules.NewProviderVersionRule("modtm", "Azure/modtm", "~> 0.3", []string{
				"0.2.99",
				"1.0.0",
			}),
			expected: helper.Issues{},
		},
		{
			desc: "no required_providers block no issue, it's not our business",
			config: `terraform {
}`,
			rule: rules.NewProviderVersionRule("modtm", "Azure/modtm", "~> 0.3", []string{
				"0.2.99",
				"1.0.0",
			}),
			expected: helper.Issues{},
		},
		{
			desc: "no modtm defined in required_providers should emit issue",
			config: `terraform{
  required_providers {
    azurerm = {
      source = "hashicorp/azurerm"
      version = "3.111.0"
    }
  }
}`,
			rule: rules.NewProviderVersionRule("modtm", "Azure/modtm", "~> 0.3", []string{
				"0.2.999",
				"1.0.0",
			}),
			expected: helper.Issues{
				{
					Rule: rules.NewProviderVersionRule("modtm", "Azure/modtm", "~> 0.3", []string{
						"0.2.999",
						"1.0.0",
					}),
					Message: "`modtm` provider should be declared in the `required_providers` block",
				},
			},
		},
		{
			desc: "modtm defined in required_providers with incorrect source emit issue",
			config: `terraform{
  required_providers {
    modtm = {
      source = "notAzure/modtm"
      version = ">=0.2.2"
    }
  }
}`,
			rule: rules.NewProviderVersionRule("modtm", "Azure/modtm", "~> 0.3", []string{
				"0.2.999",
				"1.0.0",
			}),
			expected: helper.Issues{
				{
					Rule: rules.NewProviderVersionRule("modtm", "Azure/modtm", "~> 0.3", []string{
						"0.2.999",
						"1.0.0",
					}),
					Message: "provider `modtm`'s source should be Azure/modtm, got notAzure/modtm",
				},
			},
		},
		{
			desc: "modtm defined with incorrect version constraint",
			config: `terraform {
  required_providers {
    modtm = {
      source = "Azure/modtm"
      version = ">=0.2.2"
    }
  }
}`,
			rule: rules.NewProviderVersionRule("modtm", "Azure/modtm", "~> 0.3", []string{
				"0.2.999",
				"1.0.0",
			}),
			expected: helper.Issues{
				{
					Rule: rules.NewProviderVersionRule("modtm", "Azure/modtm", "~> 0.3", []string{
						"0.2.999",
						"1.0.0",
					}),
					Message: "this module should not support provider `modtm` version 0.2.999, recommended version constraint: ~> 0.3",
				},
			},
		},
		{
			desc: "modtm defined with incorrect version constraint2",
			config: `terraform {
  required_providers {
    modtm = {
      source = "Azure/modtm"
      version = ">= 0.3.0"
    }
  }
}`,
			rule: rules.NewProviderVersionRule("modtm", "Azure/modtm", "~> 0.3", []string{
				"0.2.999",
				"1.0.0",
			}),
			expected: helper.Issues{
				{
					Rule: rules.NewProviderVersionRule("modtm", "Azure/modtm", "~> 0.3", []string{
						"0.2.999",
						"1.0.0",
					}),
					Message: "this module should not support provider `modtm` version 1.0.0, recommended version constraint: ~> 0.3",
				},
			},
		},
		{
			desc: "modtm defined with correct version constraint",
			config: `terraform {
  required_providers {
    modtm = {
      source = "Azure/modtm"
      version = "~> 0.3.0"
    }
  }
}`,
			rule: rules.NewProviderVersionRule("modtm", "Azure/modtm", "~> 0.3", []string{
				"0.2.999",
				"1.0.0",
			}),
			expected: helper.Issues{},
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"main.tf": c.config})
			stub := gostub.Stub(&attrvalue.AppFs, mockFs(c.config))
			defer stub.Reset()
			if err := c.rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			helper.AssertIssuesWithoutRange(t, c.expected, runner.Issues)
		})
	}
}

func mockFs(c string) afero.Afero {
	fs := afero.NewMemMapFs()
	_ = afero.WriteFile(fs, "main.tf", []byte(c), os.ModePerm)
	return afero.Afero{Fs: fs}
}

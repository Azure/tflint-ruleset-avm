package rules_test

import (
	"testing"

	"github.com/Azure/tflint-ruleset-avm/rules"
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
			desc:     "no terraform block no issue, it's not our business",
			config:   ``,
			rule:     rules.NewProviderVersionRule("modtm", "Azure/modtm", "0.3.0", "~> 0.3", true),
			expected: helper.Issues{},
		},
		{
			desc: "no required_providers block no issue, it's not our business",
			config: `terraform {
}`,
			rule:     rules.NewProviderVersionRule("modtm", "Azure/modtm", "0.3.0", "~> 0.3", true),
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
			rule: rules.NewProviderVersionRule("modtm", "Azure/modtm", "0.3.0", "~> 0.3", true),
			expected: helper.Issues{
				{
					Rule:    rules.NewProviderVersionRule("modtm", "Azure/modtm", "0.3.0", "~> 0.3", true),
					Message: "`modtm` provider should be declared in the `required_providers` block",
				},
			},
		},
		{
			desc: "no modtm defined in required_providers but not mandatory - ok",
			config: `terraform{
  required_providers {
    azurerm = {
      source = "hashicorp/azurerm"
      version = "3.111.0"
    }
  }
}`,
			rule:     rules.NewProviderVersionRule("modtm", "Azure/modtm", "0.3.0", "~> 0.3", false),
			expected: helper.Issues{},
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
			rule: rules.NewProviderVersionRule("modtm", "Azure/modtm", "0.3.0", "~> 0.3", true),
			expected: helper.Issues{
				{
					Rule:    rules.NewProviderVersionRule("modtm", "Azure/modtm", "0.3.0", "~> 0.3", true),
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
      version = "~> 0.2.2"
    }
  }
}`,
			rule: rules.NewProviderVersionRule("modtm", "Azure/modtm", "0.3.0", "~> 0.3", true),
			expected: helper.Issues{
				{
					Rule:    rules.NewProviderVersionRule("modtm", "Azure/modtm", "0.3.0", "~> 0.3", true),
					Message: "provider `modtm`'s version should satisfy 0.3.0, got ~> 0.2.2. Recommended version constraint `~> 0.3`",
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
			rule: rules.NewProviderVersionRule("modtm", "Azure/modtm", "0.2.0", "~> 0.2", true),
			expected: helper.Issues{
				{
					Rule:    rules.NewProviderVersionRule("modtm", "Azure/modtm", "1.0.0", "~> 0.3", true),
					Message: "provider `modtm`'s version should satisfy 0.2.0, got >= 0.3.0. Recommended version constraint `~> 0.2`",
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
			rule:     rules.NewProviderVersionRule("modtm", "Azure/modtm", "0.3.2", "~> 0.3", true),
			expected: helper.Issues{},
		},
		{
			desc: "modtm defined with correct version constraint but different case",
			config: `terraform {
  required_providers {
    modtm = {
      source = "azure/modtm"
      version = "~> 0.3.0"
    }
  }
}`,
			rule:     rules.NewProviderVersionRule("modtm", "Azure/modtm", "0.3.2", "~> 0.3", true),
			expected: helper.Issues{},
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"main.tf": c.config})
			if err := c.rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			helper.AssertIssuesWithoutRange(t, c.expected, runner.Issues)
		})
	}
}

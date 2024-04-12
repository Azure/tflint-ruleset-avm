package requiredversion_test

import (
	"testing"

	"github.com/Azure/tflint-ruleset-avm/requiredversion"
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestRequiredVersion(t *testing.T) {
	cases := []struct {
		desc   string
		files  map[string]string
		issues helper.Issues
	}{
		{
			desc: "required_version = '~> xx.xx.xx', ok",
			files: map[string]string{
				"terraform.tf": `terraform {
		 required_version = "~> 0.12.29"
		 required_providers {
		   azurerm = {
		     version = ">= 3.98.0"
		     source = "hashicorp/azurerm"
		   }
		 }
		}`,
			},
			issues: helper.Issues{},
		},
		{
			desc: "required_version = '>= xx.xx.xx, < xx.xx.xx', ok",
			files: map[string]string{
				"terraform.tf": `terraform {
		 required_version = ">= 1.6.0, < 1.7.4"
		 required_providers {
		   azurerm = {
		     version = ">= 3.98.0"
		     source = "hashicorp/azurerm"
		   }
		 }
		}`,
			},
			issues: helper.Issues{},
		},
		{
			desc: "required_version = 'xx.xx.xx', not ok",
			files: map[string]string{
				"terraform.tf": `terraform {
  required_version = "0.12.29"
  required_providers {
    azurerm = {
      version = ">= 3.98.0"
      source = "hashicorp/azurerm"
    }
  }
}`,
			},
			issues: helper.Issues{
				{
					Rule:    requiredversion.NewRequiredVersionRule("required_version", ""),
					Message: "The `required_version` property constraint can use the ~> #.# or the >= #.#.#, < #.#.# format",
					Range: hcl.Range{
						Filename: "terraform.tf",
					},
				},
			},
		},
		{
			desc: "required_version = '>= xx.xx.xx', not ok",
			files: map[string]string{
				"terraform.tf": `terraform {
  required_version = ">= 0.12.29"
  required_providers {
    azurerm = {
      version = ">= 3.98.0"
      source = "hashicorp/azurerm"
    }
  }
}`,
			},
			issues: helper.Issues{
				{
					Rule:    requiredversion.NewRequiredVersionRule("required_version", ""),
					Message: "The `required_version` property constraint can use the ~> #.# or the >= #.#.#, < #.#.# format",
					Range: hcl.Range{
						Filename: "terraform.tf",
					},
				},
			},
		},
		{
			desc: "required_version = '< xx.xx.xx', not ok",
			files: map[string]string{
				"terraform.tf": `terraform {
  required_version = "< 1.7.4"
  required_providers {
    azurerm = {
      version = ">= 3.98.0"
      source = "hashicorp/azurerm"
    }
  }
}`,
			},
			issues: helper.Issues{
				{
					Rule:    requiredversion.NewRequiredVersionRule("required_version", ""),
					Message: "The `required_version` property constraint can use the ~> #.# or the >= #.#.#, < #.#.# format",
					Range: hcl.Range{
						Filename: "terraform.tf",
					},
				},
			},
		},
		{
			desc: "required_version = ', xx.xx.xx', not ok",
			files: map[string]string{
				"terraform.tf": `terraform {
  required_version = ", 1.7.4"
  required_providers {
    azurerm = {
      version = ">= 3.98.0"
      source = "hashicorp/azurerm"
    }
  }
}`,
			},
			issues: helper.Issues{
				{
					Rule:    requiredversion.NewRequiredVersionRule("required_version", ""),
					Message: "The `required_version` property constraint can use the ~> #.# or the >= #.#.#, < #.#.# format",
					Range: hcl.Range{
						Filename: "terraform.tf",
					},
				},
			},
		},
		{
			desc: "no required_version set for terraform block",
			files: map[string]string{
				"terraform.tf": `terraform {
  required_providers {
    azurerm = {
      version = ">= 3.98.0"
      source = "hashicorp/azurerm"
    }
  }
}`,
			},
			issues: helper.Issues{
				{
					Rule:    requiredversion.NewRequiredVersionRule("required_version", ""),
					Message: "The `required_version` field should be declared in the `terraform` block",
					Range: hcl.Range{
						Filename: "terraform.tf",
					},
				},
			},
		},
		{
			desc: "required_version is not placed at the beginning of terraform block",
			files: map[string]string{
				"terraform.tf": `terraform {
		 required_providers {
		   azurerm = {
		     version = ">= 3.98.0"
		     source = "hashicorp/azurerm"
		   }
		 }
		 required_version = "~> 1.7.4"
		}`,
			},
			issues: helper.Issues{
				{
					Rule:    requiredversion.NewRequiredVersionRule("required_version", ""),
					Message: "The `required_version` field should be declared at the beginning of `terraform` block",
					Range: hcl.Range{
						Filename: "terraform.tf",
					},
				},
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			rule := requiredversion.NewRequiredVersionRule("required_version", "")

			runner := helper.TestRunner(t, tc.files)

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssuesWithoutRange(t, tc.issues, runner.Issues)
		})
	}
}

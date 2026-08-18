package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"testing"
)

func Test_TerraformModuleProviderDeclarationRule(t *testing.T) {
	cases := []struct {
		Name     string
		JSON     bool
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "1. empty provider block",
			Content: `
provider "azurerm" {
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformModuleProviderDeclarationRule(),
					Message: "Provider block in terraform module is expected to have and only have `alias` declared",
				},
			},
		},
		{
			Name: "2. provider block with field other than `alias` declared",
			Content: `
provider "azurerm" {
  location = "west"
}

provider "azurerm" {
  alias    = "test1"
  location = "east"
}

provider "azurerm" {
  alias    = "test2"
  features {}
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformModuleProviderDeclarationRule(),
					Message: "Provider block in terraform module is expected to have and only have `alias` declared",
				},
				{
					Rule:    NewTerraformModuleProviderDeclarationRule(),
					Message: "Provider block in terraform module is expected to have and only have `alias` declared",
				},
				{
					Rule:    NewTerraformModuleProviderDeclarationRule(),
					Message: "Provider block in terraform module is expected to have and only have `alias` declared",
				},
			},
		},
		{
			Name: "3. correct case",
			Content: `
provider "azurerm" {
  alias = "test"
}`,
			Expected: helper.Issues{},
		},
	}
	rule := NewTerraformModuleProviderDeclarationRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			filename := "config.tf"
			if tc.JSON {
				filename = "config.tf.json"
			}
			runner := helper.TestRunner(t, map[string]string{filename: tc.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}

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
			Name: "empty provider block",
			Content: `
provider "azurerm" {
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformModuleProviderDeclarationRule(),
					Message: "Provider blocks must not be declared in Terraform modules; declare aliases with `configuration_aliases` in `required_providers`",
				},
			},
		},
		{
			Name: "configured and aliased provider blocks",
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
					Message: "Provider blocks must not be declared in Terraform modules; declare aliases with `configuration_aliases` in `required_providers`",
				},
				{
					Rule:    NewTerraformModuleProviderDeclarationRule(),
					Message: "Provider blocks must not be declared in Terraform modules; declare aliases with `configuration_aliases` in `required_providers`",
				},
				{
					Rule:    NewTerraformModuleProviderDeclarationRule(),
					Message: "Provider blocks must not be declared in Terraform modules; declare aliases with `configuration_aliases` in `required_providers`",
				},
			},
		},
		{
			Name: "alias-only provider block uses legacy proxy pattern",
			Content: `
provider "azurerm" {
  alias = "test"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformModuleProviderDeclarationRule(),
					Message: "Provider blocks must not be declared in Terraform modules; declare aliases with `configuration_aliases` in `required_providers`",
				},
			},
		},
		{
			Name: "configuration aliases without provider block",
			Content: `
terraform {
  required_providers {
    azurerm = {
      source                = "hashicorp/azurerm"
      version               = "~> 4.0"
      configuration_aliases = [azurerm.alternate]
    }
  }
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

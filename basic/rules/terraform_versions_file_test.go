package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"testing"
)

func Test_TerraformVersionsFileRule(t *testing.T) {
	cases := []struct {
		Name     string
		JSON     bool
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "1. other file",
			Content: `
resource "azurerm_resource_group" "rg" {
  name     = "myTFResourceGroup"
  location = "westus2"
  tags = {
    Team = "DevOps"
    Environment = "Terraform Getting Started"
    #Team = "DevOps"
  }
}

resource "azurerm_virtual_network" "vnet" {
  name                = "myTFVnet"
  address_space       = ["10.0.0.0/16"]
  location            = "westus2"
  resource_group_name = azurerm_resource_group.rg.name
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "2. other type of blocks in versions.tf",
			Content: `
resource "azurerm_resource_group" "rg" {
  name     = "myTFResourceGroup"
  location = "westus2"
  tags = {
    Team = "DevOps"
    Environment = "Terraform Getting Started"
    #Team = "DevOps"
  }
}

resource "azurerm_virtual_network" "vnet" {
  name                = "myTFVnet"
  address_space       = ["10.0.0.0/16"]
  location            = "westus2"
  resource_group_name = azurerm_resource_group.rg.name
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformVersionsFileRule(),
					Message: "`versions.tf` should have and only have 1 `terraform` block",
				},
			},
		},
		{
			Name:    "3. no terraform block in versions.tf",
			Content: "",
			Expected: helper.Issues{
				{
					Rule:    NewTerraformVersionsFileRule(),
					Message: "`versions.tf` should have and only have 1 `terraform` block",
				},
			},
		},
	}
	rule := NewTerraformVersionsFileRule()

	for i, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			var filename string
			if i == 0 {
				filename = "config.tf"
			} else {
				filename = "versions.tf"
			}
			if tc.JSON {
				filename += ".json"
			}
			runner := helper.TestRunner(t, map[string]string{filename: tc.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}

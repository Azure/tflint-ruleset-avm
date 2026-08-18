package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AzurermArgOrderRule(t *testing.T) {

	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{

		{
			Name: "10. provider",
			Content: `
provider "azurerm" {
  features {}
  client_id = "temp"
}`,
			Expected: helper.Issues{
				{
					Rule: NewAzurermArgOrderRule(),
					Message: `Arguments are expected to be sorted in following order:
provider "azurerm" {
  client_id = "temp"
  
  features {}
}`,
				},
			},
		},
	}

	rule := NewAzurermArgOrderRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"config.tf": tc.Content})
		t.Run(tc.Name, func(t *testing.T) {
			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}
			AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}

func Test_NonAzurermResourceShouldNotBeChecked(t *testing.T) {
	code := `
resource "random_string" "key_vault_prefix" {
  length  = 6
  special = false
  upper   = false
  numeric = false
}`
	rule := NewAzurermArgOrderRule()
	runner := helper.TestRunner(t, map[string]string{"config.tf": code})
	if err := rule.Check(runner); err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}
	if len(runner.Issues) != 0 {
		t.Fatalf("unexpected issue")
	}
}

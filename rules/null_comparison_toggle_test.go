package rules_test

import (
	"testing"

	"github.com/Azure/tflint-ruleset-avm/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestNullComparisonToggle(t *testing.T) {
	cases := []struct {
		desc   string
		config string
		issues helper.Issues
	}{
		{
			desc: "object variable exists, ok",
			config: `variable "resource_group" {
  type = object({
    id = string
  })
}

resource "azurerm_resource_group" "test2" {
  count = var.resource_group == null ? 1 : 0
  name     = "acctest-rg-test02"
  location = "westeurope"
}`,
			issues: helper.Issues{},
		},
		{
			desc: "string variable exists, not ok",
			config: `variable "resource_group_id" {
  type = string
}

resource "azurerm_resource_group" "test2" {
  count = var.resource_group_id == null ? 1 : 0
  name     = "acctest-rg-test02"
  location = "westeurope"
}`,
			issues: helper.Issues{},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			rule := rules.NewNullComparisonToggleRule()
			filename := "terraform.tf"

			runner := helper.TestRunner(t, map[string]string{filename: tc.config})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssuesWithoutRange(t, tc.issues, runner.Issues)
		})
	}
}

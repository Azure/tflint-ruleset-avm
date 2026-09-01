package rules_test

import (
	"testing"

	"github.com/Azure/tflint-ruleset-avm/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestNoDoubleQuotesInIgnoreChanges(t *testing.T) {
	cases := []struct {
		desc   string
		config string
		issues helper.Issues
	}{
		{
			desc: "no double quotes, ok",
			config: `resource "azurerm_resource_group" "test" {
  name     = "acctestRG"
  location = "westeurope"

  lifecycle {
    ignore_changes = [tags, location]
  }
}`,
			issues: helper.Issues{},
		},
		{
			desc: "double quotes exist, not ok",
			config: `resource "azurerm_resource_group" "test" {
  name     = "acctestRG"
  location = "westeurope"

  lifecycle {
    ignore_changes = ["tags", "location"]
  }
}`,
			issues: helper.Issues{
				{
					Rule:    rules.NewNoDoubleQuotesInIgnoreChangesRule(),
					Message: "ignore_changes shouldn't include double quotes",
				},
			},
		},
		{
			desc: "some item includes double quotes, not ok",
			config: `resource "azurerm_resource_group" "test" {
  name     = "acctestRG"
  location = "westeurope"

  lifecycle {
    ignore_changes = ["tags", location]
  }
}`,
			issues: helper.Issues{
				{
					Rule:    rules.NewNoDoubleQuotesInIgnoreChangesRule(),
					Message: "ignore_changes shouldn't include double quotes",
				},
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			rule := rules.NewNoDoubleQuotesInIgnoreChangesRule()
			filename := "terraform.tf"

			runner := helper.TestRunner(t, map[string]string{filename: tc.config})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssuesWithoutRange(t, tc.issues, runner.Issues)
		})
	}
}

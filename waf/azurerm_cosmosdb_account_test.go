package waf_test

import (
	"testing"

	"github.com/prashantv/gostub"

	"github.com/Azure/tflint-ruleset-avm/attrvalue"
	"github.com/Azure/tflint-ruleset-avm/waf"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func TestAzurermCosmosDbAccountBackupMode(t *testing.T) {
	wafRules := waf.WafRules{}

	testCases := []struct {
		name     string
		rule     tflint.Rule
		content  string
		expected helper.Issues
	}{
		{
			name: "correct setting",
			rule: wafRules.AzurermCosmosDbAccountBackupMode(),
			content: `
	variable "backup_type" {
		type    = string
		default = "Continuous"
	}
	resource "azurerm_cosmosdb_account" "example" {
		backup {
			type = var.backup_type
		}
	}`,
			expected: helper.Issues{},
		},
		{
			name: "incorrect setting",
			rule: wafRules.AzurermCosmosDbAccountBackupMode(),
			content: `
    variable "backup_type" {
		type    = string
		default = "Periodic"
	}
	resource "azurerm_cosmosdb_account" "example" {
		backup {
			type = var.backup_type
		}
	}`,
			expected: helper.Issues{
				{
					Rule:    wafRules.AzurermCosmosDbAccountBackupMode(),
					Message: "Periodic is an invalid attribute value of `type` - expecting (one of) [Continuous]",
				},
			},
		},
		{
			name: "missing block",
			rule: wafRules.AzurermCosmosDbAccountBackupMode(),
			content: `
	resource "azurerm_cosmosdb_account" "example" {

	}`,
			expected: helper.Issues{
				{
					Rule:    wafRules.AzurermCosmosDbAccountBackupMode(),
					Message: "The attribute `type` must be specified",
				},
			},
		},
		{
			name: "missing block attribute",
			rule: wafRules.AzurermCosmosDbAccountBackupMode(),
			content: `
	resource "azurerm_cosmosdb_account" "example" {
		backup {
			
		}
	}`,
			expected: helper.Issues{
				{
					Rule:    wafRules.AzurermCosmosDbAccountBackupMode(),
					Message: "The attribute `type` must be specified",
				},
			},
		},
		{
			name: "missing resource",
			rule: wafRules.AzurermCosmosDbAccountBackupMode(),
			content: `
	resource "something_else" "example" {

	}`,
			expected: helper.Issues{},
		},
	}

	filename := "main.tf"
	for _, c := range testCases {
		tc := c
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{filename: tc.content})
			stub := gostub.Stub(&attrvalue.AppFs, mockFs(tc.content))
			defer stub.Reset()
			if err := tc.rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			helper.AssertIssuesWithoutRange(t, tc.expected, runner.Issues)
		})
	}
}

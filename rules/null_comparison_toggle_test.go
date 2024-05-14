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
			issues: helper.Issues{
				{
					Rule:    rules.NewNullComparisonToggleRule(),
					Message: "The variable should be defined as object type for the resource id",
				},
			},
		},
		{
			desc: "string local exists, ok",
			config: `variable "resource_group_id" {
				 type = string
				}
		
				locals {
				 resource_group_id = var.resource_group_id
				}
		
				resource "azurerm_resource_group" "test2" {
				 count = local.resource_group_id == null ? 1 : 0
				 name     = "acctest-rg-test02"
				 location = "westeurope"
				}`,
			issues: helper.Issues{},
		},
		{
			desc: "string variable and another condition exist, not ok",
			config: `variable "resource_group_id" {
		type = string
		}
		
		variable "resource_group_enabled" {
		type = bool
		}
		
		resource "azurerm_resource_group" "test2" {
		count = var.resource_group_enabled && (var.resource_group_id != null) ? 0 : 1
		name     = "acctest-rg-test02"
		location = "westeurope"
		}`,
			issues: helper.Issues{
				{
					Rule:    rules.NewNullComparisonToggleRule(),
					Message: "The variable should be defined as object type for the resource id",
				},
			},
		},
		{
			desc: "string variable with `!=` exists, not ok",
			config: `variable "resource_group_id" {
		type = string
		}
		
		resource "azurerm_resource_group" "test2" {
		count = var.resource_group_id != null ? 0 : 1
		name     = "acctest-rg-test02"
		location = "westeurope"
		}`,
			issues: helper.Issues{
				{
					Rule:    rules.NewNullComparisonToggleRule(),
					Message: "The variable should be defined as object type for the resource id",
				},
			},
		},
		{
			desc: "string variable exists and the expression starts with null, not ok",
			config: `variable "resource_group_id" {
		type = string
		}
		
		resource "azurerm_resource_group" "test2" {
		count = null == var.resource_group_id ? 1 : 0
		name     = "acctest-rg-test02"
		location = "westeurope"
		}`,
			issues: helper.Issues{
				{
					Rule:    rules.NewNullComparisonToggleRule(),
					Message: "The variable should be defined as object type for the resource id",
				},
			},
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

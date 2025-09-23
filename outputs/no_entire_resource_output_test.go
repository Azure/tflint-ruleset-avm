package outputs_test

import (
	"testing"

	"github.com/Azure/tflint-ruleset-avm/outputs"
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestNoEntireResourceOutputRule(t *testing.T) {
	cases := []struct {
		desc   string
		config string
		issues helper.Issues
	}{
		{
			desc: "output attribute only (ok)",
			config: `resource "azurerm_resource_group" "rg" {
  location = "westeurope"
  name     = "rg-test"
}

output "rg_id" {
  value = azurerm_resource_group.rg.id
}`,
			issues: helper.Issues{},
		},
		{
			desc: "entire resource output (not ok)",
			config: `resource "azurerm_resource_group" "rg" {
  location = "westeurope"
  name     = "rg-test"
}

output "rg" {
  value = azurerm_resource_group.rg
}`,
			issues: helper.Issues{{
				Rule:    outputs.NewNoEntireResourceOutputRule("no_entire_resource_output", ""),
				Message: "authors SHOULD NOT output entire resource objects; expose specific computed attributes instead",
				Range:   hcl.Range{Filename: "variables.tf"},
			}},
		},
		{
			desc: "entire resource index output (not ok)",
			config: `resource "azurerm_resource_group" "rg" {
  count    = 1
  location = "westeurope"
  name     = "rg-test"
}

output "rg" {
  value = azurerm_resource_group.rg[0]
}`,
			issues: helper.Issues{{
				Rule:    outputs.NewNoEntireResourceOutputRule("no_entire_resource_output", ""),
				Message: "authors SHOULD NOT output entire resource objects; expose specific computed attributes instead",
				Range:   hcl.Range{Filename: "variables.tf"},
			}},
		},
		{
			desc: "for_each attribute only (ok)",
			config: `resource "azurerm_resource_group" "rg" {
  for_each = { a = "one", b = "two" }
  location = "westeurope"
  name     = "rg-${each.key}"
}

output "rg_ids" {
  value = [for v in azurerm_resource_group.rg : v.id]
}`,
			issues: helper.Issues{},
		},
		{
			desc: "for_each whole instance (not ok)",
			config: `resource "azurerm_resource_group" "rg" {
  for_each = { a = "one" }
  location = "westeurope"
  name     = "rg-${each.key}"
}

output "rg" {
  value = azurerm_resource_group.rg["a"]
}`,
			issues: helper.Issues{{
				Rule:    outputs.NewNoEntireResourceOutputRule("no_entire_resource_output", ""),
				Message: "authors SHOULD NOT output entire resource objects; expose specific computed attributes instead",
				Range:   hcl.Range{Filename: "variables.tf"},
			}},
		},
		{
			desc: "for_each whole collection splat (not ok)",
			config: `resource "azurerm_resource_group" "rg" {
  for_each = { a = "one", b = "two" }
  location = "westeurope"
  name     = "rg-${each.key}"
}

output "rgs" {
  value = azurerm_resource_group.rg[*]
}`,
			issues: helper.Issues{{
				Rule:    outputs.NewNoEntireResourceOutputRule("no_entire_resource_output", ""),
				Message: "authors SHOULD NOT output entire resource objects; expose specific computed attributes instead",
				Range:   hcl.Range{Filename: "variables.tf"},
			}},
		},
		{
			desc: "for_each base map reference (not ok)",
			config: `resource "azurerm_resource_group" "rg" {
  for_each = { a = "one", b = "two" }
  location = "westeurope"
  name     = "rg-${each.key}"
}

output "rgs" {
  value = azurerm_resource_group.rg
}`,
			issues: helper.Issues{{
				Rule:    outputs.NewNoEntireResourceOutputRule("no_entire_resource_output", ""),
				Message: "authors SHOULD NOT output entire resource objects; expose specific computed attributes instead",
				Range:   hcl.Range{Filename: "variables.tf"},
			}},
		},
		{
			desc: "for_each splat attribute access (ok)",
			config: `resource "azurerm_resource_group" "rg" {
  for_each = { a = "one", b = "two" }
  location = "westeurope"
  name     = "rg-${each.key}"
}

output "rg_ids" {
  value = azurerm_resource_group.rg[*].id
}`,
			issues: helper.Issues{},
		},
	}

	for _, tc := range cases {
		c := tc
		// nolint:scopelint // parallel subtests closure capture
		t.Run(c.desc, func(t *testing.T) {
			t.Parallel()
			rule := outputs.NewNoEntireResourceOutputRule("no_entire_resource_output", "")
			filename := "variables.tf"
			runner := helper.TestRunner(t, map[string]string{filename: c.config})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssuesWithoutRange(t, c.issues, runner.Issues)
		})
	}
}

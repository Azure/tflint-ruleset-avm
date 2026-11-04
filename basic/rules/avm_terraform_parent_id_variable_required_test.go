package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

const parentIDVariable = `
variable "parent_id" {
  type        = string
  description = "The resource ID of the parent scope."
  nullable    = false
}
`

func Test_TerraformParentIDVariable(t *testing.T) {
	rule := NewTerraformParentIDVariableRule()
	runRuleTestCases(t, rule, []ruleTestCase{
		{
			Name: "compliant AzAPI resource module",
			Files: map[string]string{
				"variables.tf": parentIDVariable,
				"main.tf": `
resource "azapi_resource" "this" {
  type      = "Microsoft.KeyVault/vaults@2023-07-01"
  name      = var.name
  parent_id = var.parent_id
}`,
			},
			Expected: helper.Issues{},
		},
		{
			Name: "modules without AzAPI managed resources are not checked",
			Files: map[string]string{
				"variables.tf": `
variable "resource_group_name" {
  type = string
}`,
				"main.tf": `
resource "azurerm_key_vault" "this" {
  resource_group_name = var.resource_group_name
}

data "azapi_resource" "existing" {
  type = "Microsoft.Resources/resourceGroups@2024-03-01"
}`,
			},
			Expected: helper.Issues{},
		},
		{
			Name: "missing parent_id variable",
			Files: map[string]string{
				"main.tf": `
resource "azapi_resource" "this" {
  type      = "Microsoft.KeyVault/vaults@2023-07-01"
  parent_id = "/subscriptions/0000/resourceGroups/rg"
}`,
			},
			Expected: helper.Issues{
				{Rule: rule, Message: "Variable `parent_id` must be declared to expose the parent scope of the primary AzAPI resource"},
				{Rule: rule, Message: "`azapi_resource.this` must set `parent_id = var.parent_id`"},
			},
		},
		{
			Name: "parent_id must be a required string",
			Files: map[string]string{
				"variables.tf": `
variable "parent_id" {
  type     = any
  default  = null
  nullable = false
}`,
				"main.tf": `
resource "azapi_resource" "this" {
  type      = "Microsoft.KeyVault/vaults@2023-07-01"
  parent_id = var.parent_id
}`,
			},
			Expected: helper.Issues{
				{Rule: rule, Message: "Variable `parent_id` must be of type `string`"},
				{Rule: rule, Message: "Variable `parent_id` must be required and must not declare a default value"},
			},
		},
		{
			Name: "parent_id without a type is flagged",
			Files: map[string]string{
				"variables.tf": `
variable "parent_id" {}`,
				"main.tf": `
resource "azapi_resource" "this" {
  parent_id = var.parent_id
}`,
			},
			Expected: helper.Issues{
				{Rule: rule, Message: "Variable `parent_id` must be of type `string`"},
				{Rule: rule, Message: "Variable `parent_id` must set `nullable = false`"},
			},
		},
		{
			Name: "parent_id without nullable is flagged",
			Files: map[string]string{
				"variables.tf": `
variable "parent_id" {
  type = string
}`,
				"main.tf": `
resource "azapi_resource" "this" {
  parent_id = var.parent_id
}`,
			},
			Expected: helper.Issues{
				{Rule: rule, Message: "Variable `parent_id` must set `nullable = false`"},
			},
		},
		{
			Name: "parent_id with nullable true is flagged",
			Files: map[string]string{
				"variables.tf": `
variable "parent_id" {
  type     = string
  nullable = true
}`,
				"main.tf": `
resource "azapi_resource" "this" {
  parent_id = var.parent_id
}`,
			},
			Expected: helper.Issues{
				{Rule: rule, Message: "Variable `parent_id` must set `nullable = false`"},
			},
		},
		{
			Name: "parent_id with non-literal nullable is flagged",
			Files: map[string]string{
				"variables.tf": `
variable "parent_id" {
  type     = string
  nullable = local.parent_id_nullable
}`,
				"main.tf": `
resource "azapi_resource" "this" {
  parent_id = var.parent_id
}`,
			},
			Expected: helper.Issues{
				{Rule: rule, Message: "Variable `parent_id` must set `nullable = false`"},
			},
		},
		{
			Name: "parent-scope-specific variables are flagged",
			Files: map[string]string{
				"variables.tf": parentIDVariable + `
variable "resource_group_name" {
  type = string
}

variable "resource_group_resource_id" {
  type = string
}`,
				"main.tf": `
resource "azapi_resource" "this" {
  parent_id = var.parent_id
}`,
			},
			Expected: helper.Issues{
				{Rule: rule, Message: "Variable `resource_group_name` must not be declared; expose the parent scope through the `parent_id` variable instead"},
				{Rule: rule, Message: "Variable `resource_group_resource_id` must not be declared; expose the parent scope through the `parent_id` variable instead"},
			},
		},
		{
			Name: "azapi_resource.this must forward var.parent_id directly",
			Files: map[string]string{
				"variables.tf": parentIDVariable,
				"main.tf": `
resource "azapi_resource" "this" {
  parent_id = "${var.parent_id}/providers/Microsoft.Network/virtualNetworks/vnet"
}`,
			},
			Expected: helper.Issues{
				{Rule: rule, Message: "`azapi_resource.this` must set `parent_id = var.parent_id`"},
			},
		},
		{
			Name: "azapi_resource.this without parent_id is flagged",
			Files: map[string]string{
				"variables.tf": parentIDVariable,
				"main.tf": `
resource "azapi_resource" "this" {
  type = "Microsoft.KeyVault/vaults@2023-07-01"
}`,
			},
			Expected: helper.Issues{
				{Rule: rule, Message: "`azapi_resource.this` must set `parent_id = var.parent_id`"},
			},
		},
		{
			Name: "satellite resources may derive parent_id from the primary resource",
			Files: map[string]string{
				"variables.tf": parentIDVariable,
				"main.tf": `
resource "azapi_resource" "this" {
  parent_id = var.parent_id
}

resource "azapi_resource" "lock" {
  parent_id = azapi_resource.this.id
}`,
			},
			Expected: helper.Issues{},
		},
	})
}

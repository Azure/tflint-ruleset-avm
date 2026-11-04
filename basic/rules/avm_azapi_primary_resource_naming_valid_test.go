package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AzapiPrimaryResourceNaming(t *testing.T) {
	rule := NewAzapiPrimaryResourceNamingRule()
	runRuleTestCases(t, rule, []ruleTestCase{
		{
			Name: "primary resource named this with named satellites",
			Files: map[string]string{"main.tf": `
resource "azapi_resource" "this" {
  type = "Microsoft.KeyVault/vaults@2023-07-01"
}

resource "azapi_resource" "lock" {
  type = "Microsoft.Authorization/locks@2020-05-01"
}

resource "azapi_resource" "role_assignments" {
  for_each = var.role_assignments
  type     = "Microsoft.Authorization/roleAssignments@2022-04-01"
}`},
			Expected: helper.Issues{},
		},
		{
			Name: "other AzAPI resource types may be the primary resource",
			Files: map[string]string{"main.tf": `
resource "azapi_update_resource" "this" {
  type = "Microsoft.Network/virtualNetworks@2024-05-01"
}`},
			Expected: helper.Issues{},
		},
		{
			Name: "modules without AzAPI managed resources are not checked",
			Files: map[string]string{"main.tf": `
resource "azurerm_key_vault" "example" {}

data "azapi_resource" "existing" {}

data "azapi_resource_list" "all" {}`},
			Expected: helper.Issues{},
		},
		{
			Name: "missing primary resource named this",
			Files: map[string]string{
				"b.tf": `
resource "azapi_resource" "vault" {}`,
				"a.tf": `
resource "azapi_resource" "lock" {}`,
			},
			Expected: helper.Issues{
				{Rule: rule, Message: "The primary AzAPI resource must be named `this`, for example `azapi_resource.this`"},
			},
		},
		{
			Name: "satellites named this are flagged",
			Files: map[string]string{"main.tf": `
resource "azapi_update_resource" "this" {}

resource "azapi_resource" "this" {}

resource "azapi_resource_action" "this" {}`},
			Expected: helper.Issues{
				{Rule: rule, Message: "`azapi_update_resource.this` must not be named `this` because `azapi_resource.this` is the primary resource; name it after what it represents"},
				{Rule: rule, Message: "`azapi_resource_action.this` must not be named `this` because `azapi_resource.this` is the primary resource; name it after what it represents"},
			},
		},
		{
			Name: "first resource named this is primary when azapi_resource.this is absent",
			Files: map[string]string{"main.tf": `
resource "azapi_update_resource" "this" {}

resource "azapi_resource_action" "this" {}`},
			Expected: helper.Issues{
				{Rule: rule, Message: "`azapi_resource_action.this` must not be named `this` because `azapi_update_resource.this` is the primary resource; name it after what it represents"},
			},
		},
	})
}

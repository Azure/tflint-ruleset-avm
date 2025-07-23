package rules_test

import (
	"testing"

	"github.com/Azure/tflint-ruleset-avm/rules"
	"github.com/stretchr/testify/require"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestAzApiNoPreviewVersion(t *testing.T) {
	cases := []struct {
		desc   string
		config string
		issues helper.Issues
	}{
		// azapi_resource tests
		{
			desc: "azapi_resource with non-preview API version, ok",
			config: `resource "azapi_resource" "test" {
  type     = "Microsoft.Storage/storageAccounts@2023-01-01"
  name     = "test"
  location = "westeurope"
  parent_id = "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/test"
}`,
			issues: helper.Issues{},
		},
		{
			desc: "azapi_resource with preview API version, not ok",
			config: `resource "azapi_resource" "test" {
  type     = "Microsoft.Storage/storageAccounts@2023-01-01-preview"
  name     = "test"
  location = "westeurope"
  parent_id = "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/test"
}`,
			issues: helper.Issues{
				{
					Rule:    &rules.AzApiNoPreviewVersionRule{},
					Message: "Resource type `Microsoft.Storage/storageAccounts@2023-01-01-preview` is using a preview API version, which is not recommended.",
				},
			},
		},
		// azapi_update_resource tests
		{
			desc: "azapi_update_resource with non-preview API version, ok",
			config: `resource "azapi_update_resource" "test" {
  type      = "Microsoft.Storage/storageAccounts@2023-01-01"
  resource_id = "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/test/providers/Microsoft.Storage/storageAccounts/test"
}`,
			issues: helper.Issues{},
		},
		{
			desc: "azapi_update_resource with preview API version, not ok",
			config: `resource "azapi_update_resource" "test" {
  type      = "Microsoft.Storage/storageAccounts@2023-01-01-preview"
  resource_id = "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/test/providers/Microsoft.Storage/storageAccounts/test"
}`,
			issues: helper.Issues{
				{
					Rule:    &rules.AzApiNoPreviewVersionRule{},
					Message: "Resource type `Microsoft.Storage/storageAccounts@2023-01-01-preview` is using a preview API version, which is not recommended.",
				},
			},
		},
		// azapi_resource_action tests
		{
			desc: "azapi_resource_action with non-preview API version, ok",
			config: `resource "azapi_resource_action" "test" {
  type        = "Microsoft.Storage/storageAccounts@2023-01-01"
  resource_id = "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/test/providers/Microsoft.Storage/storageAccounts/test"
  action      = "listKeys"
}`,
			issues: helper.Issues{},
		},
		{
			desc: "azapi_resource_action with preview API version, not ok",
			config: `resource "azapi_resource_action" "test" {
  type        = "Microsoft.Storage/storageAccounts@2023-01-01-preview"
  resource_id = "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/test/providers/Microsoft.Storage/storageAccounts/test"
  action      = "listKeys"
}`,
			issues: helper.Issues{
				{
					Rule:    &rules.AzApiNoPreviewVersionRule{},
					Message: "Resource type `Microsoft.Storage/storageAccounts@2023-01-01-preview` is using a preview API version, which is not recommended.",
				},
			},
		},
		// azapi_data_plane_resource tests
		{
			desc: "azapi_data_plane_resource with non-preview API version, ok",
			config: `resource "azapi_data_plane_resource" "test" {
  type      = "Microsoft.Storage/storageAccounts@2023-01-01"
  name      = "test"
  parent_id = "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/test"
}`,
			issues: helper.Issues{},
		},
		{
			desc: "azapi_data_plane_resource with preview API version, not ok",
			config: `resource "azapi_data_plane_resource" "test" {
  type      = "Microsoft.Storage/storageAccounts@2023-01-01-preview"
  name      = "test"
  parent_id = "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/test"
}`,
			issues: helper.Issues{
				{
					Rule:    &rules.AzApiNoPreviewVersionRule{},
					Message: "Resource type `Microsoft.Storage/storageAccounts@2023-01-01-preview` is using a preview API version, which is not recommended.",
				},
			},
		},
		// Data source tests
		{
			desc: "azapi_resource data source with non-preview API version, ok",
			config: `data "azapi_resource" "test" {
  type      = "Microsoft.Storage/storageAccounts@2023-01-01"
  name      = "test"
  parent_id = "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/test"
}`,
			issues: helper.Issues{},
		},
		{
			desc: "azapi_resource data source with preview API version, not ok",
			config: `data "azapi_resource" "test" {
  type      = "Microsoft.Storage/storageAccounts@2023-01-01-preview"
  name      = "test"
  parent_id = "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/test"
}`,
			issues: helper.Issues{
				{
					Rule:    &rules.AzApiNoPreviewVersionRule{},
					Message: "Resource type `Microsoft.Storage/storageAccounts@2023-01-01-preview` is using a preview API version, which is not recommended.",
				},
			},
		},
		{
			desc: "azapi_resource_action data source with non-preview API version, ok",
			config: `data "azapi_resource_action" "test" {
  type        = "Microsoft.Storage/storageAccounts@2023-01-01"
  resource_id = "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/test/providers/Microsoft.Storage/storageAccounts/test"
  action      = "listKeys"
}`,
			issues: helper.Issues{},
		},
		{
			desc: "azapi_resource_action data source with preview API version, not ok",
			config: `data "azapi_resource_action" "test" {
  type        = "Microsoft.Storage/storageAccounts@2023-01-01-preview"
  resource_id = "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/test/providers/Microsoft.Storage/storageAccounts/test"
  action      = "listKeys"
}`,
			issues: helper.Issues{
				{
					Rule:    &rules.AzApiNoPreviewVersionRule{},
					Message: "Resource type `Microsoft.Storage/storageAccounts@2023-01-01-preview` is using a preview API version, which is not recommended.",
				},
			},
		},
		{
			desc: "azapi_resource_list data source with non-preview API version, ok",
			config: `data "azapi_resource_list" "test" {
  type      = "Microsoft.Storage/storageAccounts@2023-01-01"
  parent_id = "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/test"
}`,
			issues: helper.Issues{},
		},
		{
			desc: "azapi_resource_list data source with preview API version, not ok",
			config: `data "azapi_resource_list" "test" {
  type      = "Microsoft.Storage/storageAccounts@2023-01-01-preview"
  parent_id = "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/test"
}`,
			issues: helper.Issues{
				{
					Rule:    &rules.AzApiNoPreviewVersionRule{},
					Message: "Resource type `Microsoft.Storage/storageAccounts@2023-01-01-preview` is using a preview API version, which is not recommended.",
				},
			},
		},
		// Ephemeral tests
		{
			desc: "azapi_resource ephemeral with non-preview API version, ok",
			config: `ephemeral "azapi_resource" "test" {
  type      = "Microsoft.Storage/storageAccounts@2023-01-01"
  name      = "test"
  parent_id = "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/test"
}`,
			issues: helper.Issues{},
		},
		{
			desc: "azapi_resource ephemeral with preview API version, not ok",
			config: `ephemeral "azapi_resource" "test" {
  type      = "Microsoft.Storage/storageAccounts@2023-01-01-preview"
  name      = "test"
  parent_id = "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/test"
}`,
			issues: helper.Issues{
				{
					Rule:    &rules.AzApiNoPreviewVersionRule{},
					Message: "Resource type `Microsoft.Storage/storageAccounts@2023-01-01-preview` is using a preview API version, which is not recommended.",
				},
			},
		},
		{
			desc: "azapi_resource_action ephemeral with non-preview API version, ok",
			config: `ephemeral "azapi_resource_action" "test" {
  type        = "Microsoft.Storage/storageAccounts@2023-01-01"
  resource_id = "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/test/providers/Microsoft.Storage/storageAccounts/test"
  action      = "listKeys"
}`,
			issues: helper.Issues{},
		},
		{
			desc: "azapi_resource_action ephemeral with preview API version, not ok",
			config: `ephemeral "azapi_resource_action" "test" {
  type        = "Microsoft.Storage/storageAccounts@2023-01-01-preview"
  resource_id = "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/test/providers/Microsoft.Storage/storageAccounts/test"
  action      = "listKeys"
}`,
			issues: helper.Issues{
				{
					Rule:    &rules.AzApiNoPreviewVersionRule{},
					Message: "Resource type `Microsoft.Storage/storageAccounts@2023-01-01-preview` is using a preview API version, which is not recommended.",
				},
			},
		},
		// Edge cases and ignored scenarios
		{
			desc: "non-azapi resource, should be ignored",
			config: `resource "azurerm_storage_account" "test" {
  name                     = "test"
  resource_group_name      = "test"
  location                 = "westeurope"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}`,
			issues: helper.Issues{},
		},
		{
			desc: "azapi_resource without type attribute, should be ignored",
			config: `resource "azapi_resource" "test" {
  name      = "test"
  location  = "westeurope"
  parent_id = "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/test"
}`,
			issues: helper.Issues{},
		},
		{
			desc:   "azapi_client_config data source without type attribute, should be ignored",
			config: `data "azapi_client_config" "current" {}`,
			issues: helper.Issues{},
		},
		{
			desc: "azapi_resource_id data source without type attribute, should be ignored",
			config: `data "azapi_resource_id" "test" {
  name      = "test"
  parent_id = "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/test"
  resource_type = "Microsoft.Storage/storageAccounts"
}`,
			issues: helper.Issues{},
		},
		// Multiple resources mixed scenario
		{
			desc: "multiple azapi resources with mixed API versions",
			config: `resource "azapi_resource" "test1" {
  type     = "Microsoft.Storage/storageAccounts@2023-01-01"
  name     = "test1"
  location = "westeurope"
  parent_id = "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/test"
}

resource "azapi_resource" "test2" {
  type     = "Microsoft.Compute/virtualMachines@2023-07-01-preview"
  name     = "test2"
  location = "westeurope"
  parent_id = "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/test"
}

data "azapi_resource_action" "test3" {
  type        = "Microsoft.Network/virtualNetworks@2023-02-01-preview"
  resource_id = "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/test/providers/Microsoft.Network/virtualNetworks/test"
  action      = "checkIPAddressAvailability"
}`,
			issues: helper.Issues{
				{
					Rule:    &rules.AzApiNoPreviewVersionRule{},
					Message: "Resource type `Microsoft.Compute/virtualMachines@2023-07-01-preview` is using a preview API version, which is not recommended.",
				},
				{
					Rule:    &rules.AzApiNoPreviewVersionRule{},
					Message: "Resource type `Microsoft.Network/virtualNetworks@2023-02-01-preview` is using a preview API version, which is not recommended.",
				},
			},
		},
		// Different preview version formats
		{
			desc: "different preview version format variations, all should be detected",
			config: `resource "azapi_resource" "test1" {
  type = "Microsoft.Storage/storageAccounts@2023-01-01-PREVIEW"
  name = "test1"
  location = "westeurope"
  parent_id = "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/test"
}`,
			issues: helper.Issues{
				{
					Rule:    &rules.AzApiNoPreviewVersionRule{},
					Message: "Resource type `Microsoft.Storage/storageAccounts@2023-01-01-PREVIEW` is using a preview API version, which is not recommended.",
				},
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			rule := &rules.AzApiNoPreviewVersionRule{}
			filename := "terraform.tf"

			runner := helper.TestRunner(t, map[string]string{filename: tc.config})

			require.NoError(t, rule.Check(runner))
			helper.AssertIssuesWithoutRange(t, tc.issues, runner.Issues)
		})
	}
}

package interfaces_test

import (
	"fmt"
	"testing"

	"github.com/Azure/tflint-ruleset-avm/interfaces"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

const customerManagedKeyV1Fixture = `variable "customer_managed_key" {
  type = object({
    key_vault_resource_id = string
    key_name              = string
    key_version           = optional(string, null)
    user_assigned_identity = optional(object({
      resource_id = string
    }), null)
  })
  default = null

  validation {
    condition     = var.customer_managed_key == null || can(provider::azapi::parse_resource_id("Microsoft.KeyVault/vaults", var.customer_managed_key.key_vault_resource_id))
    error_message = "` + "`customer_managed_key.key_vault_resource_id`" + ` must be a valid Azure Key Vault resource ID."
  }
  validation {
    condition     = var.customer_managed_key == null || var.customer_managed_key.user_assigned_identity == null || can(provider::azapi::parse_resource_id("Microsoft.ManagedIdentity/userAssignedIdentities", var.customer_managed_key.user_assigned_identity.resource_id))
    error_message = "` + "`customer_managed_key.user_assigned_identity.resource_id`" + ` must be a valid user-assigned managed identity resource ID."
  }
}`

const customerManagedKeyV2Fixture = `variable "customer_managed_key" {
  type = object({
    key_vault_key_uri = string
    user_assigned_identity = optional(object({
      client_id = string
    }), null)
  })
  default = null

  validation {
    condition     = var.customer_managed_key == null || can(regex("^https://[^/]+/keys/[^/]+(/[^/]+)?$", var.customer_managed_key.key_vault_key_uri))
    error_message = "` + "`customer_managed_key.key_vault_key_uri`" + ` must be a Key Vault or Managed HSM key URI, in the form ` + "`https://{vaultHost}/keys/{keyName}`" + ` or ` + "`https://{vaultHost}/keys/{keyName}/{keyVersion}`" + `."
  }
  validation {
    condition     = var.customer_managed_key == null || var.customer_managed_key.user_assigned_identity == null || can(regex("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$", var.customer_managed_key.user_assigned_identity.client_id))
    error_message = "` + "`customer_managed_key.user_assigned_identity.client_id`" + ` must be a valid GUID."
  }
}`

// TestCustomerManagedKeyInterface tests both customer managed key interface variants.
func TestCustomerManagedKeyInterface(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		JSON     bool
		Expected helper.Issues
	}{
		{
			Name:     "correct",
			Content:  customerManagedKeyV1Fixture,
			Expected: helper.Issues{},
		},
		{
			Name:     "correct v2",
			Content:  customerManagedKeyV2Fixture,
			Expected: helper.Issues{},
		},
		{
			Name: "incorrect shape",
			Content: `
variable "customer_managed_key" {
	type = object({
		key_vault_key_uri                = string
		user_assigned_identity_client_id = optional(string, null)
	})
	default = null
}`,
			Expected: helper.Issues{
				{
					Rule:    interfaces.NewVarCheckRuleFromAvmInterface(interfaces.CustomerManagedKeyV2),
					Message: fmt.Sprintf("variable type does not comply with the interface specification:\n\n%s", interfaces.CustomerManagedKeyV2TypeString),
				},
			},
		},
	}

	rule := registeredInterfaceRule(t, "customer_managed_key")

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			filename := "variables.tf"
			if tc.JSON {
				filename += ".json"
			}

			runner := helper.TestRunner(t, map[string]string{filename: tc.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssuesWithoutRange(t, tc.Expected, runner.Issues)
		})
	}
}

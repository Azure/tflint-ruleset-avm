package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

const azapiTerraformBlock = `
terraform {
  required_providers {
    azapi = {
      source  = "Azure/azapi"
      version = "~> 2.4"
    }
  }
}
`

func Test_TerraformResourceIDVariableValidation(t *testing.T) {
	rule := NewTerraformResourceIDVariableValidationRule()
	runRuleTestCases(t, rule, []ruleTestCase{
		{
			Name: "spec-compliant validations are allowed",
			Files: map[string]string{
				"terraform.tf": azapiTerraformBlock,
				"variables.tf": `
variable "parent_id" {
  type     = string
  nullable = false

  validation {
    condition     = can(provider::azapi::parse_resource_id("Microsoft.Resources/resourceGroups", var.parent_id))
    error_message = "Must be a resource group ID."
  }
}

variable "workspace_resource_id" {
  type    = string
  default = null

  validation {
    condition     = var.workspace_resource_id == null || can(provider::azapi::parse_resource_id("Microsoft.OperationalInsights/workspaces", var.workspace_resource_id))
    error_message = "Must be a workspace ID."
  }
}

variable "user_assigned_resource_ids" {
  type     = set(string)
  default  = []
  nullable = false

  validation {
    condition = alltrue([
      for id in var.user_assigned_resource_ids :
      can(provider::azapi::parse_resource_id("Microsoft.ManagedIdentity/userAssignedIdentities", id))
    ])
    error_message = "Must be identity IDs."
  }

  validation {
    condition     = length(var.user_assigned_resource_ids) <= 10
    error_message = "At most 10 identities."
  }
}

variable "virtual_wan_resource_id" {
  type    = string
  default = null

  validation {
    condition     = var.virtual_wan_resource_id != null ? can(provider::azapi::parse_resource_id("Microsoft.Network/virtualWans", var.virtual_wan_resource_id).name) : true
    error_message = "Must be a virtual WAN ID."
  }
}`,
			},
			Expected: helper.Issues{},
		},
		{
			Name: "modules without the AzAPI provider are skipped",
			Files: map[string]string{"variables.tf": `
variable "parent_id" {
  type = string
}`},
			Expected: helper.Issues{},
		},
		{
			Name: "unrelated names and types are ignored",
			Files: map[string]string{
				"terraform.tf": azapiTerraformBlock,
				"variables.tf": `
variable "subscription_id" {
  type = string
}

variable "subnet_resource_id" {
  type = object({ id = string })
}

variable "private_dns_zone_resource_ids" {
  type = list(object({ id = string }))
}

variable "untyped_resource_id" {}`,
			},
			Expected: helper.Issues{},
		},
		{
			Name: "the AzAPI provider is detected by source",
			Files: map[string]string{
				"terraform.tf": `
terraform {
  required_providers {
    restapi = {
      source = "registry.terraform.io/azure/azapi"
    }
  }
}`,
				"variables.tf": `
variable "parent_id" {
  type     = string
  nullable = false
}`,
			},
			Expected: helper.Issues{
				{Rule: rule, Message: "Variable `parent_id` holds a resource ID and must be validated with `can(provider::azapi::parse_resource_id(\"<resource type>\", ...))`"},
			},
		},
		{
			Name: "missing and hand-rolled validations are flagged",
			Files: map[string]string{
				"terraform.tf": azapiTerraformBlock,
				"variables.tf": `
variable "parent_id" {
  type     = string
  nullable = false

  validation {
    condition     = can(regex("^/subscriptions/[^/]+/resourceGroups/[^/]+$", var.parent_id))
    error_message = "Must be a resource group ID."
  }
}

variable "key_vault_resource_id" {
  type     = string
  nullable = false

  validation {
    condition     = startswith(var.key_vault_resource_id, "/subscriptions/") && length(var.key_vault_resource_id) > 0
    error_message = "Must be a Key Vault ID."
  }
}

variable "subnet_resource_ids" {
  type     = map(string)
  default  = {}
  nullable = false

  validation {
    condition     = alltrue([for id in var.subnet_resource_ids : length(split("/", id)) == 11])
    error_message = "Must be subnet IDs."
  }
}`,
			},
			Expected: helper.Issues{
				{Rule: rule, Message: "Variable `parent_id` must not validate its resource ID with `regex`; use `provider::azapi::parse_resource_id` instead"},
				{Rule: rule, Message: "Variable `parent_id` holds a resource ID and must be validated with `can(provider::azapi::parse_resource_id(\"<resource type>\", ...))`"},
				{Rule: rule, Message: "Variable `key_vault_resource_id` must not validate its resource ID with `startswith`; use `provider::azapi::parse_resource_id` instead"},
				{Rule: rule, Message: "Variable `key_vault_resource_id` must not validate its resource ID with `length`; use `provider::azapi::parse_resource_id` instead"},
				{Rule: rule, Message: "Variable `key_vault_resource_id` holds a resource ID and must be validated with `can(provider::azapi::parse_resource_id(\"<resource type>\", ...))`"},
				{Rule: rule, Message: "Variable `subnet_resource_ids` must not validate its resource ID with `split`; use `provider::azapi::parse_resource_id` instead"},
				{Rule: rule, Message: "Variable `subnet_resource_ids` holds a resource ID and must be validated with `can(provider::azapi::parse_resource_id(\"<resource type>\", ...))`"},
			},
		},
		{
			Name: "malformed parse_resource_id calls are flagged",
			Files: map[string]string{
				"terraform.tf": azapiTerraformBlock,
				"variables.tf": `
variable "parent_id" {
  type     = string
  nullable = false

  validation {
    condition     = provider::azapi::parse_resource_id(local.parent_type, var.parent_id).resource_group_name != ""
    error_message = "Must be a resource group ID."
  }
}

variable "workspace_resource_id" {
  type    = string
  default = null

  validation {
    condition     = can(provider::azapi::parse_resource_id("Microsoft.${var.namespace}/workspaces", var.workspace_resource_id))
    error_message = "Must be a workspace ID."
  }
}

variable "private_dns_zone_resource_ids" {
  type     = set(string)
  default  = []
  nullable = false

  validation {
    condition     = can(provider::azapi::parse_resource_id("Microsoft.Network/privateDnsZones", var.private_dns_zone_resource_ids))
    error_message = "Must be private DNS zone IDs."
  }
}

variable "settings" {
  type = map(object({ workspace_id = optional(string) }))

  validation {
    condition = alltrue([
      for _, v in var.settings :
      v.workspace_id == null || provider::azapi::parse_resource_id(v.type, v.workspace_id) != null
    ])
    error_message = "Must be workspace IDs."
  }
}`,
			},
			Expected: helper.Issues{
				{Rule: rule, Message: "The resource type passed to `provider::azapi::parse_resource_id` in the validation of variable `parent_id` must be a literal string"},
				{Rule: rule, Message: "`provider::azapi::parse_resource_id` in the validation of variable `parent_id` must be wrapped in `can(...)`"},
				{Rule: rule, Message: "The resource type passed to `provider::azapi::parse_resource_id` in the validation of variable `workspace_resource_id` must be a literal string"},
				{Rule: rule, Message: "The validation of optional variable `workspace_resource_id` must allow `null`, for example `var.workspace_resource_id == null || can(provider::azapi::parse_resource_id(...))`"},
				{Rule: rule, Message: "The validation of variable `private_dns_zone_resource_ids` must check each element with `alltrue([for ... : can(provider::azapi::parse_resource_id(...))])`"},
				{Rule: rule, Message: "The resource type passed to `provider::azapi::parse_resource_id` in the validation of variable `settings` must be a literal string"},
				{Rule: rule, Message: "`provider::azapi::parse_resource_id` in the validation of variable `settings` must be wrapped in `can(...)`"},
			},
		},
	})
}

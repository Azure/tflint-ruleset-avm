package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformDeprecatedVariableFile(t *testing.T) {
	rule := NewTerraformDeprecatedVariableFileRule()
	runRuleTestCases(t, rule, []ruleTestCase{
		{
			Name: "deprecated variables in deprecated_variables.tf are allowed",
			Files: map[string]string{
				"variables.tf": `
variable "network_security_group_enabled" {
  type        = bool
  description = "Whether to generate a network security group."
}`,
				"deprecated_variables.tf": `
variable "create_nsg" {
  type        = bool
  description = "DEPRECATED, use ` + "`network_security_group_enabled`" + ` instead; Whether to generate a network security group."
}

variable "old_name" {
  type        = string
  description = <<DESCRIPTION
DEPRECATED, use ` + "`name`" + ` instead.
DESCRIPTION
}`,
			},
			Expected: helper.Issues{},
		},
		{
			Name: "deprecated variable outside deprecated_variables.tf is flagged",
			Files: map[string]string{"variables.tf": `
variable "create_nsg" {
  type        = bool
  description = "DEPRECATED, use ` + "`network_security_group_enabled`" + ` instead."
}`},
			Expected: helper.Issues{
				{Rule: rule, Message: "Deprecated variable `create_nsg` must be declared in `deprecated_variables.tf`"},
			},
		},
		{
			Name: "variable in deprecated_variables.tf without DEPRECATED is flagged",
			Files: map[string]string{"deprecated_variables.tf": `
variable "no_description" {
  type = string
}

variable "lowercase" {
  type        = string
  description = "Deprecated, use ` + "`name`" + ` instead."
}`},
			Expected: helper.Issues{
				{Rule: rule, Message: "Variable `no_description` in `deprecated_variables.tf` must have a description starting with `DEPRECATED`"},
				{Rule: rule, Message: "Variable `lowercase` in `deprecated_variables.tf` must have a description starting with `DEPRECATED`"},
			},
		},
		{
			Name: "DEPRECATED in the middle of a description is ignored",
			Files: map[string]string{"variables.tf": `
variable "name" {
  type        = string
  description = "The name. Not DEPRECATED."
}`},
			Expected: helper.Issues{},
		},
	})
}

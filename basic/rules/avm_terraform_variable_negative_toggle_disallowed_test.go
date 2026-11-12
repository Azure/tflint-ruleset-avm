package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformVariableNegativeToggle(t *testing.T) {
	rule := NewTerraformVariableNegativeToggleRule()
	runRuleTestCases(t, rule, []ruleTestCase{
		{
			Name: "positive toggles are allowed",
			Files: map[string]string{"variables.tf": `
variable "public_network_access_enabled" {
  type = bool
}

variable "enable_telemetry" {
  type = bool
}`},
			Expected: helper.Issues{},
		},
		{
			Name: "negative bool toggles are flagged",
			Files: map[string]string{"variables.tf": `
variable "public_network_access_disabled" {
  type = bool
}

variable "disable_local_auth" {
  type    = bool
  default = false
}`},
			Expected: helper.Issues{
				{Rule: rule, Message: "Feature switch variable `public_network_access_disabled` should use a positive name such as `xxx_enabled` instead of `xxx_disabled`"},
				{Rule: rule, Message: "Feature switch variable `disable_local_auth` should use a positive name such as `xxx_enabled` instead of `xxx_disabled`"},
			},
		},
		{
			Name: "non-bool or untyped variables are ignored",
			Files: map[string]string{"variables.tf": `
variable "disable_features" {
  type = list(string)
}

variable "rules_disabled" {
  default = []
}`},
			Expected: helper.Issues{},
		},
		{
			Name: "names that only contain disabled are ignored",
			Files: map[string]string{"variables.tf": `
variable "disabled_rules_count" {
  type = bool
}

variable "is_disabled_flag" {
  type = bool
}`},
			Expected: helper.Issues{},
		},
	})
}

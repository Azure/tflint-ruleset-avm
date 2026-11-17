package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformModuleControlVariable(t *testing.T) {
	rule := NewTerraformModuleControlVariableRule()
	runRuleTestCases(t, rule, []ruleTestCase{
		{
			Name: "feature toggles are allowed",
			Files: map[string]string{"variables.tf": `
variable "enable_telemetry" {
  type    = bool
  default = true
}

variable "public_network_access_enabled" {
  type    = bool
  default = false
}`},
			Expected: helper.Issues{},
		},
		{
			Name: "module control variables are flagged",
			Files: map[string]string{"variables.tf": `
variable "enabled" {
  type    = bool
  default = true
}

variable "module_depends_on" {
  type    = list(any)
  default = []
}`},
			Expected: helper.Issues{
				{Rule: rule, Message: "Variable `enabled` must not be declared to control the entire module's operation; use feature-specific toggles instead"},
				{Rule: rule, Message: "Variable `module_depends_on` must not be declared to control the entire module's operation; use feature-specific toggles instead"},
			},
		},
		{
			Name: "outputs with the same name are ignored",
			Files: map[string]string{"outputs.tf": `
output "enabled" {
  value = true
}`},
			Expected: helper.Issues{},
		},
	})
}

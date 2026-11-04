package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformSensitiveFalse(t *testing.T) {
	rule := NewTerraformSensitiveFalseRule()
	runRuleTestCases(t, rule, []ruleTestCase{
		{
			Name: "omitted or true sensitive is allowed",
			Files: map[string]string{"main.tf": `
variable "name" {
  type = string
}

variable "password" {
  type      = string
  sensitive = true
}

output "secret" {
  value     = var.password
  sensitive = true
}`},
			Expected: helper.Issues{},
		},
		{
			Name: "sensitive false on variable and output",
			Files: map[string]string{"main.tf": `
variable "name" {
  type      = string
  sensitive = false
}

output "name" {
  value     = var.name
  sensitive = false
}`},
			Expected: helper.Issues{
				{Rule: rule, Message: "`sensitive = false` must not be declared on variable `name`; values are not sensitive by default"},
				{Rule: rule, Message: "`sensitive = false` must not be declared on output `name`; values are not sensitive by default"},
			},
		},
		{
			Name: "non-literal sensitive is ignored",
			Files: map[string]string{"main.tf": `
output "name" {
  value     = var.name
  sensitive = var.is_sensitive
}`},
			Expected: helper.Issues{},
		},
		{
			Name: "other blocks are ignored",
			Files: map[string]string{"main.tf": `
resource "foo" "bar" {
  sensitive = false
}`},
			Expected: helper.Issues{},
		},
	})
}

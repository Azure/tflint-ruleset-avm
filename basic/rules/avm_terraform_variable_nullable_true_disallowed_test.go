package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformVariableNullableTrue(t *testing.T) {
	rule := NewTerraformVariableNullableTrueRule()
	runRuleTestCases(t, rule, []ruleTestCase{
		{
			Name: "omitted or false nullable is allowed",
			Files: map[string]string{"variables.tf": `
variable "name" {
  type = string
}

variable "tags" {
  type     = map(string)
  default  = {}
  nullable = false
}`},
			Expected: helper.Issues{},
		},
		{
			Name: "nullable true is flagged",
			Files: map[string]string{"variables.tf": `
variable "name" {
  type     = string
  nullable = true
}`},
			Expected: helper.Issues{
				{Rule: rule, Message: "`nullable = true` must not be declared on variable `name`"},
			},
		},
		{
			Name: "non-literal nullable is ignored",
			Files: map[string]string{"variables.tf": `
variable "name" {
  type     = string
  nullable = local.nullable
}`},
			Expected: helper.Issues{},
		},
	})
}

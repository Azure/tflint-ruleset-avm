package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformDeprecatedOutputFile(t *testing.T) {
	rule := NewTerraformDeprecatedOutputFileRule()
	runRuleTestCases(t, rule, []ruleTestCase{
		{
			Name: "deprecated outputs in deprecated_outputs.tf are allowed",
			Files: map[string]string{
				"outputs.tf": `
output "resource_id" {
  value       = "id"
  description = "The resource ID."
}`,
				"deprecated_outputs.tf": `
output "id" {
  value       = "id"
  description = "DEPRECATED, use ` + "`resource_id`" + ` instead."
}`,
			},
			Expected: helper.Issues{},
		},
		{
			Name: "deprecated output outside deprecated_outputs.tf is flagged",
			Files: map[string]string{"outputs.tf": `
output "id" {
  value       = "id"
  description = "DEPRECATED, use ` + "`resource_id`" + ` instead."
}`},
			Expected: helper.Issues{
				{Rule: rule, Message: "Deprecated output `id` must be declared in `deprecated_outputs.tf`"},
			},
		},
		{
			Name: "output without description is ignored",
			Files: map[string]string{"outputs.tf": `
output "id" {
  value = "id"
}`},
			Expected: helper.Issues{},
		},
	})
}

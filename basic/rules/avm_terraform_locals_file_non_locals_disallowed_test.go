package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformLocalsFile(t *testing.T) {
	rule := NewTerraformLocalsFileRule()
	runRuleTestCases(t, rule, []ruleTestCase{
		{
			Name: "multiple locals blocks are allowed",
			Files: map[string]string{"locals.tf": `
locals {
  a = 1
}

locals {
  b = 2
}`},
			Expected: helper.Issues{},
		},
		{
			Name: "non-locals blocks in locals.tf are flagged",
			Files: map[string]string{"locals.tf": `
locals {
  a = 1
}

resource "foo" "bar" {}

variable "baz" {
  type = string
}`},
			Expected: helper.Issues{
				{Rule: rule, Message: "`locals.tf` must only contain `locals` blocks, found `resource` block"},
				{Rule: rule, Message: "`locals.tf` must only contain `locals` blocks, found `variable` block"},
			},
		},
		{
			Name: "locals next to resources in other files are allowed",
			Files: map[string]string{"main.tf": `
locals {
  a = 1
}

resource "foo" "bar" {}`},
			Expected: helper.Issues{},
		},
		{
			Name: "locals.tf in a subdirectory path is checked",
			Files: map[string]string{"modules/child/locals.tf": `
data "foo" "bar" {}`},
			Expected: helper.Issues{
				{Rule: rule, Message: "`locals.tf` must only contain `locals` blocks, found `data` block"},
			},
		},
	})
}

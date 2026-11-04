package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformCollectionVariableNullable(t *testing.T) {
	rule := NewTerraformCollectionVariableNullableRule()
	runRuleTestCases(t, rule, []ruleTestCase{
		{
			Name: "collections with nullable false are allowed",
			Files: map[string]string{"variables.tf": `
variable "names" {
  type     = list(string)
  default  = []
  nullable = false
}

variable "subnets" {
  type     = map(object({ name = string }))
  default  = {}
  nullable = false
}`},
			Expected: helper.Issues{},
		},
		{
			Name: "scalars, objects, and untyped variables are ignored",
			Files: map[string]string{"variables.tf": `
variable "name" {
  type = string
}

variable "lock" {
  type = object({
    kind = string
  })
  default = null
}

variable "untyped" {}`},
			Expected: helper.Issues{},
		},
		{
			Name: "tags keeps its nullable map(string) interface",
			Files: map[string]string{"variables.tf": `
variable "tags" {
  type    = map(string)
  default = null
}`},
			Expected: helper.Issues{},
		},
		{
			Name: "collections without nullable are flagged",
			Files: map[string]string{"variables.tf": `
variable "names" {
  type    = list(string)
  default = []
}

variable "ids" {
  type = set(string)
}

variable "subnets" {
  type    = map(object({ name = string }))
  default = {}
}

variable "pair" {
  type = tuple([string, number])
}

variable "legacy" {
  type = map
}`},
			Expected: helper.Issues{
				{Rule: rule, Message: "Collection variable `names` must set `nullable = false`"},
				{Rule: rule, Message: "Collection variable `ids` must set `nullable = false`"},
				{Rule: rule, Message: "Collection variable `subnets` must set `nullable = false`"},
				{Rule: rule, Message: "Collection variable `pair` must set `nullable = false`"},
				{Rule: rule, Message: "Collection variable `legacy` must set `nullable = false`"},
			},
		},
		{
			Name: "explicit nullable values are left to other rules",
			Files: map[string]string{"variables.tf": `
variable "names" {
  type     = list(string)
  nullable = true
}

variable "ids" {
  type     = set(string)
  nullable = local.nullable
}`},
			Expected: helper.Issues{},
		},
	})
}

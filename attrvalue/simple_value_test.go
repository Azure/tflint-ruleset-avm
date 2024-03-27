package attrvalue_test

import (
	"testing"

	"github.com/Azure/tflint-ruleset-avm/attrvalue"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func TestSimpleValueRule(t *testing.T) {
	testCases := []struct {
		name     string
		rule     tflint.Rule
		content  string
		expected helper.Issues
	}{
		{
			name: "correct string",
			rule: attrvalue.NewSimpleRule("foo", "bar", []string{"bar", "bat"}),
			content: `
	variable "test" {
		type    = string
		default = "bat"
	}
	resource "foo" "example" {
		bar = var.test
	}`,
			expected: helper.Issues{},
		},
		{
			name: "incorrect string",
			rule: attrvalue.NewSimpleRule("foo", "bar", []string{"bar", "bat"}),
			content: `
	variable "test" {
		type    = string
		default = "fiz"
	}
	resource "foo" "example" {
		bar = var.test
	}`,
			expected: helper.Issues{
				{
					Rule:    attrvalue.NewSimpleRule("foo", "bar", []string{"bar", "bat"}),
					Message: "fiz is an invalid attribute value of `bar` - expecting (one of) [bar bat]",
				},
			},
		},
		{
			name: "correct number",
			rule: attrvalue.NewSimpleRule("foo", "bar", []int{1, 2}),
			content: `
	variable "test" {
		type    = number
		default = 2
	}
	resource "foo" "example" {
		bar = var.test
	}`,
			expected: helper.Issues{},
		},
		{
			name: "correct number float",
			rule: attrvalue.NewSimpleRule("foo", "bar", []float64{1.2, 2.1}),
			content: `
variable "test" {
	type    = number
	default = 2.1
}
resource "foo" "example" {
	bar = var.test
}`,
			expected: helper.Issues{},
		},
		{
			name: "incorrect number",
			rule: attrvalue.NewSimpleRule("foo", "bar", []int{1, 2}),
			content: `
	variable "test" {
		type    = number
		default = 3
	}
	resource "foo" "example" {
		bar = var.test
	}`,
			expected: helper.Issues{
				{
					Rule:    attrvalue.NewSimpleRule("foo", "bar", []int{1, 2}),
					Message: "3 is an invalid attribute value of `bar` - expecting (one of) [1 2]",
				},
			},
		},
		{
			name: "incorrect number float",
			rule: attrvalue.NewSimpleRule("foo", "bar", []float64{1.1, 2.2}),
			content: `
	variable "test" {
		type    = number
		default = 2.1
	}
	resource "foo" "example" {
		bar = var.test
	}`,
			expected: helper.Issues{
				{
					Rule:    attrvalue.NewSimpleRule("foo", "bar", []float64{1.1, 2.2}),
					Message: "2.1 is an invalid attribute value of `bar` - expecting (one of) [1.1 2.2]",
				},
			},
		},
		{
			name: "correct bool",
			rule: attrvalue.NewSimpleRule("foo", "bar", []bool{true}),
			content: `
	variable "test" {
		type    = bool
		default = true
	}
	resource "foo" "example" {
		bar = var.test
	}`,
			expected: helper.Issues{},
		},
		{
			name: "incorrect bool",
			rule: attrvalue.NewSimpleRule("foo", "bar", []bool{true}),
			content: `
	variable "test" {
		type    = bool
		default = false
	}
	resource "foo" "example" {
		bar = var.test
	}`,
			expected: helper.Issues{
				{
					Rule:    attrvalue.NewSimpleRule("foo", "bar", []bool{true}),
					Message: "false is an invalid attribute value of `bar` - expecting (one of) [true]",
				},
			},
		},
		{
			name: "null value",
			rule: attrvalue.NewSimpleRule("foo", "bar", []bool{true}),
			content: `
	variable "test" {
		type    = bool
		default = null
	}
	resource "foo" "example" {
		bar = var.test
	}`,
			expected: helper.Issues{},
		},
		{
			name: "optional value",
			rule: attrvalue.NewSimpleRule("foo", "bar", []bool{true}),
			content: `
	variable "test" {
		type    = object({
			optional = optional(bool)
		})
		default = {}
	}
	resource "foo" "example" {
		bar = var.test.optional
	}`,
			expected: helper.Issues{},
		},
	}

	filename := "main.tf"
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.name, func(t *testing.T) {
			t.Parallel()
			runner := helper.TestRunner(t, map[string]string{filename: tC.content})
			if err := tC.rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			helper.AssertIssuesWithoutRange(t, tC.expected, runner.Issues)
		})
	}
}

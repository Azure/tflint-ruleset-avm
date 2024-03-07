package attrvalue_test

import (
	"testing"

	"github.com/Azure/tflint-ruleset-avm/attrvalue"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func TestStringValueRule(t *testing.T) {
	testCases := []struct {
		name     string
		rule     tflint.Rule
		content  string
		expected helper.Issues
	}{
		{
			name: "incorrect",
			rule: attrvalue.NewStringRule("foo", "bar", []string{"baz", "bat"}),
			content: `
	resource "foo" "example" {
		bar = "fiz"
	}`,
			expected: helper.Issues{
				{
					Rule:    attrvalue.NewStringRule("foo", "bar", []string{"baz"}),
					Message: "\"fiz\" is an invalid attribute value of `bar` - expecting (one of) [baz bat]",
				},
			},
		},
		{
			name: "correct",
			rule: attrvalue.NewStringRule("foo", "bar", []string{"baz", "bat"}),
			content: `
	resource "foo" "example" {
		bar = "baz"
	}`,
			expected: helper.Issues{},
		},
		{
			name: "correct second value",
			rule: attrvalue.NewStringRule("foo", "bar", []string{"baz", "bat"}),
			content: `
	resource "foo" "example" {
		bar = "bat"
	}`,
			expected: helper.Issues{},
		},
		{
			name: "not applicable",
			rule: attrvalue.NewStringRule("foo", "bar", []string{"baz", "bat"}),
			content: `
	resource "zzz" "example" {
		bar = "fiz"
	}`,
			expected: helper.Issues{},
		},
		{
			name: "traversal",
			rule: attrvalue.NewStringRule("foo", "bar", []string{"baz", "bat"}),
			content: `
	variable "test" {
		type = string
		default = "baz"
	}
	resource "foo" "example" {
		bar = var.test
	}`,
			expected: helper.Issues{},
		},
		{
			name: "incorrect value type",
			rule: attrvalue.NewStringRule("foo", "bar", []string{"baz", "bat"}),
			content: `
	variable "test" {
		type = number
		default = null
	}
	resource "foo" "example" {
		bar = var.test
	}`,
			expected: helper.Issues{},
		},
		// This will need to be implemented in an integration test due to the use of the optional function
		// not being supported in helper.TestRunner
		// 	{
		// 		name: "traversal with optional func",
		// 		rule: attrvalue.NewStringRule("foo", "bar", []string{"baz", "bat"}),
		// 		content: `
		// variable "test" {
		// 	type = object({
		// 		value = optional(string, "baz")
		// 	})
		// 	default = {}
		// }
		// resource "foo" "example" {
		// 	bar = var.test.value
		// }`,
		// 		expected: helper.Issues{},
		// 	},
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

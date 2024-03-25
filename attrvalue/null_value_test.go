package attrvalue_test

import (
	"testing"

	"github.com/Azure/tflint-ruleset-avm/attrvalue"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func TestNullRule(t *testing.T) {
	testCases := []struct {
		name     string
		rule     tflint.Rule
		content  string
		expected helper.Issues
	}{
		{
			name: "correct",
			rule: attrvalue.NewNullRule("foo", "bar"),
			content: `
	variable "test" {
		type    = string
		default = null
	}
	resource "foo" "example" {
		bar = var.test
	}`,
			expected: helper.Issues{},
		},
		{
			name: "undefined",
			rule: attrvalue.NewNullRule("foo", "bar"),
			content: `
	variable "test" {
		type    = string
	}
	resource "foo" "example" {
		bar = var.test
	}`,
			expected: helper.Issues{},
		},
		{
			name: "not null (string)",
			rule: attrvalue.NewNullRule("foo", "bar"),
			content: `
	variable "test" {
		type    = string
		default = "something"
	}
	resource "foo" "example" {
		bar = var.test
	}`,
			expected: helper.Issues{
				{
					Rule:    attrvalue.NewNullRule("foo", "bar"),
					Message: "invalid attribute value of `bar` - expecting null",
				},
			},
		},
		{
			name: "not null (number)",
			rule: attrvalue.NewNullRule("foo", "bar"),
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
					Rule:    attrvalue.NewNullRule("foo", "bar"),
					Message: "invalid attribute value of `bar` - expecting null",
				},
			},
		},
		{
			name: "not null (bool)",
			rule: attrvalue.NewNullRule("foo", "bar"),
			content: `
	variable "test" {
		type    = bool
		default = true
	}
	resource "foo" "example" {
		bar = var.test
	}`,
			expected: helper.Issues{
				{
					Rule:    attrvalue.NewNullRule("foo", "bar"),
					Message: "invalid attribute value of `bar` - expecting null",
				},
			},
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

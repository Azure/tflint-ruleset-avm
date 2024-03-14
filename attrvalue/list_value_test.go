package attrvalue_test

import (
	"testing"

	"github.com/Azure/tflint-ruleset-avm/attrvalue"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func TestListNumberValueRule(t *testing.T) {
	testCases := []struct {
		name     string
		rule     tflint.Rule
		content  string
		expected helper.Issues
	}{
		{
			name: "incorrect",
			rule: attrvalue.NewListRule("foo", "bar", [][]int{{1, 2, 3}}),
			content: `
	variable "test" {
		type    = list(number)
		default = [3]
	}
	resource "foo" "example" {
		bar = var.test
	}`,
			expected: helper.Issues{
				{
					Rule:    attrvalue.NewListRule("foo", "bar", [][]int{{1, 2, 3}}),
					Message: "\"&[3]\" is an invalid attribute value of `bar` - expecting (one of) [[1 2 3]]",
				},
			},
		},
		{
			name: "correct",
			rule: attrvalue.NewListRule("foo", "bar", [][]int{{1, 2, 3}}),
			content: `
	variable "test" {
		type    = list(number)
		default = [1, 2, 3]
	}
	resource "foo" "example" {
		bar = var.test
	}`,
			expected: helper.Issues{},
		},
		{
			name: "correct but wrong order",
			rule: attrvalue.NewListRule("foo", "bar", [][]int{{1, 2, 3}}),
			content: `
	variable "test" {
		type    = list(number)
		default = [2, 3, 1]
	}
	resource "foo" "example" {
		bar = var.test
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

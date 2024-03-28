package attrvalue_test

import (
	"testing"

	"github.com/prashantv/gostub"

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
			rule: attrvalue.NewSetRule("foo", "bar", [][]int{{1, 2, 3}}, ""),
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
					Rule:    attrvalue.NewSetRule("foo", "bar", [][]int{{1, 2, 3}}, ""),
					Message: "\"[3]\" is an invalid attribute value of `bar` - expecting (one of) [[1 2 3]]",
				},
			},
		},
		{
			name: "correct",
			rule: attrvalue.NewSetRule("foo", "bar", [][]int{{1, 2, 3}}, ""),
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
			name: "correct with string list",
			rule: attrvalue.NewSetRule("foo", "bar", [][]string{{"1", "2", "3"}}, ""),
			content: `
	variable "test" {
		type    = list(string)
		default = ["1", "2", "3"]
	}
	resource "foo" "example" {
		bar = var.test
	}`,
			expected: helper.Issues{},
		},
		{
			name: "correct but different order",
			rule: attrvalue.NewSetRule("foo", "bar", [][]int{{1, 2, 3}}, ""),
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
		{
			name: "variable without default",
			rule: attrvalue.NewSetRule("foo", "bar", [][]int{{1, 2, 3}}, ""),
			content: `
		variable "test" {
			type    = list(number)
		}
		resource "foo" "example" {
			bar = var.test
		}`,
			expected: helper.Issues{},
		},
		{
			name: "variable with convertable element type",
			rule: attrvalue.NewSetRule("foo", "bar", [][]int{{1, 2, 3}}, ""),
			content: `
	variable "test" {
		type    = list(string)
		default = ["1", "2", "3"]
	}
	resource "foo" "example" {
		bar = var.test
	}`,
			expected: helper.Issues{},
		},
	}

	filename := "main.tf"
	for _, c := range testCases {
		tc := c
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{filename: tc.content})
			stub := gostub.Stub(&attrvalue.AppFs, mockFs(tc.content))
			defer stub.Reset()
			if err := tc.rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			helper.AssertIssuesWithoutRange(t, tc.expected, runner.Issues)
		})
	}
}

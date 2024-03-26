package attrvalue_test

import (
	"testing"

	"github.com/Azure/tflint-ruleset-avm/attrvalue"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func TestNestedBlockValueRule(t *testing.T) {
	testCases := []struct {
		name     string
		rule     tflint.Rule
		content  string
		expected helper.Issues
	}{
		{
			name: "correct string",
			rule: attrvalue.NewSimpleNestedBlockRule("foo", "fiz", "bar", []string{"biz", "bat"}),
			content: `
	variable "test" {
		type    = string
		default = "biz"
	}
	resource "foo" "example" {
		fiz {
			bar = var.test
		}
	}`,
			expected: helper.Issues{},
		},
		{
			name: "incorrect string",
			rule: attrvalue.NewSimpleNestedBlockRule("foo", "fiz", "bar", []string{"biz", "bat"}),
			content: `
	variable "test" {
		type    = string
		default = "baz"
	}
	resource "foo" "example" {
		fiz {
			bar = var.test
		}
	}`,
			expected: helper.Issues{
				{
					Rule:    attrvalue.NewSimpleNestedBlockRule("foo", "fiz", "bar", []string{"biz", "bat"}),
					Message: "baz is an invalid attribute value of `bar` - expecting (one of) [biz bat]",
				},
			},
		},
		{
			name: "no nested block of that type",
			rule: attrvalue.NewSimpleNestedBlockRule("foo", "fiz", "bar", []string{"biz", "bat"}),
			content: `
	variable "test" {
		type    = string
		default = "baz"
	}
	resource "foo" "example" {
		fuz {
			bar = var.test
		}
	}`,
			expected: helper.Issues{},
		},
		{
			name: "multiple blocks correct",
			rule: attrvalue.NewSimpleNestedBlockRule("foo", "fiz", "bar", []string{"biz", "bat"}),
			content: `
	variable "test" {
		type    = string
		default = "biz"
	}
	resource "foo" "example" {
		fiz {
			bar = var.test
		}
		fiz {
			bar = var.test
		}
	}`,
			expected: helper.Issues{},
		},
		{
			name: "multiple blocks partially correct",
			rule: attrvalue.NewSimpleNestedBlockRule("foo", "fiz", "bar", []string{"biz", "bat"}),
			content: `
	variable "test" {
		type    = string
		default = "biz"
	}
	variable "test2" {
		type    = string
		default = "incorrect"
	}
	resource "foo" "example" {
		fiz {
			bar = var.test
		}
		fiz {
			bar = var.test2
		}
	}`,
			expected: helper.Issues{
				{
					Rule:    attrvalue.NewSimpleNestedBlockRule("foo", "fiz", "bar", []string{"biz", "bat"}),
					Message: "incorrect is an invalid attribute value of `bar` - expecting (one of) [biz bat]",
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

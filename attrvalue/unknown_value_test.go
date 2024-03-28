package attrvalue_test

import (
	"testing"

	"github.com/prashantv/gostub"

	"github.com/Azure/tflint-ruleset-avm/attrvalue"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func TestUnknownValueRule(t *testing.T) {
	testCases := []struct {
		name     string
		rule     tflint.Rule
		content  string
		expected helper.Issues
	}{
		{
			name: "unknown string value",
			rule: attrvalue.NewUnknownValueRule("foo", "bar", ""),
			content: `
	variable "test" {
		type    = string
	}
	resource "foo" "example" {
		foo = "name"
		bar = var.test
	}`,
			expected: helper.Issues{},
		},
		{
			name: "null value with local",
			rule: attrvalue.NewUnknownValueRule("foo", "bar", ""),
			content: `
		locals {
		 bar = null
		}
		resource "foo" "example" {
			foo = "name"
			bar = local.bar
		}`,
			expected: helper.Issues{
				{
					Rule:    attrvalue.NewUnknownValueRule("foo", "bar", ""),
					Message: "invalid attribute value of `bar` - expecting unknown",
				},
			},
		},
		{
			name: "default value",
			rule: attrvalue.NewUnknownValueRule("foo", "bar", ""),
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
					Rule:    attrvalue.NewUnknownValueRule("foo", "bar", ""),
					Message: "invalid attribute value of `bar` - expecting unknown",
				},
			},
		},
		{
			name: "unknown number value",
			rule: attrvalue.NewUnknownValueRule("foo", "bar", ""),
			content: `
		variable "test" {
			type    = number
		}
		resource "foo" "example" {
			bar = var.test
		}`,
			expected: helper.Issues{},
		},
		{
			name: "not null (number)",
			rule: attrvalue.NewUnknownValueRule("foo", "bar", ""),
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
					Rule:    attrvalue.NewUnknownValueRule("foo", "bar", ""),
					Message: "invalid attribute value of `bar` - expecting unknown",
				},
			},
		},
		{
			name: "attribute not found",
			rule: attrvalue.NewUnknownValueRule("foo", "bar", ""),
			content: `
		variable "test" {
			type    = bool
			default = true
		}
		resource "foo" "example" {
		}`,
			expected: helper.Issues{},
		},
		{
			name:     "empty config",
			rule:     attrvalue.NewUnknownValueRule("foo", "bar", ""),
			content:  ``,
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

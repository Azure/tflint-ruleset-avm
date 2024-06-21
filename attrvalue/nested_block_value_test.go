package attrvalue_test

import (
	"os"
	"testing"

	"github.com/prashantv/gostub"
	"github.com/spf13/afero"

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
			rule: attrvalue.NewSimpleNestedBlockRule("foo", "fiz", "bar", []string{"biz", "bat"}, "", false, ""),
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
			rule: attrvalue.NewSimpleNestedBlockRule("foo", "fiz", "bar", []string{"biz", "bat"}, "", false, ""),
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
					Rule:    attrvalue.NewSimpleNestedBlockRule("foo", "fiz", "bar", []string{"biz", "bat"}, "", false, ""),
					Message: "baz is an invalid attribute value of `bar` - expecting (one of) [biz bat]",
				},
			},
		},
		{
			name: "incorrect resource with correct block",
			rule: attrvalue.NewSimpleNestedBlockRule("foo", "fiz", "bar", []string{"biz", "bat"}, "", false, ""),
			content: `
	variable "test" {
		type    = string
		default = "baz"
	}
	resource "fuz" "example" {
		fiz {
			bar = var.test
		}
	}`,
			expected: helper.Issues{},
		},
		{
			name: "no nested block of that type",
			rule: attrvalue.NewSimpleNestedBlockRule("foo", "fiz", "bar", []string{"biz", "bat"}, "", false, ""),
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
			name: "no nested block of that type with must exist set",
			rule: attrvalue.NewSimpleNestedBlockRule("foo", "fiz", "bar", []string{"biz", "bat"}, "", true, ""),
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
			expected: helper.Issues{
				{
					Rule:    attrvalue.NewSimpleNestedBlockRule("foo", "fiz", "bar", []string{"biz", "bat"}, "", true, ""),
					Message: "The attribute `bar` must be specified",
				},
			},
		},
		{
			name: "no nested block attribute of that type with must exist set",
			rule: attrvalue.NewSimpleNestedBlockRule("foo", "fiz", "bar", []string{"biz", "bat"}, "", true, ""),
			content: `
	variable "test" {
		type    = string
		default = "baz"
	}
	resource "foo" "example" {
		fiz {
			bor = var.test
		}
	}`,
			expected: helper.Issues{
				{
					Rule:    attrvalue.NewSimpleNestedBlockRule("foo", "fiz", "bar", []string{"biz", "bat"}, "", true, ""),
					Message: "The attribute `bar` must be specified",
				},
			},
		},
		{
			name: "multiple blocks correct",
			rule: attrvalue.NewSimpleNestedBlockRule("foo", "fiz", "bar", []string{"biz", "bat"}, "", false, ""),
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
			rule: attrvalue.NewSimpleNestedBlockRule("foo", "fiz", "bar", []string{"biz", "bat"}, "", false, ""),
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
					Rule:    attrvalue.NewSimpleNestedBlockRule("foo", "fiz", "bar", []string{"biz", "bat"}, "", false, ""),
					Message: "incorrect is an invalid attribute value of `bar` - expecting (one of) [biz bat]",
				},
			},
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

func mockFs(c string) afero.Afero {
	fs := afero.NewMemMapFs()
	_ = afero.WriteFile(fs, "main.tf", []byte(c), os.ModePerm)
	return afero.Afero{Fs: fs}
}

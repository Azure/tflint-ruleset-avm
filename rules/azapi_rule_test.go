// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rules

import (
	"os"
	"testing"

	"github.com/Azure/tflint-ruleset-avm/blockquery"
	"github.com/Azure/tflint-ruleset-avm/modulecontent"
	"github.com/prashantv/gostub"
	"github.com/spf13/afero"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

func TestAzapiRule(t *testing.T) {
	testCases := []struct {
		name     string
		rule     tflint.Rule
		content  string
		expected helper.Issues
	}{
		{
			name: "correct string",
			rule: NewAzApiRuleQueryMustExist("test", "https://example.com", "testType", "", "", "foo", blockquery.IsOneOf, blockquery.NewStringResults("fiz")...),
			content: `
		resource "azapi_resource" "test" {
		  type = "testType@0000-00-00"
		  body = {
			  foo = "fiz"
				bar = "biz"
			}
		}`,
			expected: helper.Issues{},
		},
		{
			name: "correct string with multiple values",
			rule: NewAzApiRuleQueryMustExist("test", "https://example.com", "testType", "", "", "foo", blockquery.IsOneOf, blockquery.NewStringResults("fuz", "fiz")...),
			content: `
		resource "azapi_resource" "test" {
		  type = "testType@0000-00-00"
		  body = {
			  foo = "fiz"
				bar = "biz"
			}
		}`,
			expected: helper.Issues{},
		},
		{
			name: "incorrect string",
			rule: NewAzApiRuleQueryMustExist("test", "https://example.com", "testType", "", "", "foo", blockquery.IsOneOf, blockquery.NewStringResults("not_fiz")...),
			content: `
		resource "azapi_resource" "test" {
		  type = "testType@0000-00-00"
		  body = {
			  foo = "fiz"
				bar = "biz"
			}
		}`,
			expected: helper.Issues{
				{
					Rule:    NewAzApiRuleQueryMustExist("test", "https://example.com", "testType", "", "", "foo", blockquery.IsOneOf, blockquery.NewStringResults("not_fiz")...),
					Message: "returned value `fiz` not in expected values `[not_fiz]`",
				},
			},
		},
		{
			name: "string not present but doesn't need to exist",
			rule: NewAzApiRuleQueryOptionalExist("test", "https://example.com", "testType", "", "", "bat", blockquery.IsOneOf, blockquery.NewStringResults()...),
			content: `
		resource "azapi_resource" "test" {
		  type = "testType@0000-00-00"
		  body = {
			  foo = "fiz"
				bar = "biz"
			}
		}`,
			expected: helper.Issues{},
		},
		{
			name: "string present but not expected",
			rule: NewAzApiRuleQueryOptionalExist("test", "https://example.com", "testType", "", "", "foo", blockquery.IsNull),
			content: `
		resource "azapi_resource" "test" {
		  type = "testType@0000-00-00"
		  body = {
			  foo = "fiz"
				bar = "biz"
			}
		}`,
			expected: helper.Issues{
				{
					Rule:    NewAzApiRuleQueryMustExist("test", "https://example.com", "testType", "", "", "foo", blockquery.EachIsOneOf, blockquery.NewStringResults()...),
					Message: "returned value is not null but expected to be",
				},
			},
		},
		{
			name: "string not present but not expected",
			rule: NewAzApiRuleQueryOptionalExist("test", "https://example.com", "testType", "", "", "notExist", blockquery.IsNull),
			content: `
		resource "azapi_resource" "test" {
		  type = "testType@0000-00-00"
		  body = {
			  foo = "fiz"
				bar = "biz"
			}
		}`,
			expected: helper.Issues{},
		},
		{
			name: "correct number",
			rule: NewAzApiRuleQueryMustExist("test", "https://example.com", "testType", "", "", "foo", blockquery.IsOneOf, blockquery.NewIntResults(2)...),
			content: `
		resource "azapi_resource" "test" {
		  type = "testType@0000-00-00"
		  body = {
			  foo = 2
				bar = "biz"
			}
		}`,
			expected: helper.Issues{},
		},
		{
			name: "incorrect number",
			rule: NewAzApiRuleQueryMustExist("test", "https://example.com", "testType", "", "", "foo", blockquery.IsOneOf, blockquery.NewIntResults(0)...),
			content: `
		resource "azapi_resource" "test" {
		  type = "testType@0000-00-00"
		  body = {
			  foo = 2
				bar = "biz"
			}
		}`,
			expected: helper.Issues{
				{
					Rule:    NewAzApiRuleQueryMustExist("test", "https://example.com", "testType", "", "", "foo", blockquery.IsOneOf, blockquery.NewIntResults(0)...),
					Message: "returned value `2` not in expected values `[0]`",
				},
			},
		},
		{
			name: "correct bool",
			rule: NewAzApiRuleQueryMustExist("test", "https://example.com", "testType", "", "", "foo", blockquery.IsOneOf, blockquery.NewBoolResult(true)),
			content: `
		resource "azapi_resource" "test" {
		  type = "testType@0000-00-00"
		  body = {
			  foo = true
				bar = "biz"
			}
		}`,
			expected: helper.Issues{},
		},
		{
			name: "incorrect bool",
			rule: NewAzApiRuleQueryMustExist("test", "https://example.com", "testType", "", "", "foo", blockquery.IsOneOf, blockquery.NewBoolResult(true)),
			content: `
		resource "azapi_resource" "test" {
		  type = "testType@0000-00-00"
		  body = {
			  foo = false
				bar = "biz"
			}
		}`,
			expected: helper.Issues{
				{
					Rule:    NewAzApiRuleQueryMustExist("test", "https://example.com", "testType", "", "", "foo", blockquery.EachIsOneOf, blockquery.NewBoolResult(true)),
					Message: "returned value `false` not in expected values `[true]`",
				},
			},
		},
		{
			name: "correct list",
			rule: NewAzApiRuleQueryMustExist("test", "https://example.com", "testType", "", "", "foo", blockquery.IsOneOf, blockquery.NewListResults([]cty.Value{cty.NumberIntVal(1), cty.NumberIntVal(2), cty.NumberIntVal(3)})...),
			content: `
		resource "azapi_resource" "test" {
		  type = "testType@0000-00-00"
		  body = {
			  foo = [1, 2, 3]
				bar = "biz"
			}
		}`,
			expected: helper.Issues{},
		},
		{
			name: "incorrect list",
			rule: NewAzApiRuleQueryMustExist("test", "https://example.com", "testType", "", "", "foo", blockquery.IsOneOf, blockquery.NewListResults([]cty.Value{cty.NumberIntVal(4), cty.NumberIntVal(5), cty.NumberIntVal(6)})...),
			content: `
		resource "azapi_resource" "test" {
		  type = "testType@0000-00-00"
		  body = {
			  foo = [1, 2, 3]
				bar = "biz"
			}
		}`,
			expected: helper.Issues{
				{
					Rule:    NewAzApiRuleQueryMustExist("test", "https://example.com", "testType", "", "", "foo", blockquery.EachIsOneOf, blockquery.NewListResults([]cty.Value{cty.NumberIntVal(4), cty.NumberIntVal(5), cty.NumberIntVal(6)})...),
					Message: "returned value `[1 2 3]` not in expected values `[[4 5 6]]`",
				},
			},
		},
		{
			name: "nested list",
			rule: NewAzApiRuleQueryMustExist("test", "https://example.com", "testType", "", "", "foo.#.bar", blockquery.EachIsOneOf, blockquery.NewListResults([]cty.Value{cty.NumberIntVal(1), cty.NumberIntVal(2), cty.NumberIntVal(3)})...),
			content: `
resource "azapi_resource" "test" {
	type = "testType@0000-00-00"
	body = {
			foo = [
				{
					bar = [1, 2, 3]
				},
				{
					bar = [1, 2, 3]
				},
				{
					bar = [1, 2, 3]
				}
			]
	}
}`,
			expected: helper.Issues{},
		},
		{
			name: "nested list incorrect",
			rule: NewAzApiRuleQueryMustExist("test", "https://example.com", "testType", "", "", "foo.#.bar", blockquery.EachIsOneOf, blockquery.NewListResults([]cty.Value{cty.NumberIntVal(1), cty.NumberIntVal(2), cty.NumberIntVal(3)})...),
			content: `
resource "azapi_resource" "test" {
type = "testType@0000-00-00"
body = {
		foo = [
			{
				bar = [1, 2, 3]
			},
			{
				bar = [4, 5, 6]
			},
			{
				bar = [1, 2, 3]
			}
		]
}
}`,
			expected: helper.Issues{
				{
					Rule:    NewAzApiRuleQueryMustExist("test", "https://example.com", "testType", "", "", "foo.#.bar", blockquery.EachIsOneOf, blockquery.NewListResults([]cty.Value{cty.NumberIntVal(1), cty.NumberIntVal(2), cty.NumberIntVal(3)})...),
					Message: "returned values `[[1 2 3] [4 5 6] [1 2 3]]` not in expected values `[[1 2 3]]`",
				},
			},
		},
		{
			name: "query return no results but does not need to exist",
			rule: NewAzApiRuleQueryOptionalExist("test", "https://example.com", "testType", "", "", "notexist", blockquery.EachIsOneOf, blockquery.NewStringResults("fiz")...),
			content: `
		resource "azapi_resource" "test" {
		  type = "testType@0000-00-00"
		  body = {
			  foo = "fiz"
				bar = "biz"
			}
		}`,
			expected: helper.Issues{},
		},
		{
			name: "query return no results and need to exist",
			rule: NewAzApiRuleQueryMustExist("test", "https://example.com", "testType", "", "", "notexist", blockquery.EachIsOneOf, blockquery.NewStringResults("fiz")...),
			content: `
		resource "azapi_resource" "test" {
		  type = "testType@0000-00-00"
		  body = {
			  foo = "fiz"
				bar = "biz"
			}
		}`,
			expected: helper.Issues{
				{
					Rule:    NewAzApiRuleQueryMustExist("test", "https://example.com", "testType", "", "", "notexist", blockquery.EachIsOneOf, blockquery.NewStringResults("fiz")...),
					Message: "attribute not found: notexist",
				},
			},
		},
		{
			name: "no azapi_resources - no error expected",
			rule: NewAzApiRuleQueryMustExist("test", "https://example.com", "testType", "", "", "query", blockquery.IsNotNull),
			content: `
		resource "not_azapi_resource" "test" {
			biz = "baz"
			buz = "fuz"
		}`,
			expected: helper.Issues{},
		},
		{
			name: "no type attribute",
			rule: NewAzApiRuleQueryMustExist("test", "https://example.com", "testType", "", "", "query", blockquery.IsNotNull),
			content: `
		resource "azapi_resource" "test" {
			not_type = "baz"
			body     = {}
		}`,
			expected: helper.Issues{
				{
					Rule:    NewAzApiRuleQueryMustExist("test", "https://example.com", "testType", "", "", "query", blockquery.IsNotNull),
					Message: "Resource does not have a `type` attribute",
				},
			},
		},
		{
			name: "no body attribute",
			rule: NewAzApiRuleQueryMustExist("test", "https://example.com", "testType", "", "", "query", blockquery.IsNotNull),
			content: `
resource "azapi_resource" "test" {
	type     = "testType@0000-00-00"
	not_body = {}
}`,
			expected: helper.Issues{
				{
					Rule:    NewAzApiRuleQueryMustExist("test", "https://example.com", "testType", "", "", "query", blockquery.IsNotNull),
					Message: "Resource does not have a `body` attribute",
				},
			},
		},
		{
			name: "object array query",
			rule: NewAzApiRuleQueryMustExist("test", "https://example.com", "testType", "", "", "objectarray.#.attr", blockquery.EachIsOneOf, blockquery.NewStringResults("val")...),
			content: `
resource "azapi_resource" "test" {
	type     = "testType@0000-00-00"
	body = {
		objectarray = [
		  {
				attr = "val"
			},
		  {
				attr = "val"
			}
		]
	}
}`,
			expected: helper.Issues{},
		},
		{
			name: "unknown value",
			rule: NewAzApiRuleQueryMustExist("test", "https://example.com", "testType", "", "", "key", blockquery.IsNotKnown),
			content: `
variable "unknown" {
  type = string
}

resource "azapi_resource" "test" {
	type     = "testType@0000-00-00"
	body = {
		key = var.unknown
	}
}`,
			expected: helper.Issues{},
		},
	}

	filename := "main.tf"
	for _, c := range testCases {
		tc := c
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{filename: tc.content})
			stub := gostub.Stub(&modulecontent.AppFs, mockFs(tc.content))
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

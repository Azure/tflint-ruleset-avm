package rules_test

import (
	"testing"

	"github.com/Azure/tflint-ruleset-avm/rules"
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestTerraformLockInterface(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		JSON     bool
		Expected helper.Issues
	}{
		{
			Name: "not lock variable",
			Content: `
variable "not_lock" {
	default = "default"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "lock variable correct",
			Content: `
variable "lock" {
	default = null
	type = object({
		kind = string
		name = optional(string, null)
	})
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "lock variable too many attributes in object",
			Content: `
			variable "lock" {
				default = null
				type = object({
					kind     = string
					name     = optional(string, null)
					unwanted = string
				})
			}`,
			Expected: helper.Issues{
				&helper.Issue{
					Rule:    rules.NewTerraformLockInterfaceRule(),
					Message: "`var.lock`: expression function object() argument does not have exactly two attributes",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 2, Column: 4},
						End:      hcl.Pos{Line: 2, Column: 19},
					},
				},
			},
		},
		{
			Name: "lock variable missing kind attribute in object, but has correct number of attributes",
			Content: `
			variable "lock" {
				default = null
				type = object({
					# kind is missing
					name     = optional(string, null)
					unwanted = string
				})
			}`,
			Expected: helper.Issues{
				&helper.Issue{
					Rule:    rules.NewTerraformLockInterfaceRule(),
					Message: "`var.lock`: expression function object() argument attribute `unwanted` is not valid",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 2, Column: 4},
						End:      hcl.Pos{Line: 2, Column: 19},
					},
				},
			},
		},
		{
			Name: "lock variable kind attribute incorrect type",
			Content: `
			variable "lock" {
				default = null
				type = object({
					kind = number
					name = optional(string, null)
				})
			}`,
			Expected: helper.Issues{
				&helper.Issue{
					Rule:    rules.NewTerraformLockInterfaceRule(),
					Message: "`var.lock`: expression function object() argument attribute `kind` value is not `string`",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 2, Column: 4},
						End:      hcl.Pos{Line: 2, Column: 19},
					},
				},
			},
		},
		{
			Name: "lock variable name attribute incorrect type",
			Content: `
			variable "lock" {
				default = null
				type = object({
					kind = string
					name = optional(number, null)
				})
			}`,
			Expected: helper.Issues{
				&helper.Issue{
					Rule:    rules.NewTerraformLockInterfaceRule(),
					Message: "`var.lock`: expression function object() argument attribute `name` value `optional()` first argument is not `string`",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 2, Column: 4},
						End:      hcl.Pos{Line: 2, Column: 19},
					},
				},
			},
		},
		{
			Name: "lock variable name attribute incorrect optional default",
			Content: `
			variable "lock" {
				default = null
				type = object({
					kind = string
					name = optional(string, "")
				})
			}`,
			Expected: helper.Issues{
				&helper.Issue{
					Rule:    rules.NewTerraformLockInterfaceRule(),
					Message: "`var.lock`: expression function object() argument attribute `name` value `optional()` second argument is not `null`",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 2, Column: 4},
						End:      hcl.Pos{Line: 2, Column: 19},
					},
				},
			},
		},
		{
			Name: "lock variable no default",
			Content: `
			variable "lock" {
				type = object({
					kind = string
					name = optional(string, null)
				})
			}`,
			Expected: helper.Issues{
				&helper.Issue{
					Rule:    rules.NewTerraformLockInterfaceRule(),
					Message: "`var.lock`: default not declared",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 2, Column: 4},
						End:      hcl.Pos{Line: 2, Column: 19},
					},
				},
			},
		},
		{
			Name: "lock variable incorrect default",
			Content: `
			variable "lock" {
				default = {
					kind = "CanNotDelete"
				}
				type = object({
					kind = string
					name = optional(string, null)
				})
			}`,
			Expected: helper.Issues{
				&helper.Issue{
					Rule:    rules.NewTerraformLockInterfaceRule(),
					Message: "`var.lock`: default value is not `null`",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 2, Column: 4},
						End:      hcl.Pos{Line: 2, Column: 19},
					},
				},
			},
		},
	}

	rule := rules.NewTerraformLockInterfaceRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			filename := "variables.tf"
			if tc.JSON {
				filename += ".json"
			}

			runner := helper.TestRunner(t, map[string]string{filename: tc.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}

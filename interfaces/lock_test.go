package interfaces_test

import (
	"fmt"
	"testing"

	"github.com/Azure/tflint-ruleset-avm/interfaces"
	"github.com/Azure/tflint-ruleset-avm/rules"
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

// TestLockTerraformVar tests Lock interface.
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
			Content: fmt.Sprintf(`variable "lock" {
	default = null
	type = %s
}`, interfaces.LockTypeString),
			Expected: helper.Issues{},
		},
		{
			Name: "lock variable incorrect nullable true",
			Content: fmt.Sprintf(`
variable "lock" {
	default = null
	type = %s
	nullable = true
}`, interfaces.LockTypeString),
			Expected: helper.Issues{
				&helper.Issue{
					Rule:    rules.NewVarCheckRuleFromAvmInterface(interfaces.Lock),
					Message: "nullable should not be set.",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 8, Column: 2},
						End:      hcl.Pos{Line: 8, Column: 17},
					},
				},
			},
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
					Rule:    rules.NewVarCheckRuleFromAvmInterface(interfaces.Lock),
					Message: fmt.Sprintf("variable type does not comply with the interface specification:\n\n%s", interfaces.LockTypeString),
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 4, Column: 5},
						End:      hcl.Pos{Line: 8, Column: 7},
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
					Rule:    rules.NewVarCheckRuleFromAvmInterface(interfaces.Lock),
					Message: fmt.Sprintf("variable type does not comply with the interface specification:\n\n%s", interfaces.LockTypeString),
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 4, Column: 5},
						End:      hcl.Pos{Line: 8, Column: 7},
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
					Rule:    rules.NewVarCheckRuleFromAvmInterface(interfaces.Lock),
					Message: fmt.Sprintf("variable type does not comply with the interface specification:\n\n%s", interfaces.LockTypeString),
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 4, Column: 5},
						End:      hcl.Pos{Line: 7, Column: 7},
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
					Rule:    rules.NewVarCheckRuleFromAvmInterface(interfaces.Lock),
					Message: fmt.Sprintf("variable type does not comply with the interface specification:\n\n%s", interfaces.LockTypeString),
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 4, Column: 5},
						End:      hcl.Pos{Line: 7, Column: 7},
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
					Rule:    rules.NewVarCheckRuleFromAvmInterface(interfaces.Lock),
					Message: fmt.Sprintf("variable type does not comply with the interface specification:\n\n%s", interfaces.LockTypeString),
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 4, Column: 5},
						End:      hcl.Pos{Line: 7, Column: 7},
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
					Rule:    rules.NewVarCheckRuleFromAvmInterface(interfaces.Lock),
					Message: "default not declared",
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
					Rule:    rules.NewVarCheckRuleFromAvmInterface(interfaces.Lock),
					Message: fmt.Sprintf("default value is not correct, see: %s", interfaces.Lock.Link),
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 2, Column: 4},
						End:      hcl.Pos{Line: 2, Column: 19},
					},
				},
			},
		},
	}

	rule := rules.NewVarCheckRuleFromAvmInterface(interfaces.Lock)

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
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
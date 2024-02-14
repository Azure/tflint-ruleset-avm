package rules_test

import (
	"fmt"
	"testing"

	"github.com/Azure/tflint-ruleset-avm/interfaces"
	"github.com/Azure/tflint-ruleset-avm/rules"
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestTerraformLockInterface(t *testing.T) {
	expectedLockInterfaceIssue := &helper.Issue{
		Rule:    rules.NewAvmInterfaceLockRule(),
		Message: fmt.Sprintf("`%s` variable type does not comply with the interface specification:\n\n%s", "lock", interfaces.Lock.Type),
		Range: hcl.Range{
			Filename: "variables.tf",
			Start:    hcl.Pos{Line: 2, Column: 4},
			End:      hcl.Pos{Line: 2, Column: 19},
		},
	}
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
			Content: fmt.Sprintf(`
variable "lock" {
	default = null
	type = %s
}`, interfaces.Lock.Type),
			Expected: helper.Issues{},
		},
		{
			Name: "lock variable incorrect nullable true",
			Content: fmt.Sprintf(`
variable "lock" {
	default = null
	type = %s
	nullable = true
}`, interfaces.Lock.Type),
			Expected: helper.Issues{
				&helper.Issue{
					Rule:    rules.NewAvmInterfaceLockRule(),
					Message: fmt.Sprintf("`var.%s`: nullable is set and should not be, we require this to be true and this is the default behaviour so no need to set explicitly", interfaces.Lock.Name),
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 16},
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
				expectedLockInterfaceIssue,
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
				expectedLockInterfaceIssue,
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
				expectedLockInterfaceIssue,
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
				expectedLockInterfaceIssue,
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
				expectedLockInterfaceIssue,
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
					Rule:    rules.NewAvmInterfaceLockRule(),
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
					Rule:    rules.NewAvmInterfaceLockRule(),
					Message: fmt.Sprintf("`var.%s`: default value is not correct, see: %s", interfaces.Lock.Name, interfaces.Lock.Link),
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 2, Column: 4},
						End:      hcl.Pos{Line: 2, Column: 19},
					},
				},
			},
		},
	}

	rule := rules.NewAvmInterfaceLockRule()

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

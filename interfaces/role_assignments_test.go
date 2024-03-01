package interfaces_test

import (
	"fmt"
	"testing"

	"github.com/Azure/tflint-ruleset-avm/interfaces"
	"github.com/Azure/tflint-ruleset-avm/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestRoleAssignmentsInterface(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		JSON     bool
		Expected helper.Issues
	}{
		{
			Name: "not role_assignments variable",
			Content: `
variable "not_role_assignments" {
	default = "default"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "role_assignments variable correct",
			Content: fmt.Sprintf(`
variable "role_assignments" {
  type = %s
  default     = {}
  nullable    = false
  }`, interfaces.RoleAssignmentsTypeString),
			Expected: helper.Issues{},
		},
	}

	rule := rules.NewVarCheckRuleFromAvmInterface(interfaces.RoleAssignments)

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

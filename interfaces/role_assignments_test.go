package interfaces_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Azure/tflint-ruleset-avm/interfaces"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestRoleAssignmentsInterface(t *testing.T) {
	const roleAssignmentsV2 = `variable "role_assignments" {
  type = map(object({
    name                                   = optional(string, null)
    role_definition_id_or_name             = string
    principal_id                           = string
    description                            = optional(string, null)
    skip_service_principal_aad_check       = optional(bool, false)
    condition                              = optional(string, null)
    condition_version                      = optional(string, null)
    delegated_managed_identity_resource_id = optional(string, null)
    principal_type                         = optional(string, null)
  }))
  default  = {}
  nullable = false
}`
	roleAssignmentsV1 := strings.Replace(roleAssignmentsV2, "    name                                   = optional(string, null)\n", "", 1)

	compatibilityRule := registeredInterfaceRule(t, "role_assignments")
	deprecationRule := registeredInterfaceRule(t, "deprecated_role_assignments_interface")
	cases := []struct {
		Name                string
		Content             string
		Expected            helper.Issues
		ExpectedDeprecation helper.Issues
	}{
		{
			Name:    "correct v1",
			Content: roleAssignmentsV1,
			ExpectedDeprecation: helper.Issues{
				{
					Rule:    deprecationRule,
					Message: deprecatedRoleAssignmentsMessage,
				},
			},
		},
		{
			Name:     "correct v2",
			Content:  roleAssignmentsV2,
			Expected: helper.Issues{},
		},
		{
			Name: "incorrect name type",
			Content: strings.Replace(
				roleAssignmentsV2,
				"name                                   = optional(string, null)",
				"name                                   = optional(number, null)",
				1,
			),
			Expected: helper.Issues{
				{
					Rule:    interfaces.NewVarCheckRuleFromAvmInterface(interfaces.RoleAssignmentsV2),
					Message: fmt.Sprintf("variable type does not comply with the interface specification:\n\n%s", interfaces.RoleAssignmentsV2TypeString),
				},
			},
		},
		{
			Name:    "missing variable",
			Content: "",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			runner := helper.TestRunner(t, map[string]string{"variables.tf": tc.Content})

			if err := compatibilityRule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}
			expected := tc.Expected
			if expected == nil {
				expected = helper.Issues{}
			}
			helper.AssertIssuesWithoutRange(t, expected, runner.Issues)

			runner = helper.TestRunner(t, map[string]string{"variables.tf": tc.Content})
			if err := deprecationRule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}
			expected = tc.ExpectedDeprecation
			if expected == nil {
				expected = helper.Issues{}
			}
			helper.AssertIssuesWithoutRange(t, expected, runner.Issues)
		})
	}
}

const deprecatedRoleAssignmentsMessage = "role_assignments uses deprecated interface variant 1; migrate to variant 2 by adding name = optional(string, null). Support for variant 1 will be removed in the release following the v0.19.0 migration window."

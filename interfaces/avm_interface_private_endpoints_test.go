package interfaces_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Azure/tflint-ruleset-avm/interfaces"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestPrivateEndpoints(t *testing.T) {
	const canonical = `variable "private_endpoints" {
  type = map(object({
    name = optional(string, null)
    role_assignments = optional(map(object({
      name                                   = optional(string, null)
      role_definition_id_or_name             = string
      principal_id                           = string
      description                            = optional(string, null)
      skip_service_principal_aad_check       = optional(bool, false)
      condition                              = optional(string, null)
      condition_version                      = optional(string, null)
      delegated_managed_identity_resource_id = optional(string, null)
      principal_type                         = optional(string, null)
    })), {})
    lock = optional(object({
      kind  = string
      name  = optional(string, null)
      notes = optional(string, null)
    }), null)
    tags                                    = optional(map(string), null)
    subnet_resource_id                      = string
    subresource_name                        = optional(string, null)
    private_dns_zone_group_name             = optional(string, "default")
    private_dns_zone_resource_ids           = optional(set(string), [])
    application_security_group_associations = optional(map(string), {})
    private_service_connection_name         = optional(string, null)
    network_interface_name                  = optional(string, null)
    location                                = optional(string, null)
    resource_group_name                     = optional(string, null)
    ip_configurations = optional(map(object({
      name               = string
      private_ip_address = string
      member_name        = optional(string)
    })), {})
  }))
  default  = {}
  nullable = false
}`
	compatibilityRule := registeredInterfaceRule(t, "avm_interface_private_endpoints")
	deprecationRule := registeredInterfaceRule(t, "avm_interface_private_endpoints_deprecated")
	type testCase struct {
		Name                string
		Content             string
		Expected            helper.Issues
		ExpectedDeprecation helper.Issues
	}
	cases := []testCase{}

	for _, includeLockNotes := range []bool{true, false} {
		for _, includeMemberName := range []bool{true, false} {
			for _, includeRoleAssignmentName := range []bool{true, false} {
				content := canonical
				if !includeLockNotes {
					content = strings.Replace(content, "    notes = optional(string, null)\n", "", 1)
				}
				if !includeMemberName {
					content = strings.Replace(content, "    member_name        = optional(string)\n", "", 1)
				}
				if !includeRoleAssignmentName {
					content = strings.Replace(content, "    name                                   = optional(string, null)\n", "", 1)
				}
				cases = append(cases, testCase{
					Name: fmt.Sprintf(
						"correct lock notes %t member name %t role assignment name %t",
						includeLockNotes,
						includeMemberName,
						includeRoleAssignmentName,
					),
					Content: content,
				})
				if !includeLockNotes || !includeMemberName || !includeRoleAssignmentName {
					cases[len(cases)-1].ExpectedDeprecation = helper.Issues{
						{
							Rule:    deprecationRule,
							Message: deprecatedPrivateEndpointsMessage,
						},
					}
				}
			}
		}
	}

	cases = append(cases,
		testCase{
			Name:     "incorrect required member name",
			Content:  strings.Replace(canonical, "member_name        = optional(string)", "member_name        = string", 1),
			Expected: privateEndpointsIssues(),
		},
		testCase{
			Name:     "incorrect lock notes type",
			Content:  strings.Replace(canonical, "notes = optional(string, null)", "notes = optional(number, null)", 1),
			Expected: privateEndpointsIssues(),
		},
		testCase{
			Name:     "incorrect role assignment name type",
			Content:  strings.Replace(canonical, "name                                   = optional(string, null)", "name                                   = optional(number, null)", 1),
			Expected: privateEndpointsIssues(),
		},
		testCase{
			Name:     "incorrect required subresource name",
			Content:  strings.Replace(canonical, "subresource_name                        = optional(string, null)", "subresource_name                        = string", 1),
			Expected: privateEndpointsIssues(),
		},
		testCase{
			Name:    "missing variable",
			Content: "",
		},
	)

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

func privateEndpointsIssues() helper.Issues {
	return helper.Issues{
		{
			Rule:    interfaces.NewVarCheckRuleFromAvmInterface(interfaces.PrivateEndpoints),
			Message: fmt.Sprintf("variable type does not comply with the interface specification:\n\n%s", interfaces.PrivateEndpointTypeString),
		},
	}
}

const deprecatedPrivateEndpointsMessage = "private_endpoints uses a deprecated transitional schema; migrate nested lock, role_assignments, and ip_configurations to the canonical v0.19.0 schema. Transitional schemas will be removed in the release following the v0.19.0 migration window."

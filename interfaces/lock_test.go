package interfaces_test

import (
	"fmt"
	"testing"

	"github.com/Azure/tflint-ruleset-avm/interfaces"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

// TestLockTerraformVar tests Lock interface.
func TestTerraformLockInterface(t *testing.T) {
	const lockV1 = `variable "lock" {
  type = object({
    kind = string
    name = optional(string, null)
  })
  default = null
}`
	const lockV2 = `variable "lock" {
  type = object({
    kind  = string
    name  = optional(string, null)
    notes = optional(string, null)
  })
  default = null
}`

	compatibilityRule := registeredInterfaceRule(t, "lock")
	deprecationRule := registeredInterfaceRule(t, "deprecated_lock_interface")
	cases := []struct {
		Name                string
		Content             string
		Expected            helper.Issues
		ExpectedDeprecation helper.Issues
	}{
		{
			Name:    "correct v1",
			Content: lockV1,
			ExpectedDeprecation: helper.Issues{
				{
					Rule:    deprecationRule,
					Message: deprecatedLockMessage,
				},
			},
		},
		{
			Name:     "correct v2",
			Content:  lockV2,
			Expected: helper.Issues{},
		},
		{
			Name: "incorrect notes",
			Content: `
variable "lock" {
  type = object({
    kind  = string
    name  = optional(string, null)
    notes = string
  })
  default = null
}`,
			Expected: helper.Issues{
				{
					Rule:    interfaces.NewVarCheckRuleFromAvmInterface(interfaces.LockV2),
					Message: fmt.Sprintf("variable type does not comply with the interface specification:\n\n%s", interfaces.LockV2TypeString),
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

const deprecatedLockMessage = "lock uses deprecated interface variant 1; migrate to variant 2 by adding notes = optional(string, null). Support for variant 1 will be removed in the release following the v0.19.0 migration window."

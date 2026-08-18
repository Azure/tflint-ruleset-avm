package interfaces_test

import (
	"fmt"
	"testing"

	"github.com/Azure/tflint-ruleset-avm/interfaces"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

// TestCustomerManagedKeyV2Interface tests the customer managed key v2 interface.
func TestCustomerManagedKeyV2Interface(t *testing.T) {
	rule := interfaces.NewVarCheckRuleFromAvmInterface(interfaces.CustomerManagedKeyV2)
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name:     "correct",
			Content:  customerManagedKeyV2Fixture,
			Expected: helper.Issues{},
		},
		{
			Name: "incorrect type",
			Content: `
variable "customer_managed_key" {
	type = object({
		key_vault_key_uri = string
		user_assigned_identity = optional(object({
			client_id = number
		}), null)
	})
	default = null
}`,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: fmt.Sprintf("variable type does not comply with the interface specification:\n\n%s", interfaces.CustomerManagedKeyV2TypeString),
				},
			},
		},
		{
			Name:     "missing variable",
			Content:  "",
			Expected: helper.Issues{},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			runner := helper.TestRunner(t, map[string]string{"variables.tf": tc.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssuesWithoutRange(t, tc.Expected, runner.Issues)
		})
	}
}

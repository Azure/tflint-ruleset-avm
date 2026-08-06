package interfaces_test

import (
	"fmt"
	"testing"

	"github.com/Azure/tflint-ruleset-avm/interfaces"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

// TestCustomerManagedKeyInterface tests both customer managed key interface variants.
func TestCustomerManagedKeyInterface(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		JSON     bool
		Expected helper.Issues
	}{
		{
			Name:     "correct",
			Content:  toTerraformVarType(interfaces.CustomerManagedKey),
			Expected: helper.Issues{},
		},
		{
			Name:     "correct v2",
			Content:  toTerraformVarType(interfaces.CustomerManagedKeyV2),
			Expected: helper.Issues{},
		},
		{
			Name: "incorrect shape",
			Content: `
variable "customer_managed_key" {
	type = object({
		key_vault_key_uri = number
	})
	default = null
}`,
			Expected: helper.Issues{
				{
					Rule:    interfaces.NewVarCheckRuleFromAvmInterface(interfaces.CustomerManagedKey),
					Message: fmt.Sprintf("variable type does not comply with the interface specification:\n\n%s", interfaces.CustomerManagedKeyTypeString),
				},
			},
		},
	}

	rule := interfaces.Rules[0]

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

			helper.AssertIssuesWithoutRange(t, tc.Expected, runner.Issues)
		})
	}
}

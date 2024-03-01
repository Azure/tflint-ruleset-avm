package interfaces_test

import (
	"fmt"
	"testing"

	"github.com/Azure/tflint-ruleset-avm/interfaces"
	"github.com/Azure/tflint-ruleset-avm/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

// TestDiagnosticSettingsInterface tests the diagnostic settings interface.
func TestManagedIdentitiesInterface(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		JSON     bool
		Expected helper.Issues
	}{
		{
			Name: "not managed_identities variable",
			Content: `
variable "not_managed_identities" {
	default = "default"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "managed_identities variable correct",
			Content: fmt.Sprintf(`
variable "managed_identities" {
  type = %s
  default     = {}
  nullable    = false
  }`, interfaces.ManagedIdentitiesTypeString),
			Expected: helper.Issues{},
		},
	}

	rule := rules.NewVarCheckRuleFromAvmInterface(interfaces.ManagedIdentities)

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

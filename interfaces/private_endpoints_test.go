package interfaces_test

import (
	"github.com/Azure/tflint-ruleset-avm/interfaces"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"testing"
)

func TestPrivateEndpoints(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		JSON     bool
		Expected helper.Issues
	}{
		{
			Name:     "correct",
			Content:  toTerraformVarType(interfaces.PrivateEndpoints),
			Expected: helper.Issues{},
		},
	}

	rule := interfaces.NewVarCheckRuleFromAvmInterface(interfaces.PrivateEndpoints)

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

package rules_test

import (
	"fmt"
	"testing"

	"github.com/Azure/tflint-ruleset-avm/interfaces"
	"github.com/Azure/tflint-ruleset-avm/rules"
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestDiagnosticSettingsInterface(t *testing.T) {
	expectedDiagnosticSettingsInterfaceIssue := &helper.Issue{
		Rule:    rules.NewAvmInterfaceDiagnosticSettingsRule(),
		Message: fmt.Sprintf("`%s` variable type does not comply with the interface specification:\n\n%s", "lock", interfaces.DiagnosticSettings.Type),
		Range: hcl.Range{
			Filename: "variables.tf",
			Start:    hcl.Pos{Line: 2, Column: 4},
			End:      hcl.Pos{Line: 2, Column: 19},
		},
	}
	_ = expectedDiagnosticSettingsInterfaceIssue
	cases := []struct {
		Name     string
		Content  string
		JSON     bool
		Expected helper.Issues
	}{
		{
			Name: "not diagnostic_settings variable",
			Content: `
variable "not_diagnostic_settings" {
	default = "default"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "diagnostic_settings variable correct",
			Content: fmt.Sprintf(`
variable "diagnostic_settings" {
	default = {}
	nullable = %t
	type = %s
}`, interfaces.DiagnosticSettings.Nullable, interfaces.DiagnosticSettings.Type),
			Expected: helper.Issues{},
		},
	}

	rule := rules.NewAvmInterfaceDiagnosticSettingsRule()

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

package interfaces_test

import (
	"fmt"
	"testing"

	"github.com/Azure/tflint-ruleset-avm/interfaces"
	"github.com/Azure/tflint-ruleset-avm/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

// TestDiagnosticSettingsInterface tests the diagnostic settings interface.
func TestDiagnosticSettingsInterface(t *testing.T) {
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
  nullable = false
  type = %s
}
`, interfaces.DiagnosticTypeString),
			Expected: helper.Issues{},
		},
	}

	rule := rules.NewVarCheckRuleFromAvmInterface(interfaces.DiagnosticSettings)

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

//// TerraformVar returns a string that represents the interface as the
//// minimum required Terraform variable definition for testing.
//func TerraformVar(i interfaces.AvmInterface, t *testing.T) string {
//	f := hclwrite.NewEmptyFile()
//	rootBody := f.Body()
//	varBlock := rootBody.AppendNewBlock("variable", []string{i.Name})
//	varBody := varBlock.Body()
//	var exp hcl.Expression
//	varBody.SetAttributeRaw(ex)
//	// I couldn't get the hclwrite to work with the type constraint so I'm just adding it as a string
//	// using SetSAttributeRaw and hclWrite.Token.
//	varBody.SetAttributeRaw("type",  exp.)
//	varBody.SetAttributeValue("default", i.Default)
//	// If the interface is not nullable, set the nullable attribute to false.
//	// the default is true so we only need to set it if it's false.
//	if !i.Nullable {
//		varBody.SetAttributeValue("nullable", cty.False)
//	}
//	return string(f.Bytes())
//}

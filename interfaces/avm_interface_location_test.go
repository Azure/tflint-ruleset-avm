package interfaces_test

import (
	"testing"

	"github.com/Azure/tflint-ruleset-avm/interfaces"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestTerraformLocationInterface(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		JSON     bool
		Expected helper.Issues
	}{
		{
			Name:     "correct",
			Content:  toTerraformVarType(interfaces.Location),
			Expected: helper.Issues{},
		},
		{
			Name: "telemetry only",
			Content: `variable "location" {
  type    = string
  default = null

  validation {
    condition     = !var.enable_telemetry || var.location != null
    error_message = "The location must not be null when telemetry is enabled."
  }
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "telemetry only with explicit nullable",
			Content: `variable "location" {
  type     = string
  default  = null
  nullable = true

  validation {
    condition     = var.location != null || !var.enable_telemetry
    error_message = "The location must not be null when telemetry is enabled."
  }
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "telemetry only without telemetry guard",
			Content: `variable "location" {
  type    = string
  default = null
}`,
		},
	}

	rule := registeredInterfaceRule(t, "avm_interface_location")

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

			if tc.Name == "telemetry only without telemetry guard" {
				if len(runner.Issues) == 0 {
					t.Fatal("registered location rule accepted a telemetry-only location without a telemetry guard")
				}
				return
			}

			helper.AssertIssuesWithoutRange(t, tc.Expected, runner.Issues)
		})
	}
}

func TestTelemetryLocationInterfaceRejectsInvalidDeclarations(t *testing.T) {
	cases := []struct {
		name     string
		content  string
		expected helper.Issues
	}{
		{
			name: "missing telemetry guard",
			content: `variable "location" {
  type    = string
  default = null
}`,
			expected: helper.Issues{
				{
					Rule:    interfaces.NewVarCheckRuleFromAvmInterface(interfaces.TelemetryLocation),
					Message: "location must include a validation condition that requires a value when telemetry is enabled",
				},
			},
		},
		{
			name: "nonnullable",
			content: `variable "location" {
  type     = string
  default  = null
  nullable = false

  validation {
    condition     = !var.enable_telemetry || var.location != null
    error_message = "The location must not be null when telemetry is enabled."
  }
}`,
			expected: helper.Issues{
				{
					Rule:    interfaces.NewVarCheckRuleFromAvmInterface(interfaces.TelemetryLocation),
					Message: "nullable should not be set.",
				},
			},
		},
	}

	rule := interfaces.NewVarCheckRuleFromAvmInterface(interfaces.TelemetryLocation)
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"variables.tf": tc.content})
			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}
			helper.AssertIssuesWithoutRange(t, tc.expected, runner.Issues)
		})
	}
}

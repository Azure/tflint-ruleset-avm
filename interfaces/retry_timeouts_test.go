package interfaces_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// Canonical TFFR7 snapshot:
// https://github.com/Azure/Azure-Verified-Modules/blob/005b4f94e9197d1e23ea4213c92b4263bf636a2e/docs/content/specs-defs/includes/terraform/shared/functional/TFFR7.md
const validRetry = `resource "azapi_resource" "example" {
  type = "Microsoft.Example/widgets@2024-01-01"
}

variable "retry" {
  type = object({
    error_message_regex  = optional(list(string))
    interval_seconds     = optional(number)
    max_interval_seconds = optional(number)
  })
  default = null
}`

const validTimeouts = `resource "azapi_resource" "example" {
  type = "Microsoft.Example/widgets@2024-01-01"
}

variable "timeouts" {
  type = object({
    create = optional(string)
    read   = optional(string)
    update = optional(string)
    delete = optional(string)
  })
  default = null
}`

func TestRetryAcceptsPermittedDefaults(t *testing.T) {
	rule := requiredInterfaceRule(t, "retry")
	cases := map[string]string{
		"canonical null": validRetry,
		"empty object":   strings.Replace(validRetry, "default = null", "default = {}", 1),
		"module field defaults": strings.Replace(
			validRetry,
			"default = null",
			`default = {
    error_message_regex  = ["ScopeLocked"]
    interval_seconds     = 10
    max_interval_seconds = 60
  }`,
			1,
		),
		"optional field defaults": strings.NewReplacer(
			"optional(list(string))", `optional(list(string), ["ScopeLocked"])`,
			"optional(number)", "optional(number, 10)",
		).Replace(validRetry),
		"explicit nullable true": strings.Replace(
			validRetry,
			"  default = null\n",
			"  default  = null\n  nullable = true\n",
			1,
		),
	}

	for name, content := range cases {
		content := content
		t.Run(name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"main.tf": content})
			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			helper.AssertIssues(t, helper.Issues{}, runner.Issues)
		})
	}
}

func TestRetryRejectsMalformedTypeDefaultAndNullability(t *testing.T) {
	rule := requiredInterfaceRule(t, "retry")
	cases := []struct {
		name        string
		content     string
		attribute   string
		expectation string
	}{
		{
			name:        "required field",
			content:     strings.Replace(validRetry, "interval_seconds     = optional(number)", "interval_seconds     = number", 1),
			attribute:   "type",
			expectation: "must use the TFFR7 retry object shape with optional fields",
		},
		{
			name:        "wrong field type",
			content:     strings.Replace(validRetry, "max_interval_seconds = optional(number)", "max_interval_seconds = optional(string)", 1),
			attribute:   "type",
			expectation: "must use the TFFR7 retry object shape with optional fields",
		},
		{
			name:        "missing field",
			content:     strings.Replace(validRetry, "    max_interval_seconds = optional(number)\n", "", 1),
			attribute:   "type",
			expectation: "must use the TFFR7 retry object shape with optional fields",
		},
		{
			name:        "extra field",
			content:     strings.Replace(validRetry, "    max_interval_seconds = optional(number)\n", "    max_interval_seconds = optional(number)\n    unexpected           = optional(number)\n", 1),
			attribute:   "type",
			expectation: "must use the TFFR7 retry object shape with optional fields",
		},
		{
			name:        "illegal optional field default",
			content:     strings.Replace(validRetry, "optional(number)", `optional(number, ["invalid"])`, 1),
			attribute:   "type",
			expectation: "must use the TFFR7 retry object shape with optional fields",
		},
		{
			name:        "missing module default",
			content:     strings.Replace(validRetry, "  default = null\n", "", 1),
			attribute:   "default",
			expectation: "must be `null`, `{}`, or an object containing legal retry field defaults",
		},
		{
			name:        "unknown module default field",
			content:     strings.Replace(validRetry, "default = null", "default = { unexpected = true }", 1),
			attribute:   "default",
			expectation: "must be `null`, `{}`, or an object containing legal retry field defaults",
		},
		{
			name:        "invalid module default field value",
			content:     strings.Replace(validRetry, "default = null", "default = { interval_seconds = [10] }", 1),
			attribute:   "default",
			expectation: "must be `null`, `{}`, or an object containing legal retry field defaults",
		},
		{
			name:        "nullable false",
			content:     strings.Replace(validRetry, "  default = null\n", "  default  = null\n  nullable = false\n", 1),
			attribute:   "nullable",
			expectation: "must permit `null`",
		},
	}

	runMalformedVariableCases(t, rule, cases)
}

func TestTimeoutsAcceptsPermittedDefaults(t *testing.T) {
	rule := requiredInterfaceRule(t, "timeouts")
	cases := map[string]string{
		"canonical null": validTimeouts,
		"empty object":   strings.Replace(validTimeouts, "default = null", "default = {}", 1),
		"module field defaults": strings.Replace(
			validTimeouts,
			"default = null",
			`default = {
    create = "30m"
    read   = "5m"
    update = "30m"
    delete = "30m"
  }`,
			1,
		),
		"optional field defaults": strings.ReplaceAll(validTimeouts, "optional(string)", `optional(string, "30m")`),
		"explicit nullable true": strings.Replace(
			validTimeouts,
			"  default = null\n",
			"  default  = null\n  nullable = true\n",
			1,
		),
	}

	for name, content := range cases {
		content := content
		t.Run(name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"main.tf": content})
			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			helper.AssertIssues(t, helper.Issues{}, runner.Issues)
		})
	}
}

func TestTimeoutsRejectsMalformedTypeDefaultAndNullability(t *testing.T) {
	rule := requiredInterfaceRule(t, "timeouts")
	cases := []struct {
		name        string
		content     string
		attribute   string
		expectation string
	}{
		{
			name:        "required field",
			content:     strings.Replace(validTimeouts, "create = optional(string)", "create = string", 1),
			attribute:   "type",
			expectation: "must use the TFFR7 timeouts object shape with optional fields",
		},
		{
			name:        "wrong field type",
			content:     strings.Replace(validTimeouts, "delete = optional(string)", "delete = optional(number)", 1),
			attribute:   "type",
			expectation: "must use the TFFR7 timeouts object shape with optional fields",
		},
		{
			name:        "missing field",
			content:     strings.Replace(validTimeouts, "    delete = optional(string)\n", "", 1),
			attribute:   "type",
			expectation: "must use the TFFR7 timeouts object shape with optional fields",
		},
		{
			name:        "extra field",
			content:     strings.Replace(validTimeouts, "    delete = optional(string)\n", "    delete     = optional(string)\n    unexpected = optional(string)\n", 1),
			attribute:   "type",
			expectation: "must use the TFFR7 timeouts object shape with optional fields",
		},
		{
			name:        "illegal optional field default",
			content:     strings.Replace(validTimeouts, "optional(string)", `optional(string, ["30m"])`, 1),
			attribute:   "type",
			expectation: "must use the TFFR7 timeouts object shape with optional fields",
		},
		{
			name:        "missing module default",
			content:     strings.Replace(validTimeouts, "  default = null\n", "", 1),
			attribute:   "default",
			expectation: "must be `null`, `{}`, or an object containing legal timeout field defaults",
		},
		{
			name:        "unknown module default field",
			content:     strings.Replace(validTimeouts, "default = null", "default = { unexpected = \"5m\" }", 1),
			attribute:   "default",
			expectation: "must be `null`, `{}`, or an object containing legal timeout field defaults",
		},
		{
			name:        "invalid module default field value",
			content:     strings.Replace(validTimeouts, "default = null", "default = { create = [\"5m\"] }", 1),
			attribute:   "default",
			expectation: "must be `null`, `{}`, or an object containing legal timeout field defaults",
		},
		{
			name:        "nullable false",
			content:     strings.Replace(validTimeouts, "  default = null\n", "  default  = null\n  nullable = false\n", 1),
			attribute:   "nullable",
			expectation: "must permit `null`",
		},
	}

	runMalformedVariableCases(t, rule, cases)
}

func runMalformedVariableCases(
	t *testing.T,
	rule tflint.Rule,
	cases []struct {
		name        string
		content     string
		attribute   string
		expectation string
	},
) {
	t.Helper()
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"main.tf": tc.content})
			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			helper.AssertIssuesWithoutRange(
				t,
				expectedVariableAttributeIssue(rule, tc.attribute, tc.expectation, hcl.Range{}),
				runner.Issues,
			)
		})
	}
}

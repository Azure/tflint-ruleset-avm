package interfaces_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

// Canonical TFFR6 snapshot:
// https://github.com/Azure/Azure-Verified-Modules/blob/005b4f94e9197d1e23ea4213c92b4263bf636a2e/docs/content/specs-defs/includes/terraform/shared/functional/TFFR6.md
const validResourceTypes = `resource "azapi_resource" "example" {
  type = var.resource_types.widget
}

variable "resource_types" {
  type = object({
    widget = optional(string, "Microsoft.Example/widgets@2024-01-01")
    child = optional(object({
      child_widget = optional(string)
      grandchild = optional(object({
        grandchild_widget = optional(string)
      }), {})
    }), {})
  })
  default  = {}
  nullable = false
}`

func TestResourceTypesValidRecursiveShapes(t *testing.T) {
	rule := requiredInterfaceRule(t, "resource_types")
	cases := map[string]string{
		"recursive": validResourceTypes,
		"multiple owned resources": strings.Replace(
			validResourceTypes,
			`    widget = optional(string, "Microsoft.Example/widgets@2024-01-01")`,
			`    widget       = optional(string, "Microsoft.Example/widgets@2024-01-01")
    widget_child = optional(string, "Microsoft.Example/widgets/children@2024-01-01")`,
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

func TestResourceTypesRejectMalformedRecursiveShapes(t *testing.T) {
	rule := requiredInterfaceRule(t, "resource_types")
	expectation := "must recursively contain optional string resource leaves and optional object submodule slots with canonical defaults"
	cases := map[string]string{
		"top level map": `resource "azapi_resource" "example" {}

variable "resource_types" {
  type     = map(string)
  default  = {}
  nullable = false
}`,
		"required owned resource leaf": strings.Replace(
			validResourceTypes,
			`widget = optional(string, "Microsoft.Example/widgets@2024-01-01")`,
			`widget = string`,
			1,
		),
		"owned resource leaf without API default": strings.Replace(
			validResourceTypes,
			`widget = optional(string, "Microsoft.Example/widgets@2024-01-01")`,
			`widget = optional(string)`,
			1,
		),
		"owned resource leaf with empty API default": strings.Replace(
			validResourceTypes,
			`widget = optional(string, "Microsoft.Example/widgets@2024-01-01")`,
			`widget = optional(string, "")`,
			1,
		),
		"inherited resource leaf with parent default": strings.Replace(
			validResourceTypes,
			`child_widget = optional(string)`,
			`child_widget = optional(string, "Microsoft.Example/children@2024-01-01")`,
			1,
		),
		"required child object": strings.NewReplacer(
			"grandchild = optional(object({",
			"grandchild = object({",
			"      }), {})",
			"      })",
		).Replace(validResourceTypes),
		"child object without empty default": strings.Replace(
			validResourceTypes,
			"}), {})",
			"}))",
			1,
		),
		"child object with null default": strings.Replace(
			validResourceTypes,
			"}), {})",
			"}), null)",
			1,
		),
		"non string leaf": strings.Replace(
			validResourceTypes,
			`child_widget = optional(string)`,
			`child_widget = optional(number)`,
			1,
		),
		"no direct owned resource leaf": strings.Replace(
			validResourceTypes,
			`    widget = optional(string, "Microsoft.Example/widgets@2024-01-01")
`,
			"",
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
			helper.AssertIssuesWithoutRange(
				t,
				expectedVariableAttributeIssue(rule, "type", expectation, hcl.Range{}),
				runner.Issues,
			)
		})
	}
}

func TestResourceTypesRejectDefaultAndNullability(t *testing.T) {
	rule := requiredInterfaceRule(t, "resource_types")
	cases := []struct {
		name        string
		content     string
		attribute   string
		expectation string
	}{
		{
			name:        "missing default",
			content:     strings.Replace(validResourceTypes, "  default  = {}\n", "", 1),
			attribute:   "default",
			expectation: "must be `{}`",
		},
		{
			name:        "null default",
			content:     strings.Replace(validResourceTypes, "default  = {}", "default  = null", 1),
			attribute:   "default",
			expectation: "must be `{}`",
		},
		{
			name:        "missing nullable",
			content:     strings.Replace(validResourceTypes, "  nullable = false\n", "", 1),
			attribute:   "nullable",
			expectation: "must be `false`",
		},
		{
			name:        "nullable true",
			content:     strings.Replace(validResourceTypes, "nullable = false", "nullable = true", 1),
			attribute:   "nullable",
			expectation: "must be `false`",
		},
	}

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

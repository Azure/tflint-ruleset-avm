package interfaces_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

// Canonical TFFR8 snapshot:
// https://github.com/Azure/Azure-Verified-Modules/blob/005b4f94e9197d1e23ea4213c92b4263bf636a2e/docs/content/specs-defs/includes/terraform/shared/functional/TFFR8.md
const validIgnoreBodyChanges = `resource "azapi_resource" "example" {
  type = "Microsoft.Example/widgets@2024-01-01"
}

variable "ignore_body_changes" {
  type = object({
    widget = optional(list(string), [])
    child = optional(object({
      child_widget = optional(list(string), [])
      grandchild = optional(object({
        grandchild_widget = optional(list(string), [])
      }), {})
    }), {})
  })
  default  = {}
  nullable = false
}`

func TestIgnoreBodyChangesValidRecursiveShape(t *testing.T) {
	rule := requiredInterfaceRule(t, "ignore_body_changes")
	cases := map[string]string{
		"empty resource leaf defaults": validIgnoreBodyChanges,
		"nonempty resource leaf default": strings.Replace(
			validIgnoreBodyChanges,
			"widget = optional(list(string), [])",
			`widget = optional(list(string), ["properties.status"])`,
			1,
		),
	}

	for name, content := range cases {
		t.Run(name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"main.tf": content})
			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			helper.AssertIssues(t, helper.Issues{}, runner.Issues)
		})
	}
}

func TestIgnoreBodyChangesRejectMalformedRecursiveShapes(t *testing.T) {
	rule := requiredInterfaceRule(t, "ignore_body_changes")
	expectation := "must recursively contain optional list(string) resource leaves with valid defaults and optional object submodule slots defaulting to `{}`"
	cases := map[string]string{
		"top level map": `resource "azapi_resource" "example" {}

variable "ignore_body_changes" {
  type     = map(list(string))
  default  = {}
  nullable = false
}`,
		"required resource leaf": strings.Replace(
			validIgnoreBodyChanges,
			"widget = optional(list(string), [])",
			"widget = list(string)",
			1,
		),
		"resource leaf without empty default": strings.Replace(
			validIgnoreBodyChanges,
			"widget = optional(list(string), [])",
			"widget = optional(list(string))",
			1,
		),
		"resource leaf with empty string default": strings.Replace(
			validIgnoreBodyChanges,
			"widget = optional(list(string), [])",
			`widget = optional(list(string), [""])`,
			1,
		),
		"resource leaf with wrong element type": strings.Replace(
			validIgnoreBodyChanges,
			"child_widget = optional(list(string), [])",
			"child_widget = optional(list(number), [])",
			1,
		),
		"required child object": strings.NewReplacer(
			"grandchild = optional(object({",
			"grandchild = object({",
			"      }), {})",
			"      })",
		).Replace(validIgnoreBodyChanges),
		"child object without empty default": strings.Replace(
			validIgnoreBodyChanges,
			"}), {})",
			"}))",
			1,
		),
		"child object with null default": strings.Replace(
			validIgnoreBodyChanges,
			"}), {})",
			"}), null)",
			1,
		),
		"no direct resource leaf": strings.Replace(
			validIgnoreBodyChanges,
			"    widget = optional(list(string), [])\n",
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

func TestIgnoreBodyChangesRejectDefaultAndNullability(t *testing.T) {
	rule := requiredInterfaceRule(t, "ignore_body_changes")
	cases := []struct {
		name        string
		content     string
		attribute   string
		expectation string
	}{
		{
			name:        "missing default",
			content:     strings.Replace(validIgnoreBodyChanges, "  default  = {}\n", "", 1),
			attribute:   "default",
			expectation: "must be `{}`",
		},
		{
			name:        "nonempty default",
			content:     strings.Replace(validIgnoreBodyChanges, "default  = {}", "default  = { widget = [\"properties.status\"] }", 1),
			attribute:   "default",
			expectation: "must be `{}`",
		},
		{
			name:        "missing nullable",
			content:     strings.Replace(validIgnoreBodyChanges, "  nullable = false\n", "", 1),
			attribute:   "nullable",
			expectation: "must be `false`",
		},
		{
			name:        "nullable true",
			content:     strings.Replace(validIgnoreBodyChanges, "nullable = false", "nullable = true", 1),
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

package interfaces_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

const validResourceTags = `variable "resource_tags" {
  type = object({
    resources = optional(object({
      this      = optional(map(string))
      secondary = optional(map(string))
    }))
    modules = optional(object({
      child = optional(object({
        resources = optional(object({
          this = optional(map(string))
        }))
      }))
    }))
  })
  default = null
}`

func TestResourceTagsAcceptsCanonicalRecursiveShapes(t *testing.T) {
	rule := requiredInterfaceRule(t, "resource_tags")
	cases := map[string]string{
		"resources and modules": validResourceTags,
		"resources only": `variable "resource_tags" {
  type = object({
    resources = optional(object({
      this = optional(map(string))
    }))
  })
  default = null
}`,
		"modules only": `variable "resource_tags" {
  type = object({
    modules = optional(object({
      child = optional(object({
        resources = optional(object({
          this = optional(map(string))
        }))
      }))
    }))
  })
  default = null
}`,
		"explicit nullable true": strings.Replace(
			validResourceTags,
			"  default = null\n",
			"  default  = null\n  nullable = true\n",
			1,
		),
	}

	for name, content := range cases {
		content := content
		t.Run(name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"variables.tf": content})
			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			helper.AssertIssues(t, helper.Issues{}, runner.Issues)
		})
	}
}

func TestResourceTagsIsOptional(t *testing.T) {
	rule := requiredInterfaceRule(t, "resource_tags")
	runner := helper.TestRunner(t, map[string]string{
		"variables.tf": `variable "tags" {
  type    = map(string)
  default = null
}`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	helper.AssertIssues(t, helper.Issues{}, runner.Issues)
}

func TestResourceTagsRejectsMalformedRecursiveShapes(t *testing.T) {
	rule := requiredInterfaceRule(t, "resource_tags")
	expectation := "must use optional `resources` and `modules` namespaces with optional `map(string)` resource leaves and recursive module objects"
	cases := map[string]string{
		"top level map": `variable "resource_tags" {
  type    = map(map(string))
  default = null
}`,
		"unknown namespace": strings.Replace(validResourceTags, "    modules = optional(object({", "    children = optional(object({", 1),
		"flat resource leaf": strings.Replace(
			validResourceTags,
			`    resources = optional(object({
      this      = optional(map(string))
      secondary = optional(map(string))
    }))`,
			`    this = optional(map(string))`,
			1,
		),
		"required resources namespace": `variable "resource_tags" {
  type = object({
    resources = object({
      this = optional(map(string))
    })
  })
  default = null
}`,
		"resources namespace default": strings.Replace(validResourceTags, "    }))\n    modules", "    }), {})\n    modules", 1),
		"required resource leaf":      strings.Replace(validResourceTags, "this      = optional(map(string))", "this      = map(string)", 1),
		"wrong resource leaf type":    strings.Replace(validResourceTags, "secondary = optional(map(string))", "secondary = optional(list(string))", 1),
		"resource leaf default":       strings.Replace(validResourceTags, "this      = optional(map(string))", "this      = optional(map(string), {})", 1),
		"required module": `variable "resource_tags" {
  type = object({
    modules = optional(object({
      child = object({
        resources = optional(object({
          this = optional(map(string))
        }))
      })
    }))
  })
  default = null
}`,
		"module default": strings.Replace(validResourceTags, "      }))\n    }))", "      }), {})\n    }))", 1),
		"flattened child shape": strings.Replace(
			validResourceTags,
			`      child = optional(object({
        resources = optional(object({
          this = optional(map(string))
        }))
      }))`,
			`      child = optional(object({
        this = optional(map(string))
      }))`,
			1,
		),
		"empty resources namespace": strings.Replace(
			validResourceTags,
			`    resources = optional(object({
      this      = optional(map(string))
      secondary = optional(map(string))
    }))`,
			`    resources = optional(object({}))`,
			1,
		),
		"empty top level object": `variable "resource_tags" {
  type    = object({})
  default = null
}`,
	}

	for name, content := range cases {
		content := content
		t.Run(name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"variables.tf": content})
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

func TestResourceTagsRejectsDefaultAndNullability(t *testing.T) {
	rule := requiredInterfaceRule(t, "resource_tags")
	cases := []struct {
		name        string
		content     string
		attribute   string
		expectation string
	}{
		{
			name:        "missing default",
			content:     strings.Replace(validResourceTags, "  default = null\n", "", 1),
			attribute:   "default",
			expectation: "must be `null`",
		},
		{
			name:        "empty object default",
			content:     strings.Replace(validResourceTags, "default = null", "default = {}", 1),
			attribute:   "default",
			expectation: "must be `null`",
		},
		{
			name:        "nullable false",
			content:     strings.Replace(validResourceTags, "  default = null\n", "  default  = null\n  nullable = false\n", 1),
			attribute:   "nullable",
			expectation: "must permit `null`",
		},
	}

	runMalformedVariableCases(t, rule, cases)
}

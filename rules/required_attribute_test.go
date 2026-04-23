package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TestAzapiRequiredAttributeRule validates the generic rule behaviour using an artificial attribute.
func TestAzapiRequiredAttributeRule(t *testing.T) {
	rule := NewRequiredAttributeRule(
		"azapi_resource_required_fake_attr", // name
		"https://example.invalid/docs/fake", // link
		"resource",                          // block type
		[]string{"azapi_resource"},          // first label allow list
		"fake_required",                     // attribute name
		"\"default\"",                       // suggestion text (string literal)
		tflint.ERROR,
	)

	cases := []struct {
		desc     string
		config   string
		expected helper.Issues
	}{
		{
			desc:     "no matching resources => no issue",
			config:   `resource "random_pet" "example" {}`,
			expected: helper.Issues{},
		},
		{
			desc: "matching resource missing attribute => issue",
			config: `resource "azapi_resource" "ex" {
  type      = "Microsoft.Example/examples@2023-06-01"
  name      = "ex1"
  location  = "westus2"
  parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1"
}`,
			expected: helper.Issues{{
				Rule:    rule,
				Message: "`fake_required` attribute must be specified (suggested default: \"default\")",
			}},
		},
		{
			desc: "matching resource with attribute => ok",
			config: `resource "azapi_resource" "ex" {
  type      = "Microsoft.Example/examples@2023-06-01"
  name      = "ex1"
  location  = "westus2"
  parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1"
  fake_required = "custom"
}`,
			expected: helper.Issues{},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.desc, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"main.tf": c.config})
			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			helper.AssertIssuesWithoutRange(t, c.expected, runner.Issues)
		})
	}
}

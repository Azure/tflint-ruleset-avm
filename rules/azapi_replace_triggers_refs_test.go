package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/jmespath/go-jmespath"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestAzapiReplaceTriggersRefsRule(t *testing.T) {
	tests := []struct {
		name          string
		resourceType  string
		body          string
		triggers      string
		expectedIssue string
	}{
		{
			name:         "omitted",
			resourceType: "azapi_resource",
			body:         `{ properties = { sku = { name = "Standard" } } }`,
		},
		{
			name:          "empty",
			resourceType:  "azapi_resource",
			body:          `{ properties = {} }`,
			triggers:      "[]",
			expectedIssue: "omit `replace_triggers_refs`",
		},
		{
			name:         "valid static path",
			resourceType: "azapi_resource",
			body:         `{ properties = { sku = { name = "Standard" } } }`,
			triggers:     `["properties.sku.name"]`,
		},
		{
			name:         "valid data plane path",
			resourceType: "azapi_data_plane_resource",
			body:         `{ properties = { enabled = true } }`,
			triggers:     `["properties.enabled"]`,
		},
		{
			name:          "blank path",
			resourceType:  "azapi_resource",
			body:          `{ properties = {} }`,
			triggers:      `[""]`,
			expectedIssue: "must not be blank",
		},
		{
			name:          "surrounding whitespace",
			resourceType:  "azapi_resource",
			body:          `{ properties = { sku = "Standard" } }`,
			triggers:      `[" properties.sku"]`,
			expectedIssue: "must not contain surrounding whitespace",
		},
		{
			name:          "duplicate path",
			resourceType:  "azapi_resource",
			body:          `{ properties = { sku = "Standard" } }`,
			triggers:      `["properties.sku", "properties.sku"]`,
			expectedIssue: "is duplicated",
		},
		{
			name:          "name is redundant",
			resourceType:  "azapi_resource",
			body:          `{ properties = {} }`,
			triggers:      `["name"]`,
			expectedIssue: "is redundant",
		},
		{
			name:          "location is redundant",
			resourceType:  "azapi_resource",
			body:          `{ properties = {} }`,
			triggers:      `["location"]`,
			expectedIssue: "is redundant",
		},
		{
			name:          "invalid JMESPath",
			resourceType:  "azapi_resource",
			body:          `{ properties = {} }`,
			triggers:      `["properties.["]`,
			expectedIssue: "is not valid JMESPath",
		},
		{
			name:          "missing static path",
			resourceType:  "azapi_resource",
			body:          `{ properties = { sku = "Standard" } }`,
			triggers:      `["properties.kind"]`,
			expectedIssue: "does not resolve against the static resource body",
		},
		{
			name:          "missing projected field",
			resourceType:  "azapi_resource",
			body:          `{ properties = { items = [{ existing = true }] } }`,
			triggers:      `["properties.items[].missing"]`,
			expectedIssue: "does not resolve against the static resource body",
		},
		{
			name:         "projection over empty collection is valid",
			resourceType: "azapi_resource",
			body:         `{ properties = { items = [] } }`,
			triggers:     `["properties.items[].name"]`,
		},
		{
			name:          "omitted body cannot resolve declared path",
			resourceType:  "azapi_resource",
			triggers:      `["properties.sku"]`,
			expectedIssue: "does not resolve against the static resource body",
		},
		{
			name:          "explicit null body cannot resolve declared path",
			resourceType:  "azapi_resource",
			body:          `null`,
			triggers:      `["properties.sku"]`,
			expectedIssue: "does not resolve against the static resource body",
		},
		{
			name:          "JMESPath type error becomes issue",
			resourceType:  "azapi_resource",
			body:          `{ properties = { sku = "Standard" } }`,
			triggers:      `["sort(properties.sku)"]`,
			expectedIssue: "cannot be evaluated against the static resource body",
		},
		{
			name:          "JMESPath null type error becomes issue",
			resourceType:  "azapi_resource",
			body:          `{ properties = { sku = null } }`,
			triggers:      `["length(properties.sku)"]`,
			expectedIssue: "cannot be evaluated against the static resource body",
		},
		{
			name:         "static null path resolves",
			resourceType: "azapi_resource",
			body:         `{ properties = { sku = null } }`,
			triggers:     `["properties.sku"]`,
		},
		{
			name:         "dynamic body skips resolution",
			resourceType: "azapi_resource",
			body:         `var.body`,
			triggers:     `["properties.sku"]`,
		},
		{
			name:          "dynamic trigger list is rejected",
			resourceType:  "azapi_resource",
			body:          `{ properties = {} }`,
			triggers:      `var.replace_triggers_refs`,
			expectedIssue: "must be a statically known list",
		},
		{
			name:          "object trigger value is rejected",
			resourceType:  "azapi_resource",
			body:          `{ properties = {} }`,
			triggers:      `{ path = "properties.sku" }`,
			expectedIssue: "must be a statically known list",
		},
		{
			name:         "unrelated resource is ignored",
			resourceType: "random_string",
			body:         `{}`,
			triggers:     `[]`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			runner := helper.TestRunner(t, map[string]string{
				"main.tf": replaceTriggersResource(test.resourceType, test.body, test.triggers),
			})
			rule := NewAzapiReplaceTriggersRefsRule()
			require.NoError(t, rule.Check(runner))
			if test.expectedIssue == "" {
				assert.Empty(t, runner.Issues)
				return
			}
			require.Len(t, runner.Issues, 1)
			assert.Contains(t, runner.Issues[0].Message, test.expectedIssue)
		})
	}
}

func TestStaticJSONValue_preservesNullProperties(t *testing.T) {
	expression, diagnostics := hclsyntax.ParseExpression(
		[]byte(`{ properties = { sku = null } }`),
		"test.tf",
		hcl.InitialPos,
	)
	require.False(t, diagnostics.HasErrors())

	value, known := staticJSONValue(expression)
	require.True(t, known)
	result, err := jmespath.Search("properties.sku", value)
	require.NoError(t, err)
	assert.Nil(t, result)
	assert.True(t, simpleObjectPathExists("properties.sku", value))
	assert.False(t, simpleObjectPathExists("properties.missing", value))
}

func replaceTriggersResource(resourceType, body, triggers string) string {
	bodyAttribute := ""
	if body != "" {
		bodyAttribute = "\n  body = " + body
	}
	triggerAttribute := ""
	if triggers != "" {
		triggerAttribute = "\n  replace_triggers_refs = " + triggers
	}
	return `resource "` + resourceType + `" "this" {` + bodyAttribute + triggerAttribute + `
}`
}

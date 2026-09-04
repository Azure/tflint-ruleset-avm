package rules

import (
	"errors"
	"testing"

	"github.com/Azure/tflint-ruleset-avm/internal/tagcapability"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestAzapiResourceTagRule(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		status        tagcapability.Status
		resolveErr    error
		expectedIssue string
	}{
		{
			name: "supported resource propagates tags",
			content: azapiResource(`Microsoft.Resources/resourceGroups@2021-04-01`, `
  tags = var.tags`),
			status: tagcapability.StatusWritable,
		},
		{
			name:          "supported resource is missing tags",
			content:       azapiResource(`Microsoft.Resources/resourceGroups@2021-04-01`, ""),
			status:        tagcapability.StatusWritable,
			expectedIssue: "must set `tags`",
		},
		{
			name: "supported resource uses per-resource tags",
			content: azapiResource(`Microsoft.Resources/resourceGroups@2021-04-01`, `
  tags = try(var.resource_tags.resources.this, null) != null ? var.resource_tags.resources.this : var.tags`),
			status: tagcapability.StatusWritable,
		},
		{
			name: "supported resource uses local tags expression",
			content: azapiResource(`Microsoft.Resources/resourceGroups@2021-04-01`, `
  tags = local.resource_tags`),
			status: tagcapability.StatusWritable,
		},
		{
			name:    "unsupported resource omits tags",
			content: azapiResource(`Microsoft.Authorization/roleAssignments@2022-04-01`, ""),
			status:  tagcapability.StatusUnsupported,
		},
		{
			name: "unsupported resource sets tags",
			content: azapiResource(`Microsoft.Authorization/roleAssignments@2022-04-01`, `
  tags = var.tags`),
			status:        tagcapability.StatusUnsupported,
			expectedIssue: "must not set `tags`",
		},
		{
			name: "read-only tags are rejected",
			content: azapiResource(`Microsoft.Example/widgets@2026-01-01`, `
  tags = var.tags`),
			status:        tagcapability.StatusReadOnly,
			expectedIssue: "must not set `tags`",
		},
		{
			name:       "unavailable schema is skipped",
			content:    azapiResource(`Microsoft.Unknown/widgets@2026-01-01`, ""),
			resolveErr: tagcapability.ErrUnknownResource,
		},
		{
			name: "dynamic resource type is skipped",
			content: `
resource "azapi_resource" "this" {
  type = var.resource_type
}`,
			status: tagcapability.StatusWritable,
		},
		{
			name: "non-AzAPI resource is ignored",
			content: `
resource "random_string" "this" {
  length = 8
}`,
			status: tagcapability.StatusWritable,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			calls := 0
			rule := newAzapiResourceTagRule(func(resourceType string) (tagcapability.Status, error) {
				calls++
				return test.status, test.resolveErr
			})
			runner := helper.TestRunner(t, map[string]string{"main.tf": test.content})

			err := rule.Check(runner)
			require.NoError(t, err)
			if test.expectedIssue == "" {
				assert.Empty(t, runner.Issues)
			} else {
				require.Len(t, runner.Issues, 1)
				assert.Contains(t, runner.Issues[0].Message, test.expectedIssue)
			}
			if test.name == "dynamic resource type is skipped" || test.name == "non-AzAPI resource is ignored" {
				assert.Zero(t, calls)
			} else {
				assert.Equal(t, 1, calls)
			}
		})
	}
}

func TestAzapiResourceTagRule_unexpectedResolverError(t *testing.T) {
	rule := newAzapiResourceTagRule(func(string) (tagcapability.Status, error) {
		return "", errors.New("malformed embedded snapshot")
	})
	runner := helper.TestRunner(t, map[string]string{
		"main.tf": azapiResource(`Microsoft.Example/widgets@2026-01-01`, ""),
	})

	err := rule.Check(runner)
	require.ErrorContains(t, err, "resolve tags support")
}

func TestAzapiResourceTagRule_unexpectedCapabilityStatus(t *testing.T) {
	rule := newAzapiResourceTagRule(func(string) (tagcapability.Status, error) {
		return "unknown", nil
	})
	runner := helper.TestRunner(t, map[string]string{
		"main.tf": azapiResource(`Microsoft.Example/widgets@2026-01-01`, ""),
	})

	err := rule.Check(runner)
	require.ErrorContains(t, err, "unexpected property status")
}

func TestAzapiResourceTagRule_embeddedSnapshot(t *testing.T) {
	rule := NewAzapiResourceTagRule()
	runner := helper.TestRunner(t, map[string]string{
		"main.tf": azapiResource(`Microsoft.Resources/resourceGroups@2021-04-01`, `
  tags = var.tags`),
	})

	require.NoError(t, rule.Check(runner))
	assert.Empty(t, runner.Issues)
}

func azapiResource(resourceType, body string) string {
	return `resource "azapi_resource" "this" {
  type = "` + resourceType + `"` + body + `
}`
}

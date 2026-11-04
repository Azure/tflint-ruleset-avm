package rules

import (
	"fmt"
	"testing"

	"github.com/Azure/tflint-ruleset-avm/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformCountForEachUsage(t *testing.T) {
	rule := NewTerraformCountForEachUsageRule()
	runRuleTestCases(t, rule, []ruleTestCase{
		{
			Name: "conditional count and map or set for_each are allowed",
			Files: map[string]string{
				"variables.tf": `
variable "subnets" {
  type     = map(string)
  default  = {}
  nullable = false
}

variable "names" {
  type     = list(string)
  default  = []
  nullable = false
}`,
				"main.tf": `
resource "azapi_resource" "lock" {
  count = length(var.names) > 0 ? 1 : 0
}

resource "azapi_resource" "subnet" {
  for_each = var.subnets
}

resource "azapi_resource" "named" {
  for_each = toset(var.names)
}

resource "azapi_resource" "by_name" {
  for_each = { for i, name in var.names : name => i }
}

module "child" {
  source   = "./child"
  for_each = { for k, v in var.subnets : k => v }
}

data "azapi_resource" "existing" {
  count = length(var.names)
}`,
			},
			Expected: helper.Issues{},
		},
		{
			Name: "count from length is flagged",
			Files: map[string]string{"main.tf": `
resource "azapi_resource" "subnet" {
  count = length(var.names)
}

module "child" {
  source = "./child"
  count  = var.enabled ? (length(local.names)) : 0
}`},
			Expected: helper.Issues{
				{Rule: rule, Message: "`count` must not be set from `length(...)`; use `for_each` with a map or set so instances keep stable keys"},
				{Rule: rule, Message: "`count` must not be set from `length(...)`; use `for_each` with a map or set so instances keep stable keys"},
			},
		},
		{
			Name: "for_each over list variables is flagged",
			Files: map[string]string{
				"variables.tf": `
variable "names" {
  type     = list(string)
  nullable = false
}

variable "pairs" {
  type     = tuple([string, string])
  nullable = false
}`,
				"main.tf": `
resource "azapi_resource" "direct" {
  for_each = var.names
}

resource "azapi_resource" "indexed" {
  for_each = { for i, name in var.pairs : i => name }
}`,
			},
			Expected: helper.Issues{
				{Rule: rule, Message: "`for_each` must not iterate list variable `names`; use a map or set variable instead"},
				{Rule: rule, Message: "`for_each` must not key instances by the position of list variable `pairs`; key them by a stable value instead"},
			},
		},
	})
}

func Test_TerraformCountForEachUsageIgnoredResources(t *testing.T) {
	rule := NewTerraformCountForEachUsageRule()
	files := map[string]string{"main.tf": `
resource "azapi_resource" "subnet" {
  count = length(var.names)
}`}
	config := fmt.Sprintf("rule %q {\n  enabled           = true\n  ignored_resources = [\"azapi_resource.subnet\"]\n}\n", rule.Name())
	runner := helper.TestRunner(t, withRuleConfig(files, config))

	require.NoError(t, common.NewConfigurableRule(rule).Check(runner))

	assert.Empty(t, runner.Issues)
}

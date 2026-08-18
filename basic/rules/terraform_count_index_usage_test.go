package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"testing"
)

func Test_TerraformCountIndexUsageRule(t *testing.T) {
	cases := []struct {
		Name     string
		JSON     bool
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "1. simple improper count.index usage example",
			Content: `
resource "azurerm_resource_group" "default" {
  count = length(var.my_list)

  name     = var.my_list[count.index]
  location = "west"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformCountIndexUsageRule(),
					Message: "`count.index` is not recommended to be used as the subscript of list/map, use for_each instead",
				},
			},
		},
		{
			Name: "2. complex improper count.index usage example1",
			Content: `
resource "azurerm_resource_group" "default" {
  count = length(var.my_list)

  name     = local.a_map["${max(3, count.index)}"]
  location = "west"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformCountIndexUsageRule(),
					Message: "`count.index` is not recommended to be used as the subscript of list/map, use for_each instead",
				},
			},
		},
		{
			Name: "3. complex improper count.index usage example2",
			Content: `
resource "azurerm_resource_group" "default" {
  count = length(var.my_list)

  name     = local.a_map[max(var.my_list[count.index], 3)*2 + local.set[0]]
  location = "west"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformCountIndexUsageRule(),
					Message: "`count.index` is not recommended to be used as the subscript of list/map, use for_each instead",
				},
			},
		},
		{
			Name: "4. multiple places with improper count.index usage",
			Content: `
resource "azurerm_resource_group" "default1" {
  count = length(var.my_list)

  name     = var.my_list[count.index]
  location = "west"
}

resource "azurerm_resource_group" "default2" {
  count = length(var.my_list)

  name     = local.a_map["${max(3, count.index)}"]
  location = "west"
}

resource "azurerm_resource_group" "default3" {
  count = length(var.my_list)

  name     = local.a_map[max(var.my_list[count.index], 3)*2 + local.set[0]]
  location = "west"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformCountIndexUsageRule(),
					Message: "`count.index` is not recommended to be used as the subscript of list/map, use for_each instead",
				},
				{
					Rule:    NewTerraformCountIndexUsageRule(),
					Message: "`count.index` is not recommended to be used as the subscript of list/map, use for_each instead",
				},
				{
					Rule:    NewTerraformCountIndexUsageRule(),
					Message: "`count.index` is not recommended to be used as the subscript of list/map, use for_each instead",
				},
			},
		},
		{
			Name: "5. proper count.index usage",
			Content: `
resource "azurerm_resource_group" "default1" {
  count = length(var.my_list)

  name     = "my resource ${count.index}"
  location = "west"
}

resource "azurerm_resource_group" "default2" {
  count = length(var.my_list)

  name     = local.a_map["my_resource"]
  location = "west"
}`,
			Expected: helper.Issues{},
		},
	}
	rule := NewTerraformCountIndexUsageRule()
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			filename := "config.tf"
			if tc.JSON {
				filename = "config.tf.json"
			}
			runner := helper.TestRunner(t, map[string]string{filename: tc.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}
			AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}

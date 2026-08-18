package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"testing"
)

func Test_TerraformHeredocUsageRule(t *testing.T) {
	cases := []struct {
		Name     string
		JSON     bool
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "1. use heredoc to parse JSON",
			Content: `
resource "azurerm_resource_group_policy_assignment" "example" {
  name                 = "example"
  resource_group_id    = azurerm_resource_group.example.id
  policy_definition_id = azurerm_policy_definition.example.id

  parameters = <<-JSON
{
    "BigIntSupported": 995815895020119788889,
    "date": "20180322",
    "message": "Success !",
    "status": 200,
    "city": "北京",
    "count": 632,
    "data": {
        "shidu": "34%",
        "pm25": 73,
        "pm10": 91,
        "quality": "良",
        "wendu": "5",
        "ganmao": "极少数敏感人群应减少户外活动",
        "yesterday": {
            "date": "21日星期三",
            "sunrise": "06:19",
            "high": "高温 11.0℃",
            "low": "低温 1.0℃",
            "sunset": "18:26",
            "aqi": 85,
            "fx": "南风",
            "fl": "<3级",
            "type": "多云",
            "notice": "阴晴之间，谨防紫外线侵扰"
        },
        "forecast": [
            {
                "date": "22日星期四",
                "sunrise": "06:17",
                "high": "高温 17.0℃",
                "low": "低温 1.0℃",
                "sunset": "18:27",
                "aqi": 98,
                "fx": "西南风",
                "fl": "<3级",
                "type": "晴",
                "notice": "愿你拥有比阳光明媚的心情"
            },
            {
                "date": "23日星期五",
                "sunrise": "06:16",
                "high": "高温 18.0℃",
                "low": "低温 5.0℃",
                "sunset": "18:28",
                "aqi": 118,
                "fx": "无持续风向",
                "fl": "<3级",
                "type": "多云",
                "notice": "阴晴之间，谨防紫外线侵扰"
            },
            {
                "date": "24日星期六",
                "sunrise": "06:14",
                "high": "高温 21.0℃",
                "low": "低温 7.0℃",
                "sunset": "18:29",
                "aqi": 52,
                "fx": "西南风",
                "fl": "<3级",
                "type": "晴",
                "notice": "愿你拥有比阳光明媚的心情"
            },
            {
                "date": "25日星期日",
                "sunrise": "06:13",
                "high": "高温 22.0℃",
                "low": "低温 7.0℃",
                "sunset": "18:30",
                "aqi": 71,
                "fx": "西南风",
                "fl": "<3级",
                "type": "晴",
                "notice": "愿你拥有比阳光明媚的心情"
            },
            {
                "date": "26日星期一",
                "sunrise": "06:11",
                "high": "高温 21.0℃",
                "low": "低温 8.0℃",
                "sunset": "18:31",
                "aqi": 97,
                "fx": "西南风",
                "fl": "<3级",
                "type": "多云",
                "notice": "阴晴之间，谨防紫外线侵扰"
            }
        ]
    }
}
  JSON
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformHeredocUsageRule(),
					Message: "for JSON, instead of HEREDOC, use a combination of a `local` and the `jsonencode` function",
				},
			},
		},
		{
			Name: "2. use HEREDOC to parse YAML",
			Content: `
resource "azurerm_resource_group_policy_assignment" "example" {
  name                 = "example"
  resource_group_id    = azurerm_resource_group.example.id
  policy_definition_id = azurerm_policy_definition.example.id

  parameters = <<-YAML
name: release

on:
  push:
    branches:
    - '!*'
    tags:
    - v*.*.*

jobs:
  goreleaser:
    runs-on: ubuntu-latest
    steps:
    - name: Checkout
      uses: actions/checkout@v3
      with:
        fetch-depth: 0
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.18
    - name: Import GPG key
      id: import_gpg
      uses: crazy-max/ghaction-import-gpg@v4
    - name: Run GoReleaser
      uses: goreleaser/goreleaser-action@v3
      with:
        version: v0.178.0
        args: release --rm-dist
  YAML
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformHeredocUsageRule(),
					Message: "for YAML, instead of HEREDOC, use a combination of a `local` and the `yamlencode` function",
				},
			},
		},
		{
			Name: "3. simple non-JSON/YAML text",
			Content: `
resource "azurerm_resource_group_policy_assignment" "example1" {
  name                 = "example"
  resource_group_id    = azurerm_resource_group.example.id
  policy_definition_id = azurerm_policy_definition.example.id

  parameters = <<-EMPTY
  EMPTY
}

resource "azurerm_resource_group_policy_assignment" "example2" {
  name                 = "example"
  resource_group_id    = azurerm_resource_group.example.id
  policy_definition_id = azurerm_policy_definition.example.id

  parameters = <<-SENTENCE
  hello, world!
  SENTENCE
}

resource "azurerm_resource_group_policy_assignment" "example3" {
  name                 = "example"
  resource_group_id    = azurerm_resource_group.example.id
  policy_definition_id = azurerm_policy_definition.example.id

  parameters = <<-EXP
  tmp = "test"
  EXP
}

resource "azurerm_resource_group_policy_assignment" "example4" {
  name                 = "example"
  resource_group_id    = azurerm_resource_group.example.id
  policy_definition_id = azurerm_policy_definition.example.id

  parameters = <<-MULTILINES
line1
  line2
  MULTILINES
}`,
			Expected: helper.Issues{},
		},
	}
	rule := NewTerraformHeredocUsageRule()
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

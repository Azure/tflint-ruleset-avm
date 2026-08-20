package rules

import (
	"fmt"
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestTerraformHeredocUsageRule(t *testing.T) {
	tests := []struct {
		name    string
		content string
		message string
	}{
		{
			name:    "literal JSON",
			content: heredoc(`{"name":"example"}`),
			message: "Use `jsonencode` instead of a heredoc for literal JSON",
		},
		{
			name: "literal YAML mapping",
			content: heredoc(`name: example
enabled: true`),
			message: "Use `yamlencode` instead of a heredoc for literal YAML",
		},
		{
			name:    "plain text",
			content: heredoc("plain text"),
		},
		{
			name:    "JSON scalar",
			content: heredoc("true"),
			message: "Use `jsonencode` instead of a heredoc for literal JSON",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			runner := helper.TestRunner(t, map[string]string{"main.tf": test.content})
			rule := NewTerraformHeredocUsageRule()
			if err := rule.Check(runner); err != nil {
				t.Fatal(err)
			}
			if test.message == "" {
				AssertIssues(t, helper.Issues{}, runner.Issues)
				return
			}
			AssertIssues(t, helper.Issues{
				{
					Rule:    rule,
					Message: test.message,
					Range:   runner.Issues[0].Range,
				},
			}, runner.Issues)
		})
	}
}

func heredoc(content string) string {
	return fmt.Sprintf("value = <<-EOT\n%s\nEOT\n", content)
}

package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformOutputSeperateRule(t *testing.T) {
	cases := []struct {
		Name     string
		JSON     bool
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "1. no output",
			Content: `
terraform {}`,
			Expected: helper.Issues{},
		},
		{
			Name: "2. putting output and other blocks together in the same file",
			Content: `
terraform {}

output "api_base_url" {
  value = "https://${aws_instance.example.private_dns}:8433/"

  # The EC2 instance must have an encrypted root volume.
  precondition {
    condition     = data.aws_ebs_volume.example.encrypted
    error_message = "The server's root volume is not encrypted."
  }
}

output "db_password" {
  value       = aws_db_instance.db.password
  description = "The password for logging in to the database."
  sensitive   = true
}

output "instance_ip_addr" {
  value       = aws_instance.server.private_ip
  description = "The private IP address of the main server instance."
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformOutputSeparateRule(),
					Message: "Putting outputs and other types of blocks in the same file is not recommended",
				},
			},
		},
	}
	rule := NewTerraformOutputSeparateRule()

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

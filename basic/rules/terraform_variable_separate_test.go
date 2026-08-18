package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformVariableSeperateRule(t *testing.T) {
	cases := []struct {
		Name     string
		JSON     bool
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "1. no variable",
			Content: `
terraform{}`,
			Expected: helper.Issues{},
		},
		{
			Name: "2. putting variable and other blocks together in the same file",
			Content: `
terraform{}

variable "availability_zone_names" {
  type    = list(string)
  default = ["us-west-1a"]
}

variable "docker_ports" {
  type = list(object({
    internal = number
    external = number
    protocol = string
  }))
  default = [
    {
      internal = 8300
      external = 8300
      protocol = "tcp"
    }
  ]
}

variable "image_id" {
  type = string
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformVariableSeparateRule(),
					Message: "Putting variables and other types of blocks in the same file is not recommended",
				},
			},
		},
	}
	rule := NewTerraformVariableSeparateRule()

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

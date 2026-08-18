package rules

import (
	"github.com/stretchr/testify/assert"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"testing"
)

func Test_TerraformSensitiveVariableNoDefault(t *testing.T) {
	cases := []struct {
		Name     string
		JSON     bool
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "1. correct cases",
			Content: `
variable "availability_zone_addr" {
  type = string
}

variable "availability_zone_tag" {
  type      = string
  sensitive = true
}

variable "availability_zone_names" {
  type    = list(string)
  default = ["us-west-1a"]
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "2. single sensitive varible with default value set",
			Content: `
variable "availability_zone_names" {
  type      = list(string)
  default   = ["us-west-1a"]
  sensitive = true
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformSensitiveVariableNoDefaultRule(),
					Message: "Default value is not expected to be set for sensitive variable `availability_zone_names`",
				},
			},
		},
		{
			Name: "3. multiple sensitive varibles with default value set",
			Content: `
variable "availability_zone_addr" {
  type = string
}

variable "availability_zone_names" {
  type      = list(string)
  default   = ["us-west-1a"]
  sensitive = true
}

variable "availability_zone_tag" {
  type      = string
  default   = "test"
  sensitive = true
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformSensitiveVariableNoDefaultRule(),
					Message: "Default value is not expected to be set for sensitive variable `availability_zone_names`",
				},
				{
					Rule:    NewTerraformSensitiveVariableNoDefaultRule(),
					Message: "Default value is not expected to be set for sensitive variable `availability_zone_tag`",
				},
			},
		},
		{
			Name: "4. empty list should be ok",
			Content: `
variable "availability_zone_names" {
  type      = list(string)
  default   = []
  sensitive = true
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "5. empty set should be ok",
			Content: `
variable "availability_zone_names" {
  type      = set(string)
  default   = []
  sensitive = true
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "6. empty map should be ok",
			Content: `
variable "availability_zone_names" {
  type      = map(string)
  default   = {}
  sensitive = true
}
`,
			Expected: helper.Issues{},
		},
	}
	rule := NewTerraformSensitiveVariableNoDefaultRule()

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

func Test_TerraformSensitiveVariableNoDefaultValue_NullDefaultValueShouldNotRaiseError(t *testing.T) {
	filename := "config.tf"
	code := `variable "password" {
  type      = string
  default   = null
  sensitive = true
}`
	runner := helper.TestRunner(t, map[string]string{filename: code})
	sut := NewTerraformSensitiveVariableNoDefaultRule()
	if err := sut.Check(runner); err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}

	assert.Equal(t, 0, len(runner.Issues))
}

package rules

import (
	"github.com/stretchr/testify/require"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"testing"
)

func Test_TerraformRequiredVersionDeclaration(t *testing.T) {
	cases := []struct {
		Name     string
		JSON     bool
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "1. correct case",
			Content: `
terraform {
  required_version = "~> 0.12.29"
  required_providers {
    aws = {
      version = ">= 2.7.0"
      source = "hashicorp/aws"
    }
  }
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "2. no required_version set for terraform block",
			Content: `
terraform {
  required_providers {
    aws = {
      version = ">= 2.7.0"
      source = "hashicorp/aws"
    }
  }
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformRequiredVersionDeclarationRule(),
					Message: "The `required_version` field should be declared at the beginning of `terraform` block",
				},
			},
		},
		{
			Name: "3. required_version is not placed at the beginning of terraform block",
			Content: `
terraform {
  required_providers {
    aws = {
      version = ">= 2.7.0"
      source = "hashicorp/aws"
    }
  }
  required_version = "~> 0.12.29"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformRequiredVersionDeclarationRule(),
					Message: "The `required_version` field should be declared at the beginning of `terraform` block",
				},
			},
		},
	}
	rule := NewTerraformRequiredVersionDeclarationRule()

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

func Test_TerraformRequiredVersionDeclaration_SkipOverrideFile(t *testing.T) {
	filenames := []string{"override.tf", "terraform_override.tf"}
	for _, filename := range filenames {
		t.Run(filename, func(t *testing.T) {
			rule := NewTerraformRequiredVersionDeclarationRule()
			runner := helper.TestRunner(t, map[string]string{filename: `
terraform {
  required_providers {
    aws = {
      version = ">= 2.7.0"
      source = "hashicorp/aws"
    }
  }
}`})
			require.NoError(t, rule.Check(runner))
			AssertIssues(t, helper.Issues{}, runner.Issues)
		})
	}
}

package rules

import (
	"github.com/stretchr/testify/require"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"testing"
)

func Test_TerraformRequiredProvidersDeclaration(t *testing.T) {
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
      source = "hashicorp/aws"
      version = ">= 2.7.0"
    }
	azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0.2"
    }
  }
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "2. required_providers not declared in terraform block",
			Content: `
terraform {
  required_version = "~> 0.12.29"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformRequiredProvidersDeclarationRule(),
					Message: "The `required_providers` field should be declared in `terraform` block",
				},
			},
		},
		{
			Name: "3. args in required_providers block are not sorted in alphabetic order",
			Content: `
terraform {
  required_version = "~> 0.12.29"
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0.2"
    }
    aws = {
      source = "hashicorp/aws"
      version = ">= 2.7.0"
    }
  }
}`,
			Expected: helper.Issues{
				{
					Rule: NewTerraformRequiredProvidersDeclarationRule(),
					Message: "The arguments of `required_providers` are expected to be sorted as follows:" + `
required_providers {
  aws = {
    source  = "hashicorp/aws"
    version = ">= 2.7.0"
  }
  azurerm = {
    source  = "hashicorp/azurerm"
    version = "~> 3.0.2"
  }
}`,
				},
			},
		},
		{
			Name: "4. parameters of providers are not sorted in alphabetic order in required_providers block",
			Content: `
terraform {
  required_version = "~> 0.12.29"
  required_providers {
    aws = {
      version = ">= 2.7.0"
      source = "hashicorp/aws"
    }
    azurerm = {
      version = "~> 3.0.2"
      source  = "hashicorp/azurerm"
    }
	b = {
      version = "~> 3.0.2"
      source  = "hashicorp/azurerm"
    }
	c = {
      version = "~> 3.0.2"
      source  = "hashicorp/azurerm"
    }
	d = {
      version = "~> 3.0.2"
      source  = "hashicorp/azurerm"
    }
	e = {
      version = "~> 3.0.2"
      source  = "hashicorp/azurerm"
    }
  }
}`,
			Expected: helper.Issues{
				{
					Rule: NewTerraformRequiredProvidersDeclarationRule(),
					Message: "Parameters of provider `aws` are expected to be sorted as follows:" + `
aws = {
  source  = "hashicorp/aws"
  version = ">= 2.7.0"
}`,
				},
				{
					Rule: NewTerraformRequiredProvidersDeclarationRule(),
					Message: "Parameters of provider `azurerm` are expected to be sorted as follows:" + `
azurerm = {
  source  = "hashicorp/azurerm"
  version = "~> 3.0.2"
}`,
				},
				{
					Rule: NewTerraformRequiredProvidersDeclarationRule(),
					Message: "Parameters of provider `b` are expected to be sorted as follows:" + `
b = {
  source  = "hashicorp/azurerm"
  version = "~> 3.0.2"
}`,
				},
				{
					Rule: NewTerraformRequiredProvidersDeclarationRule(),
					Message: "Parameters of provider `c` are expected to be sorted as follows:" + `
c = {
  source  = "hashicorp/azurerm"
  version = "~> 3.0.2"
}`,
				},
				{
					Rule: NewTerraformRequiredProvidersDeclarationRule(),
					Message: "Parameters of provider `d` are expected to be sorted as follows:" + `
d = {
  source  = "hashicorp/azurerm"
  version = "~> 3.0.2"
}`,
				},
				{
					Rule: NewTerraformRequiredProvidersDeclarationRule(),
					Message: "Parameters of provider `e` are expected to be sorted as follows:" + `
e = {
  source  = "hashicorp/azurerm"
  version = "~> 3.0.2"
}`,
				},
			},
		},
		{
			Name: "5. Mixed cases",
			Content: `
terraform {
  required_version = "~> 0.12.29"
  required_providers {
    azurerm = {
      version = "~> 3.0.2"
      source  = "hashicorp/azurerm"
    }
    aws = {
      source = "hashicorp/aws"
      version = ">= 2.7.0"
    }
  }
}`,
			Expected: helper.Issues{
				{
					Rule: NewTerraformRequiredProvidersDeclarationRule(),
					Message: "The arguments of `required_providers` are expected to be sorted as follows:" + `
required_providers {
  aws = {
    source  = "hashicorp/aws"
    version = ">= 2.7.0"
  }
  azurerm = {
    source  = "hashicorp/azurerm"
    version = "~> 3.0.2"
  }
}`,
				},
			},
		},
		{
			Name: "6. Empty required_providers block",
			Content: `
terraform {
  required_version = "~> 0.12.29"
  required_providers {}
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "7. No parameter and only 1 parameter for provider",
			Content: `
terraform {
  required_version = "~> 0.12.29"
  required_providers {
    aws = {}
    azurerm = "~> 3.0.2"
  }
}`,
			Expected: helper.Issues{},
		},
	}
	rule := NewTerraformRequiredProvidersDeclarationRule()

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

func Test_TerraformRequiredProvidersDeclaration_SkipOverrideFile(t *testing.T) {
	filenames := []string{"override.tf", "terraform_override.tf"}
	for _, filename := range filenames {
		t.Run(filename, func(t *testing.T) {
			rule := NewTerraformRequiredProvidersDeclarationRule()
			runner := helper.TestRunner(t, map[string]string{filename: `
terraform {
  required_version = "~> 1.1.9"
}`})
			require.NoError(t, rule.Check(runner))
			AssertIssues(t, helper.Issues{}, runner.Issues)
		})
	}
}

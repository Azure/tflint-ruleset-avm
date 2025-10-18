package rules_test

import (
	"testing"

	"github.com/Azure/tflint-ruleset-avm/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestRequiredProviders(t *testing.T) {
	cases := []struct {
		desc   string
		files  map[string]string
		issues helper.Issues
	}{
		{
			desc: "required_providers exists",
			files: map[string]string{
				"terraform.tf": `terraform {
  required_version = "~> 0.12.29"
  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = ">= 2.7.0, < 3.98.0"
    }
	azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0.2"
    }
  }
}`,
			},
			issues: helper.Issues{},
		},
		{
			desc: "required_providers not declared in terraform block",
			files: map[string]string{
				"terraform.tf": `terraform {
  required_version = "~> 0.12.29"
}`,
			},
			issues: helper.Issues{
				{
					Rule:    rules.NewRequiredProvidersRule(),
					Message: "The `required_providers` field should be declared in `terraform` block",
				},
			},
		},
		{
			desc: "args in required_providers block are not sorted in alphabetic order",
			files: map[string]string{
				"terraform.tf": `terraform {
  required_version = "~> 0.12.29"
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0.2"
    }
    aws = {
      source = "hashicorp/aws"
      version = ">= 2.7.0, < 3.98.0"
    }
  }
}`,
			},
			issues: helper.Issues{
				{
					Rule: rules.NewRequiredProvidersRule(),
					Message: "The arguments of `required_providers` are expected to be sorted as follows:" + `
required_providers {
  aws = {
    source  = "hashicorp/aws"
    version = ">= 2.7.0, < 3.98.0"
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
			desc: "parameters of providers are not sorted in alphabetic order in required_providers block",
			files: map[string]string{
				"terraform.tf": `terraform {
  required_version = "~> 0.12.29"
  required_providers {
    aws = {
      version = ">= 2.7.0, < 3.98.0"
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
			},
			issues: helper.Issues{
				{
					Rule: rules.NewRequiredProvidersRule(),
					Message: "Parameters of provider `aws` are expected to be sorted as follows:" + `
aws = {
  source  = "hashicorp/aws"
  version = ">= 2.7.0, < 3.98.0"
}`,
				},
				{
					Rule: rules.NewRequiredProvidersRule(),
					Message: "Parameters of provider `azurerm` are expected to be sorted as follows:" + `
azurerm = {
  source  = "hashicorp/azurerm"
  version = "~> 3.0.2"
}`,
				},
				{
					Rule: rules.NewRequiredProvidersRule(),
					Message: "Parameters of provider `b` are expected to be sorted as follows:" + `
b = {
  source  = "hashicorp/azurerm"
  version = "~> 3.0.2"
}`,
				},
				{
					Rule: rules.NewRequiredProvidersRule(),
					Message: "Parameters of provider `c` are expected to be sorted as follows:" + `
c = {
  source  = "hashicorp/azurerm"
  version = "~> 3.0.2"
}`,
				},
				{
					Rule: rules.NewRequiredProvidersRule(),
					Message: "Parameters of provider `d` are expected to be sorted as follows:" + `
d = {
  source  = "hashicorp/azurerm"
  version = "~> 3.0.2"
}`,
				},
				{
					Rule: rules.NewRequiredProvidersRule(),
					Message: "Parameters of provider `e` are expected to be sorted as follows:" + `
e = {
  source  = "hashicorp/azurerm"
  version = "~> 3.0.2"
}`,
				},
			},
		},
		{
			desc: "mixed cases",
			files: map[string]string{
				"terraform.tf": `terraform {
  required_version = "~> 0.12.29"
  required_providers {
    azurerm = {
      version = "~> 3.0.2"
      source  = "hashicorp/azurerm"
    }
    aws = {
      source = "hashicorp/aws"
      version = ">= 2.7.0, < 3.98.0"
    }
  }
}`,
			},
			issues: helper.Issues{
				{
					Rule: rules.NewRequiredProvidersRule(),
					Message: "The arguments of `required_providers` are expected to be sorted as follows:" + `
required_providers {
  aws = {
    source  = "hashicorp/aws"
    version = ">= 2.7.0, < 3.98.0"
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
			desc: "empty required_providers block",
			files: map[string]string{
				"terraform.tf": `terraform {
  required_version = "~> 0.12.29"
  required_providers {}
}`,
			},
			issues: helper.Issues{},
		},
		{
			desc: "no parameter and only 1 parameter for provider",
			files: map[string]string{
				"terraform.tf": `terraform {
  required_version = "~> 0.12.29"
  required_providers {
    aws = {}
    azurerm = "~> 3.0.2"
  }
}`,
			},
			issues: helper.Issues{},
		},
		{
			desc: "version = 'xx.xx.xx', not ok",
			files: map[string]string{
				"terraform.tf": `terraform {
  required_version = "~> 0.12.29"
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "3.98.0"
    }
  }
}`,
			},
			issues: helper.Issues{
				{
					Rule:    rules.NewRequiredProvidersRule(),
					Message: "The `version` property constraint can use the ~> #.# or the >= #.#.#, < #.#.# format",
				},
			},
		},
		{
			desc: "version = '>= xx.xx.xx', not ok",
			files: map[string]string{
				"terraform.tf": `terraform {
  required_version = "~> 0.12.29"
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">= 3.98.0"
    }
  }
}`,
			},
			issues: helper.Issues{
				{
					Rule:    rules.NewRequiredProvidersRule(),
					Message: "The `version` property constraint can use the ~> #.# or the >= #.#.#, < #.#.# format",
				},
			},
		},
		{
			desc: "version = '< xx.xx.xx', not ok",
			files: map[string]string{
				"terraform.tf": `terraform {
  required_version = "~> 0.12.29"
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "< 3.98.0"
    }
  }
}`,
			},
			issues: helper.Issues{
				{
					Rule:    rules.NewRequiredProvidersRule(),
					Message: "The `version` property constraint can use the ~> #.# or the >= #.#.#, < #.#.# format",
				},
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			rule := rules.NewRequiredProvidersRule()

			runner := helper.TestRunner(t, tc.files)

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssuesWithoutRange(t, tc.issues, runner.Issues)
		})
	}
}

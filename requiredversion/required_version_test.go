package requiredversion_test

import (
	"testing"

	"github.com/Azure/tflint-ruleset-avm/requiredversion"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestRequiredVersion(t *testing.T) {
	cases := []struct {
		desc         string
		config       string
		requiredName string
		issues       helper.Issues
	}{
		{
			desc: "required_version = '~> #.#', ok",
			config: `terraform {
  required_version = "~> 0.12.29"
  required_providers {
    azurerm = {
      version = ">= 3.98.0"
      source = "hashicorp/azurerm"
    }
  }
}`,
			requiredName: "required_version",
			issues:       helper.Issues{},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			rule := requiredversion.NewRequiredVersionRule("required_version", "")
			filename := "variables.tf"

			runner := helper.TestRunner(t, map[string]string{filename: tc.config})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssuesWithoutRange(t, tc.issues, runner.Issues)
		})
	}
}

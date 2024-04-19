package rules_test

import (
	"testing"

	"github.com/Azure/tflint-ruleset-avm/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestModules(t *testing.T) {
	cases := []struct {
		desc   string
		files  map[string]string
		issues helper.Issues
	}{
		{
			desc: "source exists, ok",
			files: map[string]string{
				"terraform.tf": `module "other-module" {
  source  = "Azure/terraform-azurerm-avm-res-keyvault-vault"
  version = "0.5.3"
}`,
			},
			issues: helper.Issues{},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			rule := rules.NewModulesRule()

			runner := helper.TestRunner(t, tc.files)

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssuesWithoutRange(t, tc.issues, runner.Issues)
		})
	}
}

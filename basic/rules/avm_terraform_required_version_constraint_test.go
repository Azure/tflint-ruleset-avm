package rules

import (
	"fmt"
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformRequiredVersion(t *testing.T) {
	rule := NewTerraformRequiredVersionRule()
	terraformTf := func(constraint string) map[string]string {
		return map[string]string{"terraform.tf": fmt.Sprintf(`
terraform {
  required_version = %q
}`, constraint)}
	}
	minMsg := func(c string) string {
		return fmt.Sprintf("`required_version` must include a minimum Terraform version, got `%s`; use `~> #.#` or `>= #.#.#, < #.#.#`", c)
	}
	maxMsg := func(c string) string {
		return fmt.Sprintf("`required_version` must include a maximum major Terraform version, got `%s`; use `~> #.#` or `>= #.#.#, < #.#.#`", c)
	}
	nextMajorMsg := func(next, c string) string {
		return fmt.Sprintf("`required_version` must exclude Terraform %s and later, got `%s`; use `~> #.#` or `>= #.#.#, < #.#.#`", next, c)
	}

	cases := []ruleTestCase{}
	for _, valid := range []string{"~> 1.6", "~> 1.7.0", ">= 1.9, < 2.0", ">= 1.9.0, < 1.12.0", "1.9.5", ">= 1.9, <= 1.99.0"} {
		cases = append(cases, ruleTestCase{
			Name:     "valid " + valid,
			Files:    terraformTf(valid),
			Expected: helper.Issues{},
		})
	}
	cases = append(cases,
		ruleTestCase{
			Name:     "no minimum version",
			Files:    terraformTf("< 2.0"),
			Expected: helper.Issues{{Rule: rule, Message: minMsg("< 2.0")}},
		},
		ruleTestCase{
			Name:     "no maximum major version",
			Files:    terraformTf(">= 1.9"),
			Expected: helper.Issues{{Rule: rule, Message: nextMajorMsg("2.0.0", ">= 1.9")}},
		},
		ruleTestCase{
			Name:     "pessimistic major-only constraint has no upper bound",
			Files:    terraformTf("~> 1"),
			Expected: helper.Issues{{Rule: rule, Message: nextMajorMsg("2.0.0", "~> 1")}},
		},
		ruleTestCase{
			Name:     "upper bound spanning a major version",
			Files:    terraformTf(">= 1.9, < 3.0"),
			Expected: helper.Issues{{Rule: rule, Message: nextMajorMsg("2.0.0", ">= 1.9, < 3.0")}},
		},
		ruleTestCase{
			Name:  "no bounds at all",
			Files: terraformTf("!= 1.5.0"),
			Expected: helper.Issues{
				{Rule: rule, Message: minMsg("!= 1.5.0")},
				{Rule: rule, Message: maxMsg("!= 1.5.0")},
			},
		},
		ruleTestCase{
			Name:  "invalid constraint",
			Files: terraformTf("latest"),
			Expected: helper.Issues{
				{Rule: rule, Message: "`required_version` must be a valid version constraint, got `latest`"},
			},
		},
		ruleTestCase{
			Name: "missing required_version",
			Files: map[string]string{"terraform.tf": `
terraform {
  required_providers {}
}`},
			Expected: helper.Issues{
				{Rule: rule, Message: "The `terraform` block in `terraform.tf` must define `required_version`"},
			},
		},
		ruleTestCase{
			Name: "required_version in a later terraform block is found",
			Files: map[string]string{"terraform.tf": `
terraform {
  required_providers {}
}

terraform {
  required_version = "~> 1.9"
}`},
			Expected: helper.Issues{},
		},
		ruleTestCase{
			Name: "terraform blocks outside terraform.tf are ignored",
			Files: map[string]string{"versions.tf": `
terraform {
  required_version = ">= 0.13"
}`},
			Expected: helper.Issues{},
		},
		ruleTestCase{
			Name: "non-literal required_version is ignored",
			Files: map[string]string{"terraform.tf": `
terraform {
  required_version = local.version
}`},
			Expected: helper.Issues{},
		},
	)
	runRuleTestCases(t, rule, cases)
}

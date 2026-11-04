package rules

import (
	"fmt"
	"testing"

	"github.com/Azure/tflint-ruleset-avm/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/terraform/addrs"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// violatingResourceModule triggers every resource-module rule when checked as a resource module.
var violatingResourceModule = map[string]string{
	"variables.tf": `
variable "resource_group_name" {
  type = string
}`,
	"main.tf": `
resource "azapi_resource" "vault" {
  parent_id = "/subscriptions/0000/resourceGroups/rg"
}

module "alz" {
  source = "Azure/avm-ptn-alz/azurerm"
}`,
}

type childModuleRunner struct {
	*helper.Runner
}

func (childModuleRunner) GetModulePath() (addrs.Module, error) {
	return addrs.Module{"child"}, nil
}

func resourceModuleRules() []tflint.Rule {
	return []tflint.Rule{
		NewTerraformParentIDVariableRule(),
		NewAzapiPrimaryResourceNamingRule(),
		NewTerraformPatternModuleSourceRule(),
	}
}

func withRuleConfig(files map[string]string, config string) map[string]string {
	result := make(map[string]string, len(files)+1)
	for name, content := range files {
		result[name] = content
	}
	if config != "" {
		result[".tflint.hcl"] = config
	}
	return result
}

func TestResourceModuleRulesModuleClass(t *testing.T) {
	tests := []struct {
		name        string
		moduleClass string
		child       bool
		wantIssues  bool
	}{
		{name: "default class is resource", wantIssues: true},
		{name: "explicit resource class", moduleClass: "resource", wantIssues: true},
		{name: "pattern class is skipped", moduleClass: "pattern"},
		{name: "utility class is skipped", moduleClass: "utility"},
		{name: "child modules are skipped", child: true},
	}

	for _, baseRule := range resourceModuleRules() {
		for _, wrapped := range []bool{true, false} {
			for _, tc := range tests {
				t.Run(fmt.Sprintf("%s/wrapped=%t/%s", baseRule.Name(), wrapped, tc.name), func(t *testing.T) {
					config := ""
					if tc.moduleClass != "" {
						config = fmt.Sprintf("rule %q {\n  enabled      = true\n  module_class = %q\n}\n", baseRule.Name(), tc.moduleClass)
					}
					runner := helper.TestRunner(t, withRuleConfig(violatingResourceModule, config))
					rule := baseRule
					if wrapped {
						rule = common.NewConfigurableRule(baseRule)
					}

					var checkRunner tflint.Runner = runner
					if tc.child {
						checkRunner = childModuleRunner{Runner: runner}
					}
					require.NoError(t, rule.Check(checkRunner))
					if !tc.wantIssues {
						assert.Empty(t, runner.Issues)
						return
					}
					require.NotEmpty(t, runner.Issues)
					for _, issue := range runner.Issues {
						assert.Equal(t, baseRule.Name(), issue.Rule.Name())
						assert.Equal(t, tflint.ERROR, issue.Rule.Severity())
					}
				})
			}
		}
	}
}

func TestResourceModuleRulesRejectInvalidModuleClass(t *testing.T) {
	for _, baseRule := range resourceModuleRules() {
		for _, wrapped := range []bool{true, false} {
			for _, expression := range []string{`""`, `"Resource"`, `"other"`} {
				t.Run(fmt.Sprintf("%s/wrapped=%t/%s", baseRule.Name(), wrapped, expression), func(t *testing.T) {
					config := fmt.Sprintf("rule %q {\n  enabled      = true\n  module_class = %s\n}\n", baseRule.Name(), expression)
					runner := helper.TestRunner(t, withRuleConfig(violatingResourceModule, config))
					rule := baseRule
					if wrapped {
						rule = common.NewConfigurableRule(baseRule)
					}

					err := rule.Check(runner)

					require.ErrorContains(t, err, `module_class must be one of "resource", "pattern", or "utility"`)
					assert.Empty(t, runner.Issues)
				})
			}
		}
	}
}

func TestResourceModuleRulesSeverityOverride(t *testing.T) {
	for _, baseRule := range resourceModuleRules() {
		t.Run(baseRule.Name(), func(t *testing.T) {
			config := fmt.Sprintf("rule %q {\n  enabled  = true\n  severity = \"warning\"\n}\n", baseRule.Name())
			runner := helper.TestRunner(t, withRuleConfig(violatingResourceModule, config))

			require.NoError(t, common.NewConfigurableRule(baseRule).Check(runner))

			require.NotEmpty(t, runner.Issues)
			for _, issue := range runner.Issues {
				assert.Equal(t, tflint.WARNING, issue.Rule.Severity())
			}
		})
	}
}

func TestResourceModuleRulesIgnoredResources(t *testing.T) {
	files := map[string]string{"main.tf": `
resource "azapi_resource" "this" {
  parent_id = "/subscriptions/0000"
}

resource "azapi_resource" "helper" {
  parent_id = "/subscriptions/0000"
}`}
	for _, baseRule := range []tflint.Rule{NewTerraformParentIDVariableRule(), NewAzapiPrimaryResourceNamingRule()} {
		t.Run(baseRule.Name(), func(t *testing.T) {
			config := fmt.Sprintf("rule %q {\n  enabled           = true\n  ignored_resources = [\"azapi_resource.this\", \"azapi_resource.helper\"]\n}\n", baseRule.Name())
			runner := helper.TestRunner(t, withRuleConfig(files, config))

			require.NoError(t, common.NewConfigurableRule(baseRule).Check(runner))

			assert.Empty(t, runner.Issues)
		})
	}
}

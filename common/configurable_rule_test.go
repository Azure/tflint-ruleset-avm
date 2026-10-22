package common_test

import (
	"fmt"
	"testing"

	"github.com/Azure/tflint-ruleset-avm/common"
	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type emittingRule struct {
	tflint.DefaultRule
	severity tflint.Severity
}

type blockEmittingRule struct {
	tflint.DefaultRule
}

func (r *blockEmittingRule) Name() string              { return "avm_test_block_rule" }
func (r *blockEmittingRule) Enabled() bool             { return true }
func (r *blockEmittingRule) Severity() tflint.Severity { return tflint.ERROR }
func (r *blockEmittingRule) Check(runner tflint.Runner) error {
	content, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{Type: "resource", LabelNames: []string{"type", "name"}, Body: &hclext.BodySchema{}},
			{Type: "data", LabelNames: []string{"type", "name"}, Body: &hclext.BodySchema{}},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, block := range content.Blocks {
		if err := runner.EmitIssue(r, "block issue", block.DefRange); err != nil {
			return err
		}
	}
	return nil
}

type moduleClassBlockEmittingRule struct {
	blockEmittingRule
}

func (r *moduleClassBlockEmittingRule) Name() string { return "avm_test_module_class_block_rule" }
func (r *moduleClassBlockEmittingRule) CheckWithModuleClass(runner tflint.Runner, _ string) error {
	return r.Check(runner)
}

type moduleScopeBlockEmittingRule struct {
	blockEmittingRule
}

func (r *moduleScopeBlockEmittingRule) Name() string { return "avm_test_module_scope_block_rule" }
func (r *moduleScopeBlockEmittingRule) CheckWithModuleScope(runner tflint.Runner, _ string) error {
	return r.Check(runner)
}

func (r *emittingRule) Name() string              { return "avm_test_rule" }
func (r *emittingRule) Enabled() bool             { return true }
func (r *emittingRule) Severity() tflint.Severity { return r.severity }
func (r *emittingRule) Check(runner tflint.Runner) error {
	return runner.EmitIssue(r, "test issue", hcl.Range{Filename: "main.tf"})
}

func TestConfigurableRuleSeverity(t *testing.T) {
	tests := []struct {
		name            string
		defaultSeverity tflint.Severity
		configuredValue string
		wantSeverity    tflint.Severity
	}{
		{
			name:            "default is preserved",
			defaultSeverity: tflint.WARNING,
			wantSeverity:    tflint.WARNING,
		},
		{
			name:            "error override",
			defaultSeverity: tflint.NOTICE,
			configuredValue: "error",
			wantSeverity:    tflint.ERROR,
		},
		{
			name:            "warning override",
			defaultSeverity: tflint.ERROR,
			configuredValue: "warning",
			wantSeverity:    tflint.WARNING,
		},
		{
			name:            "notice override",
			defaultSeverity: tflint.ERROR,
			configuredValue: "notice",
			wantSeverity:    tflint.NOTICE,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			underlying := &emittingRule{severity: test.defaultSeverity}
			rule := common.NewConfigurableRule(underlying)
			assert.Equal(t, underlying.Name(), rule.Name())
			config := ""
			if test.configuredValue != "" {
				config = fmt.Sprintf(`
rule "avm_test_rule" {
  enabled  = true
  severity = %q
}`, test.configuredValue)
			}
			runner := helper.TestRunner(t, map[string]string{".tflint.hcl": config})

			require.NoError(t, rule.Check(runner))
			require.Len(t, runner.Issues, 1)
			assert.Equal(t, "avm_test_rule", runner.Issues[0].Rule.Name())
			assert.Equal(t, test.wantSeverity, runner.Issues[0].Rule.Severity())
		})
	}
}

func TestConfigurableRuleRejectsInvalidSeverity(t *testing.T) {
	rule := common.NewConfigurableRule(&emittingRule{severity: tflint.ERROR})
	runner := helper.TestRunner(t, map[string]string{
		".tflint.hcl": `
rule "avm_test_rule" {
  enabled  = true
  severity = "critical"
}`,
	})

	err := rule.Check(runner)

	require.ErrorContains(t, err, `severity must be one of "error", "warning", or "notice", got "critical"`)
	assert.Empty(t, runner.Issues)
}

func TestConfigurableRuleRejectsModuleClassForOtherRules(t *testing.T) {
	rule := common.NewConfigurableRule(&emittingRule{severity: tflint.ERROR})
	runner := helper.TestRunner(t, map[string]string{
		".tflint.hcl": `
rule "avm_test_rule" {
  enabled      = true
  module_class = "pattern"
}`,
	})

	err := rule.Check(runner)

	require.ErrorContains(t, err, "Unsupported argument")
	assert.Empty(t, runner.Issues)
}

func TestConfigurableRuleRejectsModuleScopeForOtherRules(t *testing.T) {
	rule := common.NewConfigurableRule(&emittingRule{severity: tflint.ERROR})
	runner := helper.TestRunner(t, map[string]string{
		".tflint.hcl": `
rule "avm_test_rule" {
  enabled      = true
  module_scope = "submodule"
}`,
	})

	err := rule.Check(runner)

	require.ErrorContains(t, err, "Unsupported argument")
	assert.Empty(t, runner.Issues)
}

func TestConfigurableRuleFiltersIgnoredResourceAndDataAddresses(t *testing.T) {
	tests := []struct {
		name       string
		selectors  string
		wantIssues int
	}{
		{
			name:       "absent",
			wantIssues: 2,
		},
		{
			name:       "short managed resource address",
			selectors:  `ignored_resources = ["azapi_resource.telemetry"]`,
			wantIssues: 1,
		},
		{
			name:       "explicit managed resource address",
			selectors:  `ignored_resources = ["resource.azapi_resource.telemetry"]`,
			wantIssues: 1,
		},
		{
			name:       "data source address",
			selectors:  `ignored_resources = ["data.azapi_resource.telemetry"]`,
			wantIssues: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rule := common.NewConfigurableRule(&blockEmittingRule{})
			runner := helper.TestRunner(t, map[string]string{
				".tflint.hcl": fmt.Sprintf(`rule "avm_test_block_rule" {
  enabled = true
  %s
}`, test.selectors),
				"main.tf": `resource "azapi_resource" "telemetry" {
  count = 1
}

data "azapi_resource" "telemetry" {
  for_each = {}
}`,
			})

			require.NoError(t, rule.Check(runner))
			assert.Len(t, runner.Issues, test.wantIssues)
		})
	}
}

func TestConfigurableRuleRejectsInvalidIgnoredResourceAddress(t *testing.T) {
	for _, selector := range []string{
		"azapi_resource",
		"var.example",
		"local.example",
		"module.example",
		"data.azapi_resource.telemetry[0]",
		"azapi_resource.telemetry[0]",
		"azapi_resource.*",
	} {
		t.Run(selector, func(t *testing.T) {
			rule := common.NewConfigurableRule(&blockEmittingRule{})
			runner := helper.TestRunner(t, map[string]string{
				".tflint.hcl": fmt.Sprintf(`rule "avm_test_block_rule" {
  enabled           = true
  ignored_resources = [%q]
}`, selector),
			})

			err := rule.Check(runner)

			require.ErrorContains(t, err, "ignored_resources entry")
			assert.Empty(t, runner.Issues)
		})
	}
}

func TestConfigurableRuleSupportsIgnoredResourcesForEveryConfigurationShape(t *testing.T) {
	tests := []struct {
		name   string
		rule   tflint.Rule
		config string
	}{
		{
			name: "standard rule",
			rule: &blockEmittingRule{},
			config: `rule "avm_test_block_rule" {
  enabled           = true
  ignored_resources = ["azapi_resource.telemetry"]
}`,
		},
		{
			name: "module class rule",
			rule: &moduleClassBlockEmittingRule{},
			config: `rule "avm_test_module_class_block_rule" {
  enabled           = true
  module_class      = "resource"
  ignored_resources = ["azapi_resource.telemetry"]
}`,
		},
		{
			name: "module scope rule",
			rule: &moduleScopeBlockEmittingRule{},
			config: `rule "avm_test_module_scope_block_rule" {
  enabled           = true
  module_scope      = "root"
  ignored_resources = ["azapi_resource.telemetry"]
}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rule := common.NewConfigurableRule(test.rule)
			runner := helper.TestRunner(t, map[string]string{
				".tflint.hcl": test.config,
				"main.tf":     `resource "azapi_resource" "telemetry" {}`,
			})

			require.NoError(t, rule.Check(runner))
			assert.Empty(t, runner.Issues)
		})
	}
}

func TestConfigurableRuleDoesNotHideUnrelatedDiagnostics(t *testing.T) {
	rule := common.NewConfigurableRule(&emittingRule{severity: tflint.ERROR})
	runner := helper.TestRunner(t, map[string]string{
		".tflint.hcl": `rule "avm_test_rule" {
  enabled           = true
  ignored_resources = ["azapi_resource.telemetry"]
}`,
	})

	require.NoError(t, rule.Check(runner))
	assert.Len(t, runner.Issues, 1)
}

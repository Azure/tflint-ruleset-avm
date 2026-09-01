package common_test

import (
	"fmt"
	"testing"

	"github.com/Azure/tflint-ruleset-avm/common"
	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type emittingRule struct {
	tflint.DefaultRule
	severity tflint.Severity
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

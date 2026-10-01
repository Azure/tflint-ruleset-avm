package outputs_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Azure/tflint-ruleset-avm/common"
	"github.com/Azure/tflint-ruleset-avm/outputs"
	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/terraform/addrs"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func TestRequiredOutput(t *testing.T) {
	cases := []struct {
		desc         string
		config       string
		requiredName string
		issues       helper.Issues
	}{
		{
			desc: "require resource_id, ok",
			config: `output "resource_id" {
  value = azurerm_kubernetes_cluster.this.id
}`,
			requiredName: "resource_id",
			issues:       helper.Issues{},
		},
		{
			desc:         "require resource_id, not ok",
			config:       ``,
			requiredName: "resource_id",
			issues: helper.Issues{
				{
					Rule:    outputs.NewRequiredOutputRule("required_output", "resource_id", ""),
					Message: "module owners MUST output the `resource_id` in their modules",
					Range: hcl.Range{
						Filename: "outputs.tf",
					},
				},
			},
		},
		{
			desc: "require resource, ok",
			config: `output "resource" {
  value = azurerm_kubernetes_cluster.this
}`,
			requiredName: "resource",
			issues:       helper.Issues{},
		},
		{
			desc:         "require resource, not ok",
			config:       ``,
			requiredName: "resource",
			issues: helper.Issues{
				{
					Rule:    outputs.NewRequiredOutputRule("required_output", "resource", ""),
					Message: "module owners MUST output the `resource` in their modules",
					Range: hcl.Range{
						Filename: "outputs.tf",
					},
				},
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			rule := outputs.NewRequiredOutputRule("required_output", tc.requiredName, "")
			filename := "variables.tf"

			runner := helper.TestRunner(t, map[string]string{filename: tc.config})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssuesWithoutRange(t, tc.issues, runner.Issues)
		})
	}
}

type scopedOutputRunner struct {
	tflint.Runner
	modulePath addrs.Module
	originalwd string
}

func (r scopedOutputRunner) GetModulePath() (addrs.Module, error) {
	return r.modulePath, nil
}

func (r scopedOutputRunner) GetOriginalwd() (string, error) {
	return r.originalwd, nil
}

func TestResourceIDRequiredRuleModuleClass(t *testing.T) {
	tests := []struct {
		name         string
		ruleConfig   string
		output       string
		modulePath   addrs.Module
		originalwd   string
		metadata     string
		wantIssue    bool
		wantSeverity tflint.Severity
	}{
		{
			name: "resource root with resource_id",
			output: `output "resource_id" {
  value = "example"
}`,
		},
		{
			name:         "resource root without resource_id",
			wantIssue:    true,
			wantSeverity: tflint.ERROR,
		},
		{
			name: "explicit resource root without resource_id",
			ruleConfig: `rule "avm_output_resource_id_required" {
  enabled      = true
  module_class = "resource"
}`,
			wantIssue:    true,
			wantSeverity: tflint.ERROR,
		},
		{
			name: "pattern root without resource_id on Windows path",
			ruleConfig: `rule "avm_output_resource_id_required" {
  enabled      = true
  module_class = "pattern"
}`,
			originalwd: `C:\modules\pattern`,
		},
		{
			name: "utility root without resource_id on Unix path",
			ruleConfig: `rule "avm_output_resource_id_required" {
  enabled      = true
  module_class = "utility"
}`,
			originalwd: "/modules/utility",
		},
		{
			name:       "resource child without resource_id",
			modulePath: addrs.Module{"child"},
		},
		{
			name: "pattern child without resource_id",
			ruleConfig: `rule "avm_output_resource_id_required" {
  enabled      = true
  module_class = "pattern"
}`,
			modulePath: addrs.Module{"child"},
		},
		{
			name:         "missing metadata does not disable resource rule",
			wantIssue:    true,
			wantSeverity: tflint.ERROR,
		},
		{
			name:         "malformed metadata does not disable resource rule",
			metadata:     `{"telemetryIdPrefix":`,
			wantIssue:    true,
			wantSeverity: tflint.ERROR,
		},
		{
			name: "severity override still applies to resource root",
			ruleConfig: `rule "avm_output_resource_id_required" {
  enabled  = true
  severity = "notice"
}`,
			wantIssue:    true,
			wantSeverity: tflint.NOTICE,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			moduleDir := t.TempDir()
			if tc.metadata != "" {
				require.NoError(t, os.WriteFile(filepath.Join(moduleDir, "metadata.json"), []byte(tc.metadata), 0600))
			}
			if tc.originalwd != "" {
				moduleDir = tc.originalwd
			}

			files := map[string]string{"outputs.tf": tc.output}
			if tc.ruleConfig != "" {
				files[".tflint.hcl"] = tc.ruleConfig
			}
			runner := helper.TestRunner(t, files)
			scopedRunner := scopedOutputRunner{
				Runner:     runner,
				modulePath: tc.modulePath,
				originalwd: moduleDir,
			}
			rule := common.NewConfigurableRule(outputs.NewResourceIDRequiredRule(""))

			require.NoError(t, rule.Check(scopedRunner))
			if !tc.wantIssue {
				assert.Empty(t, runner.Issues)
				return
			}
			require.Len(t, runner.Issues, 1)
			issue := runner.Issues[0]
			assert.Equal(t, rule.Name(), issue.Rule.Name())
			assert.Equal(t, tc.wantSeverity, issue.Rule.Severity())
			assert.Equal(t, "module owners MUST output the `resource_id` in their modules", issue.Message)
			assert.Equal(t, "outputs.tf", issue.Range.Filename)
		})
	}
}

func TestResourceIDRequiredRuleRejectsInvalidModuleClass(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		modulePath addrs.Module
	}{
		{name: "empty", expression: `""`},
		{name: "uppercase", expression: `"Resource"`},
		{name: "trailing whitespace", expression: `"pattern "`},
		{name: "unknown", expression: `"other"`},
		{name: "child with invalid class", expression: `"other"`, modulePath: addrs.Module{"child"}},
		{name: "null", expression: "null"},
		{name: "boolean", expression: "true"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			runner := helper.TestRunner(t, map[string]string{
				"outputs.tf": "",
				".tflint.hcl": fmt.Sprintf(`rule "avm_output_resource_id_required" {
  enabled      = true
  module_class = %s
}`, tc.expression),
			})
			rule := common.NewConfigurableRule(outputs.NewResourceIDRequiredRule(""))

			err := rule.Check(scopedOutputRunner{Runner: runner, modulePath: tc.modulePath})

			require.Error(t, err)
			if tc.expression != "null" && tc.expression != "true" {
				assert.ErrorContains(t, err, `module_class must be one of "resource", "pattern", or "utility"`)
			}
			assert.Empty(t, runner.Issues)
		})
	}
}

func TestResourceIDRequiredRuleRejectsInvalidSeverity(t *testing.T) {
	runner := helper.TestRunner(t, map[string]string{
		"outputs.tf": "",
		".tflint.hcl": `rule "avm_output_resource_id_required" {
  enabled      = true
  module_class = "pattern"
  severity     = "critical"
}`,
	})
	rule := common.NewConfigurableRule(outputs.NewResourceIDRequiredRule(""))

	err := rule.Check(runner)

	require.ErrorContains(t, err, `severity must be one of "error", "warning", or "notice"`)
	assert.Empty(t, runner.Issues)
}

func TestResourceIDRequiredRuleCheckWithoutWrapper(t *testing.T) {
	runner := helper.TestRunner(t, map[string]string{
		"outputs.tf": "",
		".tflint.hcl": `rule "avm_output_resource_id_required" {
  enabled      = true
  module_class = "utility"
}`,
	})

	require.NoError(t, outputs.NewResourceIDRequiredRule("").Check(runner))
	assert.Empty(t, runner.Issues)
}

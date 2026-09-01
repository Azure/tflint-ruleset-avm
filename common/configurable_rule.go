package common

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// ConfigurableRule adds validated severity configuration to a rule.
type ConfigurableRule struct {
	tflint.Rule
	name string
}

type severityConfig struct {
	Severity string `hclext:"severity,optional"`
}

type severityRunner struct {
	tflint.Runner
	rule     tflint.Rule
	severity tflint.Severity
}

type severityRule struct {
	tflint.Rule
	severity tflint.Severity
}

// NewConfigurableRule returns a rule with a canonical name and configurable severity.
func NewConfigurableRule(name string, rule tflint.Rule) *ConfigurableRule {
	return &ConfigurableRule{
		Rule: rule,
		name: name,
	}
}

// Name returns the canonical public rule name.
func (r *ConfigurableRule) Name() string {
	return r.name
}

// Check decodes configuration before delegating to the underlying rule.
func (r *ConfigurableRule) Check(runner tflint.Runner) error {
	config := severityConfig{}
	if err := runner.DecodeRuleConfig(r.Name(), &config); err != nil {
		return err
	}

	severity, err := parseSeverity(config.Severity, r.Severity())
	if err != nil {
		return fmt.Errorf("invalid configuration for rule %q: %w", r.Name(), err)
	}

	return r.Rule.Check(&severityRunner{
		Runner:   runner,
		rule:     r,
		severity: severity,
	})
}

func parseSeverity(value string, defaultSeverity tflint.Severity) (tflint.Severity, error) {
	switch value {
	case "":
		return defaultSeverity, nil
	case "error":
		return tflint.ERROR, nil
	case "warning":
		return tflint.WARNING, nil
	case "notice":
		return tflint.NOTICE, nil
	default:
		return 0, fmt.Errorf("severity must be one of %q, %q, or %q, got %q", "error", "warning", "notice", value)
	}
}

func (r *severityRunner) EmitIssue(_ tflint.Rule, message string, issueRange hcl.Range) error {
	return r.Runner.EmitIssue(&severityRule{Rule: r.rule, severity: r.severity}, message, issueRange)
}

func (r *severityRunner) EmitIssueWithFix(
	_ tflint.Rule,
	message string,
	issueRange hcl.Range,
	fixFunc func(tflint.Fixer) error,
) error {
	return r.Runner.EmitIssueWithFix(
		&severityRule{Rule: r.rule, severity: r.severity},
		message,
		issueRange,
		fixFunc,
	)
}

func (r *severityRule) Severity() tflint.Severity {
	return r.severity
}

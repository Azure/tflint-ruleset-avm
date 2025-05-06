package rules

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var _ tflint.Runner = &dummyRunner{}

// dummyRunner is used in `valid_template_interpolation` rule to run other WAF rules, but we just want to probe potential interpolation errors, not emit issues, so we wrap the actual runner with this dummy runner, swallow all emitted issues.
type dummyRunner struct {
	tflint.Runner
}

func (r *dummyRunner) EmitIssue(rule tflint.Rule, message string, issueRange hcl.Range) error {
	return nil
}

func (r *dummyRunner) EmitIssueWithFix(rule tflint.Rule, message string, issueRange hcl.Range, fixFunc func(f tflint.Fixer) error) error {
	return nil
}

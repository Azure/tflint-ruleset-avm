package rules

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var _ tflint.Runner = &issueCollectDummyRunner{}

type issueCollectDummyRunner struct {
	tflint.Runner
}

func newIssueCollectDummyRunner(runner tflint.Runner) *issueCollectDummyRunner {
	return &issueCollectDummyRunner{
		Runner: runner,
	}
}

func (r *issueCollectDummyRunner) EmitIssue(rule tflint.Rule, message string, issueRange hcl.Range) error {
	return nil
}

func (r *issueCollectDummyRunner) EmitIssueWithFix(rule tflint.Rule, message string, issueRange hcl.Range, fixFunc func(f tflint.Fixer) error) error {
	return nil
}

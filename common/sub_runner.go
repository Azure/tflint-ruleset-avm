package common

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type subRunner struct {
	tflint.Runner
	issues []issue
}

type issue struct {
	message    string
	issueRange hcl.Range
}

// RulePasses reports whether a rule emits no issues for the supplied runner.
func RulePasses(runner tflint.Runner, rule tflint.Rule) (bool, error) {
	sr := &subRunner{Runner: runner}
	if err := rule.Check(sr); err != nil {
		return false, err
	}
	return len(sr.issues) == 0, nil
}

func (e *subRunner) EmitIssue(rule tflint.Rule, message string, issueRange hcl.Range) error {
	e.issues = append(e.issues, issue{
		message:    message,
		issueRange: issueRange,
	})
	return nil
}

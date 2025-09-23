package common

import (
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var _ tflint.Rule = new(EitherCheckRule)

type EitherCheckRule struct {
	tflint.DefaultRule
	rules    []tflint.Rule
	name     string
	enabled  bool
	severity tflint.Severity
}

func NewEitherCheckRule(name string, enabled bool, severity tflint.Severity, rules ...tflint.Rule) *EitherCheckRule {
	return &EitherCheckRule{
		name:     name,
		enabled:  enabled,
		severity: severity,
		rules:    rules,
	}
}

func (e *EitherCheckRule) Name() string {
	return e.name
}

func (e *EitherCheckRule) Enabled() bool {
	return e.enabled
}

func (e *EitherCheckRule) Severity() tflint.Severity {
	return e.severity
}

func (e *EitherCheckRule) Check(runner tflint.Runner) error {
	runners := map[tflint.Rule]*subRunner{}

	issues := []issue{}
	var failingRule tflint.Rule
	for _, r := range e.rules {
		sr := &subRunner{
			Runner: runner,
		}
		runners[r] = sr
		if err := r.Check(sr); err != nil {
			return err
		}
		if len(sr.issues) == 0 {
			return nil
		}

		if len(issues) == 0 && failingRule == nil {
			issues = sr.issues
			failingRule = r
		}
	}

	sr := runners[e.rules[0]]
	for _, issue := range sr.issues {
		if err := runner.EmitIssue(failingRule, issue.message, issue.issueRange); err != nil {
			return err
		}
	}

	return nil
}

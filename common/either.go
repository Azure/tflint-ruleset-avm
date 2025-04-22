package common

import (
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var _ tflint.Rule = new(EitherCheckRule)

type EitherCheckRule struct {
	tflint.DefaultRule
	primaryRule   tflint.Rule
	secondaryRule tflint.Rule
	name          string
	enabled       bool
	severity      tflint.Severity
}

func NewEitherCheckRule(name string, enabled bool, severity tflint.Severity, primaryRule tflint.Rule, secondary tflint.Rule) *EitherCheckRule {
	return &EitherCheckRule{
		name:          name,
		enabled:       enabled,
		severity:      severity,
		primaryRule:   primaryRule,
		secondaryRule: secondary,
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

	for _, r := range []tflint.Rule{e.primaryRule, e.secondaryRule} {
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
	}
	sr := runners[e.primaryRule]
	for _, issue := range sr.issues {
		if err := runner.EmitIssue(e.primaryRule, issue.message, issue.issueRange); err != nil {
			return err
		}
	}
	return nil
}

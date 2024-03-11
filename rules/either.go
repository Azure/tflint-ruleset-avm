package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// Check interface compliance with the tflint.Rule.
var _ tflint.Rule = new(EitherCheckRule)

type EitherCheckRule struct {
	tflint.DefaultRule
	rules       []tflint.Rule
	primaryRule tflint.Rule
	name        string
	enabled     bool
	severity    tflint.Severity
}

func NewEitherTypeCheckRule(name string, enabled bool, severity tflint.Severity, primaryRule tflint.Rule, rules ...tflint.Rule) *EitherCheckRule {
	found := false
	for _, r := range rules {
		if primaryRule == r {
			found = true
			break
		}
	}
	if !found {
		panic("primaryRule must be one of rules.")
	}
	return &EitherCheckRule{
		name:        name,
		enabled:     enabled,
		severity:    severity,
		primaryRule: primaryRule,
		rules:       rules,
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

	for _, r := range e.rules {
		sr := &subRunner{
			Runner: runner,
		}
		runners[r] = sr
		err := r.Check(sr)
		if err != nil {
			return err
		}
		if len(sr.issues) == 0 {
			return nil
		}
	}
	sr := runners[e.primaryRule]
	for _, issue := range sr.issues {
		err := runner.EmitIssue(e, issue.message, issue.issueRange)
		if err != nil {
			return err
		}
	}
	return nil
}

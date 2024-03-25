package rules

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

func (e *subRunner) EmitIssue(rule tflint.Rule, message string, issueRange hcl.Range) error {
	e.issues = append(e.issues, issue{
		message:    message,
		issueRange: issueRange,
	})
	return nil
}

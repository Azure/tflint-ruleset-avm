package rules

import (
	"fmt"
	"regexp"

	"github.com/Azure/tflint-ruleset-avm/attrvalue"
	"github.com/Azure/tflint-ruleset-avm/waf"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var TemplateInterpolationErrorRegex = regexp.MustCompile("Invalid template interpolation value;")
var _ tflint.Rule = new(ValidTemplateInterpolationRule)

type ValidTemplateInterpolationRule struct {
	tflint.DefaultRule
}

func (v *ValidTemplateInterpolationRule) Name() string {
	return "valid_template_interpolation"
}

func (v *ValidTemplateInterpolationRule) Enabled() bool {
	return true
}

func (v *ValidTemplateInterpolationRule) Severity() tflint.Severity {
	return tflint.WARNING
}

func (v *ValidTemplateInterpolationRule) Check(r tflint.Runner) error {
	dummyRunner := newIssueCollectDummyRunner(r)
	for _, wafRule := range waf.GetRules() {
		p, ok := wafRule.(attrvalue.PartialContentEvaluationRule)
		if !ok {
			continue
		}
		p.EnableFailOnEvaluationError()
		if err := wafRule.Check(dummyRunner); err != nil && TemplateInterpolationErrorRegex.MatchString(err.Error()) {
			return fmt.Errorf("cannot evaluate template interpolation, usually due to null reference error or other issues: %+v", err)
		}
	}
	return nil
}

func NewValidTemplateInterpolationRule() *ValidTemplateInterpolationRule {
	return new(ValidTemplateInterpolationRule)
}

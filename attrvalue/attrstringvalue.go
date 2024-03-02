package attrvalue

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type StringRule struct {
	tflint.DefaultRule

	resourceType   string
	attributeName  string
	expectedValues []string
}

var _ tflint.Rule = (*StringRule)(nil)

func NewAttrStringValueRule(resourceType string, attributeName string, expectedValues []string) *StringRule {
	return &StringRule{
		resourceType:   resourceType,
		attributeName:  attributeName,
		expectedValues: expectedValues,
	}
}

func (r *StringRule) Name() string {
	return fmt.Sprintf("%s.%s must be: %s", r.resourceType, r.attributeName, r.expectedValues)
}

func (r *StringRule) Enabled() bool {
	return true
}

func (r *StringRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *StringRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.attributeName},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		attribute, exists := resource.Body.Attributes[r.attributeName]
		if !exists {
			continue
		}
		err := runner.EvaluateExpr(attribute.Expr, func(val string) error {
			found := false
			for _, item := range r.expectedValues {
				if item == val {
					found = true
				}
			}
			if !found {
				runner.EmitIssue(
					r,
					fmt.Sprintf("\"%s\" is an invalid attribute value of `%s` - expecting (one of) %s", val, r.attributeName, r.expectedValues),
					attribute.Expr.Range(),
				)
			}
			return nil
		}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

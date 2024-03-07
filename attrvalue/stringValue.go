package attrvalue

import (
	"fmt"
	"slices"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// StringRule checks whether a string attribute value is one of the expected values.
type StringRule struct {
	tflint.DefaultRule // Embed the default rule to reuse its implementation

	resourceType   string   // e.g. "azurerm_storage_account"
	attributeName  string   // e.g. "account_replication_type"
	expectedValues []string // e.g. []string{"ZRS"}
}

var _ tflint.Rule = (*StringRule)(nil)

// NewStringRule returns a new rule with the given resource type, attribute name, and expected values.
func NewStringRule(resourceType string, attributeName string, expectedValues []string) *StringRule {
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
			if !slices.Contains(r.expectedValues, val) {
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

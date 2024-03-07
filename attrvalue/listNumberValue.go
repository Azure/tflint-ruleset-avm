package attrvalue

import (
	"fmt"
	"slices"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// ListNumberRule checks whether a list of numbers attribute value is one of the expected values.
// It is not concerned with the order of the numbers in the list.
type ListNumberRule struct {
	tflint.DefaultRule // Embed the default rule to reuse its implementation

	resourceType   string    // e.g. "azurerm_storage_account"
	attributeName  string    // e.g. "account_replication_type"
	expectedValues [][]int32 // e.g. [][int32{1, 2, 3}]
}

var _ tflint.Rule = (*ListNumberRule)(nil)

// NewListNumberRule returns a new rule with the given resource type, attribute name, and expected values.
func NewListNumberRule(resourceType string, attributeName string, expectedValues [][]int32) *ListNumberRule {
	return &ListNumberRule{
		resourceType:   resourceType,
		attributeName:  attributeName,
		expectedValues: expectedValues,
	}
}

func (r *ListNumberRule) Name() string {
	return fmt.Sprintf("%s.%s must be: %v", r.resourceType, r.attributeName, r.expectedValues)
}

func (r *ListNumberRule) Enabled() bool {
	return true
}

func (r *ListNumberRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *ListNumberRule) Check(runner tflint.Runner) error {
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
		wantTy := cty.List(cty.Number)
		err := runner.EvaluateExpr(attribute.Expr, func(val *[]int32) error {
			slices.Sort(*val)
			for _, exp := range r.expectedValues {
				slices.Sort(exp)
				if slices.Equal(*val, exp) {
					return nil
				}
			}
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%v\" is an invalid attribute value of `%s` - expecting (one of) %v", val, r.attributeName, r.expectedValues),
				attribute.Expr.Range(),
			)
			return nil
		}, &tflint.EvaluateExprOption{
			WantType: &wantTy,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

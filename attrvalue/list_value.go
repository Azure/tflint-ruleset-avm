package attrvalue

import (
	"cmp"
	"fmt"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
	"slices"
)

// ListRule checks whether a list of numbers attribute value is one of the expected values.
// It is not concerned with the order of the numbers in the list.
type ListRule[T cmp.Ordered] struct {
	tflint.DefaultRule // Embed the default rule to reuse its implementation

	resourceType   string // e.g. "azurerm_storage_account"
	attributeName  string // e.g. "account_replication_type"
	expectedValues [][]T  // e.g. [][int{1, 2, 3}]
}

var _ tflint.Rule = (*ListRule[int])(nil)

// NewListRule returns a new rule with the given resource type, attribute name, and expected values.
func NewListRule[T cmp.Ordered](resourceType string, attributeName string, expectedValues [][]T) *ListRule[T] {
	return &ListRule[T]{
		resourceType:   resourceType,
		attributeName:  attributeName,
		expectedValues: expectedValues,
	}
}

func (r *ListRule[T]) Name() string {
	return fmt.Sprintf("%s.%s must be: %v", r.resourceType, r.attributeName, r.expectedValues)
}

func (r *ListRule[T]) Enabled() bool {
	return true
}

func (r *ListRule[T]) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *ListRule[T]) Check(runner tflint.Runner) error {
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
		if err := runner.EvaluateExpr(attribute.Expr, func(val *[]T) error {
			slices.Sort(*val)
			for _, exp := range r.expectedValues {
				slices.Sort(exp)
				if slices.Equal(*val, exp) {
					return nil
				}
			}
			return runner.EmitIssue(
				r,
				fmt.Sprintf("\"%v\" is an invalid attribute value of `%s` - expecting (one of) %v", val, r.attributeName, r.expectedValues),
				attribute.Expr.Range(),
			)
		}, &tflint.EvaluateExprOption{
			WantType: &wantTy,
		}); err != nil {
			return err
		}
	}
	return nil
}

package attrvalue

import (
	"fmt"
	"reflect"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// SimpleRule checks whether a string attribute value is one of the expected values.
// Is can be used to check string, number, and bool attributes.
// Supply the cty.Type and reflect.Type of the attribute value, and a slice of expected values as `[]any`.
// To provide the reflect.Type of a value, use `reflect.TypeOf()` and provide a parameter,
// e.g. `""` for string, `0` for number, or `true` for bool.
// To provide the cty.Type, use `cty.String`, `cty.Number`, or `cty.Bool` for string, number, and bool, respectively.
type SimpleRule[T any] struct {
	tflint.DefaultRule // Embed the default rule to reuse its implementation

	resourceType   string // e.g. "azurerm_storage_account"
	attributeName  string // e.g. "account_replication_type"
	expectedValues []T    // e.g. []string{"ZRS"}
}

var _ tflint.Rule = (*SimpleRule[any])(nil)

// NewSimpleRule returns a new rule with the given resource type, attribute name, and expected values.
func NewSimpleRule[T any](resourceType string, attributeName string, expectedValues []T) *SimpleRule[T] {
	return &SimpleRule[T]{
		resourceType:   resourceType,
		attributeName:  attributeName,
		expectedValues: expectedValues,
	}
}

func (r *SimpleRule[T]) Name() string {
	return fmt.Sprintf("%s.%s must be: %+v", r.resourceType, r.attributeName, r.expectedValues)
}

func (r *SimpleRule[T]) Enabled() bool {
	return true
}

func (r *SimpleRule[T]) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *SimpleRule[T]) Check(runner tflint.Runner) error {
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
		var dt T
		val := toStrongTypePtr(dt)
		ctyType, err := toCtyType(dt)
		if err != nil {
			return err
		}
		if err = runner.EvaluateExpr(attribute.Expr, val, &tflint.EvaluateExprOption{
			WantType: &ctyType,
		}); err != nil {
			return err
		}
		for _, expected := range r.expectedValues {
			exp := toStrongTypePtr(expected)
			if reflect.DeepEqual(val, exp) {
				return nil
			}
		}
		v := reflect.Indirect(reflect.ValueOf(val))
		return runner.EmitIssue(
			r,
			fmt.Sprintf("%v is an invalid attribute value of `%s` - expecting (one of) %v", v, r.attributeName, r.expectedValues),
			attribute.Range,
		)
	}
	return nil
}

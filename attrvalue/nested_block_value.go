package attrvalue

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// SimpleRule checks whether a string attribute value is one of the expected values.
// It can be used to check string, number, and bool attributes.
type NestedBlockRule[T any] struct {
	tflint.DefaultRule // Embed the default rule to reuse its implementation

	resourceType    string // e.g. "azurerm_application_gateway"
	nestedBlockType string // e.g. "sku"
	attributeName   string // e.g. "name
	expectedValues  []T    // e.g. []string{"Standard_V2"}
}

var _ tflint.Rule = (*SimpleRule[any])(nil)

// NewSimpleRule returns a new rule with the given resource type, attribute name, and expected values.
func NewNestedBlockRule[T any](resourceType string, nestedBlockType string, attributeName string, expectedValues []T) *NestedBlockRule[T] {
	return &NestedBlockRule[T]{
		resourceType:    resourceType,
		nestedBlockType: nestedBlockType,
		attributeName:   attributeName,
		expectedValues:  expectedValues,
	}
}

func (r *NestedBlockRule[T]) Name() string {
	return fmt.Sprintf("%s.%s must be: %+v", r.resourceType, r.attributeName, r.expectedValues)
}

func (r *NestedBlockRule[T]) Enabled() bool {
	return true
}

func (r *NestedBlockRule[T]) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *NestedBlockRule[T]) Check(runner tflint.Runner) error {
	attrs, err := getNestedBlockAttrs(runner, r.resourceType, r.nestedBlockType, r.attributeName)
	if err != nil {
		return err
	}

	for _, attr := range attrs {
		var dt T
		val := toPtr(dt)
		ctyType, err := toCtyType(dt)
		if err != nil {
			return err
		}
		if err = runner.EvaluateExpr(attr.Expr, val, &tflint.EvaluateExprOption{
			WantType: &ctyType,
		}); err != nil {
			if slices.Contains([]string{
				"null value found",           // generated by tflint.
				"null value is not allowed"}, // generated by the test runner.
				err.Error()) {
				continue
			}
		}
		found := false
		for _, expected := range r.expectedValues {
			exp := toPtr(expected)
			if reflect.DeepEqual(val, exp) {
				found = true
			}
		}
		if !found {
			v := reflect.Indirect(reflect.ValueOf(val))
			return runner.EmitIssue(
				r,
				fmt.Sprintf("%v is an invalid attribute value of `%s` - expecting (one of) %v", v, r.attributeName, r.expectedValues),
				attr.Range,
			)
		}
	}
	return nil
}

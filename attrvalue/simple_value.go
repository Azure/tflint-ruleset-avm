package attrvalue

import (
	"fmt"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/zclconf/go-cty/cty"

	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty/gocty"
)

// SimpleRule checks whether a string attribute value is one of the expected values.
// It can be used to check string, number, and bool attributes.
type SimpleRule[T any] struct {
	tflint.DefaultRule // Embed the default rule to reuse its implementation
	baseValue
	expectedValues []T // e.g. []string{"ZRS"}
}

var _ tflint.Rule = (*SimpleRule[any])(nil)
var _ AttrValueRule = (*SimpleRule[any])(nil)

// NewSimpleRule returns a new rule with the given resource type, attribute name, and expected values.
func NewSimpleRule[T any](resourceType, attributeName string, expectedValues []T) *SimpleRule[T] {
	return &SimpleRule[T]{
		baseValue:      newBaseValue(resourceType, nil, attributeName),
		expectedValues: expectedValues,
	}
}

// NewSimpleRule returns a new rule with the given resource type, attribute name, and expected values.
func NewSimpleNestedBlockRule[T any](resourceType, nestedBlockType, attributeName string, expectedValues []T) *SimpleRule[T] {
	return &SimpleRule[T]{
		baseValue:      newBaseValue(resourceType, &nestedBlockType, attributeName),
		expectedValues: expectedValues,
	}
}

func (r *SimpleRule[T]) Name() string {
	return fmt.Sprintf("%s.%s must be: %+v", r.resourceType, r.attributeName, r.expectedValues)
}

func (r *SimpleRule[T]) Check(runner tflint.Runner) error {
	var dt T
	ctyType, err := toCtyType(dt)
	if err != nil {
		return err
	}
	return r.checkAttributes(runner, cty.DynamicPseudoType, func(attr *hclext.Attribute, val cty.Value) error {
		if val.IsNull() {
			return nil
		}
		found := false
		for _, exp := range r.expectedValues {
			ctyExp, err := gocty.ToCtyValue(exp, ctyType)
			if err != nil {
				return err
			}
			found = ctyExp.Equals(val).True()
			if found {
				break
			}
		}
		if !found {
			goVal := new(T)
			_ = gocty.FromCtyValue(val, goVal)
			return runner.EmitIssue(
				r,
				fmt.Sprintf("%v is an invalid attribute value of `%s` - expecting (one of) %v", *goVal, r.attributeName, r.expectedValues),
				attr.Range,
			)
		}
		return nil
	})
}

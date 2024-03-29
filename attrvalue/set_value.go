package attrvalue

import (
	"cmp"
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

// SetRule checks whether a list of numbers attribute value is one of the expected values.
// It is not concerned with the order of the numbers in the list.
type SetRule[T cmp.Ordered] struct {
	tflint.DefaultRule // Embed the default rule to reuse its implementation
	baseValue
	expectedValues [][]T // e.g. [][int{1, 2, 3}]
}

var _ tflint.Rule = (*SetRule[int])(nil)
var _ AttrValueRule = (*SimpleRule[any])(nil)

// NewSetRule returns a new rule with the given resource type, attribute name, and expected values.
func NewSetRule[T cmp.Ordered](resourceType string, attributeName string, expectedValues [][]T, link string) *SetRule[T] {
	return &SetRule[T]{
		baseValue:      newBaseValue(resourceType, nil, attributeName, true, link, tflint.ERROR),
		expectedValues: expectedValues,
	}
}

func (r *SetRule[T]) Name() string {
	if r.nestedBlockType != nil {
		return fmt.Sprintf("%s.%s.%s_must_be_%+v", r.resourceType, *r.nestedBlockType, r.attributeName, r.expectedValues)
	}
	return fmt.Sprintf("%s.%s_must_be_%v", r.resourceType, r.attributeName, r.expectedValues)
}

func (r *SetRule[T]) Check(runner tflint.Runner) error {
	var dts []T
	var dt T
	ctyTypeS, err := toCtyType(dts)
	if err != nil {
		return err
	}
	ctyType, err := toCtyType(dt)
	if err != nil {
		return err
	}
	return r.checkAttributes(runner, ctyTypeS, func(attr *hclext.Attribute, val cty.Value) error {
		if val.IsNull() || !val.IsKnown() {
			return nil
		}
		actual := val.AsValueSet()
		for _, exp := range r.expectedValues {
			expectedValue, err := gocty.ToCtyValue(exp, cty.Set(ctyType))
			if err != nil {
				return err
			}
			if cty.SetValFromValueSet(actual).Equals(expectedValue).True() {
				return nil
			}
		}
		goVal := new([]T)
		_ = gocty.FromCtyValue(val, goVal)
		return runner.EmitIssue(
			r,
			fmt.Sprintf("\"%v\" is an invalid attribute value of `%s` - expecting (one of) %v", *goVal, r.attributeName, r.expectedValues),
			attr.Expr.Range(),
		)
	})
}

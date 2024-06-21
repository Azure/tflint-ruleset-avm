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
	mustExist      bool
	ruleName       string
}

var _ tflint.Rule = (*SimpleRule[any])(nil)
var _ AttrValueRule = (*SimpleRule[any])(nil)

// NewSimpleRule returns a new rule with the given resource type, attribute name, and expected values.
func NewSimpleRule[T any](resourceType, attributeName string, expectedValues []T, link string, mustExist bool, ruleName string) *SimpleRule[T] {
	return &SimpleRule[T]{
		baseValue:      newBaseValue(resourceType, nil, attributeName, true, link, tflint.ERROR),
		expectedValues: expectedValues,
		mustExist:      mustExist,
		ruleName:       ruleName,
	}
}

// NewSimpleNestedBlockRule returns a new rule with the given resource type, attribute name, and expected values.
func NewSimpleNestedBlockRule[T any](resourceType, nestedBlockType, attributeName string, expectedValues []T, link string, mustExist bool, ruleName string) *SimpleRule[T] {
	return &SimpleRule[T]{
		baseValue:      newBaseValue(resourceType, &nestedBlockType, attributeName, true, link, tflint.ERROR),
		expectedValues: expectedValues,
		mustExist:      mustExist,
		ruleName:       ruleName,
	}
}

func (r *SimpleRule[T]) Link() string {
	return r.link
}

func (r *SimpleRule[T]) Name() string {
	if r.ruleName != "" {
		return r.ruleName
	}

	if r.nestedBlockType != nil {
		return fmt.Sprintf("%s.%s.%s", r.resourceType, *r.nestedBlockType, r.attributeName)
	}
	return fmt.Sprintf("%s.%s", r.resourceType, r.attributeName)
}

func (r *SimpleRule[T]) Check(runner tflint.Runner) error {
	var dt T
	ctyType, err := toCtyType(dt)
	if err != nil {
		return err
	}

	if r.mustExist {
		exists, resource, err := r.attributeExistsWhereResourceIsSpecified(runner)
		if err != nil {
			return err
		}

		if !exists {
			return runner.EmitIssue(
				r,
				fmt.Sprintf("The attribute `%s` must be specified", r.attributeName),
				resource.DefRange,
			)
		}
	}

	return r.checkAttributes(runner, cty.DynamicPseudoType, func(attr *hclext.Attribute, val cty.Value) error {
		if val.IsNull() || !val.IsKnown() {
			return nil
		}
		for _, exp := range r.expectedValues {
			ctyExp, err := gocty.ToCtyValue(exp, ctyType)
			if err != nil {
				return err
			}
			if ctyExp.Equals(val).True() {
				return nil
			}
		}
		goVal := new(T)
		_ = gocty.FromCtyValue(val, goVal)
		return runner.EmitIssue(
			r,
			fmt.Sprintf("%v is an invalid attribute value of `%s` - expecting (one of) %v", *goVal, r.attributeName, r.expectedValues),
			attr.Range,
		)
	})
}

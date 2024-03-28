package attrvalue

import (
	"fmt"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"

	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// UnknownValueRule checks whether an attribute value is null or part of a variable with no default value.
type UnknownValueRule struct {
	tflint.DefaultRule // Embed the default rule to reuse its implementation
	baseValue
}

var _ tflint.Rule = (*UnknownValueRule)(nil)
var _ AttrValueRule = (*UnknownValueRule)(nil)

func (r *UnknownValueRule) GetNestedBlockType() *string {
	return r.nestedBlockType
}

// NewSimpleRule returns a new rule with the given resource type, and attribute name
func NewUnknownValueRule(resourceType string, attributeName string) *UnknownValueRule {
	return &UnknownValueRule{
		baseValue: newBaseValue(resourceType, nil, attributeName),
	}
}

// NewUnknownValueNestedBlockRule returns a new rule with the given resource type, nested block type, and attribute name
func NewUnknownValueNestedBlockRule(resourceType, nestedBlockType, attributeName string) *UnknownValueRule {
	return &UnknownValueRule{
		baseValue: newBaseValue(resourceType, &nestedBlockType, attributeName),
	}
}

func (r *UnknownValueRule) Name() string {
	return fmt.Sprintf("%s.%s must be null", r.resourceType, r.attributeName)
}

func (r *UnknownValueRule) Check(runner tflint.Runner) error {
	return r.checkAttributes(runner, cty.DynamicPseudoType, func(attr *hclext.Attribute, val cty.Value) error {
		if val.IsKnown() {
			if err := runner.EmitIssue(
				r,
				fmt.Sprintf("invalid attribute value of `%s` - expecting unknown", r.attributeName),
				attr.Expr.Range(),
			); err != nil {
				return err
			}
		}
		return nil
	})
}

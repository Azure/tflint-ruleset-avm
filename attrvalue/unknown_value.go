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
	ruleName string
}

var _ tflint.Rule = (*UnknownValueRule)(nil)
var _ AttrValueRule = (*UnknownValueRule)(nil)

func (r *UnknownValueRule) GetNestedBlockType() *string {
	return r.nestedBlockType
}

// NewUnknownValueRule returns a new rule with the given resource type, and attribute name
func NewUnknownValueRule(resourceType, attributeName, link string, ruleName string) *UnknownValueRule {
	return &UnknownValueRule{
		baseValue: newBaseValue(resourceType, nil, attributeName, true, link, tflint.ERROR),
		ruleName: ruleName,
	}
}

func (r *UnknownValueRule) Link() string {
	return r.link
}

// NewUnknownValueNestedBlockRule returns a new rule with the given resource type, nested block type, and attribute name
func NewUnknownValueNestedBlockRule(resourceType, nestedBlockType, attributeName, link string, ruleName string) *UnknownValueRule {
	return &UnknownValueRule{
		baseValue: newBaseValue(resourceType, &nestedBlockType, attributeName, true, link, tflint.ERROR),
		ruleName: ruleName,
	}
}

func (r *UnknownValueRule) Name() string {
	if r.ruleName != "" {
		return r.ruleName
	}

	if r.nestedBlockType != nil {
		return fmt.Sprintf("%s.%s.%s", r.resourceType, *r.nestedBlockType, r.attributeName)
	}
	return fmt.Sprintf("%s.%s", r.resourceType, r.attributeName)
}

func (r *UnknownValueRule) Check(runner tflint.Runner) error {
	return r.checkAttributes(runner, cty.DynamicPseudoType, func(attr *hclext.Attribute, val cty.Value) error {
		if val.IsKnown() {
			return runner.EmitIssue(
				r,
				fmt.Sprintf("invalid attribute value of `%s` - expecting unknown", r.attributeName),
				attr.Expr.Range(),
			)
		}
		return nil
	})
}

package attrvalue

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// UnknownValueRule checks whether an attribute value is null or part of a variable with no default value.
type UnknownValueRule struct {
	tflint.DefaultRule // Embed the default rule to reuse its implementation

	resourceType    string  // e.g. "azurerm_storage_account"
	nestedBlockType *string // e.g. "sku"
	attributeName   string  // e.g. "account_replication_type"
}

var _ tflint.Rule = (*UnknownValueRule)(nil)
var _ AttrValueRule = (*UnknownValueRule)(nil)

func (r *UnknownValueRule) GetResourceType() string {
	return r.resourceType
}

func (r *UnknownValueRule) GetAttributeName() string {
	return r.attributeName
}

func (r *UnknownValueRule) GetNestedBlockType() *string {
	return r.nestedBlockType
}

// NewSimpleRule returns a new rule with the given resource type, and attribute name
func NewUnknownValueRule(resourceType string, attributeName string) *UnknownValueRule {
	return &UnknownValueRule{
		resourceType:  resourceType,
		attributeName: attributeName,
	}
}

// NewUnknownValueNestedBlockRule returns a new rule with the given resource type, nested block type, and attribute name
func NewUnknownValueNestedBlockRule(resourceType, nestedBlockType, attributeName string) *UnknownValueRule {
	return &UnknownValueRule{
		resourceType:    resourceType,
		nestedBlockType: &nestedBlockType,
		attributeName:   attributeName,
	}
}

func (r *UnknownValueRule) Name() string {
	return fmt.Sprintf("%s.%s must be null", r.resourceType, r.attributeName)
}

func (r *UnknownValueRule) Enabled() bool {
	return true
}

func (r *UnknownValueRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *UnknownValueRule) Check(runner tflint.Runner) error {
	ctx, attrs, diags := fetchAttrsAndContext(r, runner)
	if diags.HasErrors() {
		return fmt.Errorf("could not get partial content: %s", diags)
	}

	for _, attr := range attrs {
		if attr.Name != r.attributeName {
			continue
		}
		val, diags := ctx.EvaluateExpr(attr.Expr, cty.DynamicPseudoType)
		if diags.HasErrors() {
			return fmt.Errorf("could not evaluate expression: %s", diags)
		}
		if val.IsKnown() {
			if err := runner.EmitIssue(
				r,
				fmt.Sprintf("invalid attribute value of `%s` - expecting unknown", r.attributeName),
				attr.Expr.Range(),
			); err != nil {
				return err
			}
		}
	}
	return nil
}

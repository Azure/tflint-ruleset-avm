package attrvalue

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty/gocty"
)

// SimpleRule checks whether a string attribute value is one of the expected values.
// It can be used to check string, number, and bool attributes.
type SimpleRule[T any] struct {
	tflint.DefaultRule // Embed the default rule to reuse its implementation

	resourceType    string  // e.g. "azurerm_storage_account"
	nestedBlockType *string // e.g. "sku"
	attributeName   string  // e.g. "account_replication_type"
	expectedValues  []T     // e.g. []string{"ZRS"}
}

var _ tflint.Rule = (*SimpleRule[any])(nil)
var _ AttrValueRule = (*SimpleRule[any])(nil)

func (r *SimpleRule[T]) GetResourceType() string {
	return r.resourceType
}

func (r *SimpleRule[T]) GetAttributeName() string {
	return r.attributeName
}

func (r *SimpleRule[T]) GetNestedBlockType() *string {
	return r.nestedBlockType
}

// NewSimpleRule returns a new rule with the given resource type, attribute name, and expected values.
func NewSimpleRule[T any](resourceType, attributeName string, expectedValues []T) *SimpleRule[T] {
	return &SimpleRule[T]{
		resourceType:   resourceType,
		attributeName:  attributeName,
		expectedValues: expectedValues,
	}
}

// NewSimpleRule returns a new rule with the given resource type, attribute name, and expected values.
func NewSimpleNestedBlockRule[T any](resourceType, nestedBlockType, attributeName string, expectedValues []T) *SimpleRule[T] {
	return &SimpleRule[T]{
		resourceType:    resourceType,
		attributeName:   attributeName,
		nestedBlockType: &nestedBlockType,
		expectedValues:  expectedValues,
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
	ctx, attrs, diags := fetchAttrsAndContext(r, runner)
	if diags.HasErrors() {
		return fmt.Errorf("could not get partial content: %s", diags)
	}
	for _, attr := range attrs {
		var dt T
		ctyType, err := toCtyType(dt)
		if err != nil {
			return err
		}
		val, diags := ctx.EvaluateExpr(attr.Expr, ctyType)
		if diags.HasErrors() {
			return fmt.Errorf("could not evaluate expression: %s", diags)
		}
		if val.IsNull() {
			continue
		}
		found := false
		for _, exp := range r.expectedValues {
			ctyExp, err := gocty.ToCtyValue(exp, ctyType)
			if err != nil {
				return err
			}
			found = ctyExp.Equals(val).True()
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
	}
	return nil
}

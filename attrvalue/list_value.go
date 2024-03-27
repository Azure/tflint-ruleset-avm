package attrvalue

import (
	"cmp"
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
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
var _ AttrValueRule = (*SimpleRule[any])(nil)

func (r *ListRule[T]) GetResourceType() string {
	return r.resourceType
}

func (r *ListRule[T]) GetAttributeName() string {
	return r.attributeName
}

func (r *ListRule[T]) GetNestedBlockType() *string {
	return nil
}

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
	ctx, attrs, diags := fetchAttrsAndContext(r, runner)
	if diags.HasErrors() {
		return fmt.Errorf("could not get partial content: %s", diags)
	}

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
	for _, attr := range attrs {
		val, diags := ctx.EvaluateExpr(attr.Expr, ctyTypeS)
		if diags.HasErrors() {
			return fmt.Errorf("could not evaluate expression: %s", diags)
		}
		if val.IsNull() {
			continue
		}
		valSet := val.AsValueSet()
		found := false
		for _, exp := range r.expectedValues {
			ctyexp, err := gocty.ToCtyValue(exp, cty.Set(ctyType))
			if err != nil {
				return err
			}
			found = cty.SetValFromValueSet(valSet).Equals(ctyexp).True()
			if found {
				break
			}
		}
		if !found {
			goVal := new([]T)
			_ = gocty.FromCtyValue(val, goVal)
			return runner.EmitIssue(
				r,
				fmt.Sprintf("\"%v\" is an invalid attribute value of `%s` - expecting (one of) %v", *goVal, r.attributeName, r.expectedValues),
				attr.Expr.Range(),
			)
		}
	}
	return nil
}

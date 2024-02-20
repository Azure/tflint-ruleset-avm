package avmhelper

import (
	"github.com/zclconf/go-cty/cty"
	"reflect"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
)

type VariableType struct {
	Type    cty.Type
	Default *typeexpr.Defaults
}

func NewVariableTypeFromExpression(exp hcl.Expression) (VariableType, hcl.Diagnostics) {
	t, d, diags := typeexpr.TypeConstraintWithDefaults(exp)
	if diags.HasErrors() {
		return VariableType{}, diags
	}
	return VariableType{
		Type:    t,
		Default: d,
	}, nil
}

// CheckEqualTypeConstraints checks if two supplied hcl Expressions are in fact type constraints,
// and if they are that they are equal.
func CheckEqualTypeConstraints(type1, type2 VariableType) bool {
	return reflect.DeepEqual(type1, type2)
}

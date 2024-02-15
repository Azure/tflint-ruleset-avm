package avmhelper

import (
	"reflect"

	"github.com/Azure/tflint-ruleset-avm/to"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
)

// CheckTypeConstraintsAreEqual checks if two supplied hcl Expressions are in fact type constraints,
// and if they are that they are equal.
func CheckTypeConstraintsAreEqual(got, want hcl.Expression) (*bool, hcl.Diagnostics) {
	result := to.Ptr(false)
	gotTy, gotDef, diags := typeexpr.TypeConstraintWithDefaults(got)
	if diags.HasErrors() {
		return nil, diags
	}
	wantTy, wantDef, diags := typeexpr.TypeConstraintWithDefaults(want)
	if diags.HasErrors() {
		return nil, diags
	}
	if reflect.DeepEqual(gotTy, wantTy) && reflect.DeepEqual(gotDef, wantDef) {
		result = to.Ptr(true)
	}
	return result, diags
}

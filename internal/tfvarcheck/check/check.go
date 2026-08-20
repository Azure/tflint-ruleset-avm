package check

import (
	"reflect"

	"github.com/Azure/tflint-ruleset-avm/internal/tfvarcheck/varcheck"
	"github.com/zclconf/go-cty/cty"
)

// EqualCtyValue reports whether two cty values are equal.
func EqualCtyValue(got, want cty.Value) bool {
	return got.Equals(want).True()
}

// EqualTypeConstraints reports whether two parsed type constraints and defaults are equal.
func EqualTypeConstraints(type1, type2 varcheck.TypeConstraintWithDefaults) bool {
	return reflect.DeepEqual(type1, type2)
}

// Nullable reports whether a nullable attribute value matches the expected declaration.
func Nullable(got cty.Value, want bool) bool {
	if got.Type() != cty.Bool {
		return false
	}
	if want {
		return got.IsNull()
	}
	return !got.IsNull() && !got.True()
}

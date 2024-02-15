package avmhelper

import "github.com/zclconf/go-cty/cty"

// CheckNullable checks if the nullable attribute is set and if it is set to the desired value.
// If nullable is set to true in the interfaces package:
// - return true if got is null, else return false
// If nullable is set to false in the interfaces package:
// - return true if got is false null, else return false
func CheckNullable(got cty.Value, want bool) bool {
	if !got.Type().IsPrimitiveType() || got.Type() != cty.Bool {
		return false
	}
	if want {
		return got.IsNull()
	}
	if !want && got.IsNull() {
		return false
	}
	return got.False()
}

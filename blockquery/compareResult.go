// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package blockquery

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

type ResultCompareFunc func(cty.Value, ...cty.Value) (bool, string, error)

// IsKnown is a compare function that checks if the result is known.
func IsKnown(got cty.Value, _ ...cty.Value) (bool, string, error) {
	if !got.IsKnown() {
		return false, "returned value is unknown", nil
	}
	return true, "", nil
}

// IsNotKnown is a compare function that checks if the result is unknown.
func IsNotKnown(got cty.Value, _ ...cty.Value) (bool, string, error) {
	if got.IsKnown() {
		return false, "returned value is known", nil
	}
	return true, "", nil
}

// IsNotNull is a compare function that checks if the result exists.
func IsNotNull(got cty.Value, _ ...cty.Value) (bool, string, error) {
	if got.IsNull() {
		return false, "returned value is null but not expected to be", nil
	}
	return true, "", nil
}

// IsNull is a compare function that checks if the result does not exist.
func IsNull(got cty.Value, _ ...cty.Value) (bool, string, error) {
	if !got.IsNull() {
		return false, "returned value is not null but expected to be", nil
	}
	return true, "", nil
}

// EachIsOneOf is a compare function that checks if the result is one of the expected values.
// This is useful when the result is an array and you want to check if each element is one of the expected values.
func EachIsOneOf(got cty.Value, expected ...cty.Value) (bool, string, error) {
	if !got.Type().IsListType() {
		return false, "", fmt.Errorf("expected a list but got %s", got.Type().FriendlyName())
	}
	results := make([]bool, 0, got.LengthInt())
	it := got.ElementIterator()
	for it.Next() {
		_, v := it.Element()
		results = append(results, compareResults(v, expected))
	}
	if !allTrue(results...) {
		return false, fmt.Sprintf("returned values `%s` not in expected values `%v`", fmtCty(got), fmtCty(cty.ListVal(expected))), nil
	}
	return true, "", nil
}

// IsOneOf is a compare function that checks if the result is one of the expected values.
func IsOneOf(got cty.Value, expected ...cty.Value) (bool, string, error) {
	ok := compareResults(got, expected)
	if !ok {
		return false, fmt.Sprintf("returned value `%s` not in expected values `%v`", fmtCty(got), fmtCty(cty.ListVal(expected))), nil
	}
	return ok, "", nil
}

// compareResults compares the result with the expected values.
func compareResults(got cty.Value, want []cty.Value) bool {
	if !got.IsKnown() || got.IsNull() {
		return false
	}
	for _, w := range want {
		cnv, err := convert.Convert(got, w.Type())
		if err != nil {
			continue
		}
		ok := w.Equals(cnv).True()
		if ok {
			return true
		}
	}
	return false
}

// allTrue checks if all the values are true.
func allTrue(in ...bool) bool {
	for _, b := range in {
		if !b {
			return false
		}
	}
	return true
}

// fmtCty formats the cty value to a string.
func fmtCty(in cty.Value) string {
	switch in.Type() {
	case cty.Bool:
		return fmt.Sprintf("%t", in.True())
	case cty.Number:
		big := in.AsBigFloat()
		if big.IsInt() {
			i, _ := big.Int64()
			return fmt.Sprintf("%d", i)
		}
		f, _ := big.Float64()
		return fmt.Sprintf("%f", f)
	case cty.String:
		return in.AsString()
	}
	if in.Type().IsListType() || in.Type().IsTupleType() {
		res := make([]string, 0, in.LengthInt())
		it := in.ElementIterator()
		for it.Next() {
			_, v := it.Element()
			res = append(res, fmtCty(v))
		}
		return fmt.Sprintf("%s", res)
	}
	return in.GoString()
}

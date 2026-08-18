// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package blockquery

import (
	"github.com/zclconf/go-cty/cty"
)

// NewIntResults creates results of type Number,
// E.g. If you expect vales to be 1, 2, or 3 then use
// NewIntResults(1, 2, 3).
func NewIntResults(vals ...int) []cty.Value {
	results := make([]cty.Value, len(vals))
	for i, val := range vals {
		results[i] = cty.NumberIntVal(int64(val))
	}
	return results
}

// NewFloatResults creates results of type Number,
// E.g. If you expect vales to be 1.1, 2.2, or 3.3 then use
// NewFloatResults(1.1, 2.2, 3.3).
func NewFloatResults(vals ...float64) []cty.Value {
	results := make([]cty.Value, len(vals))
	for i, val := range vals {
		results[i] = cty.NumberFloatVal(val)
	}
	return results
}

// NewStringResults creates results of type String,
// E.g. If you expect vales to be "a", "b", or "c" then use
// NewStringResults("a", "b", "c").
func NewStringResults(vals ...string) []cty.Value {
	results := make([]cty.Value, len(vals))
	for i, val := range vals {
		results[i] = cty.StringVal(val)
	}
	return results
}

// NewBoolResult creates a result.
// E.g. If you expect vales to be true, then use NewBoolResult(true).
// Only specify one value, if the length of the values is not 1, it will panic.
func NewBoolResult(vals ...bool) cty.Value {
	if len(vals) != 1 {
		panic("NewBoolResult expects exactly one value")
	}
	return cty.BoolVal(vals[0])
}

// NewListResults creates results for complex types.
//
// E.g. If you expect vales to be an object: {"nested": "value"},
// then use NewListResults(map[string]any{"nested": "value"}).
//
// If you want to compare a list, then use NewListResults(cty.ListVal(cty.NumberIntVal(1), cty.NumberIntVal(2), cty.NumberIntVal(3))).
func NewListResults(vals ...[]cty.Value) []cty.Value {
	results := make([]cty.Value, len(vals))
	for i, val := range vals {
		results[i] = cty.ListVal(val)
	}
	return results
}

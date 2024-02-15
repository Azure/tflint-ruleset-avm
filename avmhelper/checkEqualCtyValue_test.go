package avmhelper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"

	"github.com/Azure/tflint-ruleset-avm/avmhelper"
)

func TestCheckEqualCtyValue(t *testing.T) {
	tests := []struct {
		Name     string
		Got      cty.Value
		Want     cty.Value
		Expected bool
	}{
		{
			Name:     "Equal values",
			Got:      cty.StringVal("hello"),
			Want:     cty.StringVal("hello"),
			Expected: true,
		},
		{
			Name:     "Not equal values",
			Got:      cty.StringVal("hello"),
			Want:     cty.StringVal("world"),
			Expected: false,
		},
		{
			Name:     "Equal list values",
			Got:      cty.ListVal([]cty.Value{cty.StringVal("hello"), cty.StringVal("world")}),
			Want:     cty.ListVal([]cty.Value{cty.StringVal("hello"), cty.StringVal("world")}),
			Expected: true,
		},
		{
			Name:     "Not equal list values",
			Got:      cty.ListVal([]cty.Value{cty.StringVal("hello"), cty.StringVal("world")}),
			Want:     cty.ListVal([]cty.Value{cty.StringVal("world"), cty.StringVal("hello")}),
			Expected: false,
		},
		{
			Name:     "Equal map values",
			Got:      cty.MapVal(map[string]cty.Value{"key1": cty.StringVal("value1"), "key2": cty.StringVal("value2")}),
			Want:     cty.MapVal(map[string]cty.Value{"key1": cty.StringVal("value1"), "key2": cty.StringVal("value2")}),
			Expected: true,
		},
		{
			Name:     "Not equal map values",
			Got:      cty.MapVal(map[string]cty.Value{"key1": cty.StringVal("value1"), "key2": cty.StringVal("value2")}),
			Want:     cty.MapVal(map[string]cty.Value{"key1": cty.StringVal("value2"), "key2": cty.StringVal("value1")}),
			Expected: false,
		},
		{
			Name:     "Equal tuple values",
			Got:      cty.TupleVal([]cty.Value{cty.StringVal("hello"), cty.StringVal("world")}),
			Want:     cty.TupleVal([]cty.Value{cty.StringVal("hello"), cty.StringVal("world")}),
			Expected: true,
		},
		{
			Name:     "Not equal tuple values",
			Got:      cty.TupleVal([]cty.Value{cty.StringVal("hello"), cty.StringVal("world")}),
			Want:     cty.TupleVal([]cty.Value{cty.StringVal("world"), cty.StringVal("hello")}),
			Expected: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			ret := avmhelper.CheckEqualCtyValue(tc.Got, tc.Want)

			assert.Equal(t, tc.Expected, ret)
		})
	}
}

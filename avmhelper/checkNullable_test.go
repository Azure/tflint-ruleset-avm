package avmhelper_test

import (
	"testing"

	"github.com/Azure/tflint-ruleset-avm/avmhelper"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestCheckNullable(t *testing.T) {
	tests := []struct {
		Name     string
		Got      cty.Value
		Want     bool
		Expected bool
	}{
		{
			Name:     "Nullable true, got should be null",
			Got:      cty.NullVal(cty.Bool),
			Want:     true,
			Expected: true,
		},
		{
			Name:     "Nullable true, got is true but should not be",
			Got:      cty.True,
			Want:     true,
			Expected: false,
		},
		{
			Name:     "Nullable false, got is null but should not be",
			Got:      cty.NullVal(cty.Bool),
			Want:     false,
			Expected: false,
		},
		{
			Name:     "Nullable fasle, got should be false",
			Got:      cty.False,
			Want:     false,
			Expected: true,
		},
		{
			Name:     "got is not a primitive type",
			Got:      cty.ListValEmpty(cty.String),
			Want:     true,
			Expected: false,
		},
		{
			Name:     "got is an incorrect primative type",
			Got:      cty.StringVal("test"),
			Want:     true,
			Expected: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			result := avmhelper.CheckNullable(tc.Got, tc.Want)
			assert.Equal(t, tc.Expected, result)
		})
	}
}

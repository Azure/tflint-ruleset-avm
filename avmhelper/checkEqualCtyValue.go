package avmhelper

import "github.com/zclconf/go-cty/cty"

func CheckEqualCtyValue(got, want cty.Value) bool {
	ret := got.Equals(want)
	return ret.True()
}

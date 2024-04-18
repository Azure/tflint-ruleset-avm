package interfaces

import (
	"github.com/matt-FFFFFF/tfvarcheck/varcheck"
	"github.com/zclconf/go-cty/cty"
)

// TagsTypeString is the type constraint string for tags.
// When updating the type constraint string, make sure to also update the two
// private endpoint interfaces (the one with subresource and the one without).
var TagsTypeString = `map(string)`

var tagsType = StringToTypeConstraintWithDefaults(TagsTypeString)

var Tags = AvmInterface{
	VarCheck:      varcheck.NewVarCheck(tagsType, cty.NullVal(cty.DynamicPseudoType), true),
	RuleName:      "tags",
	VarTypeString: TagsTypeString,
	RuleEnabled:   true,
	RuleLink:      "https://azure.github.io/Azure-Verified-Modules/specs/shared/interfaces/#tags",
}

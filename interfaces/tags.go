package interfaces

import (
	"github.com/matt-FFFFFF/tfvarcheck/varcheck"
	"github.com/zclconf/go-cty/cty"
)

var TagsTypeString = `map(string)`

var tagsType = StringToTypeConstraintWithDefaults(TagsTypeString)

var Tags = AvmInterface{
	VarCheck:      varcheck.NewVarCheck(tagsType, cty.NullVal(cty.DynamicPseudoType), true),
	RuleName:      "tags",
	VarTypeString: TagsTypeString,
	RuleEnabled:   true,
	RuleLink:      "https://azure.github.io/Azure-Verified-Modules/specs/shared/interfaces/#tags",
}

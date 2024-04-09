package interfaces

import (
	"github.com/matt-FFFFFF/tfvarcheck/varcheck"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

var LocationTypeString = `string`

var locationType = StringToTypeConstraintWithDefaults(LocationTypeString)

var Location = AvmInterface{
	VarCheck:      varcheck.NewVarCheck(locationType, cty.UnknownVal(cty.String), false),
	RuleName:      "location",
	VarTypeString: LocationTypeString,
	RuleEnabled:   true,
	RuleLink:      "https://azure.github.io/Azure-Verified-Modules/specs/shared/#id-rmnfr2---category-inputs---parametervariable-naming",
	RuleSeverity:  tflint.ERROR,
}

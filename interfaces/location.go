package interfaces

import (
	"github.com/Azure/tflint-ruleset-avm/internal/tfvarcheck/varcheck"
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
	RuleLink:      "https://azure.github.io/Azure-Verified-Modules/specs/tf/res/#id-rmnfr2---category-inputs---parametervariable-naming",
	RuleSeverity:  tflint.ERROR,
}

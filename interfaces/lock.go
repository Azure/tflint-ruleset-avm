package interfaces

import (
	"github.com/matt-FFFFFF/tfvarcheck/varcheck"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// LockTypeString is the type constraint string for lock interface.
// When updating the type constraint string, make sure to also update the two
// private endpoint interfaces (the one with subresource and the one without).
var LockTypeString = `object({
	kind = string
	name = optional(string, null)
})`

var lockType = StringToTypeConstraintWithDefaults(LockTypeString)

var Lock = AvmInterface{
	VarCheck:      varcheck.NewVarCheck(lockType, cty.NullVal(cty.DynamicPseudoType), true),
	RuleName:      "lock",
	VarTypeString: LockTypeString,
	RuleEnabled:   true,
	RuleLink:      "https://azure.github.io/Azure-Verified-Modules/specs/shared/interfaces/#resource-locks",
	RuleSeverity:  tflint.ERROR,
}

package interfaces

import (
	"github.com/matt-FFFFFF/tfvarcheck/varcheck"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

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

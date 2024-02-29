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

var lockType = stringToTypeConstraintWithDefaults(LockTypeString)

// TerraformVarLock represents the definition of the AVM Locks interface.
var Lock = AvmInterface{
	Name:          "lock",
	VarType:       varcheck.NewVarCheck(lockType, cty.NullVal(cty.DynamicPseudoType), true),
	VarTypeString: LockTypeString,
	Enabled:       true,
	Link:          "https://azure.github.io/Azure-Verified-Modules/specs/shared/interfaces/#resource-locks",
	Severity:      tflint.ERROR,
}

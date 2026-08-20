package interfaces

import (
	"github.com/Azure/tflint-ruleset-avm/internal/tfvarcheck/varcheck"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

const LockV1TypeString = `object({
  kind = string
  name = optional(string, null)
})`

const LockV2TypeString = `object({
  kind  = string
  name  = optional(string, null)
  notes = optional(string, null)
})`

const LockTypeString = LockV2TypeString

func newLockInterface(typeString string) AvmInterface {
	return AvmInterface{
		VarCheck:      varcheck.NewVarCheck(StringToTypeConstraintWithDefaults(typeString), cty.NullVal(cty.DynamicPseudoType), true),
		RuleName:      "lock",
		VarTypeString: typeString,
		RuleEnabled:   true,
		RuleLink:      "https://azure.github.io/Azure-Verified-Modules/specs/tf/interfaces/#resource-locks",
		RuleSeverity:  tflint.ERROR,
	}
}

var LockV1 = newLockInterface(LockV1TypeString)

var LockV2 = newLockInterface(LockV2TypeString)

var Lock = LockV2

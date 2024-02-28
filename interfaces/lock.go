package interfaces

import (
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// TerraformVarLock represents the definition of the AVM Locks interface.
var Lock = AvmInterface{
	Name: "lock",
	TypeStr: `object({
		kind = string
		name = optional(string, null)
	})`,
	Enabled:  true,
	Link:     "https://azure.github.io/Azure-Verified-Modules/specs/shared/interfaces/#resource-locks",
	Default:  cty.NullVal(cty.DynamicPseudoType),
	Nullable: true,
	Severity: tflint.ERROR,
}

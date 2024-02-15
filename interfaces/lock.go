package interfaces

import "github.com/zclconf/go-cty/cty"

// Lock represents the lock interface.
var Lock = AVMInterface{
	Name: "lock",
	Type: `object({
		kind = string
		name = optional(string, null)
	})`,
	Enabled:  true,
	Link:     "https://azure.github.io/Azure-Verified-Modules/specs/shared/interfaces/#resource-locks",
	Default:  cty.NullVal(cty.DynamicPseudoType),
	Nullable: true,
}

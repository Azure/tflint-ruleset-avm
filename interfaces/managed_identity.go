package interfaces

import "github.com/zclconf/go-cty/cty"

// Lock represents the lock interface.
var ManagedIdentity = AvmInterface{
	Name: "managed_identities",
	TypeStr: `object({
		kind = string
		name = optional(string, null)
	})`,
	Enabled:  true,
	Link:     "https://azure.github.io/Azure-Verified-Modules/specs/shared/interfaces/#managed-identities",
	Default:  cty.EmptyObjectVal,
	Nullable: false,
}

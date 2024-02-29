package interfaces

import (
	"github.com/matt-FFFFFF/tfvarcheck/varcheck"
	"github.com/zclconf/go-cty/cty"
)

var ManagedIdentityTypeString = `object({
		kind = string
		name = optional(string, null)
	})`

// Lock represents the lock interface.
var ManagedIdentity = AvmInterface{
	Name:          "managed_identities",
	VarType:       varcheck.NewVarCheck(stringToTypeConstraintWithDefaults(ManagedIdentityTypeString), cty.EmptyObjectVal, false),
	VarTypeString: ManagedIdentityTypeString,
	Enabled:       true,
	Link:          "https://azure.github.io/Azure-Verified-Modules/specs/shared/interfaces/#managed-identities",
}

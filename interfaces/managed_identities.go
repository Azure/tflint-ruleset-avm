package interfaces

import (
	"github.com/matt-FFFFFF/tfvarcheck/varcheck"
	"github.com/zclconf/go-cty/cty"
)

var ManagedIdentitiesTypeString = `object({
		kind = string
		name = optional(string, null)
	})`

// Lock represents the lock interface.
var ManagedIdentities = AvmInterface{
	Name:          "managed_identities",
	VarType:       varcheck.NewVarCheck(stringToTypeConstraintWithDefaults(ManagedIdentitiesTypeString), cty.EmptyObjectVal, false),
	VarTypeString: ManagedIdentitiesTypeString,
	Enabled:       true,
	Link:          "https://azure.github.io/Azure-Verified-Modules/specs/shared/interfaces/#managed-identities",
}

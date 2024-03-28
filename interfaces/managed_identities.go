package interfaces

import (
	"github.com/matt-FFFFFF/tfvarcheck/varcheck"
	"github.com/zclconf/go-cty/cty"
)

var ManagedIdentitiesTypeString = `object({
	system_assigned            = optional(bool, false)
	user_assigned_resource_ids = optional(set(string), [])
})`

var ManagedIdentities = AvmInterface{
	VarCheck:      varcheck.NewVarCheck(StringToTypeConstraintWithDefaults(ManagedIdentitiesTypeString), cty.EmptyObjectVal, false),
	RuleName:      "managed_identities",
	VarTypeString: ManagedIdentitiesTypeString,
	RuleEnabled:   true,
	RuleLink:      "https://azure.github.io/Azure-Verified-Modules/specs/shared/interfaces/#managed-identities",
}

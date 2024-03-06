package interfaces

import (
	"github.com/matt-FFFFFF/tfvarcheck/varcheck"
	"github.com/zclconf/go-cty/cty"
)

var RoleAssignmentsTypeString = `map(object({
      role_definition_id_or_name             = string
      principal_id                           = string
      description                            = optional(string, null)
      skip_service_principal_aad_check       = optional(bool, false)
      condition                              = optional(string, null)
      condition_version                      = optional(string, null)
      delegated_managed_identity_resource_id = optional(string, null)
    }))`

var roleAssignmentsType = StringToTypeConstraintWithDefaults(RoleAssignmentsTypeString)

var RoleAssignments = AvmInterface{
	VarCheck:      varcheck.NewVarCheck(roleAssignmentsType, cty.EmptyObjectVal, false),
	RuleName:      "role_assignments",
	VarTypeString: RoleAssignmentsTypeString,
	RuleEnabled:   true,
	RuleLink:      "https://azure.github.io/Azure-Verified-Modules/specs/shared/interfaces/#role-assignments",
}

package interfaces

import (
	"github.com/Azure/tflint-ruleset-avm/internal/tfvarcheck/varcheck"
	"github.com/zclconf/go-cty/cty"
)

const RoleAssignmentsV1TypeString = `map(object({
  role_definition_id_or_name             = string
  principal_id                           = string
  description                            = optional(string, null)
  skip_service_principal_aad_check       = optional(bool, false)
  condition                              = optional(string, null)
  condition_version                      = optional(string, null)
  delegated_managed_identity_resource_id = optional(string, null)
  principal_type                         = optional(string, null)
}))`

const RoleAssignmentsV2TypeString = `map(object({
  name                                   = optional(string, null)
  role_definition_id_or_name             = string
  principal_id                           = string
  description                            = optional(string, null)
  skip_service_principal_aad_check       = optional(bool, false)
  condition                              = optional(string, null)
  condition_version                      = optional(string, null)
  delegated_managed_identity_resource_id = optional(string, null)
  principal_type                         = optional(string, null)
}))`

const RoleAssignmentsTypeString = RoleAssignmentsV2TypeString

func newRoleAssignmentsInterface(typeString string) AvmInterface {
	return AvmInterface{
		VarCheck:      varcheck.NewVarCheck(StringToTypeConstraintWithDefaults(typeString), cty.EmptyObjectVal, false),
		RuleName:      "avm_interface_role_assignments",
		VariableName:  "role_assignments",
		VarTypeString: typeString,
		RuleEnabled:   true,
		RuleLink:      "https://azure.github.io/Azure-Verified-Modules/specs/tf/interfaces/#role-assignments",
	}
}

var RoleAssignmentsV1 = newRoleAssignmentsInterface(RoleAssignmentsV1TypeString)

var RoleAssignmentsV2 = newRoleAssignmentsInterface(RoleAssignmentsV2TypeString)

var RoleAssignments = RoleAssignmentsV2

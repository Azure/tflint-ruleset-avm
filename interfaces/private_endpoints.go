package interfaces

import (
	"fmt"

	"github.com/matt-FFFFFF/tfvarcheck/varcheck"
	"github.com/zclconf/go-cty/cty"
)

func privateEndpointTypeString(includeLockNotes, includeMemberName, includeRoleAssignmentName bool) string {
	lockType := `    kind = string
    name = optional(string, null)`
	if includeLockNotes {
		lockType = `    kind  = string
    name  = optional(string, null)
    notes = optional(string, null)`
	}

	memberName := ""
	if includeMemberName {
		memberName = "\n    member_name        = optional(string)"
	}

	roleAssignmentName := ""
	if includeRoleAssignmentName {
		roleAssignmentName = "    name                                   = optional(string, null)\n"
	}

	return fmt.Sprintf(`map(object({
  name = optional(string, null)
  role_assignments = optional(map(object({
%s    role_definition_id_or_name             = string
    principal_id                           = string
    description                            = optional(string, null)
    skip_service_principal_aad_check       = optional(bool, false)
    condition                              = optional(string, null)
    condition_version                      = optional(string, null)
    delegated_managed_identity_resource_id = optional(string, null)
    principal_type                         = optional(string, null)
  })), {})
  lock = optional(object({
%s
  }), null)
  tags                                    = optional(map(string), null)
  subnet_resource_id                      = string
  subresource_name                        = optional(string, null)
  private_dns_zone_group_name             = optional(string, "default")
  private_dns_zone_resource_ids           = optional(set(string), [])
  application_security_group_associations = optional(map(string), {})
  private_service_connection_name         = optional(string, null)
  network_interface_name                  = optional(string, null)
  location                                = optional(string, null)
  resource_group_name                     = optional(string, null)
  ip_configurations = optional(map(object({
    name               = string
    private_ip_address = string%s
  })), {})
}))`, roleAssignmentName, lockType, memberName)
}

func newPrivateEndpointsInterface(typeString string) AvmInterface {
	return AvmInterface{
		VarCheck:      varcheck.NewVarCheck(StringToTypeConstraintWithDefaults(typeString), cty.EmptyObjectVal, false),
		RuleName:      "private_endpoints",
		VarTypeString: typeString,
		RuleEnabled:   true,
		RuleLink:      "https://azure.github.io/Azure-Verified-Modules/specs/tf/interfaces/#private-endpoints",
	}
}

var privateEndpointVariants = func() []AvmInterface {
	variants := make([]AvmInterface, 0, 8)
	for _, includeLockNotes := range []bool{true, false} {
		for _, includeMemberName := range []bool{true, false} {
			for _, includeRoleAssignmentName := range []bool{true, false} {
				variants = append(variants, newPrivateEndpointsInterface(
					privateEndpointTypeString(includeLockNotes, includeMemberName, includeRoleAssignmentName),
				))
			}
		}
	}
	return variants
}()

var PrivateEndpointTypeString = privateEndpointVariants[0].VarTypeString

var PrivateEndpoints = privateEndpointVariants[0]

package interfaces

import (
	"github.com/matt-FFFFFF/tfvarcheck/varcheck"
	"github.com/zclconf/go-cty/cty"
)

var PrivateEndpointTypeString = `map(object({
  name               = optional(string, null)
  role_assignments   = optional(map(object({
    role_definition_id_or_name             = string
    principal_id                           = string
    description                            = optional(string, null)
    skip_service_principal_aad_check       = optional(bool, false)
    condition                              = optional(string, null)
    condition_version                      = optional(string, null)
    delegated_managed_identity_resource_id = optional(string, null)
  })), {})
  lock               = optional(object({
    kind = string
    name = optional(string, null)
  }), null)
  tags               = optional(map(string), null)
  subnet_resource_id = string
  private_dns_zone_group_name             = optional(string, "default")
  private_dns_zone_resource_ids           = optional(set(string), [])
  application_security_group_associations = optional(map(string), {})
  private_service_connection_name         = optional(string, null)
  network_interface_name                  = optional(string, null)
  location                                = optional(string, null)
  resource_group_name                     = optional(string, null)
  ip_configurations = optional(map(object({
    name               = string
    private_ip_address = string
  })), {})
}))`

var PrivateEndpoints = AvmInterface{
	VarCheck:      varcheck.NewVarCheck(StringToTypeConstraintWithDefaults(PrivateEndpointTypeString), cty.EmptyObjectVal, false),
	RuleName:      "private_endpoints",
	VarTypeString: PrivateEndpointTypeString,
	RuleEnabled:   true,
	RuleLink:      "https://azure.github.io/Azure-Verified-Modules/specs/shared/interfaces/#private-endpoints",
}

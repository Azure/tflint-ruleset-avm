package interfaces

import (
	"github.com/matt-FFFFFF/tfvarcheck/varcheck"
	"github.com/zclconf/go-cty/cty"
)

var PrivateEndpointWithSubresourceNameTypeString = `map(object({
      name               = optional(string, null)
      role_assignments   = optional(map(object({})), {})
      lock               = optional(object({}), {})
      tags               = optional(map(any), null)
      subnet_resource_id = string
      subresource_name                        = string
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

var PrivateEndpointsWithSubresourceName = AvmInterface{
	VarCheck:      varcheck.NewVarCheck(StringToTypeConstraintWithDefaults(PrivateEndpointWithSubresourceNameTypeString), cty.EmptyObjectVal, false),
	RuleName:      "private_endpoints",
	VarTypeString: PrivateEndpointWithSubresourceNameTypeString,
	RuleEnabled:   true,
	RuleLink:      "https://azure.github.io/Azure-Verified-Modules/specs/shared/interfaces/#private-endpoints",
}

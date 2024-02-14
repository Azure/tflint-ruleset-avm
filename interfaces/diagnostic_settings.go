package interfaces

import "github.com/zclconf/go-cty/cty"

var DiagnosticSettings = AVMInterface{
	Name: "diagnostic_settings",
	Type: `map(object({
		name                                     = optional(string, null)
		log_categories                           = optional(set(string), [])
		log_groups                               = optional(set(string), ["allLogs"])
		metric_categories                        = optional(set(string), ["AllMetrics"])
		log_analytics_destination_type           = optional(string, "Dedicated")
		workspace_resource_id                    = optional(string, null)
		storage_account_resource_id              = optional(string, null)
		event_hub_authorization_rule_resource_id = optional(string, null)
		event_hub_name                           = optional(string, null)
		marketplace_partner_resource_id          = optional(string, null)
	}))`,
	Enabled:  true,
	Link:     "https://azure.github.io/Azure-Verified-Modules/specs/shared/interfaces/#diagnostic-settings",
	Default:  cty.EmptyObjectVal, // need to use this rather than an empty map value.
	Nullable: false,
}

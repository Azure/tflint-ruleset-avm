package interfaces

import (
	"github.com/matt-FFFFFF/tfvarcheck/varcheck"
	"github.com/zclconf/go-cty/cty"
)

var DiagnosticTypeString = `map(object({
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
	}))`

var diagnosticType = stringToTypeConstraintWithDefaults(DiagnosticTypeString)

// DiagnosticSettings represents the diagnostic_settings interface.
var DiagnosticSettings = AvmInterface{
	Name:          "diagnostic_settings",
	VarType:       varcheck.NewVarCheck(diagnosticType, cty.EmptyObjectVal, false),
	VarTypeString: DiagnosticTypeString,
	Enabled:       true,
	Link:          "https://azure.github.io/Azure-Verified-Modules/specs/shared/interfaces/#diagnostic-settings",
}

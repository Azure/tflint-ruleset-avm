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

var diagnosticType = StringToTypeConstraintWithDefaults(DiagnosticTypeString)

var DiagnosticSettings = AvmInterface{
	VarCheck:      varcheck.NewVarCheck(diagnosticType, cty.EmptyObjectVal, false),
	RuleName:      "diagnostic_settings",
	VarTypeString: DiagnosticTypeString,
	RuleEnabled:   true,
	RuleLink:      "https://azure.github.io/Azure-Verified-Modules/specs/tf/interfaces/#diagnostic-settings",
}

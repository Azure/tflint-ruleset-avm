package waf

import "github.com/Azure/tflint-ruleset-avm/attrvalue"

var AzurermServicePlanZoneBalancingEnabled = func() *attrvalue.SimpleRule {
	return attrvalue.NewSimpleBoolRule(
		"azurerm_service_plan",
		"zone_balancing_enabled",
		[]any{true},
	)
}

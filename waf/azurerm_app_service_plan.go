package waf

import "github.com/Azure/tflint-ruleset-avm/attrvalue"

func (wf WafRules) AzurermAppServicePlanZoneRedundant() *attrvalue.SimpleRule[bool] {
	return attrvalue.NewSimpleRule[bool](
		"azurerm_app_service_plan",
		"zone_redundant",
		[]bool{true},
		"https://learn.microsoft.com/en-us/azure/reliability/reliability-app-service?tabs=cli#-asp-2--use-an-app-service-plan-that-supports-availability-zones",
		false,
		"",
	)
}
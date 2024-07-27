package waf

import (
	"github.com/Azure/tflint-ruleset-avm/attrvalue"
)

func (wf WafRules) AzurermServicePlanZoneBalancingEnabled() *attrvalue.SimpleRule[bool] {
	return attrvalue.NewSimpleRule[bool](
		"azurerm_service_plan",
		"zone_balancing_enabled",
		[]bool{true},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library/services/web/app-service-plan/#asp-1---migrate-app-service-to-availability-zone-support",
		false,
		"",
	)
}

/*
func (wf WafRules) AzurermServicePlanSku() *attrvalue.SimpleRule[string] {
	return attrvalue.NewSimpleNestedBlockRule[string](
		"azurerm_service_plan",
		"sku_name",
		"name",
		[]string{"", ""},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library-v2/azure-resources/Web/serverFarms/#use-standard-or-premium-tier",
		false,
		"",
	)
}

*/

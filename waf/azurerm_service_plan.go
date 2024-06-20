package waf

import (
	"github.com/Azure/tflint-ruleset-avm/attrvalue"
)

func AzurermServicePlanZoneBalancingEnabled() *attrvalue.SimpleRule[bool] {
	return attrvalue.NewSimpleRule[bool](
		"azurerm_service_plan",
		"zone_balancing_enabled",
		[]bool{true},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library/services/web/app-service-plan/#asp-1---migrate-app-service-to-availability-zone-support",
		false,
	)
}

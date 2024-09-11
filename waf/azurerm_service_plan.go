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

func (wf WafRules) AzurermServicePlanSku() *attrvalue.SimpleRule[string] {
	return attrvalue.NewSimpleNestedBlockRule[string](
		"azurerm_service_plan",
		"sku_name",
		"name",
		[]string{"I1", "I2", "I3", "I1v2", "I2v2", "I3v2", "I4V2", "I5V2", "I6V2", "I1MV2", "I2MV2", "I3MV2", "I4MV2", "I5MV2", "S1", "S2", "S3", "WS1", "WS2", "WS3", "P1V2", "P2V2", "P3V2", "P0V3", "P1V3", "P2V3", "P3V3", "P1MV3", "P2MV3", "P3MV3", "P4MV3", "P5MV3"},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library-v2/azure-resources/Web/serverFarms/#use-standard-or-premium-tier",
		false,
		"",
	)
}

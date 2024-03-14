package waf

import (
	"github.com/Azure/tflint-ruleset-avm/attrvalue"
)

func AzurermServicePlanZoneBalancingEnabled() *attrvalue.SimpleRule[bool] {
	return attrvalue.NewSimpleRule[bool](
		"azurerm_service_plan",
		"zone_balancing_enabled",
		[]bool{true},
	)
}

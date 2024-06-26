package waf

import "github.com/Azure/tflint-ruleset-avm/attrvalue"

func (wf WafRules) AzurermApplicationGatewayZones() *attrvalue.SetRule[int] {
	return attrvalue.NewSetRule(
		"azurerm_application_gateway",
		"zones",
		[][]int{{1, 2, 3}},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library/services/networking/application-gateway/#agw-1---set-a-minimum-instance-count-of-2",
		"",
	)
}

func (wf WafRules) AzurermApplicationGatewaySku() *attrvalue.SimpleRule[string] {
	return attrvalue.NewSimpleNestedBlockRule[string](
		"azurerm_application_gateway",
		"sku",
		"name",
		[]string{"Standard_v2", "WAF_v2"},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library/services/networking/application-gateway/#agw-4---use-application-gw-v2-instead-of-v1",
		false,
		"",
	)
}

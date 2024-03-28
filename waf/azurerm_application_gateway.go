package waf

import "github.com/Azure/tflint-ruleset-avm/attrvalue"

func AzurermApplicationGatewayZones() *attrvalue.SetRule[int] {
	return attrvalue.NewListRule(
		"azurerm_application_gateway",
		"zones",
		//TODO: What if there's no three zones in the given region?
		[][]int{{1, 2, 3}},
	)
}

func AzurermApplicationGatewaySku() *attrvalue.SimpleRule[string] {
	return attrvalue.NewSimpleNestedBlockRule[string](
		"azurerm_application_gateway",
		"sku",
		"name",
		[]string{"Standard_v2", "WAF_v2"},
	)
}

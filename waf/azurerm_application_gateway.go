package waf

import "github.com/Azure/tflint-ruleset-avm/attrvalue"

func AzurermApplicationGatewayZones() *attrvalue.ListRule[int] {
	return attrvalue.NewListRule(
		"azurerm_application_gateway",
		"zones",
		[][]int{{1, 2, 3}},
	)
}

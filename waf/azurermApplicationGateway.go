package waf

import "github.com/Azure/tflint-ruleset-avm/attrvalue"

var AzurermApplicationGatewayZones = func() *attrvalue.ListNumberRule {
	return attrvalue.NewListNumberRule(
		"azurerm_application_gateway",
		"zones",
		[][]int{{1, 2, 3}},
	)
}

package waf

import "github.com/Azure/tflint-ruleset-avm/attrvalue"

var AzurermPublicIpSku = func() *attrvalue.SimpleRule {
	return attrvalue.NewSimpleStringRule(
		"azurerm_public_ip",
		"sku",
		[]any{"Standard"},
	)
}

var AzurermPublicIpZones = func() *attrvalue.ListNumberRule {
	return attrvalue.NewListNumberRule(
		"azurerm_public_ip",
		"zones",
		[][]int{{1, 2, 3}},
	)
}

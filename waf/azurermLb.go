package waf

import "github.com/Azure/tflint-ruleset-avm/attrvalue"

var AzurermLbSku = func() *attrvalue.SimpleRule {
	return attrvalue.NewSimpleStringRule(
		"azurerm_lb",
		"sku",
		[]any{"Standard"},
	)
}

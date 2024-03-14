package waf

import (
	"github.com/Azure/tflint-ruleset-avm/attrvalue"
)

func AzurermLbSku() *attrvalue.SimpleRule[string] {
	return attrvalue.NewSimpleRule[string](
		"azurerm_lb",
		"sku",
		[]string{"Standard"},
	)
}

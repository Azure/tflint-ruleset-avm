package waf

import (
	"github.com/Azure/tflint-ruleset-avm/attrvalue"
)

func AzurermLbSku() *attrvalue.SimpleRule[string] {
	return attrvalue.NewSimpleRule[string](
		"azurerm_lb",
		"sku",
		[]string{"Standard"},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library/services/networking/load-balancer/#lb-1---use-standard-load-balancer-sku",
		false,
	)
}

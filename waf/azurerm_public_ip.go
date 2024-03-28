package waf

import (
	"github.com/Azure/tflint-ruleset-avm/attrvalue"
)

func AzurermPublicIpSku() *attrvalue.SimpleRule[string] {
	return attrvalue.NewSimpleRule[string](
		"azurerm_public_ip",
		"sku",
		[]string{"Standard"},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library/services/networking/public-ip/#pip-1---use-standard-sku-and-zone-redundant-ips-when-applicable",
	)
}

func AzurermPublicIpZones() *attrvalue.SetRule[int] {
	return attrvalue.NewSetRule(
		"azurerm_public_ip",
		"zones",
		[][]int{{1, 2, 3}},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library/services/networking/public-ip/#pip-1---use-standard-sku-and-zone-redundant-ips-when-applicable",
	)
}

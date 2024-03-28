package waf

import (
	"github.com/Azure/tflint-ruleset-avm/attrvalue"
)

func AzurermPublicIpSku() *attrvalue.SimpleRule[string] {
	return attrvalue.NewSimpleRule[string]("azurerm_public_ip", "sku", []string{"Standard"})
}

func AzurermPublicIpZones() *attrvalue.SetRule[int] {
	return attrvalue.NewListRule(
		"azurerm_public_ip",
		"zones",
		[][]int{{1, 2, 3}},
	)
}

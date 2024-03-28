package waf

import (
	"github.com/Azure/tflint-ruleset-avm/attrvalue"
)

func AzurermVirtualNetworkGatewaySku() *attrvalue.SimpleRule[string] {
	return attrvalue.NewSimpleRule[string](
		"azurerm_virtual_network_gateway",
		"sku",
		[]string{"ErGw1AZ", "ErGw2AZ", "ErGw3AZ", "VpnGw1AZ", "VpnGw2AZ", "VpnGw3AZ", "VpnGw4AZ", "VpnGw5AZ"},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library/services/networking/expressroute-gateway/#ergw-2---use-zone-redundant-gateway-skus",
	)
}

package waf

import (
	"github.com/Azure/tflint-ruleset-avm/attrvalue"
)

func (wf WafRules) AzurermVirtualNetworkGatewaySku() *attrvalue.SimpleRule[string] {
	return attrvalue.NewSimpleRule[string](
		"azurerm_virtual_network_gateway",
		"sku",
		[]string{"ErGw1AZ", "ErGw2AZ", "ErGw3AZ", "VpnGw1AZ", "VpnGw2AZ", "VpnGw3AZ", "VpnGw4AZ", "VpnGw5AZ"},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library/services/networking/expressroute-gateway/#ergw-2---use-zone-redundant-gateway-skus",
		false,
		"",
	)
}

func (wf WafRules) AzurermVirtualNetworkGatewayVpnActiveActive() *attrvalue.SimpleRule[bool] {
	return attrvalue.NewSimpleRule[bool](
		"azurerm_virtual_network_gateway",
		"active_active",
		[]bool{true},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library-v2/azure-resources/Network/virtualNetworkGateways/#plan-for-active-active-mode-with-vpn-gateways",
		"",
	)
}
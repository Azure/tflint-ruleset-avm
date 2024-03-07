package waf

import "github.com/Azure/tflint-ruleset-avm/attrvalue"

var AzurermVirtualNetworkGatewaySku = func() *attrvalue.SimpleRule {
	return attrvalue.NewSimpleStringRule(
		"azurerm_virtual_network_gateway",
		"sku",
		[]any{"ErGw1AZ", "ErGw2AZ", "ErGw3AZ", "VpnGw1AZ", "VpnGw2AZ", "VpnGw3AZ", "VpnGw4AZ", "VpnGw5AZ"},
	)
}

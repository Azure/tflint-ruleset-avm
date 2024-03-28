package waf

import (
	"github.com/Azure/tflint-ruleset-avm/attrvalue"
)

func AzurermVirtualNetworkGatewaySku() *attrvalue.SimpleRule[string] {
	return attrvalue.NewSimpleRule[string](
		"azurerm_virtual_network_gateway",
		"sku",
		[]string{"ErGw1AZ", "ErGw2AZ", "ErGw3AZ", "VpnGw1AZ", "VpnGw2AZ", "VpnGw3AZ", "VpnGw4AZ", "VpnGw5AZ"},
		"",
	)
}

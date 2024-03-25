package waf

import "github.com/Azure/tflint-ruleset-avm/attrvalue"

func AzurermVirtualMachineZoneNull() *attrvalue.NullRule {
	return attrvalue.NewNullRule(
		"azurerm_virtual_machine",
		"zone",
	)
}

func AzurermVirtualMachineZonesNull() *attrvalue.NullRule {
	return attrvalue.NewNullRule(
		"azurerm_virtual_machine",
		"zones",
	)
}

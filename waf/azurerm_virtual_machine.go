package waf

import "github.com/Azure/tflint-ruleset-avm/attrvalue"

func AzurermVirtualMachineZoneUnknown() *attrvalue.UnknownValueRule {
	return attrvalue.NewUnknownValueRule(
		"azurerm_virtual_machine",
		"zone",
	)
}

func AzurermVirtualMachineZonesUnknown() *attrvalue.UnknownValueRule {
	return attrvalue.NewUnknownValueRule(
		"azurerm_virtual_machine",
		"zones",
	)
}

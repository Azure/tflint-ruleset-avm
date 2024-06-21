package waf

import "github.com/Azure/tflint-ruleset-avm/attrvalue"

func (wf WafRules) AzurermVirtualMachineZoneUnknown() *attrvalue.UnknownValueRule {
	return attrvalue.NewUnknownValueRule(
		"azurerm_virtual_machine",
		"zone",
		"https://azure.github.io/Azure-Proactive-Resiliency-Library/services/compute/virtual-machines/#vm-2---deploy-vms-across-availability-zones",
		"",
	)
}

func (wf WafRules) AzurermVirtualMachineZonesUnknown() *attrvalue.UnknownValueRule {
	return attrvalue.NewUnknownValueRule(
		"azurerm_virtual_machine",
		"zones",
		"https://azure.github.io/Azure-Proactive-Resiliency-Library/services/compute/virtual-machines/#vm-2---deploy-vms-across-availability-zones",
		"",
	)
}

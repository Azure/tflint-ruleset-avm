package waf

import "github.com/Azure/tflint-ruleset-avm/attrvalue"

func (wf WafRules) AzurermVirtualMachineUseMaintenanceConfiguration1() *attrvalue.UnknownValueRule {
	return attrvalue.NewUnknownValueRule(
		"azurerm_maintenance_assignment_virtual_machine",
		"maintenance_configuration_id",
		"https://azure.github.io/Azure-Proactive-Resiliency-Library-v2/azure-resources/Compute/virtualMachines/#use-maintenance-configurations-for-the-vms",
		"",
	)
}

func (wf WafRules) AzurermVirtualMachineUseMaintenanceConfiguration2() *attrvalue.UnknownValueRule {
	return attrvalue.NewUnknownValueRule(
		"azurerm_maintenance_assignment_virtual_machine",
		"virtual_machine_id",
		"https://azure.github.io/Azure-Proactive-Resiliency-Library-v2/azure-resources/Compute/virtualMachines/#use-maintenance-configurations-for-the-vms",
		"",
	)
}
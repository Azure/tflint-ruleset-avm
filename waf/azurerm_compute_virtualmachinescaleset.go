package waf

import (
	"github.com/Azure/tflint-ruleset-avm/attrvalue"
)

func (wf WafRules) AutomaticRepairsPolicyEnabledLinux() *attrvalue.SimpleRule[bool] {
	return attrvalue.NewSimpleRule[bool](
		"azurerm_linux_virtual_machine_scale_set",
		"automatic_instance_repair ",
		[]bool{true},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library-v2/azure-resources/Compute/virtualMachineScaleSets/#enable-automatic-repair-policy-on-azure-virtual-machine-scale-sets",
		false,
		"",
	)
}

func (wf WafRules) AutomaticRepairsPolicyEnabledWindows() *attrvalue.SimpleRule[bool] {
	return attrvalue.NewSimpleRule[bool](
		"azurerm_windows_virtual_machine_scale_set",
		"automatic_instance_repair ",
		[]bool{true},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library-v2/azure-resources/Compute/virtualMachineScaleSets/#enable-automatic-repair-policy-on-azure-virtual-machine-scale-sets",
		false,
		"",
	)
}

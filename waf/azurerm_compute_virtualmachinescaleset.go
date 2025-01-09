package waf

import (
	"github.com/Azure/tflint-ruleset-avm/attrvalue"
)

func (wf WafRules) AutomaticRepairsPolicyEnabledLinux() *attrvalue.SimpleRule[bool] {
	return attrvalue.NewSimpleRule[bool](
		"azurerm_linux_virtual_machine_scale_set",
		"automatic_instance_repair",
		[]bool{true},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library-v2/azure-resources/Compute/virtualMachineScaleSets/#enable-automatic-repair-policy-on-azure-virtual-machine-scale-sets",
		false,
		"",
	)
}

func (wf WafRules) AutomaticRepairsPolicyEnabledWindows() *attrvalue.SimpleRule[bool] {
	return attrvalue.NewSimpleRule[bool](
		"azurerm_windows_virtual_machine_scale_set",
		"automatic_instance_repair",
		[]bool{true},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library-v2/azure-resources/Compute/virtualMachineScaleSets/#enable-automatic-repair-policy-on-azure-virtual-machine-scale-sets",
		false,
		"",
	)
}

func (wf WafRules) AutoScaleEnabled() *attrvalue.SimpleRule[bool] {
	return attrvalue.NewSimpleRule[bool](
		"azurerm_monitor_autoscale_setting",
		"enabled",
		[]bool{true},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library-v2/azure-resources/Compute/virtualMachineScaleSets/#configure-vmss-autoscale-to-custom-and-configure-the-scaling-metrics",
		false,
		"",
	)
}

func (wf WafRules) ZoneBalanceDisabledLinux() *attrvalue.SimpleRule[bool] {
	return attrvalue.NewSimpleRule[bool](
		"azurerm_linux_virtual_machine_scale_set",
		"zone_balance",
		[]bool{false},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library-v2/azure-resources/Compute/virtualMachineScaleSets/#disable-force-strictly-even-balance-across-zones-to-avoid-scale-in-and-out-fail-attempts",
		false,
		"",
	)
}

func (wf WafRules) ZoneBalanceDisabledWindows() *attrvalue.SimpleRule[bool] {
	return attrvalue.NewSimpleRule[bool](
		"azurerm_windows_virtual_machine_scale_set",
		"zone_balance",
		[]bool{false},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library-v2/azure-resources/Compute/virtualMachineScaleSets/#disable-force-strictly-even-balance-across-zones-to-avoid-scale-in-and-out-fail-attempts",
		false,
		"",
	)
}

func (wf WafRules) DeployAcrossAvailabilityZonesLinux() *attrvalue.SetRule[int] {
	return attrvalue.NewSetRule(
		"azurerm_linux_virtual_machine_scale_set",
		"zones",
		[][]int{{1, 2, 3}},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library-v2/azure-resources/Compute/virtualMachineScaleSets/#deploy-vmss-across-availability-zones-with-vmss-flex",
		"",
	)
}

func (wf WafRules) DeployAcrossAvailabilityZonesWindows() *attrvalue.SetRule[int] {
	return attrvalue.NewSetRule(
		"azurerm_windows_virtual_machine_scale_set",
		"zones",
		[][]int{{1, 2, 3}},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library-v2/azure-resources/Compute/virtualMachineScaleSets/#deploy-vmss-across-availability-zones-with-vmss-flex",
		"",
	)
}

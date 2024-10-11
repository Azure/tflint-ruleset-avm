package waf

import "github.com/Azure/tflint-ruleset-avm/attrvalue"

func (wf WafRules) AzurermVirtualMachineZoneUnknown() *attrvalue.UnknownValueRule {
	return attrvalue.NewUnknownValueRule(
		"azurerm_virtual_machine",
		"zone",
		"https://azure.github.io/Azure-Proactive-Resiliency-Library-v2/azure-resources/Compute/virtualMachines/#deploy-vms-across-availability-zones",
		"",
	)
}

func (wf WafRules) AzurermVirtualMachineZonesUnknown() *attrvalue.UnknownValueRule {
	return attrvalue.NewUnknownValueRule(
		"azurerm_virtual_machine",
		"zones",
		"https://azure.github.io/Azure-Proactive-Resiliency-Library-v2/azure-resources/Compute/virtualMachines/#deploy-vms-across-availability-zones",
		"",
	)
}

// This test checks for the use of resource type azurerm_virtual_machine since the azurerm_windows_virtual_machine and azurerm_linux_virtual_machine resources don't support unmanaged disks
// Since a test doesn't exist for checking the use of disallowed resource types, this uses an unknown value check on a required variable (name) to flag use of this resource.
func (wf WafRules) AzurermLegacyVirtualMachineNotAllowed() *attrvalue.UnknownValueRule {
	return attrvalue.NewUnknownValueRule(
		"azurerm_virtual_machine",
		"name",
		"https://azure.github.io/Azure-Proactive-Resiliency-Library-v2/azure-resources/Compute/virtualMachines/#use-managed-disks-for-vm-disks",
		"",
	)
}

// This test checks to see if a windows virtual machine's OS disk is one of the premium sku's
func (wf WafRules) AzurermWindowsVirtualMachineOSDiskDefaultSSD() *attrvalue.SimpleRule[string] {
	return attrvalue.NewSimpleNestedBlockRule[string](
		"azurerm_windows_virtual_machine",
		"os_disk",
		"storage_account_type",
		[]string{"Premium_LRS", "Premium_ZRS", "PremiumV2_LRS"},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library-v2/azure-resources/Compute/virtualMachines/#mission-critical-workloads-should-consider-using-premium-or-ultra-disks",
		true,
		"",
	)
}

// This test checks to see if a linux virtual machine's OS disk is one of the premium sku's
func (wf WafRules) AzurermLinuxVirtualMachineOSDiskDefaultSSD() *attrvalue.SimpleRule[string] {
	return attrvalue.NewSimpleNestedBlockRule[string](
		"azurerm_linux_virtual_machine",
		"os_disk",
		"storage_account_type",
		[]string{"Premium_LRS", "Premium_ZRS", "PremiumV2_LRS"},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library-v2/azure-resources/Compute/virtualMachines/#mission-critical-workloads-should-consider-using-premium-or-ultra-disks",
		true,
		"",
	)
}

// This test validates where managed disk resource types are either premium or ultra. TODO: Ensure that this doesn't conflict with other module outcomes.
func (wf WafRules) AzurermManagedDiskStorageAccountType() *attrvalue.SimpleRule[string] {
	return attrvalue.NewSimpleRule[string](
		"azurerm_managed_disk",
		"storage_account_type",
		[]string{"Premium_LRS", "Premium_ZRS", "PremiumV2_LRS", "UltraSSD_LRS"},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library-v2/azure-resources/Compute/virtualMachines/#mission-critical-workloads-should-consider-using-premium-or-ultra-disks",
		true,
		"",
	)
}

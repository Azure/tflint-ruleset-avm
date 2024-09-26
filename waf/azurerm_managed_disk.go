package waf

import (
	"github.com/Azure/tflint-ruleset-avm/attrvalue"
)

func (wf WafRules) AzurermManagedDiskStorageAccountTypeIsZRS() *attrvalue.SimpleRule[string] {
	return attrvalue.NewSimpleRule[string](
		"azurerm_amanaged_disk",
		"sku",
		[]string{"StandardSSD_ZRS", "Premium_ZRS"},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library-v2/welcome/",
		false,
		"",
	)
}
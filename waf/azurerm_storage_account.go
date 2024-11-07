package waf

import (
	"github.com/Azure/tflint-ruleset-avm/attrvalue"
)

func (wf WafRules) AzurermStorageAccountAccountReplicationType() *attrvalue.SimpleRule[string] {
	return attrvalue.NewSimpleRule[string](
		"azurerm_storage_account",
		"account_replication_type",
		[]string{"GRS"},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library/services/storage/storage-account/#st-1---ensure-that-storage-accounts-are-zone-or-region-redundant",
		false,
		"",
	)
}

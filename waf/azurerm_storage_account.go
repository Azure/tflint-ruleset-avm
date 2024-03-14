package waf

import (
	"github.com/Azure/tflint-ruleset-avm/attrvalue"
)

func AzurermStorageAccountAccountReplicationType() *attrvalue.SimpleRule[string] {
	return attrvalue.NewSimpleRule[string](
		"azurerm_storage_account",
		"account_replication_type",
		[]string{"ZRS"},
	)
}

package waf

import (
	"github.com/Azure/tflint-ruleset-avm/attrvalue"
)

// AzurermStorageAccountAccountReplicationType checks whether the account_replication_type is set to ZRS.
var AzurermStorageAccountAccountReplicationType = func() *attrvalue.SimpleRule {
	return attrvalue.NewSimpleStringRule(
		"azurerm_storage_account",
		"account_replication_type",
		[]any{"ZRS"},
	)
}

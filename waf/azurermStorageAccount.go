package waf

import (
	"github.com/Azure/tflint-ruleset-avm/attrvalue"
)

var AzurermStorageAccountAccountReplicationType = func() *attrvalue.SimpleRule {
	return attrvalue.NewSimpleStringRule(
		"azurerm_storage_account",
		"account_replication_type",
		[]any{"ZRS"},
	)
}

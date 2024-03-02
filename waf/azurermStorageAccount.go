package waf

import "github.com/Azure/tflint-ruleset-avm/attrvalue"

// AzurermStorageAccountAccountReplicationType checks whether the account_replication_type is set to ZRS.
var AzurermStorageAccountAccountReplicationType = func() *attrvalue.StringRule {
	return attrvalue.NewAttrStringValueRule(
		"azurerm_storage_account",
		"account_replication_type",
		[]string{"ZRS"},
	)
}

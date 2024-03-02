package waf

import "github.com/Azure/tflint-ruleset-avm/attrvalue"

var AzurermStorageAccountAccountReplicationType = func() *attrvalue.AttrStringValueRule {
	return attrvalue.NewAttrStringValueRule(
		"azurerm_key_vault",
		"account_replication_type",
		[]string{"ZRS"},
	)
}

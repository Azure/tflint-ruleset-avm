package waf

import (
	"reflect"

	"github.com/Azure/tflint-ruleset-avm/attrvalue"
	"github.com/zclconf/go-cty/cty"
)

// AzurermStorageAccountAccountReplicationType checks whether the account_replication_type is set to ZRS.
var AzurermStorageAccountAccountReplicationType = func() *attrvalue.SimpleRule {
	return attrvalue.NewSimpleRule(
		"azurerm_storage_account",
		"account_replication_type",
		cty.String,
		reflect.TypeOf(""),
		[]any{"ZRS"},
	)
}

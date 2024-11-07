package waf

import (
	"github.com/Azure/tflint-ruleset-avm/attrvalue"
)

func (wf WafRules) AzurermCosmosDbAccountBackupMode() *attrvalue.SimpleRule[string] {
	return attrvalue.NewSimpleNestedBlockRule[string](
		"azurerm_cosmosdb_account",
		"backup",
		"type",
		[]string{"Continuous"},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library-v2/azure-resources/DocumentDB/databaseAccounts/#configure-continuous-backup-mode",
		true,
		"",
	)
}

func (wf WafRules) AzurermCosmosDbAccountFailoverEnabled() *attrvalue.SimpleRule[bool] {
	return attrvalue.NewSimpleRule[bool](
		"azurerm_cosmosdb_account",
		"automatic_failover_enabled",
		[]bool{true},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library-v2/azure-resources/ContainerService/managedClusters/#enable-the-cluster-auto-scaler-on-an-existing-cluster",
		false,
		"",
	)
}


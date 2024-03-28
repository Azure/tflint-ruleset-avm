package waf

import "github.com/Azure/tflint-ruleset-avm/attrvalue"

func AzurermKubernetesClusterZones() *attrvalue.SetRule[int] {
	return attrvalue.NewListRule(
		"azurerm_kubernetes_cluster",
		"zones",
		[][]int{{1, 2, 3}},
	)
}

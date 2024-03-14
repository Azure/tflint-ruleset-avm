package waf

import "github.com/Azure/tflint-ruleset-avm/attrvalue"

func AzurermKubernetesClusterZones() *attrvalue.ListRule[int] {
	return attrvalue.NewListRule(
		"azurerm_kubernetes_cluster",
		"zones",
		[][]int{{1, 2, 3}},
	)
}

package waf

import "github.com/Azure/tflint-ruleset-avm/attrvalue"

func AzurermKubernetesClusterZones() *attrvalue.SetRule[int] {
	return attrvalue.NewSetRule(
		"azurerm_kubernetes_cluster",
		"zones",
		[][]int{{1, 2, 3}},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library/services/container/aks/#aks-1---deploy-aks-cluster-across-availability-zones",
	)
}

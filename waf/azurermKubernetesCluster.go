package waf

import "github.com/Azure/tflint-ruleset-avm/attrvalue"

var AzurermKubernetesClusterZones = func() *attrvalue.ListNumberRule {
	return attrvalue.NewListNumberRule(
		"azurerm_kubernetes_cluster",
		"zones",
		[][]int{{1, 2, 3}},
	)
}

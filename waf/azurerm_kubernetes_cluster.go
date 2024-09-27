package waf

import "github.com/Azure/tflint-ruleset-avm/attrvalue"

func (wf WafRules) AzurermKubernetesClusterZones() *attrvalue.SetRule[int] {
	return attrvalue.NewSetRule(
		"azurerm_kubernetes_cluster",
		"zones",
		[][]int{{1, 2, 3}},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library/services/container/aks/#aks-1---deploy-aks-cluster-across-availability-zones",
		"",
	)
}

func (wf WafRules) AzurermKubernetesClusterSkuTier() *attrvalue.SimpleRule[string] {
	return attrvalue.NewSimpleRule[string](
		"azurerm_kubernetes_cluster",
		"sku_tier",
		[]string{"Standard", "Premium"},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library-v2/azure-resources/ContainerService/managedClusters/#update-aks-tier-to-standard",
		false,
		"",
	)
}

func (wf WafRules) AzurermKubernetesClusterAutoScalingEnabled() *attrvalue.SimpleRule[bool] {
	return attrvalue.NewSimpleRule[bool](
		"azurerm_kubernetes_cluster",
		"auto_scaling_enabled",
		[]bool{true},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library-v2/azure-resources/ContainerService/managedClusters/#enable-the-cluster-auto-scaler-on-an-existing-cluster",
		false,
		"",
	)
}

func (wf WafRules) AzurermKubernetesClusterOMSAgentUnconfigured() *attrvalue.UnknownValueRule {
	return attrvalue.NewUnknownValueNestedBlockRule(
		"azurerm_kubernetes_cluster",
		"oms_agent",
		"log_analytics_workspace_id",
		"https://azure.github.io/Azure-Proactive-Resiliency-Library-v2/azure-resources/Compute/virtualMachines/#deploy-vms-across-availability-zones",
		"",
	)
}
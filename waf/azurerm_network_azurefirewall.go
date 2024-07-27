package waf

import (
	"github.com/Azure/tflint-ruleset-avm/attrvalue"
)

func (wf WafRules) DeployAcrossAvailabilityZones() *attrvalue.SetRule[int] {
	return attrvalue.NewSetRule(
		"azurerm_firewall",
		"zones",
		[][]int{{1, 2, 3}},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library-v2/azure-resources/Network/azureFirewalls/#deploy-azure-firewall-across-multiple-availability-zones",
		"",
	)
}

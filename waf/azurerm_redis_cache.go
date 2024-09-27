package waf

import (
	"github.com/Azure/tflint-ruleset-avm/attrvalue"
)

func (wf WafRules) AzurermRedisCacheZoneRedundancyEnabled() *attrvalue.SetRule[int] {
	return attrvalue.NewSetRule(
		"azurerm_redis_cache",
		"zones",
		[][]int{{1, 2, 3}},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library-v2/azure-resources/Cache/Redis/#enable-zone-redundancy-for-azure-cache-for-redis",
		"",
	)
}
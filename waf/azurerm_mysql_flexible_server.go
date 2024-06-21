package waf

import (
	"github.com/Azure/tflint-ruleset-avm/attrvalue"
)

func (wf WafRules) AzurermMySqlFlexibleServerZoneRedundancy() *attrvalue.SimpleRule[string] {
	return attrvalue.NewSimpleNestedBlockRule[string](
		"azurerm_mysql_flexible_server",
		"high_availability",
		"mode",
		[]string{"ZoneRedundant"},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library-v2/azure-resources/DBforMySQL/flexibleServers/#enable-ha-with-zone-redundancy",
		true,
		"",
	)
}

func (wf WafRules) AzurermMySqlFlexibleServerCustomMaintenanceSchedule() *attrvalue.SimpleRule[string] {
	return attrvalue.NewSimpleNestedBlockRule[string](
		"azurerm_mysql_flexible_server",
		"maintenance_window",
		"day_of_week",
		[]string{"0", "1", "2", "3", "4", "5", "6"},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library-v2/azure-resources/DBforMySQL/flexibleServers/#enable-custom-maintenance-schedule",
		true,
		"",
	)
}

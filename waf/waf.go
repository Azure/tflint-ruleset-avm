// Package waf contains the rules for Well Architected Alignment.
// To add a new rule, create a new file and add a new function that returns a new rule.
// Then add the rule to the Rules slice.
package waf

import "github.com/terraform-linters/tflint-plugin-sdk/tflint"

// Rules is the list of rules for Well Architected Alignment.
// Make sure to add any new rules to this list.
var Rules = []tflint.Rule{
	AzurermStorageAccountAccountReplicationType(),
	AzurermKubernetesClusterZones(),
	AzurermPublicIpSku(),
	AzurermPublicIpZones(),
	AzurermApplicationGatewayZones(),
	AzurermVirtualNetworkGatewaySku(),
	AzurermServicePlanZoneBalancingEnabled(),
	AzurermLbSku(),
}

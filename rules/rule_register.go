package rules

import (
	"slices"

	basic "github.com/Azure/tflint-ruleset-avm/basic/rules"
	"github.com/Azure/tflint-ruleset-avm/common"
	"github.com/Azure/tflint-ruleset-avm/interfaces"
	"github.com/Azure/tflint-ruleset-avm/outputs"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var Rules = func() []tflint.Rule {
	return registerRules(slices.Concat(
		basic.Rules,
		[]tflint.Rule{
			NewTerraformTfFileRule(),
			NewAzapiResourceTagRule(),
			NewAzapiReplaceTriggersRefsRule(),
			NewModuleSourceRule(),
			NewNoDoubleQuotesInIgnoreChangesRule(),
			NewDisallowedProviderRule("azurerm", "hashicorp/azurerm"),
			NewProviderVersionRule("modtm", "Azure/modtm", "0.3.0", "~> 0.3", true),
			NewProviderVersionRule("azapi", "Azure/azapi", "2.999.0", "~> 2.0", false),
			NewProviderVersionRule("azurerm", "hashicorp/azurerm", "4.999.0", "~> 4.0", false),
			// Generic azapi required attribute rules replacing bespoke implementations
			NewRequiredAttributeRule(
				"avm_azapi_response_export_values_required",
				"https://azure.github.io/Azure-Verified-Modules/specs/tf/azapi/#response_export_values-required",
				"resource",
				[]string{"azapi_resource", "azapi_resource_action", "azapi_update_resource"},
				"response_export_values",
				"[]",
				tflint.ERROR,
				DisallowWildcardList("response_export_values"),
			),
			NewRequiredAttributeRule(
				"avm_azapi_data_response_export_values_required",
				"https://azure.github.io/Azure-Verified-Modules/specs/tf/azapi/#response_export_values-required",
				"data",
				[]string{"azapi_resource", "azapi_resource_action", "azapi_resource_list"},
				"response_export_values",
				"[]",
				tflint.ERROR,
				DisallowWildcardList("response_export_values"),
			),
		},
		interfaces.Rules,
		outputs.Rules,
	))
}()

func registerRules(rawRules []tflint.Rule) []tflint.Rule {
	registered := make([]tflint.Rule, 0, len(rawRules))
	for _, rule := range rawRules {
		registered = append(registered, common.NewConfigurableRule(rule))
	}
	return registered
}

package rules

import (
	"fmt"
	"regexp"
	"slices"

	basic "github.com/Azure/tflint-ruleset-avm/basic/rules"
	"github.com/Azure/tflint-ruleset-avm/common"
	"github.com/Azure/tflint-ruleset-avm/interfaces"
	"github.com/Azure/tflint-ruleset-avm/outputs"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// RuleNameMigration records a breaking public rule rename.
type RuleNameMigration struct {
	Old string
	New string
}

// RuleNameMigrations is the complete migration from pre-v1 rule names.
var RuleNameMigrations = []RuleNameMigration{
	{Old: "azapi_data_response_export_values", New: "avm_azapi_data_response_export_values_required"},
	{Old: "azapi_replace_triggers_refs", New: "avm_azapi_replace_triggers_refs_valid"},
	{Old: "azapi_resource_tag", New: "avm_azapi_resource_tags_required"},
	{Old: "azapi_response_export_values", New: "avm_azapi_response_export_values_required"},
	{Old: "customer_managed_key", New: "avm_interface_customer_managed_key"},
	{Old: "deprecated_lock_interface", New: "avm_interface_lock_deprecated"},
	{Old: "deprecated_private_endpoints_interface", New: "avm_interface_private_endpoints_deprecated"},
	{Old: "deprecated_role_assignments_interface", New: "avm_interface_role_assignments_deprecated"},
	{Old: "diagnostic_settings", New: "avm_interface_diagnostic_settings"},
	{Old: "ignore_body_changes", New: "avm_interface_ignore_body_changes"},
	{Old: "location", New: "avm_interface_location"},
	{Old: "lock", New: "avm_interface_lock"},
	{Old: "managed_identities", New: "avm_interface_managed_identities"},
	{Old: "no_entire_resource_output_tffr2", New: "avm_output_entire_resource_disallowed"},
	{Old: "private_endpoints", New: "avm_interface_private_endpoints"},
	{Old: "private_endpoints_manage_dns_zone_group", New: "avm_interface_private_endpoints_manage_dns_zone_group"},
	{Old: "provider_azapi_version_constraint", New: "avm_provider_azapi_version_constraint"},
	{Old: "provider_azurerm_disallowed", New: "avm_provider_azurerm_disallowed"},
	{Old: "provider_azurerm_version_constraint", New: "avm_provider_azurerm_version_constraint"},
	{Old: "provider_modtm_version_constraint", New: "avm_provider_modtm_version_constraint"},
	{Old: "required_module_source_tffr1", New: "avm_terraform_module_source_required"},
	{Old: "required_module_source_tfnfr10", New: "avm_terraform_ignore_changes_unquoted_references"},
	{Old: "required_output_rmfr7", New: "avm_output_resource_id_required"},
	{Old: "resource_types", New: "avm_interface_resource_types"},
	{Old: "retry", New: "avm_interface_retry"},
	{Old: "role_assignments", New: "avm_interface_role_assignments"},
	{Old: "tags", New: "avm_interface_tags"},
	{Old: "terraform_heredoc_usage", New: "avm_terraform_literal_heredoc_disallowed"},
	{Old: "terraform_module_provider_declaration", New: "avm_terraform_provider_block_disallowed"},
	{Old: "terraform_sensitive_variable_no_default", New: "avm_terraform_sensitive_variable_default_disallowed"},
	{Old: "terraform_tf_file", New: "avm_terraform_configuration_file_required"},
	{Old: "timeouts", New: "avm_interface_timeouts"},
}

var canonicalRuleNamePattern = regexp.MustCompile(`^avm_[a-z0-9]+(?:_[a-z0-9]+)*$`)

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
				"azapi_response_export_values",
				"https://azure.github.io/Azure-Verified-Modules/specs/tf/azapi/#response_export_values-required",
				"resource",
				[]string{"azapi_resource", "azapi_resource_action", "azapi_update_resource"},
				"response_export_values",
				"[]",
				tflint.ERROR,
				DisallowWildcardList("response_export_values"),
			),
			NewRequiredAttributeRule(
				"azapi_data_response_export_values",
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
	migrations := make(map[string]string, len(RuleNameMigrations))
	for _, migration := range RuleNameMigrations {
		if _, exists := migrations[migration.Old]; exists {
			panic(fmt.Sprintf("duplicate legacy rule name migration: %s", migration.Old))
		}
		if !canonicalRuleNamePattern.MatchString(migration.New) {
			panic(fmt.Sprintf("invalid canonical rule name: %s", migration.New))
		}
		migrations[migration.Old] = migration.New
	}

	registered := make([]tflint.Rule, 0, len(rawRules))
	for _, rule := range rawRules {
		canonicalName, ok := migrations[rule.Name()]
		if !ok {
			panic(fmt.Sprintf("missing rule name migration for %s", rule.Name()))
		}
		delete(migrations, rule.Name())
		registered = append(registered, common.NewConfigurableRule(canonicalName, rule))
	}
	if len(migrations) != 0 {
		panic(fmt.Sprintf("rule name migrations without registered rules: %v", migrations))
	}

	return registered
}

# Rules Reference

This document lists all rules currently registered in this ruleset. AVM rules are enabled by default and support `error`, `warning`, or `notice` severity overrides in TFLint configuration.

| Name | Enabled | Severity | Link |
| ---- | ------- | -------- | ---- |
| avm_azapi_data_response_export_values_required | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/tf/a...](https://azure.github.io/Azure-Verified-Modules/specs/tf/azapi/#response_export_values-required) |
| avm_azapi_replace_triggers_refs_valid | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/spec/TFFR5/](https://azure.github.io/Azure-Verified-Modules/spec/TFFR5/) |
| avm_azapi_resource_tags_required | true | ERROR | [https://aka.ms/avm/spec/TFFR9](https://aka.ms/avm/spec/TFFR9) |
| avm_azapi_response_export_values_required | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/tf/a...](https://azure.github.io/Azure-Verified-Modules/specs/tf/azapi/#response_export_values-required) |
| avm_interface_customer_managed_key | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/tf/i...](https://azure.github.io/Azure-Verified-Modules/specs/tf/interfaces/#customer-managed-keys) |
| avm_interface_diagnostic_settings | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/tf/i...](https://azure.github.io/Azure-Verified-Modules/specs/tf/interfaces/#diagnostic-settings) |
| avm_interface_ignore_body_changes | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/spec/TFFR8/](https://azure.github.io/Azure-Verified-Modules/spec/TFFR8/) |
| avm_interface_location | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/tf/r...](https://azure.github.io/Azure-Verified-Modules/specs/tf/res/#id-rmnfr2---category-inputs---parametervariable-naming) |
| avm_interface_lock | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/tf/i...](https://azure.github.io/Azure-Verified-Modules/specs/tf/interfaces/#resource-locks) |
| avm_interface_lock_deprecated | true | NOTICE | [https://azure.github.io/Azure-Verified-Modules/specs/tf/i...](https://azure.github.io/Azure-Verified-Modules/specs/tf/interfaces/#resource-locks) |
| avm_interface_managed_identities | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/tf/i...](https://azure.github.io/Azure-Verified-Modules/specs/tf/interfaces/#managed-identities) |
| avm_interface_private_endpoints | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/tf/i...](https://azure.github.io/Azure-Verified-Modules/specs/tf/interfaces/#private-endpoints) |
| avm_interface_private_endpoints_deprecated | true | NOTICE | [https://azure.github.io/Azure-Verified-Modules/specs/tf/i...](https://azure.github.io/Azure-Verified-Modules/specs/tf/interfaces/#private-endpoints) |
| avm_interface_private_endpoints_manage_dns_zone_group | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/includes/i...](https://azure.github.io/Azure-Verified-Modules/includes/interfaces/tf/int.pe.schema.tf) |
| avm_interface_resource_tags | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/tf/i...](https://azure.github.io/Azure-Verified-Modules/specs/tf/interfaces/#tags) |
| avm_interface_resource_types | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/spec/TFFR6/](https://azure.github.io/Azure-Verified-Modules/spec/TFFR6/) |
| avm_interface_retry | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/spec/TFFR7/](https://azure.github.io/Azure-Verified-Modules/spec/TFFR7/) |
| avm_interface_role_assignments | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/tf/i...](https://azure.github.io/Azure-Verified-Modules/specs/tf/interfaces/#role-assignments) |
| avm_interface_role_assignments_deprecated | true | NOTICE | [https://azure.github.io/Azure-Verified-Modules/specs/tf/i...](https://azure.github.io/Azure-Verified-Modules/specs/tf/interfaces/#role-assignments) |
| avm_interface_tags | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/tf/i...](https://azure.github.io/Azure-Verified-Modules/specs/tf/interfaces/#tags) |
| avm_interface_timeouts | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/spec/TFFR7/](https://azure.github.io/Azure-Verified-Modules/spec/TFFR7/) |
| avm_output_entire_resource_disallowed | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/tf/r...](https://azure.github.io/Azure-Verified-Modules/specs/tf/res/#id-tffr2---category-outputs---additional-terraform-outputs) |
| avm_output_resource_id_required | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/shar...](https://azure.github.io/Azure-Verified-Modules/specs/shared/#id-rmfr7---category-outputs---minimum-required-outputs) |
| avm_provider_azapi_version_constraint | true | ERROR | - |
| avm_provider_azurerm_disallowed | true | ERROR | - |
| avm_provider_azurerm_version_constraint | true | ERROR | - |
| avm_provider_modtm_version_constraint | true | ERROR | - |
| avm_terraform_configuration_file_required | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/spec/TFNFR39/](https://azure.github.io/Azure-Verified-Modules/spec/TFNFR39/) |
| avm_terraform_ignore_changes_unquoted_references | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/terr...](https://azure.github.io/Azure-Verified-Modules/specs/terraform/#id-tfnfr10---category-code-style---no-double-quotes-in-ignore_changes) |
| avm_terraform_literal_heredoc_disallowed | true | NOTICE | [https://aka.ms/avm/spec/TFNFR40](https://aka.ms/avm/spec/TFNFR40) |
| avm_terraform_module_source_required | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/terr...](https://azure.github.io/Azure-Verified-Modules/specs/terraform/#id-tffr1---category-composition---cross-referencing-modules) |
| avm_terraform_provider_block_disallowed | true | WARNING | [https://azure.github.io/Azure-Verified-Modules/spec/TFNFR27/](https://azure.github.io/Azure-Verified-Modules/spec/TFNFR27/) |
| avm_terraform_sensitive_variable_default_disallowed | true | WARNING | [https://azure.github.io/Azure-Verified-Modules/spec/TFNFR23/](https://azure.github.io/Azure-Verified-Modules/spec/TFNFR23/) |

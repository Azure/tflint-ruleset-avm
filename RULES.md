# Rules Reference

This document lists all rules currently registered in this ruleset. AVM rules are enabled by default and can be disabled explicitly in TFLint configuration.

| Name | Enabled | Severity | Link |
| ---- | ------- | -------- | ---- |
| azapi_data_response_export_values | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/tf/a...](https://azure.github.io/Azure-Verified-Modules/specs/tf/azapi/#response_export_values-required) |
| azapi_replace_triggers_refs | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/spec/TFFR5/](https://azure.github.io/Azure-Verified-Modules/spec/TFFR5/) |
| azapi_resource_tag | true | ERROR | [https://aka.ms/avm/spec/TFFR9](https://aka.ms/avm/spec/TFFR9) |
| azapi_response_export_values | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/tf/a...](https://azure.github.io/Azure-Verified-Modules/specs/tf/azapi/#response_export_values-required) |
| customer_managed_key | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/tf/i...](https://azure.github.io/Azure-Verified-Modules/specs/tf/interfaces/#customer-managed-keys) |
| deprecated_lock_interface | true | NOTICE | [https://azure.github.io/Azure-Verified-Modules/specs/tf/i...](https://azure.github.io/Azure-Verified-Modules/specs/tf/interfaces/#resource-locks) |
| deprecated_private_endpoints_interface | true | NOTICE | [https://azure.github.io/Azure-Verified-Modules/specs/tf/i...](https://azure.github.io/Azure-Verified-Modules/specs/tf/interfaces/#private-endpoints) |
| deprecated_role_assignments_interface | true | NOTICE | [https://azure.github.io/Azure-Verified-Modules/specs/tf/i...](https://azure.github.io/Azure-Verified-Modules/specs/tf/interfaces/#role-assignments) |
| diagnostic_settings | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/tf/i...](https://azure.github.io/Azure-Verified-Modules/specs/tf/interfaces/#diagnostic-settings) |
| ignore_body_changes | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/spec/TFFR8/](https://azure.github.io/Azure-Verified-Modules/spec/TFFR8/) |
| location | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/tf/r...](https://azure.github.io/Azure-Verified-Modules/specs/tf/res/#id-rmnfr2---category-inputs---parametervariable-naming) |
| lock | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/tf/i...](https://azure.github.io/Azure-Verified-Modules/specs/tf/interfaces/#resource-locks) |
| managed_identities | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/tf/i...](https://azure.github.io/Azure-Verified-Modules/specs/tf/interfaces/#managed-identities) |
| no_entire_resource_output_tffr2 | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/tf/r...](https://azure.github.io/Azure-Verified-Modules/specs/tf/res/#id-tffr2---category-outputs---additional-terraform-outputs) |
| private_endpoints | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/tf/i...](https://azure.github.io/Azure-Verified-Modules/specs/tf/interfaces/#private-endpoints) |
| private_endpoints_manage_dns_zone_group | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/includes/i...](https://azure.github.io/Azure-Verified-Modules/includes/interfaces/tf/int.pe.schema.tf) |
| provider_azapi_version_constraint | true | ERROR | - |
| provider_azurerm_disallowed | true | ERROR | - |
| provider_azurerm_version_constraint | true | ERROR | - |
| provider_modtm_version_constraint | true | ERROR | - |
| required_module_source_tffr1 | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/terr...](https://azure.github.io/Azure-Verified-Modules/specs/terraform/#id-tffr1---category-composition---cross-referencing-modules) |
| required_module_source_tfnfr10 | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/terr...](https://azure.github.io/Azure-Verified-Modules/specs/terraform/#id-tfnfr10---category-code-style---no-double-quotes-in-ignore_changes) |
| required_output_rmfr7 | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/shar...](https://azure.github.io/Azure-Verified-Modules/specs/shared/#id-rmfr7---category-outputs---minimum-required-outputs) |
| resource_types | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/spec/TFFR6/](https://azure.github.io/Azure-Verified-Modules/spec/TFFR6/) |
| retry | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/spec/TFFR7/](https://azure.github.io/Azure-Verified-Modules/spec/TFFR7/) |
| role_assignments | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/tf/i...](https://azure.github.io/Azure-Verified-Modules/specs/tf/interfaces/#role-assignments) |
| tags | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/tf/i...](https://azure.github.io/Azure-Verified-Modules/specs/tf/interfaces/#tags) |
| terraform_heredoc_usage | true | NOTICE | [https://aka.ms/avm/spec/TFNFR40](https://aka.ms/avm/spec/TFNFR40) |
| terraform_module_provider_declaration | true | WARNING | [https://azure.github.io/Azure-Verified-Modules/spec/TFNFR27/](https://azure.github.io/Azure-Verified-Modules/spec/TFNFR27/) |
| terraform_sensitive_variable_no_default | true | WARNING | [https://azure.github.io/Azure-Verified-Modules/spec/TFNFR23/](https://azure.github.io/Azure-Verified-Modules/spec/TFNFR23/) |
| terraform_tf_file | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/spec/TFNFR39/](https://azure.github.io/Azure-Verified-Modules/spec/TFNFR39/) |
| timeouts | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/spec/TFFR7/](https://azure.github.io/Azure-Verified-Modules/spec/TFFR7/) |

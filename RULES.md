# Rules Reference

This document lists all rules currently registered in this ruleset. The Enabled column reflects the default state (some consolidated rules are wrapped to be disabled by default).

| Name | Enabled | Severity | Link |
| ---- | ------- | -------- | ---- |
| azapi_data_response_export_values | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/tf/a...](https://azure.github.io/Azure-Verified-Modules/specs/tf/azapi/#response_export_values-required) |
| azapi_replace_triggers_refs | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/tf/a...](https://azure.github.io/Azure-Verified-Modules/specs/tf/azapi/#replace_triggers_refs) |
| azapi_response_export_values | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/tf/a...](https://azure.github.io/Azure-Verified-Modules/specs/tf/azapi/#response_export_values-required) |
| azurerm_arg_order | false | NOTICE | [https://github.com/Azure/tflint-ruleset-avm/blob/main/doc...](https://github.com/Azure/tflint-ruleset-avm/blob/main/docs/azurerm/azurerm_arg_order.md) |
| azurerm_resource_tag | false | NOTICE | [https://github.com/Azure/tflint-ruleset-avm/blob/main/doc...](https://github.com/Azure/tflint-ruleset-avm/blob/main/docs/azurerm/azurerm_resource_tag.md) |
| customer_managed_key | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/tf/i...](https://azure.github.io/Azure-Verified-Modules/specs/tf/interfaces/#customer-managed-keys) |
| diagnostic_settings | true | ERROR | - |
| location | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/tf/r...](https://azure.github.io/Azure-Verified-Modules/specs/tf/res/#id-rmnfr2---category-inputs---parametervariable-naming) |
| lock | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/tf/i...](https://azure.github.io/Azure-Verified-Modules/specs/tf/interfaces/#resource-locks) |
| managed_identities | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/tf/i...](https://azure.github.io/Azure-Verified-Modules/specs/tf/interfaces/#managed-identities) |
| no_entire_resource_output_tffr2 | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/tf/r...](https://azure.github.io/Azure-Verified-Modules/specs/tf/res/#id-tffr2---category-outputs---additional-terraform-outputs) |
| private_endpoints | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/tf/i...](https://azure.github.io/Azure-Verified-Modules/specs/tf/interfaces/#private-endpoints) |
| provider_azapi_version_constraint | true | ERROR | - |
| provider_azurerm_disallowed | true | ERROR | - |
| provider_azurerm_version_constraint | true | ERROR | - |
| provider_modtm_version_constraint | true | ERROR | - |
| required_module_source_tffr1 | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/terr...](https://azure.github.io/Azure-Verified-Modules/specs/terraform/#id-tffr1---category-composition---cross-referencing-modules) |
| required_module_source_tfnfr10 | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/terr...](https://azure.github.io/Azure-Verified-Modules/specs/terraform/#id-tfnfr10---category-code-style---no-double-quotes-in-ignore_changes) |
| required_output_rmfr7 | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/shar...](https://azure.github.io/Azure-Verified-Modules/specs/shared/#id-rmfr7---category-outputs---minimum-required-outputs) |
| role_assignments | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/tf/i...](https://azure.github.io/Azure-Verified-Modules/specs/tf/interfaces/#role-assignments) |
| tags | true | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/tf/i...](https://azure.github.io/Azure-Verified-Modules/specs/tf/interfaces/#tags) |
| terraform_count_index_usage | false | WARNING | - |
| terraform_heredoc_usage | false | NOTICE | - |
| terraform_locals_order | false | NOTICE | [https://github.com/Azure/tflint-ruleset-avm/blob/main/doc...](https://github.com/Azure/tflint-ruleset-avm/blob/main/docs/basic/terraform_locals_order.md) |
| terraform_module_provider_declaration | false | WARNING | - |
| terraform_output_order | false | NOTICE | [https://github.com/Azure/tflint-ruleset-avm/blob/main/doc...](https://github.com/Azure/tflint-ruleset-avm/blob/main/docs/basic/terraform_output_order.md) |
| terraform_output_separate | false | NOTICE | [https://github.com/Azure/tflint-ruleset-avm/blob/main/doc...](https://github.com/Azure/tflint-ruleset-avm/blob/main/docs/basic/terraform_output_separate.md) |
| terraform_required_providers_declaration | false | NOTICE | - |
| terraform_required_version_declaration | false | NOTICE | - |
| terraform_resource_data_arg_layout | false | NOTICE | [https://github.com/Azure/tflint-ruleset-avm/blob/main/doc...](https://github.com/Azure/tflint-ruleset-avm/blob/main/docs/basic/terraform_resource_data_arg_layout.md) |
| terraform_sensitive_variable_no_default | false | WARNING | - |
| terraform_variable_nullable_false | false | NOTICE | [https://github.com/Azure/tflint-ruleset-avm/blob/main/doc...](https://github.com/Azure/tflint-ruleset-avm/blob/main/docs/basic/terraform_variable_nullable_false.md) |
| terraform_variable_order | false | NOTICE | [https://github.com/Azure/tflint-ruleset-avm/blob/main/doc...](https://github.com/Azure/tflint-ruleset-avm/blob/main/docs/basic/terraform_variable_order.md) |
| terraform_variable_separate | false | NOTICE | [https://github.com/Azure/tflint-ruleset-avm/blob/main/doc...](https://github.com/Azure/tflint-ruleset-avm/blob/main/docs/basic/terraform_variable_separate.md) |
| terraform_versions_file | false | NOTICE | - |
| tfnfr26 | false | ERROR | [https://azure.github.io/Azure-Verified-Modules/specs/terr...](https://azure.github.io/Azure-Verified-Modules/specs/terraform/#id-tfnfr26---category-code-style---providers-must-be-declared-in-the-required_providers-block-in-terraformtf-and-must-have-a-constraint-on-minimum-and-maximum-major-version) |

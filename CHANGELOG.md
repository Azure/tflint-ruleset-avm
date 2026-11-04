# Changelog

## Unreleased

### Fixed

- Reject AzureRM constraints that allow 5.0.0 or later, including ranges with individual version exclusions. AzureRM remains optional and the existing 4.x compatibility check is unchanged.

### Breaking

- Renamed every rule to the canonical `avm_` lowercase snake_case scheme. Existing rule override blocks must use the new names listed in the README migration table.

### Added

- Added rules for AVM Terraform specs TFNFR14, TFNFR16, TFNFR21, TFNFR22, TFNFR24, TFNFR25, TFNFR30, and TFNFR31: `avm_terraform_module_control_variable_disallowed`, `avm_terraform_variable_negative_toggle_disallowed`, `avm_terraform_variable_nullable_true_disallowed`, `avm_terraform_sensitive_false_disallowed`, `avm_terraform_deprecated_variable_file_required`, `avm_terraform_required_version_constraint`, `avm_terraform_deprecated_output_file_required`, and `avm_terraform_locals_file_non_locals_disallowed`. All are enabled by default. `avm_terraform_variable_negative_toggle_disallowed` uses `NOTICE` severity; the others use `ERROR`.
- Added resource-module rules for AVM specs TFRMFR1, TFRMNFR2, and RMFR8: `avm_terraform_parent_id_variable_required`, `avm_azapi_primary_resource_naming_valid`, and `avm_terraform_pattern_module_source_disallowed`. All are enabled by default with `ERROR` severity, and support `module_class` so pattern and utility roots can skip them.
- Added rules for AVM Terraform specs TFNFR38, TFNFR21 (with TFNFR20), and TFNFR7: `avm_terraform_resource_id_variable_validation_required`, `avm_terraform_collection_variable_nullable_false_required`, and `avm_terraform_count_for_each_usage_valid`. All are enabled by default with `ERROR` severity. `avm_terraform_resource_id_variable_validation_required` runs only when the module declares the `azapi` provider, and `avm_terraform_collection_variable_nullable_false_required` skips the `tags` variable.
- Added per-rule `severity` configuration. Supported values are `error`, `warning`, and `notice`; omitting the input preserves the rule's default severity.

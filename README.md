# Azure Verified Modules TFLint Ruleset

[![Build Status](https://github.com/Azure/tflint-ruleset-avm/workflows/build/badge.svg?branch=main)](https://github.com/Azure/tflint-ruleset-avm/actions)

This repository contains the TFLint ruleset for Azure Verified Modules.

## Embedded AzAPI tag capabilities

The `avm_azapi_resource_tags_required` rule uses a standalone, embedded snapshot at
`internal/tagcapability/data/azapi_tags_v1.json`; it does not download schemas
at runtime or depend on the AzAPI provider. The snapshot records its
MIT-licensed upstream sources and exact revisions. The weekly/manual
[`Update AzAPI Tag Snapshot`](.github/workflows/update-azapi-tag-snapshot.yml)
workflow regenerates it from `azure-rest-api-specs` through
`bicep-types-az`, validates it, and opens or updates a deterministic PR only
when the data changes.

For writable resource types, the rule requires a consumer-settable `tags`
argument but permits any Terraform expression so modules can implement the
standard module-wide fallback and per-resource override contract. Resource
types with read-only or unsupported tags must omit the argument.

When a module exposes per-resource or submodule tag overrides, the
`avm_interface_resource_tags` rule validates the typed recursive
`resource_tags` interface. The interface separates resource block labels under
`resources` from child module labels under `modules`; omitted or null resource
maps inherit `var.tags`, while supplied maps replace it without merging, and
`{}` deliberately removes all tags.

## Requirements

- TFLint v0.62+
- Go v1.26+

## Installation

You can install a released plugin with `tflint --init`. Declare a config in `.tflint.hcl`, replacing `X.Y.Z` with the release version without the `v` prefix:

```hcl
plugin "avm" {
  enabled = true

  version   = "X.Y.Z"
  source    = "github.com/Azure/tflint-ruleset-avm"
  signature = "attestation"
}
```

## Rules

Every rule name uses `avm_<domain>_<subject>_<requirement>`-style lowercase snake_case. Rules are enabled by default and preserve their documented severity unless overridden:

```hcl
rule "avm_interface_resource_types" {
  enabled  = true
  severity = "notice"
}
```

The supported severity values are `error`, `warning`, and `notice`. Any other value is rejected as invalid configuration.

### Resource exclusions

Every AVM rule accepts an optional `ignored_resources` array. Each selector excludes only the matching static resource or data block from that rule's resource checks; it does not suppress unrelated provider, variable, module, or file diagnostics.

Use Terraform-native managed-resource addresses by default, with `data.` for data sources:

```hcl
rule "avm_interface_resource_types" {
  enabled = true
  ignored_resources = [
    "azapi_resource.telemetry",
    "data.azapi_resource.telemetry",
  ]
}
```

`resource.azapi_resource.telemetry` is accepted as an equivalent explicit managed-resource form. Selectors must include the namespace (when using `data.`), type, and block name. They apply to static block addresses only: `count` and `for_each` blocks are evaluated as declared rather than expanded from Terraform state. Rules without this option preserve their existing behavior.

### AzureRM version constraint

`avm_provider_azurerm_version_constraint` requires declared AzureRM constraints to include `4.999.0` and exclude `5.0.0` and later. Use `~> 4.0` or a compatible range such as `>= 4.2, < 5.0`. Open-ended ranges and ranges ending at `< 6.0` fail even with individual 5.x versions excluded. AzureRM remains optional for AzAPI-only modules; source checks and the existing rejection of narrow 4.x pins are unchanged.

This upper-bound check requires a plugin release containing the fix; v1.2.0 and earlier do not enforce it. ModTM and AzAPI version policies are unchanged.

### Resource ID output rule

`avm_output_resource_id_required` applies to resource-module roots, not pattern or utility roots. It defaults to `resource` when `module_class` is omitted, including when `metadata.json` is missing or malformed. Configure non-resource roots explicitly in `.tflint.hcl`:

```hcl
rule "avm_output_resource_id_required" {
  enabled      = true
  module_class = "pattern"
}
```

Use `module_class = "utility"` for a utility root or `module_class = "resource"` to state the default explicitly. Other values are rejected. This option affects only this rule; child modules are still ignored and severity overrides continue to work. Module class is not inferred from the process directory, repository name, or optional metadata. Set the appropriate rule configuration for each linted root, including when using `--chdir` or `--recursive`.

The `module_class` option requires a plugin release containing this change. Existing releases such as v1.0.0 and v1.1.0 reject it; do not add it to configurations pinned to those versions.

### Module source rule

`avm_terraform_module_source_required` defaults to `module_scope = "root"`, allowing direct local submodules through `./modules/<name>`. When linting a direct `modules/*` scope independently, configure:

```hcl
rule "avm_terraform_module_source_required" {
  enabled      = true
  module_scope = "submodule"
}
```

This permits direct sibling sources such as `../route-table-route`, as specified by [TFRMNFR1](https://azure.github.io/Azure-Verified-Modules/spec/TFRMNFR1/#module-source-references). Local paths must contain exactly one directory name after the scope's prefix; a single trailing `/` is allowed. Nested roots, extra traversal, backslash separators, drive-qualified names, and names ending in a dot or space are rejected. Root scopes still reject sibling paths, and submodule scopes reject `./modules/<name>`. Other `module_scope` values are configuration errors; severity overrides remain supported.

Set the scope for each lint target, including with `--chdir` or `--recursive`; it is not inferred from the working directory or TFLint's module call tree. The caller is responsible for selecting only direct `modules/*` roots. Existing external AVM Registry references, including Registry package submodules, remain accepted. This rule has no repository identity with which to distinguish an external Registry submodule from a same-repository reference. Example configurations that disable this rule are unchanged.

The `module_scope` option requires a published plugin release containing this fix. Releases through v1.3.0 do not support it; do not add it to configurations pinned to those releases.

<!-- RULES:START -->
<!-- This section is auto-generated by scripts/update-readme.sh. Do not edit manually. -->

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
<!-- RULES:END -->

See [`RULES.md`](./RULES.md) for the full rules reference (same content, kept in sync automatically).

### Rule name migration

This release removes the old names without compatibility aliases.

| Old name | New name |
| ---- | ---- |
| `azapi_data_response_export_values` | `avm_azapi_data_response_export_values_required` |
| `azapi_replace_triggers_refs` | `avm_azapi_replace_triggers_refs_valid` |
| `azapi_resource_tag` | `avm_azapi_resource_tags_required` |
| `azapi_response_export_values` | `avm_azapi_response_export_values_required` |
| `customer_managed_key` | `avm_interface_customer_managed_key` |
| `deprecated_lock_interface` | `avm_interface_lock_deprecated` |
| `deprecated_private_endpoints_interface` | `avm_interface_private_endpoints_deprecated` |
| `deprecated_role_assignments_interface` | `avm_interface_role_assignments_deprecated` |
| `diagnostic_settings` | `avm_interface_diagnostic_settings` |
| `ignore_body_changes` | `avm_interface_ignore_body_changes` |
| `location` | `avm_interface_location` |
| `lock` | `avm_interface_lock` |
| `managed_identities` | `avm_interface_managed_identities` |
| `no_entire_resource_output_tffr2` | `avm_output_entire_resource_disallowed` |
| `private_endpoints` | `avm_interface_private_endpoints` |
| `private_endpoints_manage_dns_zone_group` | `avm_interface_private_endpoints_manage_dns_zone_group` |
| `provider_azapi_version_constraint` | `avm_provider_azapi_version_constraint` |
| `provider_azurerm_disallowed` | `avm_provider_azurerm_disallowed` |
| `provider_azurerm_version_constraint` | `avm_provider_azurerm_version_constraint` |
| `provider_modtm_version_constraint` | `avm_provider_modtm_version_constraint` |
| `required_module_source_tffr1` | `avm_terraform_module_source_required` |
| `required_module_source_tfnfr10` | `avm_terraform_ignore_changes_unquoted_references` |
| `required_output_rmfr7` | `avm_output_resource_id_required` |
| `resource_types` | `avm_interface_resource_types` |
| `retry` | `avm_interface_retry` |
| `role_assignments` | `avm_interface_role_assignments` |
| `tags` | `avm_interface_tags` |
| `terraform_heredoc_usage` | `avm_terraform_literal_heredoc_disallowed` |
| `terraform_module_provider_declaration` | `avm_terraform_provider_block_disallowed` |
| `terraform_sensitive_variable_no_default` | `avm_terraform_sensitive_variable_default_disallowed` |
| `terraform_tf_file` | `avm_terraform_configuration_file_required` |
| `timeouts` | `avm_interface_timeouts` |

### Regenerating the Rules Documentation

The full list of rules is generated automatically into `RULES.md` and mirrored into the `## Rules` section above.

To regenerate after adding, removing, or modifying rules, run:

```bash
./scripts/update-readme.sh
```

This script:

1. Runs `go generate ./...`, which executes `cmd/rulesdoc` to refresh `RULES.md` from `rules.Rules`.
2. Splices the rules table from `RULES.md` into the `<!-- RULES:START -->` / `<!-- RULES:END -->` markers in `README.md`.

A GitHub Actions workflow (`.github/workflows/rules-docs-check.yml`) runs the same script on pull requests that touch rule sources and fails the build if `README.md` or `RULES.md` is out of date, so drift cannot land on `main`.

## Building the plugin

Clone the repository locally and run the following command:

```bash
make
```

You can easily install the built plugin with the following:

```bash
make install
```

You can run the built plugin like the following:

```bash
$ cat << EOS > .tflint.hcl
plugin "avm" {
  enabled = true
}
EOS
$ tflint
```

## Releasing

1. Create and publish a GitHub pre-release with a `vX.Y.Z` tag and release notes.
2. Queue the `tflint-ruleset-avm` release pipeline in `github-private/azure/Azure-Verified-Modules` with that tag.
3. The pipeline builds every target in `release/targets.json`, Authenticode-signs the Windows binaries, packages and validates the assets, replaces the pre-release assets, and promotes the pre-release.
4. Promotion triggers `sign-checksums.yml`, which downloads the published release's assets by ID, validates them against the release tag, and creates a GitHub Artifact Attestation for `checksums.txt`.

Release pipelines are manually queued and never run from a GitHub push or tag.

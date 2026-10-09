# Changelog

## Unreleased

### Fixed

- Reject AzureRM constraints that allow 5.0.0 or later, including ranges with individual version exclusions. AzureRM remains optional and the existing 4.x compatibility check is unchanged.

### Breaking

- Renamed every rule to the canonical `avm_` lowercase snake_case scheme. Existing rule override blocks must use the new names listed in the README migration table.

### Added

- Added per-rule `severity` configuration. Supported values are `error`, `warning`, and `notice`; omitting the input preserves the rule's default severity.

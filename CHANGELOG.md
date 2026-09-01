# Changelog

## Unreleased

### Breaking

- Renamed every rule to the canonical `avm_` lowercase snake_case scheme. Existing rule override blocks must use the new names listed in the README migration table.

### Added

- Added per-rule `severity` configuration. Supported values are `error`, `warning`, and `notice`; omitting the input preserves the rule's default severity.

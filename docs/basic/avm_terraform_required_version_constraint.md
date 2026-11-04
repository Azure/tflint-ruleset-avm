# avm_terraform_required_version_constraint

checks that `required_version` in `terraform.tf` sets a minimum version and a maximum major version

## Example

```hcl
# terraform.tf
terraform {
  required_version = ">= 1.9"
}
```

```
$ tflint
1 issue(s) found:

Error: `required_version` must exclude Terraform 2.0.0 and later, got `>= 1.9`; use `~> #.#` or `>= #.#.#, < #.#.#` (avm_terraform_required_version_constraint)

  on terraform.tf line 3:
   3:   required_version = ">= 1.9"

Reference: https://azure.github.io/Azure-Verified-Modules/spec/TFNFR25/

```

## Why
[TFNFR25](https://azure.github.io/Azure-Verified-Modules/spec/TFNFR25/) says `required_version` MUST include a minimum Terraform CLI version and a maximum major version. Major releases can introduce breaking changes and must be tested.

The rule only checks `terraform` blocks in `terraform.tf`. It reports:

- A missing `required_version`.
- A constraint that can't be parsed.
- A constraint with no minimum version.
- A constraint that allows the next major version after its minimum.

Other parts of TFNFR25 are covered elsewhere. `avm_terraform_configuration_file_required` checks for exactly one `terraform` block. The standard Terraform TFLint plugin checks that `required_version` comes first.

## How To Fix
Use the `~> #.#` or `>= #.#.#, < #.#.#` format, for example `required_version = "~> 1.9"`.

# avm_terraform_deprecated_variable_file_required

checks that deprecated variables are declared in `deprecated_variables.tf`

## Example

```hcl
# variables.tf
variable "create_nsg" {
  type        = bool
  description = "DEPRECATED, use `network_security_group_enabled` instead; Whether to generate a network security group."
}
```

```
$ tflint
1 issue(s) found:

Error: Deprecated variable `create_nsg` must be declared in `deprecated_variables.tf` (avm_terraform_deprecated_variable_file_required)

  on variables.tf line 2:
   2: variable "create_nsg" {

Reference: https://azure.github.io/Azure-Verified-Modules/spec/TFNFR24/

```

## Why
[TFNFR24](https://azure.github.io/Azure-Verified-Modules/spec/TFNFR24/) says a deprecated variable MUST move to an independent `deprecated_variables.tf` file. Its `description` MUST start with `DEPRECATED`.

The rule checks both directions:

- A variable whose description starts with `DEPRECATED` must be in `deprecated_variables.tf`.
- Every variable in `deprecated_variables.tf` must have a description that starts with `DEPRECATED`.

The `DEPRECATED` match is case-sensitive and ignores leading whitespace. Descriptions that aren't static strings are skipped.

## How To Fix
Move the variable to `deprecated_variables.tf`. For a variable already in that file, start its description with `DEPRECATED` and name the replacement.

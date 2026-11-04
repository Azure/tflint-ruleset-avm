# avm_terraform_resource_id_variable_validation_required

checks that resource ID variables are validated with `provider::azapi::parse_resource_id`

## Example

```hcl
variable "subnet_resource_id" {
  type     = string
  nullable = false

  validation {
    condition     = can(regex("^/subscriptions/.+/subnets/.+$", var.subnet_resource_id))
    error_message = "Must be a subnet resource ID."
  }
}
```

```
$ tflint
2 issue(s) found:

Error: Variable `subnet_resource_id` must not validate its resource ID with `regex`; use `provider::azapi::parse_resource_id` instead (avm_terraform_resource_id_variable_validation_required)

  on variables.tf line 6:
   6:     condition     = can(regex("^/subscriptions/.+/subnets/.+$", var.subnet_resource_id))

Reference: https://azure.github.io/Azure-Verified-Modules/spec/TFNFR38/

Error: Variable `subnet_resource_id` holds a resource ID and must be validated with `can(provider::azapi::parse_resource_id("<resource type>", ...))` (avm_terraform_resource_id_variable_validation_required)

  on variables.tf line 1:
   1: variable "subnet_resource_id" {

Reference: https://azure.github.io/Azure-Verified-Modules/spec/TFNFR38/

```

## Why
[TFNFR38](https://azure.github.io/Azure-Verified-Modules/spec/TFNFR38/) says variables that hold Azure resource IDs MUST be validated with `provider::azapi::parse_resource_id`. The resource type argument MUST be a literal string, and the call MUST be wrapped in `can(...)`. Hand-rolled `regex`, `startswith`, `length`, or `split` checks MUST NOT be used.

## What This Rule Checks
- The rule runs only when the module declares the `azapi` provider in `required_providers`, because `parse_resource_id` is an AzAPI provider function.
- Resource ID variables are top-level variables named `parent_id`, `resource_id`, `resource_ids`, `*_resource_id`, or `*_resource_ids`, typed `string` or `list`/`set`/`map(string)`.
- Each resource ID variable needs at least one validation that calls `parse_resource_id`.
- Collection variables must check each element inside a `for` expression, for example `alltrue([for v in var.x : can(...)])`.
- Optional string variables (`default = null` without `nullable = false`) must short-circuit on `null`, for example `var.x == null || can(...)`.
- `regex`, `regexall`, `startswith`, and `split` must not appear in the validation of a resource ID variable. `length` is flagged only for string variables, so collection size checks stay allowed.
- In any variable's validation, `parse_resource_id` calls must use a literal resource type and be wrapped in `can(...)`.
- Nested attributes inside object variables are not checked yet.

## How To Fix
```hcl
variable "subnet_resource_id" {
  type     = string
  nullable = false

  validation {
    condition     = can(provider::azapi::parse_resource_id("Microsoft.Network/virtualNetworks/subnets", var.subnet_resource_id))
    error_message = "Must be a subnet resource ID."
  }
}
```

The spec allows variables that legitimately accept IDs of more than one resource type to skip this validation. Suppress the rule for those variables with a `# tflint-ignore: avm_terraform_resource_id_variable_validation_required` comment.

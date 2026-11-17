# avm_terraform_module_control_variable_disallowed

checks that the module doesn't declare variables named `enabled` or `module_depends_on`

## Example

```hcl
variable "enabled" {
  type    = bool
  default = true
}
```

```
$ tflint
1 issue(s) found:

Error: Variable `enabled` must not be declared to control the entire module's operation; use feature-specific toggles instead (avm_terraform_module_control_variable_disallowed)

  on variables.tf line 1:
   1: variable "enabled" {

Reference: https://azure.github.io/Azure-Verified-Modules/spec/TFNFR14/

```

## Why
[TFNFR14](https://azure.github.io/Azure-Verified-Modules/spec/TFNFR14/) says module owners MUST NOT add variables like `enabled` or `module_depends_on` to control the entire module's operation. Terraform supports `count`, `for_each`, and `depends_on` on module calls for that purpose.

The rule only flags these two exact names. Boolean feature toggles such as `public_network_access_enabled` are allowed.

## How To Fix
Remove the variable. Callers can use `count`, `for_each`, or `depends_on` on the module block instead.

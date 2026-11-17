# avm_terraform_variable_negative_toggle_disallowed

checks that boolean feature switch variables use positive names

## Example

```hcl
variable "public_network_access_disabled" {
  type    = bool
  default = false
}
```

```
$ tflint
1 issue(s) found:

Notice: Feature switch variable `public_network_access_disabled` should use a positive name such as `xxx_enabled` instead of `xxx_disabled` (avm_terraform_variable_negative_toggle_disallowed)

  on variables.tf line 1:
   1: variable "public_network_access_disabled" {

Reference: https://azure.github.io/Azure-Verified-Modules/spec/TFNFR16/

```

## Why
[TFNFR16](https://azure.github.io/Azure-Verified-Modules/spec/TFNFR16/) says feature switch variables SHOULD use a positive statement: `xxx_enabled` instead of `xxx_disabled`.

The rule only flags variables with `type = bool` whose names start with `disable_` or end with `_disabled`.

## How To Fix
Rename the variable to a positive form, such as `public_network_access_enabled`, and invert its default and usage.

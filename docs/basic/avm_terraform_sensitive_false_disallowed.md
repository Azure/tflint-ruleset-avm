# avm_terraform_sensitive_false_disallowed

checks that variables and outputs don't declare `sensitive = false`

## Example

```hcl
variable "name" {
  type      = string
  sensitive = false
}
```

```
$ tflint
1 issue(s) found:

Error: `sensitive = false` must not be declared on variable `name`; values are not sensitive by default (avm_terraform_sensitive_false_disallowed)

  on main.tf line 3:
   3:   sensitive = false

Reference: https://azure.github.io/Azure-Verified-Modules/spec/TFNFR22/

```

## Why
[TFNFR22](https://azure.github.io/Azure-Verified-Modules/spec/TFNFR22/) says `sensitive = false` MUST be avoided. Variables and outputs aren't sensitive by default, so the attribute adds noise.

## How To Fix
Remove the `sensitive = false` attribute.

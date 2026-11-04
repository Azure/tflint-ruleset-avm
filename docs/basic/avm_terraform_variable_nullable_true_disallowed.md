# avm_terraform_variable_nullable_true_disallowed

checks that variables don't declare `nullable = true`

## Example

```hcl
variable "name" {
  type     = string
  nullable = true
}
```

```
$ tflint
1 issue(s) found:

Error: `nullable = true` must not be declared on variable `name` (avm_terraform_variable_nullable_true_disallowed)

  on variables.tf line 3:
   3:   nullable = true

Reference: https://azure.github.io/Azure-Verified-Modules/spec/TFNFR21/

```

## Why
[TFNFR21](https://azure.github.io/Azure-Verified-Modules/spec/TFNFR21/) says `nullable = true` MUST be avoided. Variables are nullable by default, so the attribute adds noise.

## How To Fix
Remove the `nullable = true` attribute. Use `nullable = false` if the variable must never be `null`.

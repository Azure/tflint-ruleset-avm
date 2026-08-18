# terraform_variable_nullable_false

According to the [document](https://developer.hashicorp.com/terraform/language/values/variables#disallowing-null-input-values):

>The default value for `nullable` is `true`.

So simplify the code, we can just remove `nullable = true` attribute for a variable block.

## Example

```hcl
variable "var" {
  type     = string
  nullable = false
}
```

```
$ tflint
1 issue(s) found:

Notice: `nullable` is default to `true` so we don't need to declare it explicitly. (terraform_variable_nullable_false)

  on main.tf line 3:
  3:   nullable = true

Reference: https://github.com/Azure/tflint-ruleset-avm/blob/main/docs/basic/terraform_variable_nullable_false.md
```

## Why
It helps to simplifier the Terraform code.

## How To Fix
Just remove `nullable = true`.
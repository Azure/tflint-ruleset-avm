# terraform_locals_order

Recommend proper order for variables in locals blocks
Those variables are sorted based on their names (alphabetic order)

## Example

```hcl
locals {
  service_name = "forum"
  owner        = "Community Team"
}
```

```
$ tflint
1 issue(s) found:

Notice: Recommended locals variable order:
locals {
  owner        = "Community Team"
  service_name = "forum"
} (terraform_locals_order)

  on main.tf line 1:
   1: locals {

Reference: https://github.com/Azure/tflint-ruleset-avm/blob/main/docs/basic/terraform_locals_order.md
```

## Why
It helps to improve the readability of terraform code by sorting variables in locals blocks in the order above.

## How To Fix
Just copy the text with recommended locals variable order and paste it in the tf config file to overwrite the original style of it.
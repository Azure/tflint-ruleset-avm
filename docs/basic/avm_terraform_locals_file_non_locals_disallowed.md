# avm_terraform_locals_file_non_locals_disallowed

checks that `locals.tf` only contains `locals` blocks

## Example

```hcl
# locals.tf
locals {
  name = "example"
}

resource "azurerm_resource_group" "this" {
  name     = local.name
  location = "westeurope"
}
```

```
$ tflint
1 issue(s) found:

Error: `locals.tf` must only contain `locals` blocks, found `resource` block (avm_terraform_locals_file_non_locals_disallowed)

  on locals.tf line 6:
   6: resource "azurerm_resource_group" "this" {

Reference: https://azure.github.io/Azure-Verified-Modules/spec/TFNFR31/

```

## Why
[TFNFR31](https://azure.github.io/Azure-Verified-Modules/spec/TFNFR31/) says `locals.tf` may hold multiple `locals` blocks, but only `locals` blocks are allowed.

The rule only checks files named `locals.tf`. The spec allows `locals` blocks next to `resource` or `data` blocks in other files.

## How To Fix
Move the non-`locals` block to another file, such as `main.tf`.

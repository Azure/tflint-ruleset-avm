# avm_terraform_provider_block_disallowed

Check the usage of `provider` block in terraform module

## Example

```hcl
provider "azurerm" {
  alias    = "test"
  location = "east"
}
```

```
$ tflint
1 issue(s) found:

Warning: Provider block in terraform module is expected to have and only have `alias` declared (avm_terraform_provider_block_disallowed)

  on main.tf line 1:
   1: provider "azurerm" {

Reference: https://github.com/Azure/tflint-ruleset-avm/blob/main/docs/basic/avm_terraform_provider_block_disallowed.md
```

## Why
The declaration of `provider` block in module is not expected unless it has and only has `alias` field declared to prevent bugs.
See https://www.terraform.io/language/modules/develop/providers

## How To Fix
Consider remove the `provider` block from terraform module or reformat it to have and only have `alias` field declared
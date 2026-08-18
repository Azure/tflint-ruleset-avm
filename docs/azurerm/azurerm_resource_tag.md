# azurerm_resource_tag

Check whether the tags argument is set if it's supported in a (nested block of) Azurerm resource

## Example

```hcl
resource "azurerm_container_group" "example" {
  location            = azurerm_resource_group.example.location
  name                = "example-continst"
  ip_address_type     = "Public"
  dns_name_label      = "aci-label"
  os_type             = "Linux"
  resource_group_name = azurerm_resource_group.example.name
  
  container {
    cpu    = "0.5"
    image  = "mcr.microsoft.com/azuredocs/aci-tutorial-sidecar"
    memory = "1.5"
    name   = "sidecar"
  }
}
```

```
$ tflint
1 issue(s) found:

Notice: `tags` argument is not set but supported in resource `azurerm_container_group` (azurerm_resource_tag)

  on main.tf line 1:
   1: resource "azurerm_container_group" "example" {

Reference: https://github.com/Azure/tflint-ruleset-avm/blob/main/docs/azurerm/azurerm_resource_tag.md
```

## Why

It helps users to know which resource supports tags and customize this argument based on their needs.

## How To Fix

Specify the tags argument in corresponding blocks with it supported based on your need.
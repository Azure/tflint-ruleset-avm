# avm_azapi_primary_resource_naming_valid

checks that the primary AzAPI managed resource is named `this` and that satellite AzAPI resources don't reuse that name

## Example

```hcl
resource "azapi_resource" "key_vault" {
  type = "Microsoft.KeyVault/vaults@2023-07-01"
}

resource "azapi_resource" "lock" {
  type = "Microsoft.Authorization/locks@2020-05-01"
}
```

```
$ tflint
1 issue(s) found:

Error: The primary AzAPI resource must be named `this`, for example `azapi_resource.this` (avm_azapi_primary_resource_naming_valid)

  on main.tf line 1:
   1: resource "azapi_resource" "key_vault" {

Reference: https://azure.github.io/Azure-Verified-Modules/spec/TFRMNFR2/

```

## Why
[TFRMNFR2](https://azure.github.io/Azure-Verified-Modules/spec/TFRMNFR2/) says the primary AzAPI resource of a resource module (`azapi_resource`, `azapi_update_resource`, `azapi_resource_action`, or another AzAPI resource) MUST be named `this`. Satellite resources such as locks, role assignments, diagnostic settings, and private endpoints MUST NOT be named `this`; name them after what they represent.

The rule considers AzAPI managed resources only. Data sources, including `azapi_resource_list`, are ignored, and so are resources listed in `ignored_resources`. When a module declares any AzAPI managed resource, the rule reports:

- One issue when none of them is named `this`.
- One issue for each extra resource named `this`. `azapi_resource.this` is treated as the primary resource when present; otherwise the first resource named `this` is.

The spec lets utility modules, and modules whose primary functionality is two or more peer AzAPI resources, skip the `this` name if the README explains why. tflint can't verify that, so those modules should disable this rule in their tflint configuration. Pattern and utility modules can also set `module_class`; see [Resource module rules](../../README.md#resource-module-rules).

## How To Fix
Rename the primary resource to `this` and update every reference to it:

```hcl
resource "azapi_resource" "this" {
  type = "Microsoft.KeyVault/vaults@2023-07-01"
}

resource "azapi_resource" "lock" {
  type = "Microsoft.Authorization/locks@2020-05-01"
}
```

Add a `moved` block when renaming a resource in a released module, so existing deployments don't recreate it.

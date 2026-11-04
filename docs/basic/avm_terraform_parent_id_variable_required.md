# avm_terraform_parent_id_variable_required

checks that AzAPI resource modules expose their parent scope through a required `parent_id` string variable

## Example

```hcl
# variables.tf
variable "resource_group_name" {
  type = string
}

# main.tf
resource "azapi_resource" "this" {
  type      = "Microsoft.KeyVault/vaults@2023-07-01"
  name      = var.name
  parent_id = "/subscriptions/${var.subscription_id}/resourceGroups/${var.resource_group_name}"
}
```

```
$ tflint
3 issue(s) found:

Error: `azapi_resource.this` must set `parent_id = var.parent_id` (avm_terraform_parent_id_variable_required)

  on main.tf line 4:
   4:   parent_id = "/subscriptions/${var.subscription_id}/resourceGroups/${var.resource_group_name}"

Reference: https://azure.github.io/Azure-Verified-Modules/spec/TFRMFR1/

Error: Variable `parent_id` must be declared to expose the parent scope of the primary AzAPI resource (avm_terraform_parent_id_variable_required)

  on variables.tf line 0:

Reference: https://azure.github.io/Azure-Verified-Modules/spec/TFRMFR1/

Error: Variable `resource_group_name` must not be declared; expose the parent scope through the `parent_id` variable instead (avm_terraform_parent_id_variable_required)

  on variables.tf line 1:
   1: variable "resource_group_name" {

Reference: https://azure.github.io/Azure-Verified-Modules/spec/TFRMFR1/

```

## Why
[TFRMFR1](https://azure.github.io/Azure-Verified-Modules/spec/TFRMFR1/) says a Terraform resource module MUST expose its parent scope as a single required `string` variable named `parent_id`, and MUST assign it to the `parent_id` argument of the primary AzAPI resource. Modules MUST NOT expose parent-scope-specific variables such as `resource_group_name`. [RMFR3](https://azure.github.io/Azure-Verified-Modules/spec/RMFR3/) states the same requirement for Terraform.

The rule runs only when the module declares at least one AzAPI managed resource, so AzureRM-based modules aren't flagged. It checks that:

- A `parent_id` variable exists, has `type = string`, and has no `default`.
- No variable is named `resource_group_name` or `resource_group_resource_id`.
- `azapi_resource.this`, when present, sets `parent_id = var.parent_id` directly.

Satellite resources may derive their `parent_id` from the primary resource, for example `azapi_resource.this.id`. The rule doesn't check `nullable = false`, which the spec requires only for extension-resource modules. Validation of `parent_id` with `provider::azapi::parse_resource_id` is out of scope for this rule.

The rule applies to resource modules only. See [Resource module rules](../../README.md#resource-module-rules) to configure `module_class` for pattern and utility modules.

## How To Fix
Replace parent-scope-specific variables with a required `parent_id` variable and forward it to the primary resource:

```hcl
variable "parent_id" {
  type        = string
  description = "The resource ID of the resource group to deploy the Key Vault into."
  nullable    = false
}

resource "azapi_resource" "this" {
  type      = "Microsoft.KeyVault/vaults@2023-07-01"
  name      = var.name
  parent_id = var.parent_id
}
```

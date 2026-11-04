# avm_terraform_count_for_each_usage_valid

checks that `count` and `for_each` don't create instances keyed by list position

## Example

```hcl
variable "subnet_names" {
  type = list(string)
}

resource "azapi_resource" "subnet" {
  count = length(var.subnet_names)
  # ...
}
```

```
$ tflint
1 issue(s) found:

Error: `count` must not be set from `length(...)`; use `for_each` with a map or set so instances keep stable keys (avm_terraform_count_for_each_usage_valid)

  on main.tf line 6:
   6:   count = length(var.subnet_names)

Reference: https://azure.github.io/Azure-Verified-Modules/spec/TFNFR7/

```

## Why
[TFNFR7](https://azure.github.io/Azure-Verified-Modules/spec/TFNFR7/) says `for_each` MUST use a `map` or `set` collection. Instances keyed by list position get recreated when an earlier element is removed or reordered. `count` is fine for conditional creation, such as `count = var.enabled ? 1 : 0`.

## What This Rule Checks
The rule checks `resource` and `module` blocks. It flags:
- `count` set to `length(...)`, directly or as a result of a conditional. `length(...)` inside the condition itself is allowed.
- `for_each = var.x` where `x` is declared with a `list(...)` or `tuple(...)` type.
- `for_each = { for i, v in var.x : i => ... }` where `x` is a list or tuple and the key is the index.

The spec also requires map keys and set elements to be static literals. This rule can't verify that, so it doesn't try. Data sources and `dynamic` blocks are not checked.

## How To Fix
Declare the variable as a map or set and iterate it with `for_each`:

```hcl
variable "subnet_names" {
  type     = set(string)
  nullable = false
}

resource "azapi_resource" "subnet" {
  for_each = var.subnet_names
  # ...
}
```

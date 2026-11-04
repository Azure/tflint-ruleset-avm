# avm_terraform_collection_variable_nullable_false_required

checks that collection variables declare `nullable = false`

## Example

```hcl
variable "subnets" {
  type    = map(string)
  default = {}
}
```

```
$ tflint
1 issue(s) found:

Error: Collection variable `subnets` must set `nullable = false` (avm_terraform_collection_variable_nullable_false_required)

  on variables.tf line 1:
   1: variable "subnets" {

Reference: https://azure.github.io/Azure-Verified-Modules/spec/TFNFR21/

```

## Why
[TFNFR21](https://azure.github.io/Azure-Verified-Modules/spec/TFNFR21/) says variables MUST set `nullable = false` when their type has a meaningful zero value. Its behavior-toggle exception does not extend to collection-shaped inputs. [TFNFR20](https://azure.github.io/Azure-Verified-Modules/spec/TFNFR20/) adds that collections used in loops SHOULD be non-nullable, so callers can pass an empty collection instead of `null`.

## What This Rule Checks
- Variables typed `list(...)`, `set(...)`, `map(...)`, or `tuple(...)` must declare `nullable`.
- The `tags` variable is skipped. The AVM interfaces spec requires `tags` to stay a `map(string)` with a default of `null`.
- Explicit `nullable` values are not checked here. `avm_terraform_variable_nullable_true_disallowed` already flags `nullable = true`.

## How To Fix
```hcl
variable "subnets" {
  type     = map(string)
  default  = {}
  nullable = false
}
```

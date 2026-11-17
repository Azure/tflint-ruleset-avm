# avm_terraform_deprecated_output_file_required

checks that deprecated outputs are declared in `deprecated_outputs.tf`

## Example

```hcl
# outputs.tf
output "id" {
  value       = azurerm_resource_group.this.id
  description = "DEPRECATED, use `resource_id` instead."
}
```

```
$ tflint
1 issue(s) found:

Error: Deprecated output `id` must be declared in `deprecated_outputs.tf` (avm_terraform_deprecated_output_file_required)

  on outputs.tf line 2:
   2: output "id" {

Reference: https://azure.github.io/Azure-Verified-Modules/spec/TFNFR30/

```

## Why
[TFNFR30](https://azure.github.io/Azure-Verified-Modules/spec/TFNFR30/) says a deprecated output MUST move to an independent `deprecated_outputs.tf` file. Its name MUST NOT change directly.

The rule treats an output as deprecated when its static `description` starts with `DEPRECATED`. The match is case-sensitive and ignores leading whitespace.

## How To Fix
Move the output to `deprecated_outputs.tf`, then declare the replacement output in `outputs.tf`.

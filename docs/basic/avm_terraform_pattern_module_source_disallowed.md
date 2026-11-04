# avm_terraform_pattern_module_source_disallowed

checks that resource modules don't call AVM pattern modules

## Example

```hcl
module "alz" {
  source  = "Azure/avm-ptn-alz/azurerm"
  version = "0.12.0"
}
```

```
$ tflint
1 issue(s) found:

Error: Module `alz` must not call the AVM pattern module `Azure/avm-ptn-alz/azurerm`; resource modules may only call AVM resource or utility modules (avm_terraform_pattern_module_source_disallowed)

  on main.tf line 2:
   2:   source  = "Azure/avm-ptn-alz/azurerm"

Reference: https://azure.github.io/Azure-Verified-Modules/spec/RMFR8/

```

## Why
[RMFR8](https://azure.github.io/Azure-Verified-Modules/spec/RMFR8/) says a resource module MAY call other AVM resource modules, but MUST NOT call non-AVM modules or AVM pattern modules.

`avm_terraform_module_source_required` already rejects non-AVM sources. This rule covers the remaining case: literal registry sources in the `Azure` namespace whose module name starts with `avm-ptn-`. It accepts an optional registry hostname and a `//subdir` suffix, and matches the namespace and prefix case-insensitively. AVM resource modules (`avm-res-`), AVM utility modules (`avm-utl-`), and local sources are allowed.

The rule applies to resource modules only. See [Resource module rules](../../README.md#resource-module-rules) to configure `module_class` for pattern and utility modules.

## How To Fix
Call the AVM resource modules the pattern module wraps instead, or move the composition into a pattern module.

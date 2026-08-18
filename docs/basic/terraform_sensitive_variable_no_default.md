# terraform_sensitive_variable_no_default

checks whether the default value is set for sensitive variable

## Example

```hcl
variable "availability_zone_names" {
  type      = list(string)
  default   = ["us-west-1a"]
  sensitive = true
}
```

```
$ tflint
1 issue(s) found:

Warning: Default value is not expected to be set for sensitive variable `availability_zone_names` (terraform_sensitive_variable_no_default)

  on main.tf line 3:
   3:   default   = ["us-west-1a"]

Reference: https://github.com/Azure/tflint-ruleset-avm/blob/main/docs/basic/terraform_sensitive_variable_no_default.md

```

## Why
Sensitive variable shouldn't have default value set.

## How To Fix
Change variable to insensitive or delete its default value.
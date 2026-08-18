# terraform_required_version_declaration

Check whether `required_version` is declared at the beginning of terraform setting block

## Example

```hcl
terraform {
  required_providers {
    aws = {
      version = ">= 2.7.0"
      source = "hashicorp/aws"
    }
  }
  required_version = "~> 0.12.29"
}
```

```
$ tflint
1 issue(s) found:

Notice: The `required_version` field should be declared at the beginning of `terraform` block (terraform_required_version_declaration)

  on main.tf line 8:
   8:   required_version = "~> 0.12.29"

Reference: https://github.com/Azure/tflint-ruleset-avm/blob/main/docs/basic/terraform_required_version_declaration.md
```

## Why
To better manage terraform CLI version for modules and improve readability of the code, the `required_version` field should be declared at the beginning of `terraform` block

## How To Fix
Declare the `required_version` field at the beginning of `terraform` block
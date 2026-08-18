# terraform_required_providers_declaration

Check whether `required_providers` block is declared in the terraform setting block and whether the arguments of it are sorted in alphabetic order

## Example

```hcl
terraform {
  required_version = "~> 0.12.29"
  required_providers {
    azurerm = {
      version = "~> 3.0.2"
      source  = "hashicorp/azurerm"
    }
    aws = {
      source = "hashicorp/aws"
      version = ">= 2.7.0"
    }
  }
}
```

```
$ tflint
1 issue(s) found:

Notice: The arguments of `required_providers` are expected to be sorted as follows:
required_providers {
  aws = {
    source  = "hashicorp/aws"
    version = ">= 2.7.0"
  }
  azurerm = {
    source  = "hashicorp/azurerm"
    version = "~> 3.0.2"
  }
} (terraform_required_providers_declaration)

  on versions.tf line 3:
   3:   required_providers {

Reference: https://github.com/Azure/tflint-ruleset-avm/blob/main/docs/basic/terraform_required_providers_declaration.md
```

## Why
There must be a `required_providers` block in `terraform` block to declare the version and source information of the providers used in the project, 
and it helps to improve the readability of code by sorting the arguments of this block in alphabetic order

## How To Fix
Declare the `required_providers` block in `terraform` block
Copy the text with recommended argument order of a specific code block and paste it in the tf config file to overwrite the original style of this code block.
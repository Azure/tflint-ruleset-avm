# terraform_variable_separate

checks whether the variables are declared in a file with other types of blocks declared

## Example

```hcl
terraform {}

variable "docker_ports" {
  type = list(object({
    internal = number
    external = number
    protocol = string
  }))
  default = [
    {
      internal = 8300
      external = 8300
      protocol = "tcp"
    }
  ]
}

variable "availability_zone_names" {
  type    = list(string)
  default = ["us-west-1a"]
}

variable "image_id" {
  type = string
}
```

```
$ tflint
1 issue(s) found:

Notice: Putting variables and other types of blocks in the same file is not recommended (terraform_variable_separate)

  on main.tf line 1:
   1: terraform{}

Reference: https://github.com/Azure/tflint-ruleset-avm/blob/main/docs/basic/terraform_variable_separate.md

```

## Why
It helps to improve the readability and development efficiency of terraform code by separating variables from other types of blocks.

## How To Fix
Just consider putting the variable blocks in a separate file.
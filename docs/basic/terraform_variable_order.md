# terraform_variable_order

Recommend proper order for variable blocks
The variables without default value are placed prior to those with default value set
Then the variables are sorted based on their names (alphabetic order)

## Example

```hcl
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

Notice: Recommended variable order:
variable "image_id" {
  type = string
}

variable "availability_zone_names" {
  type    = list(string)
  default = ["us-west-1a"]
}

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
} (terraform_variable_order)

  on main.tf line 1:
   1: variable "image_id" {

Reference: https://github.com/Azure/tflint-ruleset-avm/blob/main/docs/basic/terraform_variable_order.md
```

## Why
It helps to improve the readability of terraform code by sorting variable blocks in the order above.

## How To Fix
Just copy the text with recommended variable order and paste it in the tf config file to overwrite the original style of it.
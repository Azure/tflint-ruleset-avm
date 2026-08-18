# terraform_versions_file

check whether `versions.tf` has and only has 1 `terraform` block

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
```

```
$ tflint
1 issue(s) found:

Notice: `versions.tf` should have and only have 1 `terraform` block (terraform_versions_file)

  on  line 0:
   (source code not available)

Reference: https://github.com/Azure/tflint-ruleset-avm/blob/main/docs/basic/terraform_versions_file.md
```

## Why
To better manage terraform project, it's better to align with the agreement that `versions.tf` should have and only have 1 `terraform` block

## How To Fix
Clear other types of blocks in `versions.tf`
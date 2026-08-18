# terraform_output_separate

checks whether the outputs are declared in a file with other types of blocks declared

## Example

```hcl
terraform{}

output "api_base_url" {
  value = "https://${aws_instance.example.private_dns}:8433/"

  # The EC2 instance must have an encrypted root volume.
  precondition {
    condition     = data.aws_ebs_volume.example.encrypted
    error_message = "The server's root volume is not encrypted."
  }
}

output "db_password" {
  value       = aws_db_instance.db.password
  description = "The password for logging in to the database."
  sensitive   = true
}

output "instance_ip_addr" {
  value       = aws_instance.server.private_ip
  description = "The private IP address of the main server instance."
}
```

```
$ tflint
1 issue(s) found:

Notice: Putting outputs and other types of block in the same file is not recommended (terraform_output_separate)

  on main.tf line 1:
   1: terraform{}

Reference: https://github.com/Azure/tflint-ruleset-avm/blob/main/docs/basic/terraform_output_separate.md
```

## Why
It helps to improve the readability and development efficiency of terraform code by separating outputs from other types of blocks.

## How To Fix
Consider putting the output blocks in a separate file.
# terraform_output_order

Recommend proper order for output blocks
The outputs are sorted based on their names (alphabetic order)

## Example

```hcl
output "instance_ip_addr" {
  value       = aws_instance.server.private_ip
  description = "The private IP address of the main server instance."
}

output "db_password" {
  value       = aws_db_instance.db.password
  description = "The password for logging in to the database."
  sensitive   = true
}

output "api_base_url" {
  value = "https://${aws_instance.example.private_dns}:8433/"

  # The EC2 instance must have an encrypted root volume.
  precondition {
    condition     = data.aws_ebs_volume.example.encrypted
    error_message = "The server's root volume is not encrypted."
  }
}
```

```
$ tflint
1 issue(s) found:

Notice: Recommended output order:
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
} (terraform_output_order)

  on main.tf line 1:
   1: output "instance_ip_addr" {

Reference: https://github.com/Azure/tflint-ruleset-avm/blob/main/docs/basic/terraform_output_order.md
```

## Why
It helps to improve the readability of terraform code by sorting output blocks in the order above.

## How To Fix
Just copy the text with recommended output order and paste it in the tf config file to overwrite the original style of it.
variable "private_endpoints" {
  type = map(object({
    name               = optional(string, null)
  }))
  default     = {}
  nullable    = false
  description = <<DESCRIPTION
  A map of private endpoints to create on the Key Vault. The map key is deliberately arbitrary to avoid issues where map keys maybe unknown at plan time.

  - `name` - (Optional) The name of the private endpoint. One will be generated if not set.
  DESCRIPTION
}

output "resource_id" {
  # Just make var.private_endpoints an used variable, not mean to be valid.
  value = var.private_endpoints["default"].name
}
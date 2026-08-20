resource "azapi_resource" "example" {
  type                   = var.resource_types.widget
  name                   = "example"
  parent_id              = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example"
  response_export_values = []
}

variable "resource_types" {
  type = object({
    widget = optional(string, "Microsoft.Example/widgets@2024-01-01")
  })
  default  = {}
  nullable = false
}

variable "retry" {
  type = object({
    error_message_regex  = optional(list(string))
    interval_seconds     = optional(number)
    max_interval_seconds = optional(number)
  })
  default = null
}

variable "timeouts" {
  type = object({
    create = optional(string)
    read   = optional(string)
    update = optional(string)
    delete = optional(string)
  })
  default = null
}

variable "ignore_body_changes" {
  type = object({
    widget = optional(list(string), [])
  })
  default  = {}
  nullable = false
}

output "resource_id" {
  value = azapi_resource.example.id
}

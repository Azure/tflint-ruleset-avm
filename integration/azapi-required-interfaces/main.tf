resource "azapi_resource" "example" {
  type                   = var.resource_types.widget
  retry                  = var.retry
  ignore_body_changes    = var.ignore_body_changes.widget
  name                   = "example"
  parent_id              = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example"
  response_export_values = []

  timeouts {
    create = try(var.timeouts.create, null)
    read   = try(var.timeouts.read, null)
    update = try(var.timeouts.update, null)
    delete = try(var.timeouts.delete, null)
  }
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
  default = {}
}

variable "timeouts" {
  type = object({
    create = optional(string)
    read   = optional(string)
    update = optional(string)
    delete = optional(string)
  })
  default = {
    delete = "30m"
  }
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

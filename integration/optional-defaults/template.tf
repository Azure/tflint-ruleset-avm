variable "variable" {
  type = object({
    default = optional(string, "Standard")
  })
  default = {}
}

resource "azurerm_lb" "test" {
  sku = var.variable.default
}

output "resource" {
  value = null
}

output "resource_id" {
  value = null
}

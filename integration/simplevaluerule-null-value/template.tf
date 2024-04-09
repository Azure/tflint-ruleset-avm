variable "variable" {
  type    = string
  default = null
}

resource "azurerm_lb" "test" {
  sku = var.variable
}

output "resource" {
  value = null
}

output "resource_id" {
  value = null
}

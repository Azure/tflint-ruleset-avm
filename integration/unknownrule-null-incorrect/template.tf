variable "variable" {
  type    = number
  default = 1
}

resource "azurerm_virtual_machine" "test" {
  zone = var.variable
}

output "resource" {
  value = null
}

output "resource_id" {
  value = null
}

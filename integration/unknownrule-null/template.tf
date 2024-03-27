variable "variable" {
  type = number
}

resource "azurerm_virtual_machine" "test" {
  zone = var.variable
}

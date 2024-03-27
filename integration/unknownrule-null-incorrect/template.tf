variable "variable" {
  type    = number
  default = 1
}

resource "azurerm_virtual_machine" "test" {
  zone = var.variable
}

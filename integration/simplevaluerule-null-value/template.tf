variable "variable" {
  type    = string
  default = null
}

resource "azurerm_lb" "test" {
  sku = var.variable
}

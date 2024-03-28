variable "variable" {
  type = object({
    default = optional(string, "Standard")
  })
  default = {}
}

resource "azurerm_lb" "test" {
  sku = var.variable.default
}

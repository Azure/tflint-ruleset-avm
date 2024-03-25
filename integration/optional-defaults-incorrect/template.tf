variable "variable" {
  type = object({
    default = optional(string, "invalid")
  })
  default = {}
}

resource "azurerm_lb" "test" {
  sku = var.variable.default
}

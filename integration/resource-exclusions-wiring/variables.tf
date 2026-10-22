variable "resource_types" {
  type = object({
    widget = optional(string, "Microsoft.Example/widgets@2024-01-01")
  })
  default  = {}
  nullable = false
}

resource "azapi_resource" "example" {
  type                   = "Microsoft.Example/widgets@2024-01-01"
  name                   = "example"
  parent_id              = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example"
  response_export_values = []
}

output "resource_id" {
  value = azapi_resource.example.id
}

resource "azapi_resource" "telemetry" {
  type = "Microsoft.Example/telemetry@2024-01-01"
}

resource "azapi_resource" "business_first" {
  type = "Microsoft.Example/widgets@2024-01-01"
}

resource "azapi_resource" "business_second" {
  type = "Microsoft.Example/widgets@2024-01-01"
}

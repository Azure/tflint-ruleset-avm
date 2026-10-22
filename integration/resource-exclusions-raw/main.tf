resource "azapi_resource" "telemetry" {
  body = <<-JSON
{"name":"telemetry"}
JSON
}

data "azapi_resource" "telemetry" {
  body = <<-JSON
{"name":"telemetry"}
JSON
}

resource "azapi_resource" "business" {
  body = <<-JSON
{"name":"business"}
JSON
}

locals {
  business = <<-JSON
{"name":"local"}
JSON
}

plugin "terraform" {
  enabled = false
}

plugin "avm" {
  enabled = true
}

rule "avm_interface_resource_types" {
  enabled           = true
  ignored_resources = ["azapi_resource.telemetry"]
}

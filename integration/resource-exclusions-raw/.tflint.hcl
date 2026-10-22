plugin "terraform" {
  enabled = false
}

plugin "avm" {
  enabled = true
}

rule "avm_terraform_literal_heredoc_disallowed" {
  enabled = true
  ignored_resources = [
    "azapi_resource.telemetry",
    "data.azapi_resource.telemetry",
  ]
}

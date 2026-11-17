plugin "terraform" {
  enabled = false
}

plugin "avm" {
  enabled = true
}

rule "avm_interface_resource_types" {
  enabled  = true
  severity = "notice"
}

rule "avm_provider_modtm_disallowed" {
  enabled = false
}

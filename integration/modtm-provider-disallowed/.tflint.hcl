plugin "terraform" {
  enabled = false
}

plugin "avm" {
  enabled = true
}

rule "avm_provider_modtm_disallowed" {
  enabled = true
}

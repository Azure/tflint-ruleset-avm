plugin "terraform" {
  enabled = false
}

plugin "avm" {
  enabled = true
}

rule "avm_provider_azurerm_disallowed" {
  enabled = false
}

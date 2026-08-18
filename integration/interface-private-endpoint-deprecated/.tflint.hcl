plugin "terraform" {
  enabled = false
}

plugin "avm" {
  enabled = true
}

rule "provider_azurerm_disallowed" {
  enabled = false
}

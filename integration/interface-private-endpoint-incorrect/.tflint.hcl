plugin "terraform" {
  enabled = false
}

plugin "avm" {
  enabled = true
}

# Disable rules unrelated to the interface being exercised by this fixture.
rule "avm_provider_azurerm_disallowed" {
  enabled = false
}

rule "avm_provider_modtm_disallowed" {
  enabled = false
}

plugin "terraform" {
  enabled = false
}

plugin "avm" {
  enabled = true
}

# Disable rules unrelated to the interface being exercised by this fixture.
rule "provider_azurerm_disallowed" {
  enabled = false
}

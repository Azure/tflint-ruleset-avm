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

# Disable resource-module rules unrelated to the interfaces exercised by this fixture.
rule "avm_terraform_parent_id_variable_required" {
  enabled = false
}

rule "avm_azapi_primary_resource_naming_valid" {
  enabled = false
}
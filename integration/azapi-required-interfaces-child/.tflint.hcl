plugin "terraform" {
  enabled = false
}

plugin "avm" {
  enabled = true
}

rule "avm_terraform_module_source_required" {
  enabled = false
}

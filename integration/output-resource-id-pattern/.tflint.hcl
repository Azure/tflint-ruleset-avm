plugin "terraform" {
  enabled = false
}

plugin "avm" {
  enabled = true
}

rule "avm_output_resource_id_required" {
  enabled      = true
  module_class = "pattern"
}

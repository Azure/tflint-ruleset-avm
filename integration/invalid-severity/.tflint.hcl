plugin "terraform" {
  enabled = false
}

plugin "avm" {
  enabled = true
}

rule "avm_terraform_literal_heredoc_disallowed" {
  enabled  = true
  severity = "critical"
}

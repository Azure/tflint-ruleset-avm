package main

import (
	"github.com/Azure/tflint-ruleset-avm/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "template",
			Version: "v0.1.0",
			Rules:   rules.Rules,
		},
	})
}

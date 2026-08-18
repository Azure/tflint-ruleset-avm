package project

import "fmt"

// Version is ruleset version
var Version string = "0.6.0"

// ReferenceLink returns the rule reference link
func ReferenceLink(name string) string {
	return fmt.Sprintf("https://github.com/Azure/tflint-ruleset-avm/blob/main/docs/basic/%s.md", name)
}

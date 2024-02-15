package rules

import "github.com/Azure/tflint-ruleset-avm/interfaces"

// NewAvmInterfaceRule returns a new rule with the given interface.
func NewAvmInterfaceRule(i interfaces.AvmInterface) *AvmInterfaceRule {
	r := new(AvmInterfaceRule)
	r.Iface = i
	return r
}

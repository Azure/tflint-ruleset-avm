package rules

import "github.com/Azure/tflint-ruleset-avm/interfaces"

func NewAvmInterfaceLockRule() *AVMInterfaceRule {
	r := new(AVMInterfaceRule)
	r.Iface = interfaces.Lock
	return r
}

package rules

import "github.com/Azure/tflint-ruleset-avm/interfaces"

func NewAvmInterfaceDisgnosticSettingsRule() *AVMInterfaceRule {
	r := new(AVMInterfaceRule)
	r.Iface = interfaces.DiagnosticSettings
	return r
}

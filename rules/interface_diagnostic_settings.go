package rules

import "github.com/Azure/tflint-ruleset-avm/interfaces"

func NewAvmInterfaceDiagnosticSettingsRule() *AVMInterfaceRule {
	r := new(AVMInterfaceRule)
	r.Iface = interfaces.DiagnosticSettings
	return r
}

package interfaces

import (
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var Rules = []tflint.Rule{
	NewVarCheckRuleFromAvmInterface(CustomerManagedKey),
	NewVarCheckRuleFromAvmInterface(DiagnosticSettings),
	NewVarCheckRuleFromAvmInterface(Location),
	NewVarCheckRuleFromAvmInterface(Lock),
	NewVarCheckRuleFromAvmInterface(ManagedIdentities),
	NewVarCheckRuleFromAvmInterface(RoleAssignments),
	NewVarCheckRuleFromAvmInterface(Tags),
	NewVarCheckRuleFromAvmInterface(PrivateEndpoints),
}

package interfaces

import (
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var Rules = []tflint.Rule{
	NewVarCheckRuleFromAvmInterface(CustomerManagedKey),
	NewVarCheckRuleFromAvmInterface(Location),
	NewVarCheckRuleFromAvmInterface(Lock),
	NewVarCheckRuleFromAvmInterface(ManagedIdentities),
	NewVarCheckRuleFromAvmInterface(RoleAssignments),
	NewVarCheckRuleFromAvmInterface(Tags),
	NewVarCheckRuleFromAvmInterface(PrivateEndpoints),
    func() tflint.Rule {
		return common.NewEitherCheckRule("diagnostic_settings", true, tflint.ERROR,
			NewVarCheckRuleFromAvmInterface(DiagnosticSettings),
			NewVarCheckRuleFromAvmInterface(DiagnosticSettingsV2))
	}(),
}

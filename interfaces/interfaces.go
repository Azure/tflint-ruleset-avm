package interfaces

import (
	"github.com/Azure/tflint-ruleset-avm/common"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var Rules = []tflint.Rule{
	func() tflint.Rule {
		return common.NewEitherCheckRule("customer_managed_key", true, tflint.ERROR,
			NewVarCheckRuleFromAvmInterface(CustomerManagedKey),
			NewVarCheckRuleFromAvmInterface(CustomerManagedKeyV2))
	}(),
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

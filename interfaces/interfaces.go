package interfaces

import (
	"github.com/Azure/tflint-ruleset-avm/common"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var Rules = []tflint.Rule{
	NewVarCheckRuleFromAvmInterface(CustomerManagedKey),
	NewVarCheckRuleFromAvmInterface(DiagnosticSettings),
	NewVarCheckRuleFromAvmInterface(Lock),
	NewVarCheckRuleFromAvmInterface(ManagedIdentities),
	NewVarCheckRuleFromAvmInterface(RoleAssignments),
	NewVarCheckRuleFromAvmInterface(Tags),
	func() tflint.Rule {
		return common.NewEitherCheckRule("private_endpoints", true, tflint.ERROR,
			NewVarCheckRuleFromAvmInterface(PrivateEndpoints),
			NewVarCheckRuleFromAvmInterface(PrivateEndpointsWithSubresourceName))
	}(),
}

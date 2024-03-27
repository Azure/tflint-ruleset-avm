package interfaces

import (
	"github.com/Azure/tflint-ruleset-avm/common"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var Rules = []tflint.Rule{
	NewVarCheckRuleFromAvmInterface(Lock),
	NewVarCheckRuleFromAvmInterface(DiagnosticSettings),
	NewVarCheckRuleFromAvmInterface(ManagedIdentities),
	NewVarCheckRuleFromAvmInterface(RoleAssignments),
	NewVarCheckRuleFromAvmInterface(CustomerManagedKey),
	func() tflint.Rule {
		return common.NewEitherCheckRule("private_endpoints", true, tflint.ERROR,
			NewVarCheckRuleFromAvmInterface(PrivateEndpoints),
			NewVarCheckRuleFromAvmInterface(PrivateEndpointsWithSubresourceName))
	}(),
}

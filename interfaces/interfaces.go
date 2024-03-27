package interfaces

import (
	"github.com/Azure/tflint-ruleset-avm/common"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var privateEndpointsWithoutSubresourceNameRule = NewVarCheckRuleFromAvmInterface(PrivateEndpoints)
var privateEndpointsWithSubresourceNameRule = NewVarCheckRuleFromAvmInterface(PrivateEndpointsWithSubresourceName)
var PrivateEndpointsRule = common.NewEitherCheckRule("private_endpoints", true, tflint.ERROR,
	privateEndpointsWithoutSubresourceNameRule,
	privateEndpointsWithSubresourceNameRule)

var Rules = []tflint.Rule{
	NewVarCheckRuleFromAvmInterface(Lock),
	NewVarCheckRuleFromAvmInterface(DiagnosticSettings),
	NewVarCheckRuleFromAvmInterface(ManagedIdentities),
	NewVarCheckRuleFromAvmInterface(RoleAssignments),
	NewVarCheckRuleFromAvmInterface(CustomerManagedKey),
	PrivateEndpointsRule,
}

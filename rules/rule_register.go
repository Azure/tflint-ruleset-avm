package rules

import (
	"github.com/Azure/tflint-ruleset-avm/interfaces"
	azurerm "github.com/Azure/tflint-ruleset-azurerm-ext/rules"
	basic "github.com/Azure/tflint-ruleset-basic-ext/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var privateEndpointsWithoutSubresourceNameRule = NewVarCheckRuleFromAvmInterface(interfaces.PrivateEndpoints)
var privateEndpointsWithSubresourceNameRule = NewVarCheckRuleFromAvmInterface(interfaces.PrivateEndpointsWithSubresourceName)
var PrivateEndpointsRule = NewEitherCheckRule("private_endpoints", true, tflint.ERROR,
	privateEndpointsWithoutSubresourceNameRule,
	privateEndpointsWithSubresourceNameRule)

var Rules = func() []tflint.Rule {
	return []tflint.Rule{
		Wrap(basic.NewTerraformHeredocUsageRule()),
		Wrap(basic.NewTerraformModuleProviderDeclarationRule()),
		Wrap(basic.NewTerraformOutputSeparateRule()),
		Wrap(basic.NewTerraformRequiredProvidersDeclarationRule()),
		Wrap(basic.NewTerraformRequiredVersionDeclarationRule()),
		Wrap(basic.NewTerraformSensitiveVariableNoDefaultRule()),
		Wrap(basic.NewTerraformVariableNullableFalseRule()),
		Wrap(basic.NewTerraformVariableSeparateRule()),
		Wrap(azurerm.NewAzurermResourceTagRule()),
		NewTerraformDotTfRule(),
		NewVarCheckRuleFromAvmInterface(interfaces.Lock),
		NewVarCheckRuleFromAvmInterface(interfaces.DiagnosticSettings),
		NewVarCheckRuleFromAvmInterface(interfaces.ManagedIdentities),
		NewVarCheckRuleFromAvmInterface(interfaces.RoleAssignments),
		NewVarCheckRuleFromAvmInterface(interfaces.CustomerManagedKey),
		PrivateEndpointsRule,
	}
}()

type wrappedRule struct {
	tflint.Rule
}

func (*wrappedRule) Enabled() bool {
	return false
}

func Wrap(r tflint.Rule) tflint.Rule {
	return &wrappedRule{
		Rule: r,
	}
}

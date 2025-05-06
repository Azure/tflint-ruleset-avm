package rules

import (
	"slices"

	"github.com/Azure/tflint-ruleset-avm/interfaces"
	"github.com/Azure/tflint-ruleset-avm/outputs"
	azurerm "github.com/Azure/tflint-ruleset-azurerm-ext/rules"
	basic "github.com/Azure/tflint-ruleset-basic-ext/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var Rules = func() []tflint.Rule {
	return slices.Concat(
		[]tflint.Rule{
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
			NewModuleSourceRule(),
			NewNoDoubleQuotesInIgnoreChangesRule(),
			NewProviderVersionRule("modtm", "Azure/modtm", "0.3.0", "~> 0.3", true),
			NewProviderVersionRule("azapi", "Azure/azapi", "2.999.0", "~> 2.0", false),
			NewProviderVersionRule("azurerm", "hashicorp/azurerm", "4.999.0", "~> 4.0", false),
		},
		interfaces.Rules,
		outputs.Rules,
	)
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

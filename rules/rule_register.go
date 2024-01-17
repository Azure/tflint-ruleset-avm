package rules

import (
	azurerm "github.com/Azure/tflint-ruleset-azurerm-ext/rules"
	basic "github.com/Azure/tflint-ruleset-basic-ext/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	terraform "github.com/terraform-linters/tflint-ruleset-terraform/rules"
)

var Rules = []tflint.Rule{
	Wrap(terraform.NewTerraformCommentSyntaxRule()),
	Wrap(terraform.NewTerraformDeprecatedIndexRule()),
	Wrap(terraform.NewTerraformDeprecatedInterpolationRule()),
	Wrap(terraform.NewTerraformDeprecatedLookupRule()),
	Wrap(terraform.NewTerraformDocumentedOutputsRule()),
	Wrap(terraform.NewTerraformDocumentedVariablesRule()),
	Wrap(terraform.NewTerraformEmptyListEqualityRule()),
	Wrap(terraform.NewTerraformModulePinnedSourceRule()),
	Wrap(terraform.NewTerraformModuleVersionRule()),
	Wrap(terraform.NewTerraformNamingConventionRule()),
	Wrap(terraform.NewTerraformRequiredProvidersRule()),
	Wrap(terraform.NewTerraformRequiredVersionRule()),
	Wrap(terraform.NewTerraformStandardModuleStructureRule()),
	Wrap(terraform.NewTerraformTypedVariablesRule()),
	Wrap(terraform.NewTerraformUnusedDeclarationsRule()),
	Wrap(terraform.NewTerraformUnusedRequiredProvidersRule()),
	Wrap(terraform.NewTerraformWorkspaceRemoteRule()),

	Wrap(basic.NewTerraformHeredocUsageRule()),
	Wrap(basic.NewTerraformModuleProviderDeclarationRule()),
	Wrap(basic.NewTerraformOutputSeparateRule()),
	Wrap(basic.NewTerraformRequiredProvidersDeclarationRule()),
	Wrap(basic.NewTerraformRequiredVersionDeclarationRule()),
	Wrap(basic.NewTerraformSensitiveVariableNoDefaultRule()),
	Wrap(basic.NewTerraformVariableNullableFalseRule()),
	Wrap(basic.NewTerraformVariableSeparateRule()),

	Wrap(azurerm.NewAzurermResourceTagRule()),
}

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

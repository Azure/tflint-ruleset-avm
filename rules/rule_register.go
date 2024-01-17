package rules

import (
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

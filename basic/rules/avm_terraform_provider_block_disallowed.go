package rules

import (
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var _ tflint.Rule = &TerraformModuleProviderDeclarationRule{}

// TerraformModuleProviderDeclarationRule checks that modules do not declare provider blocks.
type TerraformModuleProviderDeclarationRule struct {
	tflint.DefaultRule
}

// NewTerraformModuleProviderDeclarationRule returns a new rule
func NewTerraformModuleProviderDeclarationRule() *TerraformModuleProviderDeclarationRule {
	return &TerraformModuleProviderDeclarationRule{}
}

func (r *TerraformModuleProviderDeclarationRule) Enabled() bool {
	return true
}

func (r *TerraformModuleProviderDeclarationRule) Check(runner tflint.Runner) error {
	return ForFiles(runner, r.CheckFile)
}

// Name returns the rule name
func (r *TerraformModuleProviderDeclarationRule) Name() string {
	return "avm_terraform_provider_block_disallowed"
}

// Severity returns the rule severity
func (r *TerraformModuleProviderDeclarationRule) Severity() tflint.Severity {
	return tflint.WARNING
}

func (r *TerraformModuleProviderDeclarationRule) Link() string {
	return "https://azure.github.io/Azure-Verified-Modules/spec/TFNFR27/"
}

func (r *TerraformModuleProviderDeclarationRule) CheckFile(runner tflint.Runner, file *hcl.File) error {
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		logger.Debug("skip terraform_module_provider_declaration since it's not hcl file")
		return nil
	}
	var err error
	for _, block := range body.Blocks {
		if block.Type != "provider" {
			continue
		}
		subErr := runner.EmitIssue(
			r,
			"Provider blocks must not be declared in Terraform modules; declare aliases with `configuration_aliases` in `required_providers`",
			block.DefRange(),
		)
		if subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}

package rules

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var _ tflint.Rule = new(TerraformDeprecatedVariableFileRule)

const deprecatedVariablesFile = "deprecated_variables.tf"

// TerraformDeprecatedVariableFileRule checks that deprecated variables live in `deprecated_variables.tf`,
// and that every variable in that file is annotated as `DEPRECATED`.
type TerraformDeprecatedVariableFileRule struct {
	tflint.DefaultRule
}

// NewTerraformDeprecatedVariableFileRule returns a new rule.
func NewTerraformDeprecatedVariableFileRule() *TerraformDeprecatedVariableFileRule {
	return new(TerraformDeprecatedVariableFileRule)
}

func (r *TerraformDeprecatedVariableFileRule) Name() string {
	return "avm_terraform_deprecated_variable_file_required"
}

func (r *TerraformDeprecatedVariableFileRule) Enabled() bool {
	return true
}

func (r *TerraformDeprecatedVariableFileRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *TerraformDeprecatedVariableFileRule) Link() string {
	return "https://azure.github.io/Azure-Verified-Modules/spec/TFNFR24/"
}

func (r *TerraformDeprecatedVariableFileRule) Check(runner tflint.Runner) error {
	return ForFiles(runner, r.CheckFile)
}

func (r *TerraformDeprecatedVariableFileRule) CheckFile(runner tflint.Runner, file *hcl.File) error {
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		logger.Debug("skip avm_terraform_deprecated_variable_file_required since it is not an HCL file")
		return nil
	}
	inDeprecatedFile := fileBaseName(body) == deprecatedVariablesFile
	var err error
	for _, block := range body.Blocks {
		if block.Type != "variable" || len(block.Labels) == 0 {
			continue
		}
		deprecated := hasDeprecatedDescription(block)
		var msg string
		switch {
		case deprecated && !inDeprecatedFile:
			msg = fmt.Sprintf("Deprecated variable `%s` must be declared in `%s`", block.Labels[0], deprecatedVariablesFile)
		case !deprecated && inDeprecatedFile:
			msg = fmt.Sprintf("Variable `%s` in `%s` must have a description starting with `DEPRECATED`", block.Labels[0], deprecatedVariablesFile)
		default:
			continue
		}
		if subErr := runner.EmitIssue(r, msg, block.DefRange()); subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}

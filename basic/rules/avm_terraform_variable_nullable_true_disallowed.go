package rules

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var _ tflint.Rule = new(TerraformVariableNullableTrueRule)

// TerraformVariableNullableTrueRule checks that variables do not declare `nullable = true`.
type TerraformVariableNullableTrueRule struct {
	tflint.DefaultRule
}

// NewTerraformVariableNullableTrueRule returns a new rule.
func NewTerraformVariableNullableTrueRule() *TerraformVariableNullableTrueRule {
	return new(TerraformVariableNullableTrueRule)
}

func (r *TerraformVariableNullableTrueRule) Name() string {
	return "avm_terraform_variable_nullable_true_disallowed"
}

func (r *TerraformVariableNullableTrueRule) Enabled() bool {
	return true
}

func (r *TerraformVariableNullableTrueRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *TerraformVariableNullableTrueRule) Link() string {
	return "https://azure.github.io/Azure-Verified-Modules/spec/TFNFR21/"
}

func (r *TerraformVariableNullableTrueRule) Check(runner tflint.Runner) error {
	return ForFiles(runner, r.CheckFile)
}

func (r *TerraformVariableNullableTrueRule) CheckFile(runner tflint.Runner, file *hcl.File) error {
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		logger.Debug("skip avm_terraform_variable_nullable_true_disallowed since it is not an HCL file")
		return nil
	}
	var err error
	for _, block := range body.Blocks {
		if block.Type != "variable" || len(block.Labels) == 0 {
			continue
		}
		nullable, ok := staticBool(block, "nullable")
		if !ok || !nullable {
			continue
		}
		subErr := runner.EmitIssue(
			r,
			fmt.Sprintf("`nullable = true` must not be declared on variable `%s`", block.Labels[0]),
			block.Body.Attributes["nullable"].SrcRange,
		)
		if subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}

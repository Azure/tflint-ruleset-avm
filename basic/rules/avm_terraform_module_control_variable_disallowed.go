package rules

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var _ tflint.Rule = new(TerraformModuleControlVariableRule)

// moduleControlVariableNames lists variable names that control the entire module's operation.
var moduleControlVariableNames = map[string]bool{
	"enabled":           true,
	"module_depends_on": true,
}

// TerraformModuleControlVariableRule checks that modules do not declare variables that control the entire module's operation.
type TerraformModuleControlVariableRule struct {
	tflint.DefaultRule
}

// NewTerraformModuleControlVariableRule returns a new rule.
func NewTerraformModuleControlVariableRule() *TerraformModuleControlVariableRule {
	return new(TerraformModuleControlVariableRule)
}

func (r *TerraformModuleControlVariableRule) Name() string {
	return "avm_terraform_module_control_variable_disallowed"
}

func (r *TerraformModuleControlVariableRule) Enabled() bool {
	return true
}

func (r *TerraformModuleControlVariableRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *TerraformModuleControlVariableRule) Link() string {
	return "https://azure.github.io/Azure-Verified-Modules/spec/TFNFR14/"
}

func (r *TerraformModuleControlVariableRule) Check(runner tflint.Runner) error {
	return ForFiles(runner, r.CheckFile)
}

func (r *TerraformModuleControlVariableRule) CheckFile(runner tflint.Runner, file *hcl.File) error {
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		logger.Debug("skip avm_terraform_module_control_variable_disallowed since it is not an HCL file")
		return nil
	}
	var err error
	for _, block := range body.Blocks {
		if block.Type != "variable" || len(block.Labels) == 0 || !moduleControlVariableNames[block.Labels[0]] {
			continue
		}
		subErr := runner.EmitIssue(
			r,
			fmt.Sprintf("Variable `%s` must not be declared to control the entire module's operation; use feature-specific toggles instead", block.Labels[0]),
			block.DefRange(),
		)
		if subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}

package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var _ tflint.Rule = new(TerraformVariableNegativeToggleRule)

// TerraformVariableNegativeToggleRule checks that boolean feature switches use a positive name such as `xxx_enabled`.
type TerraformVariableNegativeToggleRule struct {
	tflint.DefaultRule
}

// NewTerraformVariableNegativeToggleRule returns a new rule.
func NewTerraformVariableNegativeToggleRule() *TerraformVariableNegativeToggleRule {
	return new(TerraformVariableNegativeToggleRule)
}

func (r *TerraformVariableNegativeToggleRule) Name() string {
	return "avm_terraform_variable_negative_toggle_disallowed"
}

func (r *TerraformVariableNegativeToggleRule) Enabled() bool {
	return true
}

func (r *TerraformVariableNegativeToggleRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

func (r *TerraformVariableNegativeToggleRule) Link() string {
	return "https://azure.github.io/Azure-Verified-Modules/spec/TFNFR16/"
}

func (r *TerraformVariableNegativeToggleRule) Check(runner tflint.Runner) error {
	return ForFiles(runner, r.CheckFile)
}

func (r *TerraformVariableNegativeToggleRule) CheckFile(runner tflint.Runner, file *hcl.File) error {
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		logger.Debug("skip avm_terraform_variable_negative_toggle_disallowed since it is not an HCL file")
		return nil
	}
	var err error
	for _, block := range body.Blocks {
		if block.Type != "variable" || len(block.Labels) == 0 {
			continue
		}
		name := block.Labels[0]
		if !strings.HasPrefix(name, "disable_") && !strings.HasSuffix(name, "_disabled") {
			continue
		}
		// Only boolean variables are feature switches; other types may mirror resource arguments.
		typeAttr, typed := block.Body.Attributes["type"]
		if !typed || hcl.ExprAsKeyword(typeAttr.Expr) != "bool" {
			continue
		}
		subErr := runner.EmitIssue(
			r,
			fmt.Sprintf("Feature switch variable `%s` should use a positive name such as `xxx_enabled` instead of `xxx_disabled`", name),
			block.DefRange(),
		)
		if subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}

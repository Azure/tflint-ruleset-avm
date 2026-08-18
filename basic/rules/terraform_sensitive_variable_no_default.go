package rules

import (
	"fmt"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var _ tflint.Rule = &TerraformSensitiveVariableNoDefaultRule{}

// TerraformSensitiveVariableNoDefaultRule checks whether default value is set for sensitive variables
type TerraformSensitiveVariableNoDefaultRule struct {
	tflint.DefaultRule
}

// NewTerraformSensitiveVariableNoDefaultRule returns a new rule
func NewTerraformSensitiveVariableNoDefaultRule() *TerraformSensitiveVariableNoDefaultRule {
	return &TerraformSensitiveVariableNoDefaultRule{}
}

func (r *TerraformSensitiveVariableNoDefaultRule) Enabled() bool {
	return false
}

func (r *TerraformSensitiveVariableNoDefaultRule) Check(runner tflint.Runner) error {
	return ForFiles(runner, r.CheckFile)
}

// Name returns the rule name
func (r *TerraformSensitiveVariableNoDefaultRule) Name() string {
	return "terraform_sensitive_variable_no_default"
}

// Severity returns the rule severity
func (r *TerraformSensitiveVariableNoDefaultRule) Severity() tflint.Severity {
	return tflint.WARNING
}

func (r *TerraformSensitiveVariableNoDefaultRule) CheckFile(runner tflint.Runner, file *hcl.File) error {
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		logger.Debug("skip terraform_sensitive_variable_no_default check since it's not hcl file")
		return nil
	}
	blocks := body.Blocks
	var err error
	for _, block := range blocks {
		if block.Type != "variable" {
			continue
		}
		sensitive := false
		if attr, sensitiveSet := block.Body.Attributes["sensitive"]; sensitiveSet {
			val, diags := attr.Expr.Value(nil)
			if diags.HasErrors() {
				err = multierror.Append(err, diags)
			}
			sensitive = val.True()
		}
		if !sensitive {
			continue
		}
		nullOrEmpty, evalDefaultErr := nullOrZeroDefaultValue(block)
		if evalDefaultErr != nil {
			return evalDefaultErr
		}
		if !nullOrEmpty {
			subErr := runner.EmitIssue(
				r,
				fmt.Sprintf("Default value is not expected to be set for sensitive variable `%s`", block.Labels[0]),
				block.Body.Attributes["default"].NameRange,
			)
			if subErr != nil {
				err = multierror.Append(err, subErr)
			}
		}
	}
	return err
}

func nullOrZeroDefaultValue(b *hclsyntax.Block) (bool, error) {
	attr, set := b.Body.Attributes["default"]
	if !set {
		return true, nil
	}
	v, diag := attr.Expr.Value(&hcl.EvalContext{})
	if diag.HasErrors() {
		return false, diag
	}
	return v.Equals(cty.NullVal(cty.DynamicPseudoType)).True() ||
		(v.CanIterateElements() && v.LengthInt() == 0), nil
}

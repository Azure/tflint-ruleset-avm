package rules

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var _ tflint.Rule = new(TerraformSensitiveFalseRule)

// TerraformSensitiveFalseRule checks that variables and outputs do not declare `sensitive = false`.
type TerraformSensitiveFalseRule struct {
	tflint.DefaultRule
}

// NewTerraformSensitiveFalseRule returns a new rule.
func NewTerraformSensitiveFalseRule() *TerraformSensitiveFalseRule {
	return new(TerraformSensitiveFalseRule)
}

func (r *TerraformSensitiveFalseRule) Name() string {
	return "avm_terraform_sensitive_false_disallowed"
}

func (r *TerraformSensitiveFalseRule) Enabled() bool {
	return true
}

func (r *TerraformSensitiveFalseRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *TerraformSensitiveFalseRule) Link() string {
	return "https://azure.github.io/Azure-Verified-Modules/spec/TFNFR22/"
}

func (r *TerraformSensitiveFalseRule) Check(runner tflint.Runner) error {
	return ForFiles(runner, r.CheckFile)
}

func (r *TerraformSensitiveFalseRule) CheckFile(runner tflint.Runner, file *hcl.File) error {
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		logger.Debug("skip avm_terraform_sensitive_false_disallowed since it is not an HCL file")
		return nil
	}
	var err error
	for _, block := range body.Blocks {
		if (block.Type != "variable" && block.Type != "output") || len(block.Labels) == 0 {
			continue
		}
		sensitive, ok := staticBool(block, "sensitive")
		if !ok || sensitive {
			continue
		}
		subErr := runner.EmitIssue(
			r,
			fmt.Sprintf("`sensitive = false` must not be declared on %s `%s`; values are not sensitive by default", block.Type, block.Labels[0]),
			block.Body.Attributes["sensitive"].SrcRange,
		)
		if subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}

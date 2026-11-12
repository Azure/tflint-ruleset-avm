package rules

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var _ tflint.Rule = new(TerraformDeprecatedOutputFileRule)

const deprecatedOutputsFile = "deprecated_outputs.tf"

// TerraformDeprecatedOutputFileRule checks that outputs annotated as `DEPRECATED` live in `deprecated_outputs.tf`.
type TerraformDeprecatedOutputFileRule struct {
	tflint.DefaultRule
}

// NewTerraformDeprecatedOutputFileRule returns a new rule.
func NewTerraformDeprecatedOutputFileRule() *TerraformDeprecatedOutputFileRule {
	return new(TerraformDeprecatedOutputFileRule)
}

func (r *TerraformDeprecatedOutputFileRule) Name() string {
	return "avm_terraform_deprecated_output_file_required"
}

func (r *TerraformDeprecatedOutputFileRule) Enabled() bool {
	return true
}

func (r *TerraformDeprecatedOutputFileRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *TerraformDeprecatedOutputFileRule) Link() string {
	return "https://azure.github.io/Azure-Verified-Modules/spec/TFNFR30/"
}

func (r *TerraformDeprecatedOutputFileRule) Check(runner tflint.Runner) error {
	return ForFiles(runner, r.CheckFile)
}

func (r *TerraformDeprecatedOutputFileRule) CheckFile(runner tflint.Runner, file *hcl.File) error {
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		logger.Debug("skip avm_terraform_deprecated_output_file_required since it is not an HCL file")
		return nil
	}
	if fileBaseName(body) == deprecatedOutputsFile {
		return nil
	}
	var err error
	for _, block := range body.Blocks {
		if block.Type != "output" || len(block.Labels) == 0 || !hasDeprecatedDescription(block) {
			continue
		}
		subErr := runner.EmitIssue(
			r,
			fmt.Sprintf("Deprecated output `%s` must be declared in `%s`", block.Labels[0], deprecatedOutputsFile),
			block.DefRange(),
		)
		if subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}

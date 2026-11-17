package rules

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var _ tflint.Rule = new(TerraformLocalsFileRule)

// TerraformLocalsFileRule checks that `locals.tf` only contains `locals` blocks.
type TerraformLocalsFileRule struct {
	tflint.DefaultRule
}

// NewTerraformLocalsFileRule returns a new rule.
func NewTerraformLocalsFileRule() *TerraformLocalsFileRule {
	return new(TerraformLocalsFileRule)
}

func (r *TerraformLocalsFileRule) Name() string {
	return "avm_terraform_locals_file_non_locals_disallowed"
}

func (r *TerraformLocalsFileRule) Enabled() bool {
	return true
}

func (r *TerraformLocalsFileRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *TerraformLocalsFileRule) Link() string {
	return "https://azure.github.io/Azure-Verified-Modules/spec/TFNFR31/"
}

func (r *TerraformLocalsFileRule) Check(runner tflint.Runner) error {
	return ForFiles(runner, r.CheckFile)
}

func (r *TerraformLocalsFileRule) CheckFile(runner tflint.Runner, file *hcl.File) error {
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		logger.Debug("skip avm_terraform_locals_file_non_locals_disallowed since it is not an HCL file")
		return nil
	}
	if fileBaseName(body) != "locals.tf" {
		return nil
	}
	var err error
	for _, block := range body.Blocks {
		if block.Type == "locals" {
			continue
		}
		subErr := runner.EmitIssue(
			r,
			fmt.Sprintf("`locals.tf` must only contain `locals` blocks, found `%s` block", block.Type),
			block.DefRange(),
		)
		if subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}

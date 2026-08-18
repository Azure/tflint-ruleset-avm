package rules

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var _ tflint.Rule = &TerraformVersionsFileRule{}

// TerraformVersionsFileRule checks whether `versions.tf` only has 1 `terraform` block
type TerraformVersionsFileRule struct {
	tflint.DefaultRule
}

// NewTerraformVersionsFileRule returns a new rule
func NewTerraformVersionsFileRule() *TerraformVersionsFileRule {
	return &TerraformVersionsFileRule{}
}

func (r *TerraformVersionsFileRule) Enabled() bool {
	return false
}

func (r *TerraformVersionsFileRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

func (r *TerraformVersionsFileRule) Check(runner tflint.Runner) error {
	return ForFiles(runner, r.CheckFile)
}

// Name returns the rule name
func (r *TerraformVersionsFileRule) Name() string {
	return "terraform_versions_file"
}

func (r *TerraformVersionsFileRule) CheckFile(runner tflint.Runner, file *hcl.File) error {
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		logger.Debug("skip terraform_versions_file check since it's not hcl file")
		return nil
	}
	filename := body.Range().Filename
	if filename != "versions.tf" {
		return nil
	}
	blocks := body.Blocks
	if len(blocks) != 1 || blocks[0].Type != "terraform" {
		return runner.EmitIssue(
			r,
			"`versions.tf` should have and only have 1 `terraform` block",
			hcl.Range{},
		)
	}
	return nil
}

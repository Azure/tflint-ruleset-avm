package rules

import (
	"github.com/Azure/tflint-ruleset-avm/basic/project"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TerraformOutputSeparateRule checks whether the outputs are separated from other types of blocks
type TerraformOutputSeparateRule struct {
	tflint.DefaultRule
}

// NewTerraformOutputSeparateRule returns a new rule
func NewTerraformOutputSeparateRule() *TerraformOutputSeparateRule {
	return &TerraformOutputSeparateRule{}
}

// Name returns the rule name
func (r *TerraformOutputSeparateRule) Name() string {
	return "terraform_output_separate"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformOutputSeparateRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *TerraformOutputSeparateRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *TerraformOutputSeparateRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether the variables are separated from other types of blocks
func (r *TerraformOutputSeparateRule) Check(runner tflint.Runner) error {

	files, err := runner.GetFiles()
	if err != nil {
		return err
	}
	for _, file := range files {
		if subErr := r.checkOutputSeparate(runner, file); subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}

func (r *TerraformOutputSeparateRule) checkOutputSeparate(runner tflint.Runner, file *hcl.File) error {
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		logger.Debug("skip terraform_output_separate check since it's not hcl file")
		return nil
	}
	blocks := body.Blocks
	var firstNonOutputBlockRange *hcl.Range
	outputDefined := false
	for _, block := range blocks {
		switch block.Type {
		case "output":
			{
				outputDefined = true
			}
		default:
			if firstNonOutputBlockRange == nil {
				firstNonOutputBlockRange = ref(block.DefRange())
			}
		}
	}

	if outputDefined && firstNonOutputBlockRange != nil {
		return runner.EmitIssue(
			r,
			"Putting outputs and other types of blocks in the same file is not recommended",
			*firstNonOutputBlockRange,
		)
	}
	return nil
}

package rules

import (
	"github.com/Azure/tflint-ruleset-avm/basic/project"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TerraformVariableSeparateRule checks whether the variables are separated from other types of blocks
type TerraformVariableSeparateRule struct {
	tflint.DefaultRule
}

// NewTerraformVariableSeparateRule returns a new rule
func NewTerraformVariableSeparateRule() *TerraformVariableSeparateRule {
	return &TerraformVariableSeparateRule{}
}

// Name returns the rule name
func (r *TerraformVariableSeparateRule) Name() string {
	return "terraform_variable_separate"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformVariableSeparateRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *TerraformVariableSeparateRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *TerraformVariableSeparateRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether the variables are separated from other types of blocks
func (r *TerraformVariableSeparateRule) Check(runner tflint.Runner) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}
	for _, file := range files {
		if subErr := r.checkVariableSeparate(runner, file); subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}

func (r *TerraformVariableSeparateRule) checkVariableSeparate(runner tflint.Runner, file *hcl.File) error {
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		logger.Debug("skip terraform_variable_separate check since it's not hcl file")
		return nil
	}
	blocks := body.Blocks

	var firstNonVarBlockRange *hcl.Range
	variableDefined := false
	for _, block := range blocks {
		switch block.Type {
		case "variable":
			if !variableDefined {
				variableDefined = true
			}
		default:
			if firstNonVarBlockRange == nil {
				firstNonVarBlockRange = ref(block.DefRange())
			}
		}
	}

	if variableDefined && firstNonVarBlockRange != nil {
		return runner.EmitIssue(
			r,
			"Putting variables and other types of blocks in the same file is not recommended",
			*firstNonVarBlockRange,
		)
	}
	return nil
}

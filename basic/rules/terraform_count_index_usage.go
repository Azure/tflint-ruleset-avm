package rules

import (
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var _ tflint.Rule = &TerraformCountIndexUsageRule{}

// TerraformCountIndexUsageRule checks whether count.index is used as subscript of list/map
type TerraformCountIndexUsageRule struct {
	tflint.DefaultRule
}

func (r *TerraformCountIndexUsageRule) Enabled() bool {
	return false
}

// NewTerraformCountIndexUsageRule returns a new rule
func NewTerraformCountIndexUsageRule() *TerraformCountIndexUsageRule {
	return &TerraformCountIndexUsageRule{}
}

// Name returns the rule name
func (r *TerraformCountIndexUsageRule) Name() string {
	return "terraform_count_index_usage"
}

// Severity returns the rule severity
func (r *TerraformCountIndexUsageRule) Severity() tflint.Severity {
	return tflint.WARNING
}

func (r *TerraformCountIndexUsageRule) Check(runner tflint.Runner) error {
	return ForFiles(runner, r.CheckFile)
}

func (r *TerraformCountIndexUsageRule) CheckFile(runner tflint.Runner, file *hcl.File) error {
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		logger.Debug("skip terraform_count_index_usage since it's not hcl file")
		return nil
	}
	blocks := body.Blocks
	var err error
	for _, block := range blocks {
		if subErr := r.visitBlock(runner, block); subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}

func (r *TerraformCountIndexUsageRule) visitBlock(runner tflint.Runner, block *hclsyntax.Block) error {
	var err error
	for _, attr := range block.Body.Attributes {
		if subErr := r.visitExp(runner, attr.Expr); subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	for _, nestedBlock := range block.Body.Blocks {
		if subErr := r.visitBlock(runner, nestedBlock); subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}

func (r *TerraformCountIndexUsageRule) visitExp(runner tflint.Runner, exp hclsyntax.Expression) error {
	file, _ := runner.GetFile(exp.Range().Filename)
	tokens, diags := hclsyntax.LexExpression(exp.Range().SliceBytes(file.Bytes),
		exp.Range().Filename,
		exp.StartRange().Start)
	if diags.HasErrors() {
		return diags
	}
	var err error
	depth := 0
	for i, token := range tokens {
		switch token.Type {
		case hclsyntax.TokenOBrack:
			depth++
		case hclsyntax.TokenCBrack:
			depth--
		case hclsyntax.TokenIdent:
			if depth == 0 || i+2 >= len(tokens) {
				continue
			}
			first, second, third := string(token.Bytes), string(tokens[i+1].Bytes), string(tokens[i+2].Bytes)
			if first == "count" && second == "." && third == "index" {
				subErr := runner.EmitIssue(
					r,
					"`count.index` is not recommended to be used as the subscript of list/map, use for_each instead",
					hcl.Range{
						Filename: token.Range.Filename,
						Start:    token.Range.Start,
						End:      tokens[i+2].Range.End,
					},
				)
				if subErr != nil {
					err = multierror.Append(err, subErr)
				}
			}
		}
	}
	return err
}

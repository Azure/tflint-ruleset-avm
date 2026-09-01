package rules

import (
	"encoding/json"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"gopkg.in/yaml.v3"
)

var _ tflint.Rule = new(TerraformHeredocUsageRule)

// TerraformHeredocUsageRule checks whether literal JSON or YAML uses a heredoc.
type TerraformHeredocUsageRule struct {
	tflint.DefaultRule
}

// NewTerraformHeredocUsageRule returns a heredoc usage rule.
func NewTerraformHeredocUsageRule() *TerraformHeredocUsageRule {
	return new(TerraformHeredocUsageRule)
}

func (r *TerraformHeredocUsageRule) Name() string {
	return "avm_terraform_literal_heredoc_disallowed"
}

func (r *TerraformHeredocUsageRule) Enabled() bool {
	return true
}

func (r *TerraformHeredocUsageRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

func (r *TerraformHeredocUsageRule) Link() string {
	return "https://aka.ms/avm/spec/TFNFR40"
}

func (r *TerraformHeredocUsageRule) Check(runner tflint.Runner) error {
	return ForFiles(runner, r.CheckFile)
}

func (r *TerraformHeredocUsageRule) CheckFile(runner tflint.Runner, file *hcl.File) error {
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		logger.Debug("skip terraform_heredoc_usage since it is not an HCL file")
		return nil
	}
	tokens, diags := hclsyntax.LexConfig(file.Bytes, body.Range().Filename, hcl.InitialPos)
	if diags.HasErrors() {
		return diags
	}

	var err error
	var start hcl.Range
	var heredoc strings.Builder
	inHeredoc := false
	for _, token := range tokens {
		switch token.Type {
		case hclsyntax.TokenOHeredoc:
			inHeredoc = true
			heredoc.Reset()
			start = token.Range
		case hclsyntax.TokenCHeredoc:
			inHeredoc = false
			if subErr := r.checkLiteral(runner, heredoc.String(), start); subErr != nil {
				err = multierror.Append(err, subErr)
			}
		case hclsyntax.TokenStringLit:
			if inHeredoc {
				heredoc.Write(token.Bytes)
			}
		}
	}
	return err
}

func (r *TerraformHeredocUsageRule) checkLiteral(runner tflint.Runner, heredoc string, sourceRange hcl.Range) error {
	if strings.TrimSpace(heredoc) == "" {
		return nil
	}
	if json.Valid([]byte(heredoc)) {
		return runner.EmitIssue(
			r,
			"Use `jsonencode` instead of a heredoc for literal JSON",
			sourceRange,
		)
	}
	var value any
	if yaml.Unmarshal([]byte(heredoc), &value) == nil {
		if _, isMapping := value.(map[string]any); isMapping {
			return runner.EmitIssue(
				r,
				"Use `yamlencode` instead of a heredoc for literal YAML",
				sourceRange,
			)
		}
	}
	return nil
}

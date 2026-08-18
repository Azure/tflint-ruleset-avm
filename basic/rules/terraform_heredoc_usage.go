package rules

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"gopkg.in/yaml.v3"
	"strings"
)

var _ tflint.Rule = &TerraformHeredocUsageRule{}

// TerraformHeredocUsageRule checks whether HEREDOC is used for JSON/YAML
type TerraformHeredocUsageRule struct {
	tflint.DefaultRule
}

// NewTerraformHeredocUsageRule returns a new rule
func NewTerraformHeredocUsageRule() *TerraformHeredocUsageRule {
	return &TerraformHeredocUsageRule{}
}

func (r *TerraformHeredocUsageRule) Enabled() bool {
	return false
}

func (r *TerraformHeredocUsageRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

func (r *TerraformHeredocUsageRule) Check(runner tflint.Runner) error {
	return ForFiles(runner, r.CheckFile)
}

// Name returns the rule name
func (r *TerraformHeredocUsageRule) Name() string {
	return "terraform_heredoc_usage"
}

func (r *TerraformHeredocUsageRule) CheckFile(runner tflint.Runner, file *hcl.File) error {
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		logger.Debug("skip terraform_heredoc_usage since it's not hcl file")
		return nil
	}
	fileName := body.Range().Filename
	tokens, diags := hclsyntax.LexConfig(file.Bytes, fileName, hcl.InitialPos)
	if diags.HasErrors() {
		return diags
	}
	var err error
	var hereDocStartRange hcl.Range
	var heredoc string
	var inHeredoc bool
	for _, token := range tokens {
		switch token.Type {
		case hclsyntax.TokenOHeredoc:
			inHeredoc = true
			heredoc = ""
			hereDocStartRange = token.Range
		case hclsyntax.TokenCHeredoc:
			inHeredoc = false
			if subErr := r.checkHeredocIsJSONOrYAML(runner, heredoc, hereDocStartRange); subErr != nil {
				err = multierror.Append(err, subErr)
			}
		case hclsyntax.TokenStringLit:
			if inHeredoc {
				heredoc = fmt.Sprintf("%s%s", heredoc, string(token.Bytes))
			}
		}
	}
	return err
}

func (r *TerraformHeredocUsageRule) checkHeredocIsJSONOrYAML(runner tflint.Runner, heredoc string, heredocStartRange hcl.Range) error {
	prunedHereDoc := strings.ReplaceAll(heredoc, "\t", "")
	prunedHereDoc = strings.ReplaceAll(prunedHereDoc, " ", "")
	prunedHereDoc = strings.ReplaceAll(prunedHereDoc, "\n", "")
	if prunedHereDoc == "" {
		return nil
	}
	bytes := []byte(heredoc)
	if json.Valid(bytes) {
		return runner.EmitIssue(
			r,
			"for JSON, instead of HEREDOC, use a combination of a `local` and the `jsonencode` function",
			heredocStartRange,
		)
	}
	temp := map[string]interface{}{}
	if yaml.Unmarshal(bytes, &temp) == nil {
		return runner.EmitIssue(
			r,
			"for YAML, instead of HEREDOC, use a combination of a `local` and the `yamlencode` function",
			heredocStartRange,
		)
	}
	return nil
}

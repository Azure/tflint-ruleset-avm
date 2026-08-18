package rules

import (
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"strings"
)

var _ tflint.Rule = &TerraformRequiredVersionDeclarationRule{}

// TerraformRequiredVersionDeclarationRule checks whether required_version field is declared at the beginning of terraform block
type TerraformRequiredVersionDeclarationRule struct {
	tflint.DefaultRule
}

// NewTerraformRequiredVersionDeclarationRule returns a new rule
func NewTerraformRequiredVersionDeclarationRule() *TerraformRequiredVersionDeclarationRule {
	return &TerraformRequiredVersionDeclarationRule{}
}

func (r *TerraformRequiredVersionDeclarationRule) Enabled() bool {
	return false
}

func (r *TerraformRequiredVersionDeclarationRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

func (r *TerraformRequiredVersionDeclarationRule) Check(runner tflint.Runner) error {
	return ForFiles(runner, r.CheckFile)
}

// Name returns the rule name
func (r *TerraformRequiredVersionDeclarationRule) Name() string {
	return "terraform_required_version_declaration"
}

func (r *TerraformRequiredVersionDeclarationRule) CheckFile(runner tflint.Runner, file *hcl.File) error {
	var err error
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		logger.Debug("skip terraform_required_version_declaration check since it's not hcl file")
		return nil
	}
	filename := body.Range().Filename
	if isOverrideTfFile(filename) {
		logger.Debug("skip terraform_required_version_declaration check since it's override file")
		return nil
	}
	blocks := body.Blocks
	for _, block := range blocks {
		if block.Type != "terraform" {
			continue
		}
		if subErr := r.checkBlock(runner, block); subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}

func isOverrideTfFile(filename string) bool {
	return strings.HasSuffix(filename, "_override.tf") || filename == "override.tf"
}

func (r *TerraformRequiredVersionDeclarationRule) checkBlock(runner tflint.Runner, block *hclsyntax.Block) error {
	msg := "The `required_version` field should be declared at the beginning of `terraform` block"
	versionAttr, defined := block.Body.Attributes["required_version"]
	if !defined {
		return runner.EmitIssue(
			r,
			msg,
			block.DefRange(),
		)
	}
	comparePos := func(pos1 hcl.Pos, pos2 hcl.Pos) bool {
		if pos1.Line != pos2.Line {
			return pos1.Line < pos2.Line
		}
		return pos1.Column < pos2.Column
	}
	for _, attr := range block.Body.Attributes {
		if attr.Name != "required_version" && comparePos(attr.Range().Start, versionAttr.Range().Start) {
			return runner.EmitIssue(
				r,
				msg,
				versionAttr.NameRange,
			)
		}
	}
	for _, nestedBlock := range block.Body.Blocks {
		if comparePos(nestedBlock.Range().Start, versionAttr.Range().Start) {
			return runner.EmitIssue(
				r,
				msg,
				versionAttr.NameRange,
			)
		}
	}
	return nil
}

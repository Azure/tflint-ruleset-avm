package rules

import (
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var _ tflint.Rule = new(RequiredVersionRule)

type RequiredVersionRule struct {
	tflint.DefaultRule
}

func NewRequiredVersionRule() *RequiredVersionRule {
	return &RequiredVersionRule{}
}

func (t *RequiredVersionRule) Name() string {
	return "required_version_tfnfr25"
}

func (t *RequiredVersionRule) Link() string {
	return "https://azure.github.io/Azure-Verified-Modules/specs/terraform/#id-tfnfr25---category-code-style---verified-modules-requirements"
}

func (t *RequiredVersionRule) Enabled() bool {
	return true
}

func (t *RequiredVersionRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (t *RequiredVersionRule) Check(r tflint.Runner) error {
	tFile, err := r.GetFile("terraform.tf")
	if err != nil {
		return err
	}

	body, ok := tFile.Body.(*hclsyntax.Body)
	if !ok {
		return nil
	}

	for _, terraformBlock := range body.Blocks {
		if terraformBlock.Type != "terraform" {
			continue
		}

		requiredVersion, exists := terraformBlock.Body.Attributes["required_version"]
		if !exists {
			return r.EmitIssue(
				t,
				"The `required_version` field should be declared in the `terraform` block",
				terraformBlock.DefRange(),
			)
		}
		if err = r.EvaluateExpr(requiredVersion.Expr, t.versionHasLimitedMajorVersion(r, requiredVersion.NameRange),
			&tflint.EvaluateExprOption{ModuleCtx: tflint.RootModuleCtxType}); err != nil {
			return err
		}

		if !t.requiredVersionIsOnTheTop(terraformBlock) {
			return r.EmitIssue(
				t,
				"The `required_version` field should be declared at the beginning of `terraform` block",
				requiredVersion.NameRange,
			)
		}
	}

	return nil
}

func (t *RequiredVersionRule) requiredVersionIsOnTheTop(terraformBlock *hclsyntax.Block) bool {
	requiredVersion := terraformBlock.Body.Attributes["required_version"]
	comparePos := func(pos1 hcl.Pos, pos2 hcl.Pos) bool {
		if pos1.Line != pos2.Line {
			return pos1.Line < pos2.Line
		}
		return pos1.Column < pos2.Column
	}

	for _, attr := range terraformBlock.Body.Attributes {
		if attr.Name == "required_version" {
			continue
		}
		if comparePos(attr.Range().Start, requiredVersion.Range().Start) {
			return false
		}
	}

	for _, block := range terraformBlock.Body.Blocks {
		if comparePos(block.Range().Start, requiredVersion.Range().Start) {
			return false
		}
	}
	return true
}

func (t *RequiredVersionRule) versionHasLimitedMajorVersion(r tflint.Runner, issueRange hcl.Range) func(string) error {
	return func(requiredVersionVal string) error {
		if !strings.Contains(requiredVersionVal, "~>") && !(strings.Contains(requiredVersionVal, ">") && strings.Contains(requiredVersionVal, "<")) {
			return r.EmitIssue(
				t,
				"The `required_version` property constraint can use the ~> #.# or the >= #.#.#, < #.#.# format",
				issueRange,
			)
		}

		return nil
	}
}

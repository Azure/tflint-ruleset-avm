package requiredversion

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"strings"
)

var _ tflint.Rule = new(RequiredVersionRule)

type RequiredVersionRule struct {
	tflint.DefaultRule
	ruleName string
	link     string
}

func NewRequiredVersionRule(ruleName, link string) *RequiredVersionRule {
	return &RequiredVersionRule{
		ruleName: ruleName,
		link:     link,
	}
}

func (t *RequiredVersionRule) Name() string {
	return t.ruleName
}

func (t *RequiredVersionRule) Link() string {
	return t.link
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
		if terraformBlock.Type == "terraform" {

			if requiredVersion, exists := terraformBlock.Body.Attributes["required_version"]; exists {
				err := r.EvaluateExpr(requiredVersion.Expr, func(requiredVersionVal string) error {
					if !strings.Contains(requiredVersionVal, "~>") && !(strings.Contains(requiredVersionVal, ">=") && strings.Contains(requiredVersionVal, ",") && strings.Contains(requiredVersionVal, "<")) {
						return r.EmitIssue(
							t,
							"The `required_version` property constraint can use the ~> #.# or the >= #.#.#, < #.#.# format",
							requiredVersion.NameRange,
						)
					}

					comparePos := func(pos1 hcl.Pos, pos2 hcl.Pos) bool {
						if pos1.Line != pos2.Line {
							return pos1.Line < pos2.Line
						}
						return pos1.Column < pos2.Column
					}

					for _, attr := range terraformBlock.Body.Attributes {
						if attr.Name != "required_version" && comparePos(attr.Range().Start, requiredVersion.Range().Start) {
							return r.EmitIssue(
								t,
								"The `required_version` field should be declared at the beginning of `terraform` block",
								requiredVersion.NameRange,
							)
						}
					}

					for _, block := range terraformBlock.Body.Blocks {
						if comparePos(block.Range().Start, requiredVersion.Range().Start) {
							return r.EmitIssue(
								t,
								"The `required_version` field should be declared at the beginning of `terraform` block",
								requiredVersion.NameRange,
							)
						}
					}
					return nil
				}, &tflint.EvaluateExprOption{ModuleCtx: tflint.RootModuleCtxType})
				if err != nil {
					return err
				}
			} else {
				return r.EmitIssue(
					t,
					"The `required_version` field should be declared in the `terraform` block",
					terraformBlock.DefRange(),
				)
			}
		}
	}

	return nil
}

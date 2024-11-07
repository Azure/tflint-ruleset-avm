package rules

import (
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var noDoubleQuotesInIgnoreChangesBodySchema = &hclext.BodySchema{
	Blocks: []hclext.BlockSchema{
		{
			Type:       "resource",
			LabelNames: []string{"type", "name"},
			Body: &hclext.BodySchema{
				Blocks: []hclext.BlockSchema{
					{
						Type: "lifecycle",
						Body: &hclext.BodySchema{
							Attributes: []hclext.AttributeSchema{
								{Name: "ignore_changes"},
							},
						},
					},
				},
			},
		},
	},
}

var _ tflint.Rule = new(NoDoubleQuotesInIgnoreChangesRule)

type NoDoubleQuotesInIgnoreChangesRule struct {
	tflint.DefaultRule
}

func NewNoDoubleQuotesInIgnoreChangesRule() *NoDoubleQuotesInIgnoreChangesRule {
	return &NoDoubleQuotesInIgnoreChangesRule{}
}

func (t *NoDoubleQuotesInIgnoreChangesRule) Name() string {
	return "required_module_source_tfnfr10"
}

func (t *NoDoubleQuotesInIgnoreChangesRule) Link() string {
	return "https://azure.github.io/Azure-Verified-Modules/specs/terraform/#id-tfnfr10---category-code-style---no-double-quotes-in-ignore_changes"
}

func (t *NoDoubleQuotesInIgnoreChangesRule) Enabled() bool {
	return true
}

func (t *NoDoubleQuotesInIgnoreChangesRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (t *NoDoubleQuotesInIgnoreChangesRule) Check(r tflint.Runner) error {
	path, err := r.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() {
		return nil
	}

	body, err := r.GetModuleContent(
		noDoubleQuotesInIgnoreChangesBodySchema,
		&tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return err
	}

	var errList error
	for _, block := range body.Blocks {
		if block.Type != "resource" {
			continue
		}

		if subErr := t.checkBlock(r, block); subErr != nil {
			errList = multierror.Append(errList, subErr)
		}
	}

	return errList
}

func (t *NoDoubleQuotesInIgnoreChangesRule) checkBlock(r tflint.Runner, block *hclext.Block) error {
	for _, subBlock := range block.Body.Blocks {
		ignoreChanges, exists := subBlock.Body.Attributes["ignore_changes"]
		if !exists {
			return nil
		}

		ignoreChangesExpr, ok := ignoreChanges.Expr.(*hclsyntax.TupleConsExpr)
		if !ok {
			return nil
		}

		for _, itemExpr := range ignoreChangesExpr.Exprs {
			if _, ok := itemExpr.(*hclsyntax.ScopeTraversalExpr); !ok {
				return r.EmitIssue(
					t,
					"ignore_changes shouldn't include double quotes",
					ignoreChanges.Range,
				)
			}
		}
	}

	return nil
}

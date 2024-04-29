package rules

import (
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var nullComparisonToggleBodySchema = &hclext.BodySchema{
	Blocks: []hclext.BlockSchema{
		{
			Type:       "resource",
			LabelNames: []string{"type", "name"},
			Body: &hclext.BodySchema{
				Attributes: []hclext.AttributeSchema{
					{
						Name: "count",
					},
				},
			},
		},
	},
}

var _ tflint.Rule = new(NullComparisonToggleRule)

type NullComparisonToggleRule struct {
	tflint.DefaultRule
}

func NewNullComparisonToggleRule() *NullComparisonToggleRule {
	return &NullComparisonToggleRule{}
}

func (t *NullComparisonToggleRule) Name() string {
	return "null_comparison_toggle_tfnfr11"
}

func (t *NullComparisonToggleRule) Link() string {
	return "https://azure.github.io/Azure-Verified-Modules/specs/terraform/#id-tfnfr11---category-code-style---null-comparison-toggle"
}

func (t *NullComparisonToggleRule) Enabled() bool {
	return true
}

func (t *NullComparisonToggleRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (t *NullComparisonToggleRule) Check(r tflint.Runner) error {
	path, err := r.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() {
		return nil
	}

	body, err := r.GetModuleContent(
		nullComparisonToggleBodySchema,
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

func (t *NullComparisonToggleRule) checkBlock(r tflint.Runner, block *hclext.Block) error {
	count, exists := block.Body.Attributes["count"]

	if exists {
		conditionalExpr, ok := count.Expr.(*hclsyntax.ConditionalExpr)
		if ok {
			//todo
		}
	}

	return nil
}

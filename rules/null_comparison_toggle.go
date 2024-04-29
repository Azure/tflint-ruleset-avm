package rules

import (
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var nullComparisonToggleBodySchema = &hclext.BodySchema{
	Blocks: []hclext.BlockSchema{
		{
			Type:       "module",
			LabelNames: []string{"name"},
			Body: &hclext.BodySchema{
				Attributes: []hclext.AttributeSchema{
					{
						Name: "source",
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
		moduleSourceBodySchema,
		&tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return err
	}

	var errList error
	for _, block := range body.Blocks {
		if block.Type != "module" {
			continue
		}

		if subErr := t.checkBlock(r, block); subErr != nil {
			errList = multierror.Append(errList, subErr)
		}
	}

	return errList
}

func (t *NullComparisonToggleRule) checkBlock(r tflint.Runner, block *hclext.Block) error {
	source, exists := block.Body.Attributes["source"]
	if !exists {
		return r.EmitIssue(
			t,
			"The `source` field should be declared in the `module` block",
			block.DefRange,
		)
	}

	if err := r.EvaluateExpr(source.Expr, t.isAVMModule(r, source.NameRange), &tflint.EvaluateExprOption{ModuleCtx: tflint.RootModuleCtxType}); err != nil {
		return err
	}

	return nil
}

func (t *NullComparisonToggleRule) isAVMModule(r tflint.Runner, issueRange hcl.Range) func(string) error {
	return func(source string) error {
		if !(strings.HasPrefix(source, "Azure/") && strings.Contains(source, "avm-")) {
			return r.EmitIssue(
				t,
				"The `source` property constraint should start with `Azure/` and contain `avm-` to only involve AVM Module",
				issueRange,
			)
		}

		return nil
	}
}

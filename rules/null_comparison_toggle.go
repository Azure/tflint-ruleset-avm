package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var nullComparisonToggleBodySchema = &hclext.BodySchema{
	Blocks: []hclext.BlockSchema{
		{
			Type:       "variable",
			LabelNames: []string{"name"},
			Body: &hclext.BodySchema{
				Attributes: []hclext.AttributeSchema{
					{Name: "type"},
				},
			},
		},
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

	var resourceIdStr string
	var errList error
	for _, block := range body.Blocks {
		if block.Type != "resource" {
			continue
		}

		resourceIdString, subErr := t.checkResourceBlock(r, block)
		if subErr != nil {
			errList = multierror.Append(errList, subErr)
		}
		resourceIdStr = resourceIdString
	}

	for _, block := range body.Blocks {
		if block.Type != "variable" {
			continue
		}

		if subErr := t.checkVariableBlock(r, block, resourceIdStr); subErr != nil {
			errList = multierror.Append(errList, subErr)
		}
	}

	return errList
}

func (t *NullComparisonToggleRule) checkResourceBlock(r tflint.Runner, block *hclext.Block) (string, error) {
	count, exists := block.Body.Attributes["count"]
	if !exists {
		return "", nil
	}

	countConditionalExpr, ok := count.Expr.(*hclsyntax.ConditionalExpr)
	if !ok {
		return "", nil
	}

	var resourceIdStr string
	var err error
	for _, dynamicObj := range countConditionalExpr.Variables() {
		resourceIdStr, err = t.checkCountBlock(r, dynamicObj, count.Range)
		if err != nil {
			return resourceIdStr, err
		}
	}

	return resourceIdStr, nil
}

func (t *NullComparisonToggleRule) checkCountBlock(r tflint.Runner, dynamicObj hcl.Traversal, rangeInfo hcl.Range) (string, error) {
	for _, dynamicVal := range dynamicObj {
		if v, ok := dynamicVal.(hcl.TraverseRoot); ok && strings.HasSuffix(v.Name, "local") {
			break
		}

		if v, ok := dynamicVal.(hcl.TraverseAttr); ok && strings.HasSuffix(strings.ToLower(v.Name), "_id") {
			return v.Name, r.EmitIssue(
				t,
				"The variable should be defined as object type for the resource id",
				rangeInfo,
			)
		}
	}

	return "", nil
}

func (t *NullComparisonToggleRule) checkVariableBlock(r tflint.Runner, block *hclext.Block, resourceIdStr string) error {
	for _, label := range block.Labels {
		if label == resourceIdStr {
			attr, exists := block.Body.Attributes["type"]
			if !exists {
				return r.EmitIssue(
					t,
					fmt.Sprintf("`%s` type not declared", label),
					attr.Range,
				)
			}

			for _, dynamicVal := range attr.Expr.Variables() {
				if v := dynamicVal.RootName(); v == "string" {
					return r.EmitIssue(
						t,
						"The variable should be defined as object type for the resource id",
						dynamicVal.SourceRange(),
					)
				}
			}
		}
	}

	return nil
}

package rules

import (
	"fmt"
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

var _ tflint.Rule = new(AzApiNoPreviewVersionRule)

type AzApiNoPreviewVersionRule struct {
	tflint.DefaultRule
}

func NewAzApiNoPreviewVersionRule() *AzApiNoPreviewVersionRule {
	return new(AzApiNoPreviewVersionRule)
}

func (a *AzApiNoPreviewVersionRule) Name() string {
	return "azapi_no_preview_version_sfr1"
}

func (a *AzApiNoPreviewVersionRule) Link() string {
	return "https://azure.github.io/Azure-Verified-Modules/spec/SFR1/"
}

func (a *AzApiNoPreviewVersionRule) Enabled() bool {
	return false
}

func (a *AzApiNoPreviewVersionRule) Severity() tflint.Severity {
	return tflint.WARNING
}

func (a *AzApiNoPreviewVersionRule) Check(r tflint.Runner) error {
	types := []string{"data", "resource", "ephemeral"}
	for _, t := range types {
		var bodySchema = &hclext.BodySchema{
			Blocks: []hclext.BlockSchema{
				{
					Type:       t,
					LabelNames: []string{"type", "name"},
					Body: &hclext.BodySchema{
						Mode: 0,
						Attributes: []hclext.AttributeSchema{
							{
								Name:     "type",
								Required: false,
							},
						},
					},
				},
			},
		}
		body, err := r.GetModuleContent(bodySchema, nil)
		if err != nil {
			return err
		}
		for _, block := range body.Blocks {
			if !strings.HasPrefix(block.Labels[0], "azapi_") {
				continue
			}
			typeAttribute, ok := block.Body.Attributes["type"]
			if !ok {
				continue
			}
			var typeValue string
			if err = r.EvaluateExpr(typeAttribute.Expr, &typeValue, &tflint.EvaluateExprOption{
				WantType: &cty.String,
			}); err != nil {
				return err
			}
			if !strings.HasSuffix(strings.ToLower(typeValue), "-preview") {
				continue
			}
			if err = r.EmitIssue(
				a,
				fmt.Sprintf("Resource type `%s` is using a preview API version, which is not recommended.", typeValue),
				typeAttribute.Range,
			); err != nil {
				return err
			}
		}
	}
	return nil
}

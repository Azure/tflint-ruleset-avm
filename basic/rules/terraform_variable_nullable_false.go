package rules

import (
	"github.com/Azure/tflint-ruleset-avm/basic/project"
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var _ tflint.Rule = &TerraformVariableNullableFalseRule{}

type TerraformVariableNullableFalseRule struct {
	tflint.DefaultRule
}

func NewTerraformVariableNullableFalseRule() *TerraformVariableNullableFalseRule {
	return &TerraformVariableNullableFalseRule{}
}

func (r *TerraformVariableNullableFalseRule) Name() string {
	return "terraform_variable_nullable_false"
}

func (r *TerraformVariableNullableFalseRule) Link() string {
	return project.ReferenceLink(r.Name())
}

func (r *TerraformVariableNullableFalseRule) Enabled() bool {
	return true
}

func (r *TerraformVariableNullableFalseRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

func (r *TerraformVariableNullableFalseRule) Check(runner tflint.Runner) error {
	content, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "variable",
				LabelNames: []string{"name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{
							Name:     "nullable",
							Required: false,
						},
					},
				},
			},
		},
	}, nil)
	if err != nil {
		return err
	}
	for _, b := range content.Blocks {
		attribute, ok := b.Body.Attributes["nullable"]
		if !ok {
			continue
		}
		v, _ := attribute.Expr.Value(&hcl.EvalContext{})
		if v.False() {
			continue
		}
		err := runner.EmitIssue(r, r.errorMessage(), attribute.Range)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *TerraformVariableNullableFalseRule) errorMessage() string {
	return "`nullable` is default to `true` so we don't need to declare it explicitly."
}

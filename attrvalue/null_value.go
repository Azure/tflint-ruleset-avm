package attrvalue

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// NullRule checks whether an attribute value is null.
type NullRule struct {
	tflint.DefaultRule // Embed the default rule to reuse its implementation

	resourceType  string // e.g. "azurerm_storage_account"
	attributeName string // e.g. "account_replication_type"
}

var _ tflint.Rule = (*SimpleRule[any])(nil)

// NewSimpleRule returns a new rule with the given resource type, attribute name, and expected values.
func NewNullRule(resourceType string, attributeName string) *NullRule {
	return &NullRule{
		resourceType:  resourceType,
		attributeName: attributeName,
	}
}

func (r *NullRule) Name() string {
	return fmt.Sprintf("%s.%s must be null", r.resourceType, r.attributeName)
}

func (r *NullRule) Enabled() bool {
	return true
}

func (r *NullRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *NullRule) Check(runner tflint.Runner) error {
	variables, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "variable",
				LabelNames: []string{"name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{
							Name:     "type",
							Required: false,
						},
						{
							Name:     "default",
							Required: false,
						},
					},
				},
			},
		},
	},
		&tflint.GetModuleContentOption{
			ExpandMode: tflint.ExpandModeNone,
		})
	if err != nil {
		return err
	}

	variablesMap := make(map[string]cty.Value, len(variables.Blocks))
	for _, v := range variables.Blocks {
		if len(v.Labels) != 1 {
			return fmt.Errorf("variables should have exactly one label: %s", v.DefRange)
		}
		vName := v.Labels[0]
		defAttr, ok := v.Body.Attributes["default"]
		if !ok {
			variablesMap[vName] = cty.NullVal(cty.DynamicPseudoType)
			continue
		}
		defVal, diags := defAttr.Expr.Value(nil)
		if diags.HasErrors() {
			return fmt.Errorf("could not evaluate variable %s default value: %s", vName, diags)
		}
		variablesMap[vName] = defVal
	}

	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.attributeName},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		attribute, exists := resource.Body.Attributes[r.attributeName]
		if !exists {
			continue
		}
		attrVal, diags := attribute.Expr.Value(&hcl.EvalContext{
			Variables: map[string]cty.Value{
				"var": cty.ObjectVal(variablesMap),
			},
		})
		if diags.HasErrors() {
			return fmt.Errorf("could not evaluate attribute %s value: %s", r.attributeName, diags)
		}
		if !attrVal.IsNull() {
			runner.EmitIssue(
				r,
				fmt.Sprintf("invalid attribute value of `%s` - expecting null", r.attributeName),
				attribute.Range,
			)
		}
	}
	return nil
}

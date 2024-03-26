package attrvalue

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// getSimpleAttrs returns a slice of attributes with the given attribute name from the resources of the given resource type.
func getSimpleAttrs(runner tflint.Runner, resourceType string, attributeName string) ([]*hclext.Attribute, error) {
	resources, err := runner.GetResourceContent(resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: attributeName},
		},
	}, nil)
	if err != nil {
		return nil, err
	}
	attrs := make([]*hclext.Attribute, 0, len(resources.Blocks))
	for _, resource := range resources.Blocks {
		if attribute := getAttrFromBlock(resource, attributeName); attribute != nil {
			attrs = append(attrs, attribute)
		}
	}
	return attrs, nil
}

// getNestedBlockAttrs returns a slice of attributes with the given attribute name from the nested blocks of the given resource type.
func getNestedBlockAttrs(runner tflint.Runner, resourceType string, nestedBlockType string, attributeName string) ([]*hclext.Attribute, error) {
	resources, err := runner.GetResourceContent(resourceType, &hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type: nestedBlockType,
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{
							Name:     attributeName,
							Required: false,
						},
					},
				},
			},
		},
	}, nil)
	if err != nil {
		return nil, err
	}
	attrs := make([]*hclext.Attribute, 0, len(resources.Blocks))
	for _, resource := range resources.Blocks {
		for _, block := range resource.Body.Blocks {
			if attr := getAttrFromBlock(block, attributeName); attr != nil {
				attrs = append(attrs, attr)
			}
		}
	}
	return attrs, nil
}

// getAttrFromBlock returns the attribute with the given attribute name from the block.
func getAttrFromBlock(block *hclext.Block, attributeName string) *hclext.Attribute {
	attribute, exists := block.Body.Attributes[attributeName]
	if !exists {
		return nil
	}
	return attribute
}

// getDefaultVariableEvalContext returns the evaluation context with the default values of the variables.
func getDefaultVariableEvalContext(runner tflint.Runner) (*hcl.EvalContext, error) {
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
		return nil, err
	}

	// variablesMap is needed to construct the hcl.EvalContext.
	variablesMap := make(map[string]cty.Value, len(variables.Blocks))
	for _, v := range variables.Blocks {
		if len(v.Labels) != 1 {
			return nil, fmt.Errorf("variables should have exactly one label: %s", v.DefRange)
		}
		vName := v.Labels[0]
		defAttr, ok := v.Body.Attributes["default"]
		if !ok {
			variablesMap[vName] = cty.NullVal(cty.DynamicPseudoType)
			continue
		}
		defVal, diags := defAttr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("could not evaluate variable %s default value: %s", vName, diags)
		}
		variablesMap[vName] = defVal
	}

	ec := &hcl.EvalContext{
		Variables: map[string]cty.Value{
			"var": cty.ObjectVal(variablesMap),
		},
	}
	return ec, nil
}

package rules

import (
	"fmt"

	"github.com/Azure/tflint-ruleset-avm/avmhelper"
	"github.com/Azure/tflint-ruleset-avm/interfaces"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var _ tflint.Rule = new(AVMInterfaceRule)

// AVMInterfaceRule is the struct that represents a rule that
// check for the correct usage of an interface.
type AVMInterfaceRule struct {
	tflint.DefaultRule
	Iface interfaces.Interface
}

func (t *AVMInterfaceRule) Name() string {
	return t.Iface.Name
}

func (t *AVMInterfaceRule) Link() string {
	return t.Iface.Link
}

func (t *AVMInterfaceRule) Enabled() bool {
	return false
}

func (t *AVMInterfaceRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (t *AVMInterfaceRule) Check(r tflint.Runner) error {
	path, err := r.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() {
		// This rule does not evaluate child modules.
		return nil
	}

	// Define the schema that we want to pull out of the module content.
	body, err := r.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "variable",
				LabelNames: []string{"name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "type"},
						{Name: "default"},
					},
					Blocks: []hclext.BlockSchema{
						{
							Type:       "validation",
							LabelNames: []string{},
							Body: &hclext.BodySchema{
								Attributes: []hclext.AttributeSchema{
									{Name: "condition"},
									{Name: "error_message"},
								},
							},
						},
					},
				},
			},
		},
	}, &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return err
	}

	// Iterate over the variables and check for one called `lock`.
	for _, variable := range body.Blocks {
		if variable.Labels[0] != t.Iface.Name {
			continue
		}

		// Check if the variable has a type attribute.
		typeattr, exists := variable.Body.Attributes["type"]
		if !exists {
			if err := r.EmitIssue(
				t,
				fmt.Sprintf("`%s` variable type not declared", variable.Labels[0]),
				variable.DefRange,
			); err != nil {
				return err
			}
			continue
		}

		// Check if the type interface is correct.
		if eq, diags := avmhelper.CheckTypeConstraintsAreEqual(typeattr.AsNative().Expr, t.Iface.TypeExpression()); diags.HasErrors() || !*eq {
			r.EmitIssue(t,
				fmt.Sprintf("`%s` variable type does not comply with the interface specification:\n\n%s", variable.Labels[0], t.Iface.Type),
				variable.DefRange,
			)
		}

		// Check if the variable has a default attribute.
		defaultattr, exists := variable.Body.Attributes["default"]
		if !exists {
			if err := r.EmitIssue(
				t,
				fmt.Sprintf("`var.%s`: default not declared", variable.Labels[0]),
				variable.DefRange,
			); err != nil {
				return err
			}
			continue
		}

		// Check if the default value is correct.
		defaultval, _ := defaultattr.Expr.Value(nil)
		if defaultval != interfaces.Lock.Default {
			if err := r.EmitIssue(
				t,
				fmt.Sprintf("`var.%s`: default value is not `null`", variable.Labels[0]),
				variable.DefRange,
			); err != nil {
				return err
			}
		}
	}

	return nil
}

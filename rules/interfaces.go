package rules

import (
	"fmt"

	"github.com/Azure/tflint-ruleset-avm/avmhelper"
	"github.com/Azure/tflint-ruleset-avm/interfaces"
	"github.com/zclconf/go-cty/cty"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// Check interface compliance with the tflint.Rule.
var _ tflint.Rule = new(AvmInterfaceRule)

// AvmInterfaceRule is the struct that represents a rule that
// check for the correct usage of an interface.
// Note interfaces.AVMInterface is embedded in this struct,
// which is used by the constructor func NewAVMInterface...Rule().
type AvmInterfaceRule struct {
	tflint.DefaultRule
	Iface interfaces.AvmInterface
}

// NewAVMInterfaceRule returns a new rule with the given interface.
// The data is taken from the embedded interfaces.AVMInterface.
func (t *AvmInterfaceRule) Name() string {
	return t.Iface.Name
}

func (t *AvmInterfaceRule) Link() string {
	return t.Iface.Link
}

// Enabled returns whether the rule is enabled.
// This is sourced from the embedded interfaces.AVMInterface.
func (t *AvmInterfaceRule) Enabled() bool {
	return t.Iface.Enabled
}

// Severity returns the severity of the rule.
// Currently all interfaces have severity ERROR.
func (t *AvmInterfaceRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Check checks whether the module satisfies the interface.
// It will search for a variable with the same name as the interface.
// It will check the type, default value and nullable attributes.
// TODO: think about how we can check validation rules.
func (t *AvmInterfaceRule) Check(r tflint.Runner) error {
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
						{Name: "nullable"},
					},
					// We do not do anything with the validation data at the moment.
					Blocks: []hclext.BlockSchema{
						{
							Type: "validation",
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

	// Iterate over the variables and check for the name we are interested in.
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
			if err := r.EmitIssue(t,
				fmt.Sprintf("`%s` variable type does not comply with the interface specification:\n\n%s", variable.Labels[0], t.Iface.Type),
				typeattr.Range,
			); err != nil {
				return err
			}
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

		if defaultval.Equals(t.Iface.Default) != cty.True {
			if err := r.EmitIssue(
				t,
				fmt.Sprintf("`var.%s`: default value is not correct, see: %s", variable.Labels[0], t.Link()),
				variable.DefRange,
			); err != nil {
				return err
			}
		}

		// Check nullable
		nullableattr, nullableExists := variable.Body.Attributes["nullable"]
		switch {
		// Raise issue if nullable not set and desired is that nullable is false.
		case !nullableExists && !t.Iface.Nullable:
			if err := r.EmitIssue(
				t,
				fmt.Sprintf("`var.%s`: nullable is not set and should be set to false", variable.Labels[0]),
				variable.DefRange,
			); err != nil {
				return err
			}
			// Raise issue if nullable is set and desired is that nullable is true (default, should not explicitly set nullable to true).
		case nullableExists && t.Iface.Nullable:
			if err := r.EmitIssue(
				t,
				fmt.Sprintf("`var.%s`: nullable is set and should not be, we require this to be true and this is the default behaviour so no need to set explicitly", variable.Labels[0]),
				nullableattr.Range,
			); err != nil {
				return err
			}
		// Raise issue if nullable is set and desired is that nullable is false.
		case !t.Iface.Nullable && nullableExists:
			nullableval, _ := nullableattr.Expr.Value(nil)
			if nullableval != cty.BoolVal(false) {
				if err := r.EmitIssue(
					t,
					fmt.Sprintf("`var.%s`: nullable is set to true and should be set to false", variable.Labels[0]),
					nullableattr.Range,
				); err != nil {
					return err
				}
			}
		default:
			continue
		}
		// TODO: Check validation rules.
	}
	return nil
}

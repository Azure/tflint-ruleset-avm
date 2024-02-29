package rules

import (
	"fmt"

	"github.com/Azure/tflint-ruleset-avm/interfaces"
	"github.com/hashicorp/hcl/v2"
	"github.com/matt-FFFFFF/tfvarcheck/check"
	"github.com/matt-FFFFFF/tfvarcheck/varcheck"
	"github.com/zclconf/go-cty/cty"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// variableBodySchema is the schema for the variable block that we want to extract from the config.
var variableBodySchema = hclext.BodySchema{
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
}

// Check interface compliance with the tflint.Rule.
var _ tflint.Rule = new(InterfaceVarCheckRule)

// InterfaceVarCheckRule is the struct that represents a rule that
// check for the correct usage of an interface.
type InterfaceVarCheckRule struct {
	tflint.DefaultRule
	avmInterface interfaces.AvmInterface // This is the interface we are checking for.
	vc           varcheck.VarCheck       // This is the VarCheck object containing the strong types that we will use to check the variable. It is created from the AvmInterface.
}

// NewVarCheckRule returns a new rule with the given variable.
func NewVarCheckRuleFromAvmInterface(ifce interfaces.AvmInterface) *InterfaceVarCheckRule {
	return &InterfaceVarCheckRule{
		avmInterface: ifce,
		vc:           ifce.VarType,
	}
}

// NewAVMInterfaceRule returns a new rule with the given interface.
// The data is taken from the embedded interfaces.AVMInterface.
func (vcr *InterfaceVarCheckRule) Name() string {
	return vcr.avmInterface.Name
}

func (vcr *InterfaceVarCheckRule) Link() string {
	return vcr.avmInterface.Link
}

// Enabled returns whether the rule is enabled.
// This is sourced from the embedded interfaces.AVMInterface.
func (vcr *InterfaceVarCheckRule) Enabled() bool {
	return vcr.avmInterface.Enabled
}

// Severity returns the severity of the rule.
func (vcr *InterfaceVarCheckRule) Severity() tflint.Severity {
	return vcr.avmInterface.Severity
}

// Check checks whether the module satisfies the interface.
// It will search for a variable with the same name as the interface.
// It will check the type, default value and nullable attributes.
func (vcr *InterfaceVarCheckRule) Check(r tflint.Runner) error {
	path, err := r.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() {
		// This rule does not evaluate child modules.
		return nil
	}

	// Define the schema that we want to pull out of the module content.
	body, err := r.GetModuleContent(
		&variableBodySchema,
		&tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return err
	}

	// Iterate over the variables and check for the name we are interested in.
	for _, v := range body.Blocks {
		if v.Labels[0] != vcr.avmInterface.Name {
			continue
		}

		// Check if the variable has a type attribute.
		typeAttr, exists := v.Body.Attributes["type"]
		if !exists {
			if err := r.EmitIssue(
				vcr,
				fmt.Sprintf("`%s` variable type not declared", v.Labels[0]),
				v.DefRange,
			); err != nil {
				return err
			}
			continue
		}

		// Check if the type interface is correct.
		gotType, diags := varcheck.NewTypeConstraintWithDefaultsFromExp(typeAttr.Expr)
		if diags.HasErrors() {
			return diags
		}
		if eq := check.EqualTypeConstraints(gotType, vcr.vc.TypeConstraintWithDefs); !eq {
			if err := r.EmitIssue(vcr,
				fmt.Sprintf("variable type does not comply with the interface specification:\n\n%s", vcr.avmInterface.VarTypeString),
				typeAttr.Range,
			); err != nil {
				return err
			}
		}

		// Check if the variable has a default attribute.
		defaultAttr, exists := v.Body.Attributes["default"]
		if !exists {
			if err := r.EmitIssue(
				vcr,
				"default not declared",
				v.DefRange,
			); err != nil {
				return err
			}
			continue
		}

		// Check if the default value is correct.
		defaultVal, _ := defaultAttr.Expr.Value(nil)

		if !check.EqualCtyValue(defaultVal, vcr.vc.Default) {
			if err := r.EmitIssue(
				vcr,
				fmt.Sprintf("default value is not correct, see: %s", vcr.Link()),
				v.DefRange,
			); err != nil {
				return err
			}
		}

		// Check if the variable has a nullable attribute and fetch the value,
		// else set it to null.
		var nullableVal cty.Value
		nullableAttr, nullableExists := v.Body.Attributes["nullable"]
		if !nullableExists {
			nullableVal = cty.NullVal(cty.Bool)
		} else {
			var diags hcl.Diagnostics
			if nullableVal, diags = nullableAttr.Expr.Value(nil); diags.HasErrors() {
				if diags.HasErrors() {
					return diags
				}
			}
		}

		// Check nullable attribute.
		if ok := check.Nullable(nullableVal, vcr.vc.Nullable); !ok {
			msg := "nullable should not be set."
			if !vcr.vc.Nullable {
				msg = "nullable should be set to false"
			}
			rg := v.DefRange
			if nullableAttr != nil {
				rg = nullableAttr.Range
			}
			if err := r.EmitIssue(
				vcr,
				msg,
				rg,
			); err != nil {
				return err
			}
		}

		// TODO: Check validation rules.
	}
	return nil
}

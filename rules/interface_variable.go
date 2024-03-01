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
	interfaces.AvmInterface // This is the interface we are checking for.
}

// NewVarCheckRule returns a new rule with the given variable.
func NewVarCheckRuleFromAvmInterface(ifce interfaces.AvmInterface) *InterfaceVarCheckRule {
	return &InterfaceVarCheckRule{
		AvmInterface: ifce,
	}
}

// NewAVMInterfaceRule returns a new rule with the given interface.
// The data is taken from the embedded interfaces.AVMInterface.
func (vcr *InterfaceVarCheckRule) Name() string {
	return vcr.RuleName
}

func (vcr *InterfaceVarCheckRule) Link() string {
	return vcr.RuleLink
}

// Enabled returns whether the rule is enabled.
// This is sourced from the embedded interfaces.AVMInterface.
func (vcr *InterfaceVarCheckRule) Enabled() bool {
	return vcr.RuleEnabled
}

// Severity returns the severity of the rule.
func (vcr *InterfaceVarCheckRule) Severity() tflint.Severity {
	return vcr.RuleSeverity
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
	for _, b := range body.Blocks {
		if b.Labels[0] != vcr.RuleName {
			continue
		}

		typeAttr, c := CheckWithReturnValue(newChecker(), getTypeAttr(vcr, r, b))
		defaultAttr, c := CheckWithReturnValue(c, getDefaultAttr(vcr, r, b))
		c = c.Check(checkVarType(vcr, r, typeAttr)).
			Check(checkDefaultValue(vcr, r, b, defaultAttr)).
			Check(checkNullableValue(vcr, r, b))

		if c.err != nil {
			return c.err
		}
		// TODO: Check validation rules.
		return nil
	}
	return nil
}

func checkNullableValue(vcr *InterfaceVarCheckRule, r tflint.Runner, b *hclext.Block) func() (bool, error) {
	return func() (bool, error) {
		nullableAttr, nullableExists := b.Body.Attributes["nullable"]
		nullableVal := cty.NullVal(cty.Bool)
		if nullableExists {
			var diags hcl.Diagnostics
			if nullableVal, diags = nullableAttr.Expr.Value(nil); diags.HasErrors() {
				return false, diags
			}
		}
		// Check nullable attribute.
		if ok := check.Nullable(nullableVal, vcr.Nullable); !ok {
			msg := "nullable should not be set."
			if !vcr.Nullable {
				msg = "nullable should be set to false"
			}
			rg := b.DefRange
			if nullableAttr != nil {
				rg = nullableAttr.Range
			}
			return false, r.EmitIssue(vcr, msg, rg)
		}
		return true, nil
	}
}

func getTypeAttr(rule tflint.Rule, r tflint.Runner, b *hclext.Block) func() (*hclext.Attribute, bool, error) {
	return func() (*hclext.Attribute, bool, error) {
		// Check if the variable has a type attribute.
		typeAttr, exists := b.Body.Attributes["type"]
		if !exists {
			return typeAttr, false, r.EmitIssue(
				rule,
				fmt.Sprintf("`%s` variable type not declared", b.Labels[0]),
				b.DefRange,
			)
		}
		return typeAttr, true, nil
	}
}

func checkVarType(vcr *InterfaceVarCheckRule, r tflint.Runner, typeAttr *hclext.Attribute) func() (bool, error) {
	return func() (bool, error) {
		// Check if the type interface is correct.
		gotType, diags := varcheck.NewTypeConstraintWithDefaultsFromExp(typeAttr.Expr)
		if diags.HasErrors() {
			return false, diags
		}
		if eq := check.EqualTypeConstraints(gotType, vcr.TypeConstraintWithDefs); !eq {
			return true, r.EmitIssue(vcr,
				fmt.Sprintf("variable type does not comply with the interface specification:\n\n%s", vcr.VarTypeString),
				typeAttr.Range,
			)
		}
		return true, nil
	}
}

func getDefaultAttr(vcr tflint.Rule, r tflint.Runner, b *hclext.Block) func() (*hclext.Attribute, bool, error) {
	return func() (*hclext.Attribute, bool, error) {
		// Check if the variable has a default attribute.
		defaultAttr, exists := b.Body.Attributes["default"]
		if !exists {
			return defaultAttr, false, r.EmitIssue(
				vcr,
				"default not declared",
				b.DefRange,
			)
		}
		return defaultAttr, true, nil
	}
}

func checkDefaultValue(vcr *InterfaceVarCheckRule, r tflint.Runner, b *hclext.Block, defaultAttr *hclext.Attribute) func() (bool, error) {
	return func() (bool, error) {
		// Check if the default value is correct.
		defaultVal, _ := defaultAttr.Expr.Value(nil)
		if !check.EqualCtyValue(defaultVal, vcr.Default) {
			return true, r.EmitIssue(
				vcr,
				fmt.Sprintf("default value is not correct, see: %s", vcr.Link()),
				b.DefRange,
			)
		}
		return true, nil
	}
}

func newChecker() checker {
	return checker{
		continueCheck: true,
	}
}

type checker struct {
	continueCheck bool
	err           error
}

func (c checker) Check(check func() (bool, error)) checker {
	if c.err != nil || !c.continueCheck {
		return c
	}
	continueCheck, err := check()
	return checker{
		continueCheck: continueCheck,
		err:           err,
	}
}

func CheckWithReturnValue[TR any](c checker, check func() (TR, bool, error)) (ret TR, rc checker) {
	if c.err != nil || !c.continueCheck {
		rc = c
		return
	}
	tr, continueCheck, err := check()
	return tr, checker{
		continueCheck: continueCheck,
		err:           err,
	}
}

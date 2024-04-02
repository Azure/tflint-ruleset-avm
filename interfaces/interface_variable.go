package interfaces

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/matt-FFFFFF/tfvarcheck/check"
	"github.com/matt-FFFFFF/tfvarcheck/varcheck"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// variableBodySchema is the schema for the variable block that we want to extract from the config.
var variableBodySchema = &hclext.BodySchema{
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
	AvmInterface // This is the interface we are checking for.
}

// NewVarCheckRuleFromAvmInterface returns a new rule with the given variable.
func NewVarCheckRuleFromAvmInterface(ifce AvmInterface) *InterfaceVarCheckRule {
	return &InterfaceVarCheckRule{
		AvmInterface: ifce,
	}
}

// Name returns the rule name.
func (vcr *InterfaceVarCheckRule) Name() string {
	return vcr.RuleName
}

// Link returns the link to the rule documentation.
func (vcr *InterfaceVarCheckRule) Link() string {
	return vcr.RuleLink
}

// Enabled returns whether the rule is enabled.
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
		variableBodySchema,
		&tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return err
	}

	// Iterate over the variables and check for the name we are interested in.
	for _, b := range body.Blocks {
		if b.Labels[0] != vcr.RuleName {
			continue
		}

		typeAttr, c := CheckWithReturnValue(NewChecker(), getAttr(vcr, r, b, "type"))
		defaultAttr, c := CheckWithReturnValue(c, getAttr(vcr, r, b, "default"))
		if c = c.Check(checkVarType(vcr, r, typeAttr)).
			Check(checkDefaultValue(vcr, r, b, defaultAttr)).
			Check(checkNullableValue(vcr, r, b)); c.err != nil {
			return c.err
		}
		// TODO: Check validation rules.
		return nil
	}
	return nil
}

// getTypeAttr returns a function that will return the type attribute from a given hcl block.
// It is designed to be used with the CheckWithReturnValue function.
func getAttr(rule tflint.Rule, r tflint.Runner, b *hclext.Block, attrName string) func() (*hclext.Attribute, bool, error) {
	return func() (*hclext.Attribute, bool, error) {
		attr, exists := b.Body.Attributes[attrName]
		if !exists {
			return attr, false, r.EmitIssue(
				rule,
				fmt.Sprintf("`%s` %s not declared", b.Labels[0], attrName),
				b.DefRange,
			)
		}
		return attr, true, nil
	}
}

// checkNullableValue checks if the nullable attribute is correct.
// It is designed to be supplied to the Checker.Check() function.
func checkNullableValue(vcr *InterfaceVarCheckRule, r tflint.Runner, b *hclext.Block) func() (bool, error) {
	return func() (bool, error) {
		nullableAttr, nullableExists := b.Body.Attributes["nullable"]
		nullableVal := cty.NullVal(cty.Bool)
		var diags hcl.Diagnostics
		if nullableExists {
			nullableVal, diags = nullableAttr.Expr.Value(nil)
		}
		if diags.HasErrors() {
			return false, diags
		}
		// Check nullable attribute.
		if ok := check.Nullable(nullableVal, vcr.Nullable); ok {
			return true, nil
		}
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
}

// checkVarType checks if the type of the variable is correct.
// It is designed to be supplied to the Checker.Check() function.
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

// checkDefaultValue checks if the default value of a variable is correct.
// It is designed to be supplied to the Checker.Check() function.
func checkDefaultValue(vcr *InterfaceVarCheckRule, r tflint.Runner, b *hclext.Block, defaultAttr *hclext.Attribute) func() (bool, error) {
	return func() (bool, error) {
		// Check if the default value is correct.
		defaultVal, _ := defaultAttr.Expr.Value(nil)

		// The `{}` for `defaultVal` would be evaluated as an empty object here no matter what type is set, so we have to compare string without element type info.
		if isEmptyCollection(defaultVal) && isEmptyCollection(vcr.Default) && toString(defaultVal) == toString(vcr.Default) {
			return true, nil
		}
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

func isEmptyCollection(v cty.Value) bool {
	return (v.Type().IsCollectionType() || v.Type().IsObjectType()) && v.LengthInt() == 0
}

func toString(v cty.Value) string {
	return string(hclwrite.NewExpressionLiteral(v).BuildTokens(hclwrite.Tokens{}).Bytes())
}

// NewChecker is the constructor for the Checker type.
func NewChecker() Checker {
	return Checker{
		continueCheck: true,
	}
}

// Checker is a struct that is used to chain checks together.
type Checker struct {
	continueCheck bool
	err           error
}

// Check is a executes a supplied function that returns a bool and an error.
// The bool is a continueCheck value that is used to determine if the check should continue.
// The error is the error that is returned from the check.
//
// This function returns a new Checker, so it can be chained with other checks in a fluent style.
func (c Checker) Check(check func() (bool, error)) Checker {
	if c.err != nil || !c.continueCheck {
		return c
	}
	continueCheck, err := check()
	return Checker{
		continueCheck: continueCheck,
		err:           err,
	}
}

// CheckWithReturnValue is a generic function that runs a check func() that, as well as
// returning a bool & error, also returns a value.
// The main function will then return the value and a new Checker with the continueCheck and err.
func CheckWithReturnValue[TR any](c Checker, check func() (TR, bool, error)) (ret TR, rc Checker) {
	if c.err != nil || !c.continueCheck {
		rc = c
		return
	}
	tr, continueCheck, err := check()
	return tr, Checker{
		continueCheck: continueCheck,
		err:           err,
	}
}

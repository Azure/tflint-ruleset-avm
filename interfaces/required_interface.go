package interfaces

import (
	"fmt"

	"github.com/Azure/tflint-ruleset-avm/internal/tfvarcheck/varcheck"
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

var azapiResourceBodySchema = &hclext.BodySchema{
	Blocks: []hclext.BlockSchema{
		{
			Type:       "resource",
			LabelNames: []string{"type", "name"},
			Body:       &hclext.BodySchema{},
		},
	},
}

var mandatoryInterfaceAzapiResourceTypes = map[string]struct{}{
	"azapi_data_plane_resource": {},
	"azapi_resource":            {},
	"azapi_resource_action":     {},
	"azapi_update_resource":     {},
}

type requiredVariableApplicability func(tflint.Runner) (bool, hcl.Range, error)

type requiredVariableValidator func(*requiredVariableRule, tflint.Runner, *hclext.Block) error

type requiredVariableRule struct {
	tflint.DefaultRule

	name          string
	variableName  string
	link          string
	requirement   string
	applicability requiredVariableApplicability
	validator     requiredVariableValidator
}

func newRequiredVariableRule(
	name string,
	variableName string,
	link string,
	requirement string,
	applicability requiredVariableApplicability,
	validator requiredVariableValidator,
) *requiredVariableRule {
	return &requiredVariableRule{
		name:          name,
		variableName:  variableName,
		link:          link,
		requirement:   requirement,
		applicability: applicability,
		validator:     validator,
	}
}

func (r *requiredVariableRule) Name() string {
	return r.name
}

func (r *requiredVariableRule) Link() string {
	return r.link
}

func (r *requiredVariableRule) Enabled() bool {
	return true
}

func (r *requiredVariableRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *requiredVariableRule) Check(runner tflint.Runner) error {
	applicable, sourceRange, err := r.applicability(runner)
	if err != nil || !applicable {
		return err
	}

	variable, err := findVariableDeclaration(runner, r.variableName)
	if err != nil {
		return err
	}
	if variable == nil {
		return runner.EmitIssue(
			r,
			fmt.Sprintf("variable `%s` must be declared %s; see: %s", r.variableName, r.requirement, r.link),
			sourceRange,
		)
	}

	return r.validator(r, runner, variable)
}

func appliesToDirectAzapiResource(runner tflint.Runner) (bool, hcl.Range, error) {
	content, err := runner.GetModuleContent(
		azapiResourceBodySchema,
		&tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone},
	)
	if err != nil {
		return false, hcl.Range{}, err
	}

	for _, block := range content.Blocks {
		if len(block.Labels) != 2 {
			continue
		}
		if _, ok := mandatoryInterfaceAzapiResourceTypes[block.Labels[0]]; ok {
			return true, block.DefRange, nil
		}
	}

	return false, hcl.Range{}, nil
}

func appliesWhenVariableDeclared(name string) requiredVariableApplicability {
	return func(runner tflint.Runner) (bool, hcl.Range, error) {
		variable, err := findVariableDeclaration(runner, name)
		if err != nil || variable == nil {
			return false, hcl.Range{}, err
		}
		return true, variable.DefRange, nil
	}
}

func findVariableDeclaration(runner tflint.Runner, name string) (*hclext.Block, error) {
	content, err := runner.GetModuleContent(
		variableBodySchema,
		&tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone},
	)
	if err != nil {
		return nil, err
	}

	for _, block := range content.Blocks {
		if len(block.Labels) == 1 && block.Labels[0] == name {
			return block, nil
		}
	}

	return nil, nil
}

type nullableRequirement int

const (
	nullableMustBeFalse nullableRequirement = iota
	nullableMayBeTrue
)

type variableShape struct {
	typeDescription     string
	defaultDescription  string
	nullableDescription string
	validateType        func(varcheck.TypeConstraintWithDefaults) bool
	validateDefault     func(cty.Value, varcheck.TypeConstraintWithDefaults) bool
	nullable            nullableRequirement
}

func validateVariableShape(shape variableShape) requiredVariableValidator {
	return func(rule *requiredVariableRule, runner tflint.Runner, block *hclext.Block) error {
		typeAttr, ok := block.Body.Attributes["type"]
		if !ok {
			return emitVariableAttributeIssue(rule, runner, block.DefRange, "type", shape.typeDescription)
		}

		typeConstraint, diags := varcheck.NewTypeConstraintWithDefaultsFromExp(typeAttr.Expr)
		if diags.HasErrors() || !shape.validateType(typeConstraint) {
			return emitVariableAttributeIssue(rule, runner, typeAttr.Range, "type", shape.typeDescription)
		}

		defaultAttr, ok := block.Body.Attributes["default"]
		if !ok {
			return emitVariableAttributeIssue(rule, runner, block.DefRange, "default", shape.defaultDescription)
		}
		defaultValue, diags := defaultAttr.Expr.Value(nil)
		if diags.HasErrors() || !defaultValue.IsKnown() || !shape.validateDefault(defaultValue, typeConstraint) {
			return emitVariableAttributeIssue(rule, runner, defaultAttr.Range, "default", shape.defaultDescription)
		}

		nullableAttr, ok := block.Body.Attributes["nullable"]
		if !ok {
			if shape.nullable == nullableMayBeTrue {
				return nil
			}
			return emitVariableAttributeIssue(rule, runner, block.DefRange, "nullable", shape.nullableDescription)
		}

		nullableValue, diags := nullableAttr.Expr.Value(nil)
		if diags.HasErrors() || !nullableValue.IsKnown() || nullableValue.IsNull() || !nullableValue.Type().Equals(cty.Bool) {
			return emitVariableAttributeIssue(rule, runner, nullableAttr.Range, "nullable", shape.nullableDescription)
		}

		if shape.nullable == nullableMustBeFalse && nullableValue.True() {
			return emitVariableAttributeIssue(rule, runner, nullableAttr.Range, "nullable", shape.nullableDescription)
		}
		if shape.nullable == nullableMayBeTrue && !nullableValue.True() {
			return emitVariableAttributeIssue(rule, runner, nullableAttr.Range, "nullable", shape.nullableDescription)
		}

		return nil
	}
}

func emitVariableAttributeIssue(
	rule *requiredVariableRule,
	runner tflint.Runner,
	sourceRange hcl.Range,
	attribute string,
	expectation string,
) error {
	return runner.EmitIssue(
		rule,
		fmt.Sprintf("variable `%s` %s %s; see: %s", rule.variableName, attribute, expectation, rule.link),
		sourceRange,
	)
}

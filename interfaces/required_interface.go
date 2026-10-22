package interfaces

import (
	"fmt"
	"strings"

	"github.com/Azure/tflint-ruleset-avm/internal/tfvarcheck/varcheck"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

var azapiResourceBodySchema = &hclext.BodySchema{
	Blocks: []hclext.BlockSchema{
		{
			Type:       "resource",
			LabelNames: []string{"type", "name"},
			Body: &hclext.BodySchema{
				Attributes: []hclext.AttributeSchema{
					{Name: "type"},
					{Name: "retry"},
					{Name: "timeouts"},
					{Name: "ignore_body_changes"},
				},
			},
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

type requiredVariableValidator func(*requiredVariableRule, tflint.Runner, *hclext.Block) (bool, error)

type requiredVariableWiringValidator func(*requiredVariableRule, tflint.Runner, []*hclext.Block) error

type requiredVariableRule struct {
	tflint.DefaultRule

	name          string
	variableName  string
	link          string
	requirement   string
	applicability requiredVariableApplicability
	validator     requiredVariableValidator
	wiring        requiredVariableWiringValidator
}

func newRequiredVariableRule(
	name string,
	variableName string,
	link string,
	requirement string,
	applicability requiredVariableApplicability,
	validator requiredVariableValidator,
	wiring ...requiredVariableWiringValidator,
) *requiredVariableRule {
	var wiringValidator requiredVariableWiringValidator
	if len(wiring) > 0 {
		wiringValidator = wiring[0]
	}
	return &requiredVariableRule{
		name:          name,
		variableName:  variableName,
		link:          link,
		requirement:   requirement,
		applicability: applicability,
		validator:     validator,
		wiring:        wiringValidator,
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
		if r.wiring != nil {
			sourceRange, err = missingVariableRange(runner)
			if err != nil {
				return err
			}
		}
		return runner.EmitIssue(
			r,
			fmt.Sprintf("variable `%s` must be declared %s; see: %s", r.variableName, r.requirement, r.link),
			sourceRange,
		)
	}

	valid, err := r.validator(r, runner, variable)
	if err != nil || !valid || r.wiring == nil {
		return err
	}

	resources, err := directAzapiResources(runner)
	if err != nil {
		return err
	}
	return r.wiring(r, runner, resources)
}

func appliesToDirectAzapiResource(runner tflint.Runner) (bool, hcl.Range, error) {
	resources, err := directAzapiResources(runner)
	if err != nil {
		return false, hcl.Range{}, err
	}
	if len(resources) == 0 {
		return false, hcl.Range{}, nil
	}
	return true, resources[0].DefRange, nil
}

func directAzapiResources(runner tflint.Runner) ([]*hclext.Block, error) {
	content, err := runner.GetModuleContent(
		azapiResourceBodySchema,
		&tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone},
	)
	if err != nil {
		return nil, err
	}

	resources := make([]*hclext.Block, 0, len(content.Blocks))
	for _, block := range content.Blocks {
		if len(block.Labels) != 2 {
			continue
		}
		if _, ok := mandatoryInterfaceAzapiResourceTypes[block.Labels[0]]; ok {
			resources = append(resources, block)
		}
	}
	return resources, nil
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

func missingVariableRange(runner tflint.Runner) (hcl.Range, error) {
	file, err := runner.GetFile("variables.tf")
	if err != nil {
		if strings.Contains(err.Error(), "file not found") {
			return hcl.Range{}, nil
		}
		return hcl.Range{}, err
	}
	if file == nil {
		return hcl.Range{}, nil
	}
	return file.Body.MissingItemRange(), nil
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
	return func(rule *requiredVariableRule, runner tflint.Runner, block *hclext.Block) (bool, error) {
		typeAttr, ok := block.Body.Attributes["type"]
		if !ok {
			return false, emitVariableAttributeIssue(rule, runner, block.DefRange, "type", shape.typeDescription)
		}

		typeConstraint, diags := varcheck.NewTypeConstraintWithDefaultsFromExp(typeAttr.Expr)
		if diags.HasErrors() || !shape.validateType(typeConstraint) {
			return false, emitVariableAttributeIssue(rule, runner, typeAttr.Range, "type", shape.typeDescription)
		}

		defaultAttr, ok := block.Body.Attributes["default"]
		if !ok {
			return false, emitVariableAttributeIssue(rule, runner, block.DefRange, "default", shape.defaultDescription)
		}
		defaultValue, diags := defaultAttr.Expr.Value(nil)
		if diags.HasErrors() || !defaultValue.IsKnown() || !shape.validateDefault(defaultValue, typeConstraint) {
			return false, emitVariableAttributeIssue(rule, runner, defaultAttr.Range, "default", shape.defaultDescription)
		}

		nullableAttr, ok := block.Body.Attributes["nullable"]
		if !ok {
			if shape.nullable == nullableMayBeTrue {
				return true, nil
			}
			return false, emitVariableAttributeIssue(rule, runner, block.DefRange, "nullable", shape.nullableDescription)
		}

		nullableValue, diags := nullableAttr.Expr.Value(nil)
		if diags.HasErrors() || !nullableValue.IsKnown() || nullableValue.IsNull() || !nullableValue.Type().Equals(cty.Bool) {
			return false, emitVariableAttributeIssue(rule, runner, nullableAttr.Range, "nullable", shape.nullableDescription)
		}

		if shape.nullable == nullableMustBeFalse && nullableValue.True() {
			return false, emitVariableAttributeIssue(rule, runner, nullableAttr.Range, "nullable", shape.nullableDescription)
		}
		if shape.nullable == nullableMayBeTrue && !nullableValue.True() {
			return false, emitVariableAttributeIssue(rule, runner, nullableAttr.Range, "nullable", shape.nullableDescription)
		}

		return true, nil
	}
}

func validateAzapiInterfaceWiring(attribute string, requiresMember bool) requiredVariableWiringValidator {
	return func(rule *requiredVariableRule, runner tflint.Runner, resources []*hclext.Block) error {
		locals, err := localExpressions(runner)
		if err != nil {
			return err
		}
		for _, resource := range resources {
			value, ok := resource.Body.Attributes[attribute]
			if ok && referencesInterfaceVariable(value.Expr, rule.variableName, requiresMember, locals, map[string]bool{}) {
				continue
			}
			sourceRange := resource.DefRange
			if ok {
				sourceRange = value.Range
			}
			if err := runner.EmitIssue(
				rule,
				fmt.Sprintf("AzAPI resource must set `%s` using variable `%s`; see: %s", attribute, rule.variableName, rule.link),
				sourceRange,
			); err != nil {
				return err
			}
		}
		return nil
	}
}

func localExpressions(runner tflint.Runner) (map[string]hcl.Expression, error) {
	content, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type: "locals",
				Body: &hclext.BodySchema{
					Mode: hclext.SchemaJustAttributesMode,
				},
			},
		},
	}, &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return nil, err
	}

	locals := make(map[string]hcl.Expression)
	for _, block := range content.Blocks {
		for name, attribute := range block.Body.Attributes {
			locals[name] = attribute.Expr
		}
	}
	return locals, nil
}

func referencesInterfaceVariable(
	expression hcl.Expression,
	variableName string,
	requiresMember bool,
	locals map[string]hcl.Expression,
	visitedLocals map[string]bool,
) bool {
	return expressionForwardsInterface(
		expression,
		nil,
		variableName,
		requiresMember,
		locals,
		visitedLocals,
	)
}

func expressionForwardsInterface(
	expression hcl.Expression,
	suffix hcl.Traversal,
	variableName string,
	requiresMember bool,
	locals map[string]hcl.Expression,
	visitedLocals map[string]bool,
) bool {
	switch expression := expression.(type) {
	case *hclsyntax.ScopeTraversalExpr:
		return traversalForwardsInterface(expression.Traversal, suffix, variableName, requiresMember, locals, visitedLocals)
	case *hclsyntax.RelativeTraversalExpr:
		return expressionForwardsInterface(
			expression.Source,
			append(append(hcl.Traversal{}, expression.Traversal...), suffix...),
			variableName,
			requiresMember,
			locals,
			visitedLocals,
		)
	case *hclsyntax.ParenthesesExpr:
		return expressionForwardsInterface(expression.Expression, suffix, variableName, requiresMember, locals, visitedLocals)
	case *hclsyntax.TemplateWrapExpr:
		return expressionForwardsInterface(expression.Wrapped, suffix, variableName, requiresMember, locals, visitedLocals)
	case *hclsyntax.TemplateExpr:
		for _, part := range expression.Parts {
			if expressionForwardsInterface(part, suffix, variableName, requiresMember, locals, visitedLocals) {
				return true
			}
		}
	case *hclsyntax.ConditionalExpr:
		return expressionForwardsInterface(expression.TrueResult, suffix, variableName, requiresMember, locals, visitedLocals) ||
			expressionForwardsInterface(expression.FalseResult, suffix, variableName, requiresMember, locals, visitedLocals)
	case *hclsyntax.FunctionCallExpr:
		for index, argument := range expression.Args {
			argumentSuffix := suffix
			if strings.EqualFold(expression.Name, "lookup") && index == 0 && requiresMember {
				argumentSuffix = append(hcl.Traversal{hcl.TraverseAttr{Name: "lookup"}}, suffix...)
			}
			if expressionForwardsInterface(argument, argumentSuffix, variableName, requiresMember, locals, visitedLocals) {
				return true
			}
		}
	case *hclsyntax.ObjectConsExpr:
		if len(suffix) == 0 {
			for _, item := range expression.Items {
				if expressionForwardsInterface(item.ValueExpr, nil, variableName, requiresMember, locals, visitedLocals) {
					return true
				}
			}
			return false
		}
		member, ok := traversalMemberName(suffix[0])
		if !ok {
			return false
		}
		for _, item := range expression.Items {
			key, diagnostics := item.KeyExpr.Value(nil)
			if diagnostics.HasErrors() || !key.IsKnown() || key.IsNull() || !key.Type().Equals(cty.String) || key.AsString() != member {
				continue
			}
			return expressionForwardsInterface(item.ValueExpr, suffix[1:], variableName, requiresMember, locals, visitedLocals)
		}
	default:
		traversal, diagnostics := hcl.AbsTraversalForExpr(expression)
		if !diagnostics.HasErrors() {
			return traversalForwardsInterface(traversal, suffix, variableName, requiresMember, locals, visitedLocals)
		}
	}
	return false
}

func traversalForwardsInterface(
	traversal hcl.Traversal,
	suffix hcl.Traversal,
	variableName string,
	requiresMember bool,
	locals map[string]hcl.Expression,
	visitedLocals map[string]bool,
) bool {
	if len(traversal) < 2 {
		return false
	}
	root, ok := traversal[0].(hcl.TraverseRoot)
	if !ok {
		return false
	}

	name, ok := traversal[1].(hcl.TraverseAttr)
	if !ok {
		return false
	}
	switch root.Name {
	case "var":
		if name.Name != variableName {
			return false
		}
		return !requiresMember || len(traversal)+len(suffix) > 2
	case "local":
		if visitedLocals[name.Name] {
			return false
		}
		localExpression, ok := locals[name.Name]
		if !ok {
			return false
		}
		nextVisitedLocals := make(map[string]bool, len(visitedLocals)+1)
		for localName, visited := range visitedLocals {
			nextVisitedLocals[localName] = visited
		}
		nextVisitedLocals[name.Name] = true
		remaining := append(append(hcl.Traversal{}, traversal[2:]...), suffix...)
		return expressionForwardsInterface(localExpression, remaining, variableName, requiresMember, locals, nextVisitedLocals)
	default:
		return false
	}
}

func traversalMemberName(traverser hcl.Traverser) (string, bool) {
	switch traverser := traverser.(type) {
	case hcl.TraverseAttr:
		return traverser.Name, true
	case hcl.TraverseIndex:
		if traverser.Key.IsKnown() && !traverser.Key.IsNull() && traverser.Key.Type().Equals(cty.String) {
			return traverser.Key.AsString(), true
		}
	}
	return "", false
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

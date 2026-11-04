package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

var _ tflint.Rule = new(TerraformResourceIDVariableValidationRule)

const parseResourceIDFunction = "provider::azapi::parse_resource_id"

// handRolledResourceIDFunctions lists functions that must not be used to validate resource IDs. `length` is checked separately because it is a valid element count check on collections.
var handRolledResourceIDFunctions = map[string]bool{
	"regex":      true,
	"regexall":   true,
	"startswith": true,
	"split":      true,
}

// TerraformResourceIDVariableValidationRule checks that variables holding Azure resource IDs are validated with `provider::azapi::parse_resource_id`.
type TerraformResourceIDVariableValidationRule struct {
	tflint.DefaultRule
}

// NewTerraformResourceIDVariableValidationRule returns a new rule.
func NewTerraformResourceIDVariableValidationRule() *TerraformResourceIDVariableValidationRule {
	return new(TerraformResourceIDVariableValidationRule)
}

func (r *TerraformResourceIDVariableValidationRule) Name() string {
	return "avm_terraform_resource_id_variable_validation_required"
}

func (r *TerraformResourceIDVariableValidationRule) Enabled() bool {
	return true
}

func (r *TerraformResourceIDVariableValidationRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *TerraformResourceIDVariableValidationRule) Link() string {
	return "https://azure.github.io/Azure-Verified-Modules/spec/TFNFR38/"
}

func (r *TerraformResourceIDVariableValidationRule) Check(runner tflint.Runner) error {
	blocks, err := moduleBlocks(runner)
	if err != nil {
		return err
	}
	// parse_resource_id is an AzAPI provider function, so only modules that require the AzAPI provider can use it.
	if !requiresAzapiProvider(blocks) {
		return nil
	}

	var errList error
	emit := func(message string, issueRange hcl.Range) {
		if subErr := runner.EmitIssue(r, message, issueRange); subErr != nil {
			errList = multierror.Append(errList, subErr)
		}
	}
	for _, block := range blocks {
		if block.Type == "variable" && len(block.Labels) > 0 {
			r.checkVariable(block, emit)
		}
	}
	return errList
}

func (r *TerraformResourceIDVariableValidationRule) checkVariable(block *hclsyntax.Block, emit func(string, hcl.Range)) {
	name := block.Labels[0]
	kind := resourceIDVariableKind(block)
	needsNullCheck := kind == resourceIDString && isOptionalNullable(block)

	validated := false
	for _, validation := range block.Body.Blocks {
		if validation.Type != "validation" {
			continue
		}
		condition, ok := validation.Body.Attributes["condition"]
		if !ok {
			continue
		}
		walker := &validationConditionWalker{}
		// The walker never returns diagnostics.
		_ = hclsyntax.Walk(condition.Expr, walker)

		for _, call := range walker.parseCalls {
			if !isLiteralString(call.call.Args) {
				emit(fmt.Sprintf("The resource type passed to `%s` in the validation of variable `%s` must be a literal string", parseResourceIDFunction, name), call.call.Range())
			}
			if !call.inCan {
				emit(fmt.Sprintf("`%s` in the validation of variable `%s` must be wrapped in `can(...)`", parseResourceIDFunction, name), call.call.Range())
			}
			if kind == resourceIDCollection && !call.inFor {
				emit(fmt.Sprintf("The validation of variable `%s` must check each element with `alltrue([for ... : can(%s(...))])`", name, parseResourceIDFunction), call.call.Range())
			}
		}
		if kind == resourceIDNone {
			continue
		}
		if len(walker.parseCalls) > 0 {
			validated = true
			if needsNullCheck && !walker.comparesVariableToNull(name) {
				emit(fmt.Sprintf("The validation of optional variable `%s` must allow `null`, for example `var.%s == null || can(%s(...))`", name, name, parseResourceIDFunction), condition.SrcRange)
			}
		}
		for _, call := range walker.calls {
			if handRolledResourceIDFunctions[call.Name] || (call.Name == "length" && kind == resourceIDString) {
				emit(fmt.Sprintf("Variable `%s` must not validate its resource ID with `%s`; use `%s` instead", name, call.Name, parseResourceIDFunction), call.Range())
			}
		}
	}
	if kind != resourceIDNone && !validated {
		emit(fmt.Sprintf("Variable `%s` holds a resource ID and must be validated with `can(%s(\"<resource type>\", ...))`", name, parseResourceIDFunction), block.DefRange())
	}
}

type resourceIDKind int

const (
	resourceIDNone resourceIDKind = iota
	resourceIDString
	resourceIDCollection
)

// resourceIDVariableKind classifies a variable that holds resource IDs by its name and type. Only `string` and `list`, `set`, or `map` of `string` types are checked.
func resourceIDVariableKind(block *hclsyntax.Block) resourceIDKind {
	name := block.Labels[0]
	isIDName := name == "parent_id" ||
		name == "resource_id" || strings.HasSuffix(name, "_resource_id") ||
		name == "resource_ids" || strings.HasSuffix(name, "_resource_ids")
	if !isIDName {
		return resourceIDNone
	}
	typeName, args := variableType(block)
	switch typeName {
	case "string":
		return resourceIDString
	case "list", "set", "map":
		if len(args) == 1 && hcl.ExprAsKeyword(args[0]) == "string" {
			return resourceIDCollection
		}
	}
	return resourceIDNone
}

// isOptionalNullable reports whether a variable defaults to `null` and doesn't set `nullable = false`.
func isOptionalNullable(block *hclsyntax.Block) bool {
	if nullable, ok := staticBool(block, "nullable"); ok && !nullable {
		return false
	}
	attr, ok := block.Body.Attributes["default"]
	if !ok {
		return false
	}
	v, diags := attr.Expr.Value(nil)
	return !diags.HasErrors() && v.IsNull()
}

// isLiteralString reports whether the first argument is a literal string with no interpolation.
func isLiteralString(args []hclsyntax.Expression) bool {
	if len(args) == 0 {
		return false
	}
	switch arg := args[0].(type) {
	case *hclsyntax.TemplateExpr:
		return arg.IsStringLiteral()
	case *hclsyntax.LiteralValueExpr:
		return arg.Val.Type() == cty.String
	}
	return false
}

// requiresAzapiProvider reports whether a terraform block declares the AzAPI provider in required_providers.
func requiresAzapiProvider(blocks []*hclsyntax.Block) bool {
	for _, block := range blocks {
		if block.Type != "terraform" {
			continue
		}
		for _, nested := range block.Body.Blocks {
			if nested.Type != "required_providers" {
				continue
			}
			for name, attr := range nested.Body.Attributes {
				if name == "azapi" || providerSource(attr) == "azure/azapi" {
					return true
				}
			}
		}
	}
	return false
}

// providerSource returns the lowercase source of a required_providers entry without its registry host.
func providerSource(attr *hclsyntax.Attribute) string {
	v, diags := attr.Expr.Value(nil)
	if diags.HasErrors() || !v.IsKnown() || v.IsNull() || !v.Type().IsObjectType() || !v.Type().HasAttribute("source") {
		return ""
	}
	source := v.GetAttr("source")
	if !source.IsKnown() || source.IsNull() || !source.Type().Equals(cty.String) {
		return ""
	}
	return strings.TrimPrefix(strings.ToLower(source.AsString()), "registry.terraform.io/")
}

type parseResourceIDCall struct {
	call  *hclsyntax.FunctionCallExpr
	inCan bool
	inFor bool
}

// validationConditionWalker collects function calls and null comparisons in a validation condition, tracking the ancestors of each parse_resource_id call.
type validationConditionWalker struct {
	ancestors   []hclsyntax.Node
	calls       []*hclsyntax.FunctionCallExpr
	parseCalls  []parseResourceIDCall
	comparisons []*hclsyntax.BinaryOpExpr
}

func (w *validationConditionWalker) Enter(node hclsyntax.Node) hcl.Diagnostics {
	switch n := node.(type) {
	case *hclsyntax.FunctionCallExpr:
		w.calls = append(w.calls, n)
		if n.Name == parseResourceIDFunction {
			w.parseCalls = append(w.parseCalls, w.parseCall(n))
		}
	case *hclsyntax.BinaryOpExpr:
		if n.Op == hclsyntax.OpEqual || n.Op == hclsyntax.OpNotEqual {
			w.comparisons = append(w.comparisons, n)
		}
	}
	w.ancestors = append(w.ancestors, node)
	return nil
}

func (w *validationConditionWalker) Exit(hclsyntax.Node) hcl.Diagnostics {
	w.ancestors = w.ancestors[:len(w.ancestors)-1]
	return nil
}

func (w *validationConditionWalker) parseCall(call *hclsyntax.FunctionCallExpr) parseResourceIDCall {
	result := parseResourceIDCall{call: call}
	for _, ancestor := range w.ancestors {
		switch a := ancestor.(type) {
		case *hclsyntax.FunctionCallExpr:
			if a.Name == "can" {
				result.inCan = true
			}
		case *hclsyntax.ForExpr:
			result.inFor = true
		}
	}
	return result
}

// comparesVariableToNull reports whether the condition compares `var.<name>` with `null`.
func (w *validationConditionWalker) comparesVariableToNull(name string) bool {
	for _, comparison := range w.comparisons {
		if isVariableAndNull(comparison.LHS, comparison.RHS, name) || isVariableAndNull(comparison.RHS, comparison.LHS, name) {
			return true
		}
	}
	return false
}

func isVariableAndNull(variable, null hclsyntax.Expression, name string) bool {
	ref, ok := variableReference(variable)
	if !ok || ref != name {
		return false
	}
	literal, ok := unwrapParens(null).(*hclsyntax.LiteralValueExpr)
	return ok && literal.Val.IsNull()
}

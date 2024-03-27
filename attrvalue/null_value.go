package attrvalue

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// NullRule checks whether an attribute value is null or part of a variable with no default value.
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
	// We cannot use runner.EvaluateExpr here because it needs a concrete type to evaluate the expression into.
	// This rule will work on any value type, so we need to use (hcl.Expression).Value instead.
	// To do this we need the evaluation context to include the variables, so we need to evaluate the variables first.
	evalContext, err := getDefaultVariableEvalContext(runner)
	if err != nil {
		return err
	}

	// Get the resources and extract the attribute values.
	attrs, err := getSimpleAttrs(runner, r.resourceType, r.attributeName)
	if err != nil {
		return err
	}

	// Check the attribute values.
	for _, attr := range attrs {
		attrVal, diags := attr.Expr.Value(evalContext)
		if diags.HasErrors() {
			return fmt.Errorf("could not evaluate attribute %s value: %s", r.attributeName, diags)
		}
		if attrVal.IsNull() {
			continue
		}
		if err = runner.EmitIssue(
			r,
			fmt.Sprintf("invalid attribute value of `%s` - expecting null", r.attributeName),
			attr.Range,
		); err != nil {
			return err
		}
	}
	return nil
}

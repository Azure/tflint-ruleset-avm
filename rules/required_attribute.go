package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// RequiredAttributeRule is a generic rule to enforce that a particular attribute is present
// on blocks of a given type whose first label (usually the resource type) matches an allow list.
// This is intentionally generic so future required-attribute rules can be added with a single
// line registration instead of bespoke implementations.
//
// Example enforced block shape: resource "<firstLabel>" "name" { <attributeName> = <...> }
// Only blocks whose first label value appears in allowedFirstLabelValues are checked.
// When the attribute is missing an issue is emitted including a suggested default value text.
type RequiredAttributeRule struct {
	tflint.DefaultRule

	// blockType is the HCL block type (e.g. "resource", "data")
	blockType string
	// allowedFirstLabelValues restricts enforcement to blocks whose first label value is in the set.
	allowedFirstLabelValues map[string]struct{}
	// attributeName is the attribute that must be present.
	attributeName string
	// defaultSuggestion is text describing a default value (NOT auto-fix) shown in the issue message.
	defaultSuggestion string
	// name is the unique rule name.
	name string
	// link is documentation URL.
	link string
	// severity of the rule (defaults to ERROR when empty).
	severity tflint.Severity
	// optional value validator executed when attribute is present
	valueValidator func(run tflint.Runner, rule *RequiredAttributeRule, attr *hclext.Attribute) error
}

// NewAzapiRequiredAttributeRule constructs a new required attribute rule.
type RequiredAttributeOption func(*RequiredAttributeRule)

// WithValueValidator registers a callback to validate attribute value (only runs if attribute exists).
func WithValueValidator(fn func(run tflint.Runner, rule *RequiredAttributeRule, attr *hclext.Attribute) error) RequiredAttributeOption {
	return func(r *RequiredAttributeRule) { r.valueValidator = fn }
}

// DisallowWildcardList returns an option that forbids a single-item list containing "*" for the given attribute.
func DisallowWildcardList(attributeName string) RequiredAttributeOption {
	return WithValueValidator(func(run tflint.Runner, rule *RequiredAttributeRule, attr *hclext.Attribute) error {
		var values []string
		wantType := cty.List(cty.String)
		if err := run.EvaluateExpr(attr.Expr, &values, &tflint.EvaluateExprOption{WantType: &wantType}); err == nil {
			if len(values) == 1 && values[0] == "*" {
				msg := fmt.Sprintf("`%s` must not contain the wildcard \"*\"; explicitly list required fields or use empty list []", attributeName)
				return run.EmitIssue(rule, msg, attr.Range)
			}
		}
		return nil
	})
}

func NewRequiredAttributeRule(name, link, blockType string, allowedFirstLabels []string, attributeName, defaultSuggestion string, severity tflint.Severity, opts ...RequiredAttributeOption) *RequiredAttributeRule {
	set := make(map[string]struct{}, len(allowedFirstLabels))
	for _, v := range allowedFirstLabels {
		set[v] = struct{}{}
	}
	if severity == 0 { // treat unknown as ERROR
		severity = tflint.ERROR
	}
	r := &RequiredAttributeRule{
		blockType:               blockType,
		allowedFirstLabelValues: set,
		attributeName:           attributeName,
		defaultSuggestion:       defaultSuggestion,
		name:                    name,
		link:                    link,
		severity:                severity,
	}
	for _, o := range opts {
		o(r)
	}
	return r
}

func (r *RequiredAttributeRule) Name() string              { return r.name }
func (r *RequiredAttributeRule) Link() string              { return r.link }
func (r *RequiredAttributeRule) Enabled() bool             { return true }
func (r *RequiredAttributeRule) Severity() tflint.Severity { return r.severity }

// dynamic schema built per rule (we only care about presence, not value shape).
func (r *RequiredAttributeRule) schema() *hclext.BodySchema {
	return &hclext.BodySchema{ // e.g. resource "type" "name" {}
		Blocks: []hclext.BlockSchema{{
			Type:       r.blockType,
			LabelNames: []string{"type", "name"},
			Body: &hclext.BodySchema{Attributes: []hclext.AttributeSchema{{
				Name: r.attributeName,
			}}},
		}},
	}
}

func (r *RequiredAttributeRule) Check(run tflint.Runner) error {
	path, err := run.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() { // only enforce at root, consistent with existing azapi rules
		return nil
	}

	content, err := run.GetModuleContent(r.schema(), &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return err
	}

	var errs error
	for _, b := range content.Blocks {
		if b.Type != r.blockType || len(b.Labels) < 2 { // defensive
			continue
		}
		first := b.Labels[0]
		if _, ok := r.allowedFirstLabelValues[first]; !ok {
			continue
		}
		attr, ok := b.Body.Attributes[r.attributeName]
		if ok {
			if r.valueValidator != nil {
				if e := r.valueValidator(run, r, attr); e != nil {
					errs = multierror.Append(errs, e)
				}
			}
			continue
		}

		suggestion := r.defaultSuggestion
		if strings.TrimSpace(suggestion) != "" {
			suggestion = fmt.Sprintf(" (suggested default: %s)", suggestion)
		}
		msg := fmt.Sprintf("`%s` attribute must be specified%s", r.attributeName, suggestion)
		if e := run.EmitIssue(r, msg, b.DefRange); e != nil {
			errs = multierror.Append(errs, e)
		}
	}
	return errs
}

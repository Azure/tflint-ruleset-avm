package rules

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/jmespath/go-jmespath"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

var _ tflint.Rule = new(AzapiReplaceTriggersRefsRule)

// AzapiReplaceTriggersRefsRule validates replacement paths when they are declared.
type AzapiReplaceTriggersRefsRule struct {
	tflint.DefaultRule
}

// NewAzapiReplaceTriggersRefsRule returns an AzAPI replacement-trigger rule.
func NewAzapiReplaceTriggersRefsRule() *AzapiReplaceTriggersRefsRule {
	return new(AzapiReplaceTriggersRefsRule)
}

func (r *AzapiReplaceTriggersRefsRule) Name() string {
	return "azapi_replace_triggers_refs"
}

func (r *AzapiReplaceTriggersRefsRule) Link() string {
	return "https://azure.github.io/Azure-Verified-Modules/spec/TFFR5/"
}

func (r *AzapiReplaceTriggersRefsRule) Enabled() bool {
	return true
}

func (r *AzapiReplaceTriggersRefsRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *AzapiReplaceTriggersRefsRule) Check(runner tflint.Runner) error {
	modulePath, err := runner.GetModulePath()
	if err != nil {
		return err
	}
	if !modulePath.IsRoot() {
		return nil
	}

	content, err := runner.GetModuleContent(azapiReplaceTriggersBodySchema, &tflint.GetModuleContentOption{
		ExpandMode: tflint.ExpandModeNone,
	})
	if err != nil {
		return err
	}

	var errs error
	for _, block := range content.Blocks {
		if len(block.Labels) != 2 || !replaceTriggersResourceTypes[block.Labels[0]] {
			continue
		}
		attribute, exists := block.Body.Attributes["replace_triggers_refs"]
		if !exists {
			continue
		}

		paths, ok := staticStringList(attribute.Expr)
		if !ok {
			if err := runner.EmitIssue(
				r,
				"`replace_triggers_refs` must be a statically known list of JMESPath expressions",
				attribute.Range,
			); err != nil {
				errs = multierror.Append(errs, err)
			}
			continue
		}
		if len(paths) == 0 {
			if err := runner.EmitIssue(
				r,
				"omit `replace_triggers_refs` when no body paths require replacement",
				attribute.Range,
			); err != nil {
				errs = multierror.Append(errs, err)
			}
			continue
		}

		compiled := make(map[string]*jmespath.JMESPath, len(paths))
		seen := make(map[string]struct{}, len(paths))
		for _, path := range paths {
			normalized := strings.TrimSpace(path)
			switch {
			case normalized == "":
				errs = appendTriggerIssue(errs, runner, r, "replacement paths must not be blank", attribute.Range)
				continue
			case normalized != path:
				errs = appendTriggerIssue(errs, runner, r, fmt.Sprintf("replacement path %q must not contain surrounding whitespace", path), attribute.Range)
				continue
			case normalized == "name" || normalized == "location":
				errs = appendTriggerIssue(errs, runner, r, fmt.Sprintf("replacement path %q is redundant because AzAPI already replaces resources when it changes", normalized), attribute.Range)
				continue
			}
			if _, duplicate := seen[normalized]; duplicate {
				errs = appendTriggerIssue(errs, runner, r, fmt.Sprintf("replacement path %q is duplicated", normalized), attribute.Range)
				continue
			}
			seen[normalized] = struct{}{}

			expression, err := jmespath.Compile(normalized)
			if err != nil {
				errs = appendTriggerIssue(errs, runner, r, fmt.Sprintf("replacement path %q is not valid JMESPath: %v", normalized, err), attribute.Range)
				continue
			}
			compiled[normalized] = expression
		}

		if len(compiled) == 0 {
			continue
		}
		var body any = map[string]any{}
		if bodyAttribute, hasBody := block.Body.Attributes["body"]; hasBody {
			var known bool
			body, known = staticJSONValue(bodyAttribute.Expr)
			if !known {
				continue
			}
		}
		for path, expression := range compiled {
			result, err := expression.Search(body)
			if err != nil {
				errs = appendTriggerIssue(
					errs,
					runner,
					r,
					fmt.Sprintf("replacement path %q cannot be evaluated against the static resource body: %v", path, err),
					attribute.Range,
				)
				continue
			}
			projectionChecked, projectionResolved := simpleProjectionPathResolves(path, body)
			if (projectionChecked && !projectionResolved) ||
				(!projectionChecked && result == nil && !simpleObjectPathExists(path, body)) {
				errs = appendTriggerIssue(
					errs,
					runner,
					r,
					fmt.Sprintf("replacement path %q does not resolve against the static resource body", path),
					attribute.Range,
				)
			}
		}
	}
	return errs
}

var replaceTriggersResourceTypes = map[string]bool{
	"azapi_resource":            true,
	"azapi_data_plane_resource": true,
}

var azapiReplaceTriggersBodySchema = &hclext.BodySchema{
	Blocks: []hclext.BlockSchema{
		{
			Type:       "resource",
			LabelNames: []string{"type", "name"},
			Body: &hclext.BodySchema{
				Attributes: []hclext.AttributeSchema{
					{Name: "body"},
					{Name: "replace_triggers_refs"},
				},
			},
		},
	},
}

func staticStringList(expression hcl.Expression) ([]string, bool) {
	value, diagnostics := expression.Value(nil)
	valueType := value.Type()
	if diagnostics.HasErrors() ||
		value.IsNull() ||
		!value.IsKnown() ||
		(!valueType.IsTupleType() && !valueType.IsListType() && !valueType.IsSetType()) {
		return nil, false
	}

	values := make([]string, 0, value.LengthInt())
	iterator := value.ElementIterator()
	for iterator.Next() {
		_, item := iterator.Element()
		if item.IsNull() || !item.IsKnown() || item.Type() != cty.String {
			return nil, false
		}
		values = append(values, item.AsString())
	}
	return values, true
}

func staticJSONValue(expression hcl.Expression) (any, bool) {
	value, diagnostics := expression.Value(nil)
	if diagnostics.HasErrors() || !value.IsKnown() {
		return nil, false
	}
	if value.IsNull() {
		return map[string]any{}, true
	}
	data, err := ctyjson.Marshal(value, value.Type())
	if err != nil {
		return nil, false
	}
	var decoded any
	if err := json.Unmarshal(data, &decoded); err != nil {
		return nil, false
	}
	return decoded, true
}

func simpleObjectPathExists(path string, value any) bool {
	current := value
	for _, segment := range strings.Split(path, ".") {
		if !validSimpleIdentifier(segment) {
			return false
		}
		object, ok := current.(map[string]any)
		if !ok {
			return false
		}
		current, ok = object[segment]
		if !ok {
			return false
		}
	}
	return true
}

func simpleProjectionPathResolves(path string, value any) (bool, bool) {
	if !strings.Contains(path, "[]") {
		return false, false
	}
	segments := strings.Split(path, ".")
	return walkSimpleProjectionPath(segments, value)
}

func walkSimpleProjectionPath(segments []string, value any) (bool, bool) {
	if len(segments) == 0 {
		return true, true
	}
	segment := segments[0]
	if strings.HasSuffix(segment, "[]") {
		name := strings.TrimSuffix(segment, "[]")
		if !validSimpleIdentifier(name) {
			return false, false
		}
		object, ok := value.(map[string]any)
		if !ok {
			return true, false
		}
		child, exists := object[name]
		if !exists {
			return true, false
		}
		items, ok := child.([]any)
		if !ok {
			return true, false
		}
		for _, item := range items {
			checked, resolved := walkSimpleProjectionPath(segments[1:], item)
			if !checked {
				return false, false
			}
			if !resolved {
				return true, false
			}
		}
		return true, true
	}
	if !validSimpleIdentifier(segment) {
		return false, false
	}
	object, ok := value.(map[string]any)
	if !ok {
		return true, false
	}
	child, exists := object[segment]
	if !exists {
		return true, false
	}
	return walkSimpleProjectionPath(segments[1:], child)
}

func validSimpleIdentifier(identifier string) bool {
	if identifier == "" {
		return false
	}
	for index, character := range identifier {
		if (character >= 'a' && character <= 'z') ||
			(character >= 'A' && character <= 'Z') ||
			character == '_' ||
			(index > 0 && character >= '0' && character <= '9') {
			continue
		}
		return false
	}
	return true
}

func appendTriggerIssue(errs error, runner tflint.Runner, rule tflint.Rule, message string, sourceRange hcl.Range) error {
	if err := runner.EmitIssue(rule, message, sourceRange); err != nil {
		return multierror.Append(errs, err)
	}
	return errs
}

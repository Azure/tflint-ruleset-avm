package outputs

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// Schema to capture resource and output blocks.
var noEntireResourceOutputBodySchema = &hclext.BodySchema{
	Blocks: []hclext.BlockSchema{
		{ // resource "type" "name" {}
			Type:       "resource",
			LabelNames: []string{"type", "name"},
		},
		{ // output "name" { value = ... }
			Type:       "output",
			LabelNames: []string{"name"},
			Body:       &hclext.BodySchema{Attributes: []hclext.AttributeSchema{{Name: "value"}}},
		},
	},
}

var _ tflint.Rule = new(NoEntireResourceOutputRule)

// NoEntireResourceOutputRule checks that an output does not expose an entire resource object.
type NoEntireResourceOutputRule struct {
	tflint.DefaultRule
	ruleName string
	link     string
}

// NewNoEntireResourceOutputRule returns the rule.
func NewNoEntireResourceOutputRule(ruleName, link string) *NoEntireResourceOutputRule {
	return &NoEntireResourceOutputRule{ruleName: ruleName, link: link}
}

func (r *NoEntireResourceOutputRule) Name() string              { return r.ruleName }
func (r *NoEntireResourceOutputRule) Link() string              { return r.link }
func (r *NoEntireResourceOutputRule) Enabled() bool             { return true }
func (r *NoEntireResourceOutputRule) Severity() tflint.Severity { return tflint.ERROR }

func (r *NoEntireResourceOutputRule) Check(run tflint.Runner) error {
	path, err := run.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() { // only evaluate root module like other rules
		return nil
	}

	content, err := run.GetModuleContent(noEntireResourceOutputBodySchema, &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return err
	}

	// Build a set of resource identifiers: type.name
	resources := map[string]struct{}{}
	for _, b := range content.Blocks {
		if b.Type == "resource" && len(b.Labels) == 2 { // type, name
			resources[fmt.Sprintf("%s.%s", b.Labels[0], b.Labels[1])] = struct{}{}
		}
	}

	var errs error
	for _, b := range content.Blocks {
		if b.Type != "output" {
			continue
		}
		attr, ok := b.Body.Attributes["value"]
		if !ok { // nothing to check
			continue
		}
		if isWholeResourceExpression(attr.Expr, resources) {
			if e := run.EmitIssue(r, "authors SHOULD NOT output entire resource objects; expose specific computed attributes instead", attr.Range); e != nil {
				errs = multierror.Append(errs, e)
			}
		}
	}
	return errs
}

// isWholeResourceTraversal returns true if traversal represents a resource base object optionally followed by only index or splat traversers (i.e. still the entire object or collection of objects).
func isWholeResourceTraversal(traversal hcl.Traversal, resources map[string]struct{}) bool {
	if len(traversal) < 2 {
		return false
	}
	root, ok := traversal[0].(hcl.TraverseRoot)
	if !ok {
		return false
	}
	nameAttr, ok := traversal[1].(hcl.TraverseAttr)
	if !ok {
		return false
	}
	key := fmt.Sprintf("%s.%s", root.Name, nameAttr.Name)
	if _, exists := resources[key]; !exists {
		return false
	}
	// (debug tracing removed in finalized implementation)
	// Remaining traversers (if any) must all be index or splat; BUT if an attribute appears after a splat/index chain we treat it as narrowed (acceptable)
	encounteredAttr := false
	for _, tr := range traversal[2:] {
		switch tr.(type) {
		case hcl.TraverseIndex, hcl.TraverseSplat:
			if encounteredAttr { // once narrowed by attr, further index/splat shouldn't re-widen
				return false
			}
			continue
		case hcl.TraverseAttr:
			// attribute access narrows to a specific property => not whole resource object(s)
			return false
		default:
			return false
		}
	}
	return true
}

// isCompositeWholeResource unwraps index/splat expression layers until reaching a traversal and evaluates it.
func isWholeResourceExpression(expr hcl.Expression, resources map[string]struct{}) bool {
	switch e := expr.(type) {
	case *hclsyntax.ScopeTraversalExpr:
		return isWholeResourceTraversal(e.Traversal, resources)
	case *hclsyntax.RelativeTraversalExpr:
		// RelativeTraversalExpr has a source expression and additional traversers appended.
		if base := isWholeResourceExpression(e.Source, resources); base {
			// If base already whole, any added index/splat keeps it whole; any attr makes it partial.
			for _, tr := range e.Traversal {
				switch tr.(type) {
				case hcl.TraverseIndex, hcl.TraverseSplat:
					continue
				default:
					return false
				}
			}
			return true
		}
		// Or base not whole: attempt to treat combined traversal if source is plain
		if se, ok := e.Source.(*hclsyntax.ScopeTraversalExpr); ok {
			combined := append(se.Traversal, e.Traversal...)
			return isWholeResourceTraversal(combined, resources)
		}
		return false
	case *hclsyntax.IndexExpr:
		return isWholeResourceExpression(e.Collection, resources)
	case *hclsyntax.SplatExpr:
		// Determine if this is a full splat (exposes entire resource objects) or a projection.
		// Full splat:    collection[*]
		// Projection:    collection[*].attribute  => represented by Each being a RelativeTraversalExpr whose Source is the element symbol.
		baseWhole := isWholeResourceExpression(e.Source, resources)
		if !baseWhole {
			return false
		}
		// Classify Each
		switch each := e.Each.(type) {
		case *hclsyntax.AnonSymbolExpr:
			// Full splat: still whole
			return true
		case *hclsyntax.RelativeTraversalExpr:
			if _, ok := each.Source.(*hclsyntax.AnonSymbolExpr); ok {
				// Attribute (or further) traversal off element => projection (narrowed)
				return false
			}
			return true // Unexpected shape; be conservative and treat as whole
		default:
			// Unexpected Each expression kind; conservatively treat as whole
			return true
		}
	}
	return false
}

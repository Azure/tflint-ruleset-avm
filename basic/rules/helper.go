package rules

import (
	"maps"
	"path/filepath"
	"slices"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// ForFiles runs action for every file in the module, in file name order so issues are reported deterministically.
func ForFiles(runner tflint.Runner, action func(tflint.Runner, *hcl.File) error) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}
	for _, name := range slices.Sorted(maps.Keys(files)) {
		if subErr := action(runner, files[name]); subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}

// fileBaseName returns the base name of the file that contains body.
func fileBaseName(body *hclsyntax.Body) string {
	return filepath.Base(body.SrcRange.Filename)
}

// staticValue evaluates a literal attribute value of the wanted type. It returns false when the value is not a known, non-null literal of that type.
func staticValue(attr *hclsyntax.Attribute, want cty.Type) (cty.Value, bool) {
	v, diags := attr.Expr.Value(nil)
	if diags.HasErrors() || !v.IsKnown() || v.IsNull() || !v.Type().Equals(want) {
		return cty.NilVal, false
	}
	return v, true
}

// staticBool returns the literal boolean value of a block attribute.
func staticBool(block *hclsyntax.Block, name string) (value bool, ok bool) {
	attr, set := block.Body.Attributes[name]
	if !set {
		return false, false
	}
	v, ok := staticValue(attr, cty.Bool)
	if !ok {
		return false, false
	}
	return v.True(), true
}

// staticString returns the literal string value of a block attribute.
func staticString(block *hclsyntax.Block, name string) (value string, ok bool) {
	attr, set := block.Body.Attributes[name]
	if !set {
		return "", false
	}
	v, ok := staticValue(attr, cty.String)
	if !ok {
		return "", false
	}
	return v.AsString(), true
}

// variableReference returns the name of the input variable that expr references directly, as in `var.name`.
func variableReference(expr hcl.Expression) (string, bool) {
	traversal, diags := hcl.AbsTraversalForExpr(expr)
	if diags.HasErrors() || len(traversal) != 2 {
		return "", false
	}
	root, ok := traversal[0].(hcl.TraverseRoot)
	if !ok || root.Name != "var" {
		return "", false
	}
	attr, ok := traversal[1].(hcl.TraverseAttr)
	if !ok {
		return "", false
	}
	return attr.Name, true
}

// variableType returns the top-level type keyword or type constructor of a variable block, such as `string` or `map`, plus the constructor arguments. It returns an empty name when the type is omitted.
func variableType(block *hclsyntax.Block) (string, []hclsyntax.Expression) {
	attr, ok := block.Body.Attributes["type"]
	if !ok {
		return "", nil
	}
	if call, ok := unwrapParens(attr.Expr).(*hclsyntax.FunctionCallExpr); ok {
		return call.Name, call.Args
	}
	return hcl.ExprAsKeyword(attr.Expr), nil
}

// collectionTypeNames lists the type constructors and legacy keywords that declare collection-shaped variables.
var collectionTypeNames = map[string]bool{
	"list":  true,
	"set":   true,
	"map":   true,
	"tuple": true,
}

// unwrapParens strips any parentheses around expr.
func unwrapParens(expr hclsyntax.Expression) hclsyntax.Expression {
	for {
		parens, ok := expr.(*hclsyntax.ParenthesesExpr)
		if !ok {
			return expr
		}
		expr = parens.Expression
	}
}

// hasDeprecatedDescription reports whether a block's literal description starts with `DEPRECATED`.
func hasDeprecatedDescription(block *hclsyntax.Block) bool {
	desc, ok := staticString(block, "description")
	return ok && strings.HasPrefix(strings.TrimSpace(desc), "DEPRECATED")
}

// Package interfaces provides the standard AVM interfaces that are checked by rules.AVMInterfaceRule.
// It allows us to define the standard interfaces as a struct.
// When adding a new interface, be sure to add the corresponding NewAvmInterface* function in the rules package,
// and then add to the rule register.
package interfaces

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// Interfaces represent the standard AVM interfaces that are checked by rules.AVMInterfaceRule.
type AvmInterface struct {
	Default  cty.Value // Default value for the interface as a cty.Value
	Enabled  bool      // Whether to test this interface interface.
	Link     string    // Link to the interface documentation.
	Name     string    // Name of the interface.
	Nullable bool      // Whether the interface is nullable.
	Type     string    // String representing the type value.
	// TODO: add validation rule checks
}

// TypeExpression returns an HCL expression that represents the interface type.
// If the interface cannot be correctly parsed, this function will panic.
func (i AvmInterface) TypeExpression() hcl.Expression {
	e, d := hclsyntax.ParseExpression([]byte(i.Type), "variables.tf", hcl.Pos{})
	if d.HasErrors() {
		panic(d.Error())
	}
	return e
}

// TerrafromVar returns a string that represents the interface as the
// minimum required Terraform variable definition for testing.
func (i AvmInterface) TerrafromVar() string {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()
	varBlock := rootBody.AppendNewBlock("variable", []string{i.Name})
	varBody := varBlock.Body()
	// check the Type constraint is valid and panic if not
	if _, _, diags := typeexpr.TypeConstraintWithDefaults(i.TypeExpression()); diags.HasErrors() {
		panic(diags.Error())
	}
	// I couldn't get the hclwrite to work with the type constraint so I'm just adding it as a string
	// using SetSAttributeRaw and hclWrite.Token.
	varBody.SetAttributeRaw("type", hclwrite.Tokens{
		{
			Type:  hclsyntax.TokenStringLit,
			Bytes: []byte(" " + i.Type),
		},
	})
	varBody.SetAttributeValue("default", i.Default)
	// If the interface is not nullable, set the nullable attribute to false.
	// the default is true so we only need to set it if it's false.
	if !i.Nullable {
		varBody.SetAttributeValue("nullable", cty.False)
	}
	return string(f.Bytes())
}

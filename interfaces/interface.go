package interfaces

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/matt-FFFFFF/tfvarcheck/varcheck"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// AvmInterface represents the definition of an AVM interface.
type AvmInterface struct {
	Name     string          // Name of the interface, also the name of the variable to check.
	TypeStr  string          // The type of the interface as a string.
	Enabled  bool            // Whether the rule is enabled by default.
	Link     string          // Link to the interface specification.
	Default  cty.Value       // Default value for the interface.
	Nullable bool            // Whether the interface is nullable.
	Severity tflint.Severity // Summary of the interface.
}

// TerrafromVar returns a string that represents the interface as the
// minimum required Terraform variable definition for testing.
func (i AvmInterface) TerrafromVar() string {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()
	varBlock := rootBody.AppendNewBlock("variable", []string{i.Name})
	varBody := varBlock.Body()
	// check the Type constraint is valid and panic if not
	if _, diags := varcheck.NewTypeConstraintWithDefaultsFromBytes([]byte(i.TypeStr)); diags.HasErrors() {
		panic(diags.Error())
	}
	// I couldn't get the hclwrite to work with the type constraint so I'm just adding it as a string
	// using SetSAttributeRaw and hclWrite.Token.
	varBody.SetAttributeRaw("type", hclwrite.Tokens{
		{
			Type:  hclsyntax.TokenStringLit,
			Bytes: []byte(" " + i.TypeStr),
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

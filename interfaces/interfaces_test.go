package interfaces_test

import (
	"github.com/Azure/tflint-ruleset-avm/interfaces"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func toTerraformVarType(i interfaces.AvmInterface) string {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()
	varBlock := rootBody.AppendNewBlock("variable", []string{i.RuleName})
	varBody := varBlock.Body()

	// I couldn't get the hclwrite to work with the type constraint so I'm just adding it as a string
	// using SetSAttributeRaw and hclWrite.Token.
	varBody.SetAttributeRaw("type", hclwrite.Tokens{
		&hclwrite.Token{
			Type:         hclsyntax.TokenStringLit,
			Bytes:        []byte(i.VarTypeString),
			SpacesBefore: 1,
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

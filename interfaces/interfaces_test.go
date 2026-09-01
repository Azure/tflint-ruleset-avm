package interfaces_test

import (
	"testing"

	"github.com/Azure/tflint-ruleset-avm/interfaces"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

func toTerraformVarType(i interfaces.AvmInterface) string {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()
	variableName := i.VariableName
	if variableName == "" {
		variableName = i.RuleName
	}
	varBlock := rootBody.AppendNewBlock("variable", []string{variableName})
	varBody := varBlock.Body()

	varBody.SetAttributeRaw("type", hclwrite.Tokens{
		&hclwrite.Token{
			Type:         hclsyntax.TokenStringLit,
			Bytes:        []byte(i.VarTypeString),
			SpacesBefore: 1,
		},
	})
	if i.Default.IsKnown() {
		varBody.SetAttributeValue("default", i.Default)
	}
	if !i.Nullable {
		varBody.SetAttributeValue("nullable", cty.False)
	}
	return string(f.Bytes())
}

func registeredInterfaceRule(t *testing.T, name string) tflint.Rule {
	t.Helper()

	for _, rule := range interfaces.Rules {
		if rule.Name() == name {
			return rule
		}
	}

	t.Fatalf("interface rule %q is not registered", name)
	return nil
}

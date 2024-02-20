package avmhelper_test

import (
	"testing"

	"github.com/Azure/tflint-ruleset-avm/avmhelper"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func hclExpressionFromString(expr string) hcl.Expression {
	e, _ := hclsyntax.ParseExpression([]byte(expr), "test.tf", hcl.Pos{})
	return e
}

func TestCheckEqualTypeConstraints(t *testing.T) {
	cases := []struct {
		Name   string
		Want   hcl.Expression
		Got    hcl.Expression
		Result bool
	}{
		{
			Name:   "Same",
			Want:   hclExpressionFromString("object({kind = string, name = optional(string, null)})"),
			Got:    hclExpressionFromString("object({kind = string, name = optional(string, null)})"),
			Result: true,
		},
		{
			Name:   "Different",
			Want:   hclExpressionFromString("object({kind = string, name = optional(string, null)})"),
			Got:    hclExpressionFromString("object({kind = string, name = optional(number, null)})"),
			Result: false,
		},
		{
			Name:   "IncorrectDefaults",
			Want:   hclExpressionFromString("object({kind = string, name = optional(number, null)})"),
			Got:    hclExpressionFromString("object({kind = string, name = optional(number, 2)})"),
			Result: false,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			gotType, _ := avmhelper.NewVariableTypeFromExpression(tc.Got)
			wantType, _ := avmhelper.NewVariableTypeFromExpression(tc.Want)
			res := avmhelper.CheckEqualTypeConstraints(gotType, wantType)
			if res != tc.Result {
				t.Errorf("Test %s: Expected %v, got %v", tc.Name, tc.Result, res)
			}
		})
	}
}

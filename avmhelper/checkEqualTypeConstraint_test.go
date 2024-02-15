package avmhelper_test

import (
	"testing"

	"github.com/Azure/tflint-ruleset-avm/avmhelper"
	"github.com/Azure/tflint-ruleset-avm/to"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func hclExpressionFromString(expr string) hcl.Expression {
	e, _ := hclsyntax.ParseExpression([]byte(expr), "test.tf", hcl.Pos{})
	return e
}

func TestCheckEqualTypeConstraints(t *testing.T) {
	cases := []struct {
		Name     string
		Want     hcl.Expression
		Got      hcl.Expression
		Result   *bool
		ErrorMsg *string
	}{
		{
			Name:     "Same",
			Want:     hclExpressionFromString("object({kind = string, name = optional(string, null)})"),
			Got:      hclExpressionFromString("object({kind = string, name = optional(string, null)})"),
			Result:   to.Ptr(true),
			ErrorMsg: nil,
		},
		{
			Name:     "Different",
			Want:     hclExpressionFromString("object({kind = string, name = optional(string, null)})"),
			Got:      hclExpressionFromString("object({kind = string, name = optional(number, null)})"),
			Result:   to.Ptr(false),
			ErrorMsg: nil,
		},
		{
			Name:     "WantIsNotTypeConstraint",
			Want:     hclExpressionFromString("anotherfunc({kind = string, name = optional(string, null)})"),
			Got:      hclExpressionFromString("object({kind = string, name = optional(number, null)})"),
			Result:   nil,
			ErrorMsg: to.Ptr(`test.tf:0,0-59: Invalid type specification; Keyword "anotherfunc" is not a valid type constructor.`),
		},
		{
			Name:     "IncorrectDefaults",
			Want:     hclExpressionFromString("object({kind = string, name = optional(number, null)})"),
			Got:      hclExpressionFromString("object({kind = string, name = optional(number, 2)})"),
			Result:   to.Ptr(false),
			ErrorMsg: nil,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			res, diags := avmhelper.CheckEqualTypeConstraints(tc.Got, tc.Want)
			if diags.HasErrors() && tc.ErrorMsg == nil {
				t.Errorf("Test %s: Unexpected error: %s", tc.Name, diags.Error())
			}
			if res != nil && tc.Result != nil && *res != *tc.Result {
				t.Errorf("Test %s: Expected %v, got %v", tc.Name, *tc.Result, *res)
			}
			if tc.ErrorMsg != nil && diags.Error() != *tc.ErrorMsg {
				t.Errorf("Test %s: Expected error message %s, got %s", tc.Name, *tc.ErrorMsg, diags.Error())
			}
		})
	}
}

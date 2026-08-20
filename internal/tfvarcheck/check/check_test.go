package check_test

import (
	"testing"

	"github.com/Azure/tflint-ruleset-avm/internal/tfvarcheck/check"
	"github.com/Azure/tflint-ruleset-avm/internal/tfvarcheck/varcheck"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestEqualCtyValue(t *testing.T) {
	tests := []struct {
		name string
		got  cty.Value
		want cty.Value
		equal bool
	}{
		{name: "equal strings", got: cty.StringVal("hello"), want: cty.StringVal("hello"), equal: true},
		{name: "different strings", got: cty.StringVal("hello"), want: cty.StringVal("world")},
		{
			name:  "equal lists",
			got:   cty.ListVal([]cty.Value{cty.StringVal("hello"), cty.StringVal("world")}),
			want:  cty.ListVal([]cty.Value{cty.StringVal("hello"), cty.StringVal("world")}),
			equal: true,
		},
		{
			name: "different lists",
			got:  cty.ListVal([]cty.Value{cty.StringVal("hello"), cty.StringVal("world")}),
			want: cty.ListVal([]cty.Value{cty.StringVal("world"), cty.StringVal("hello")}),
		},
		{
			name:  "equal maps",
			got:   cty.MapVal(map[string]cty.Value{"key1": cty.StringVal("value1"), "key2": cty.StringVal("value2")}),
			want:  cty.MapVal(map[string]cty.Value{"key1": cty.StringVal("value1"), "key2": cty.StringVal("value2")}),
			equal: true,
		},
		{
			name: "different maps",
			got:  cty.MapVal(map[string]cty.Value{"key1": cty.StringVal("value1"), "key2": cty.StringVal("value2")}),
			want: cty.MapVal(map[string]cty.Value{"key1": cty.StringVal("value2"), "key2": cty.StringVal("value1")}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.equal, check.EqualCtyValue(test.got, test.want))
		})
	}
}

func TestEqualTypeConstraints(t *testing.T) {
	tests := []struct {
		name  string
		want  string
		got   string
		equal bool
	}{
		{
			name:  "same",
			want:  "object({kind = string, name = optional(string, null)})",
			got:   "object({kind = string, name = optional(string, null)})",
			equal: true,
		},
		{
			name: "different type",
			want: "object({kind = string, name = optional(string, null)})",
			got:  "object({kind = string, name = optional(number, null)})",
		},
		{
			name: "different default",
			want: "object({kind = string, name = optional(number, null)})",
			got:  "object({kind = string, name = optional(number, 2)})",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, gotDiags := varcheck.NewTypeConstraintWithDefaultsFromExp(parseExpression(test.got))
			want, wantDiags := varcheck.NewTypeConstraintWithDefaultsFromExp(parseExpression(test.want))
			assert.False(t, gotDiags.HasErrors())
			assert.False(t, wantDiags.HasErrors())
			assert.Equal(t, test.equal, check.EqualTypeConstraints(got, want))
		})
	}
}

func TestNullable(t *testing.T) {
	tests := []struct {
		name     string
		got      cty.Value
		want     bool
		expected bool
	}{
		{name: "nullable omitted", got: cty.NullVal(cty.Bool), want: true, expected: true},
		{name: "nullable explicitly true", got: cty.True, want: true},
		{name: "non-nullable omitted", got: cty.NullVal(cty.Bool), want: false},
		{name: "non-nullable false", got: cty.False, want: false, expected: true},
		{name: "wrong primitive type", got: cty.StringVal("test"), want: true},
		{name: "non-primitive type", got: cty.ListValEmpty(cty.String), want: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expected, check.Nullable(test.got, test.want))
		})
	}
}

func parseExpression(source string) hcl.Expression {
	expression, diags := hclsyntax.ParseExpression([]byte(source), "test.tf", hcl.Pos{})
	if diags.HasErrors() {
		panic(diags)
	}
	return expression
}

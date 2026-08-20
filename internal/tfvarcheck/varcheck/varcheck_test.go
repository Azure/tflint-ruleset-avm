package varcheck_test

import (
	"testing"

	"github.com/Azure/tflint-ruleset-avm/internal/tfvarcheck/varcheck"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestNewTypeConstraintWithDefaultsFromBytes(t *testing.T) {
	tests := []struct {
		name      string
		source    string
		wantType  cty.Type
		hasErrors bool
	}{
		{
			name: "simple",
			source: `object({
				foo = object({
					bar = optional(string, "baz")
				})
			})`,
			wantType: cty.Object(map[string]cty.Type{
				"foo": cty.ObjectWithOptionalAttrs(
					map[string]cty.Type{"bar": cty.String},
					[]string{"bar"},
				),
			}),
		},
		{
			name: "invalid",
			source: `obect({
				foo = object({
					bar = optional(string, "baz")
				})
			})`,
			wantType:  cty.Type{},
			hasErrors: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, diags := varcheck.NewTypeConstraintWithDefaultsFromBytes([]byte(test.source))
			assert.Equal(t, test.wantType, got.Type)
			assert.Equal(t, test.hasErrors, diags.HasErrors())
		})
	}
}

func TestNewTypeConstraintWithDefaultsFromBytes_empty(t *testing.T) {
	_, diags := varcheck.NewTypeConstraintWithDefaultsFromBytes(nil)
	assert.True(t, diags.HasErrors())
}

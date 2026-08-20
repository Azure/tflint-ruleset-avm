package varcheck

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

// VarCheck describes the expected type, default, and nullability of a Terraform variable.
type VarCheck struct {
	Default                cty.Value
	Nullable               bool
	TypeConstraintWithDefs TypeConstraintWithDefaults
}

// TypeConstraintWithDefaults pairs a Terraform type constraint with optional attribute defaults.
type TypeConstraintWithDefaults struct {
	Type    cty.Type
	Default *typeexpr.Defaults
}

// NewVarCheck creates a variable check.
func NewVarCheck(ty TypeConstraintWithDefaults, def cty.Value, nullable bool) VarCheck {
	return VarCheck{
		Default:                def,
		Nullable:               nullable,
		TypeConstraintWithDefs: ty,
	}
}

// NewTypeConstraintWithDefaultsFromExp parses a Terraform type constraint expression.
func NewTypeConstraintWithDefaultsFromExp(exp hcl.Expression) (TypeConstraintWithDefaults, hcl.Diagnostics) {
	ty, defaults, diags := typeexpr.TypeConstraintWithDefaults(exp)
	if diags.HasErrors() {
		return TypeConstraintWithDefaults{}, diags
	}
	return TypeConstraintWithDefaults{
		Type:    ty,
		Default: defaults,
	}, nil
}

// NewTypeConstraintWithDefaultsFromBytes parses a Terraform type constraint from source bytes.
func NewTypeConstraintWithDefaultsFromBytes(src []byte) (TypeConstraintWithDefaults, hcl.Diagnostics) {
	expression, diags := hclsyntax.ParseExpression(src, "variables.tf", hcl.Pos{})
	if diags.HasErrors() {
		return TypeConstraintWithDefaults{}, diags
	}
	return NewTypeConstraintWithDefaultsFromExp(expression)
}

package interfaces

import (
	"github.com/matt-FFFFFF/tfvarcheck/varcheck"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AvmInterface represents the definition of an AVM interface.
type AvmInterface struct {
	varcheck.VarCheck
	Name          string // Name of the interface, also the name of the variable to check.
	VarTypeString string
	Enabled       bool            // Whether the rule is enabled by default.
	Link          string          // Link to the interface specification.
	Severity      tflint.Severity // Summary of the interface.
}

func stringToTypeConstraintWithDefaults(c string) varcheck.TypeConstraintWithDefaults {
	v, d := varcheck.NewTypeConstraintWithDefaultsFromBytes([]byte(c))
	if d.HasErrors() {
		panic(d.Error())
	}
	return v
}

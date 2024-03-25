package interfaces

import (
	"github.com/matt-FFFFFF/tfvarcheck/varcheck"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AvmInterface represents the definition of an AVM interface,
// with additional information for use in TFLint.
type AvmInterface struct {
	varcheck.VarCheck
	RuleName      string          // RuleName of the interface, also the name of the variable to check.
	VarTypeString string          // The variable type value as a sting.
	RuleEnabled   bool            // Whether the rule is enabled by default.
	RuleLink      string          // RuleLink to the interface specification.
	RuleSeverity  tflint.Severity // Severity of the interface.
}

// StringToTypeConstraintWithDefaults converts a string to a TypeConstraintWithDefaults.
// The function will panic if the string is not valid.
func StringToTypeConstraintWithDefaults(c string) varcheck.TypeConstraintWithDefaults {
	v, d := varcheck.NewTypeConstraintWithDefaultsFromBytes([]byte(c))
	if d.HasErrors() {
		panic(d.Error())
	}
	return v
}

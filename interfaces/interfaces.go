package interfaces

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

// Interfaces represent the standard AVM interfaces that are checked by rules.AVMInterfaceRule.
type AVMInterface struct {
	Name    string    // Name of the interface.
	Link    string    // Link to the interface documentation.
	Type    string    // String representing the type value.
	Enabled bool      // Whether to test this interface interface.
	Default cty.Value // Default value for the interface in cty.
	// TODO: add validation rule checks
}

var AllInterfaces = []AVMInterface{
	Lock,
}

// TypeExpression returns an HCL expression that represents the interface type.
// If the interface cannot be correctly parsed, this function will panic.
func (i AVMInterface) TypeExpression() hcl.Expression {
	e, d := hclsyntax.ParseExpression([]byte(i.Type), "variables.tf", hcl.Pos{})
	if d.HasErrors() {
		panic(d.Error())
	}
	return e
}

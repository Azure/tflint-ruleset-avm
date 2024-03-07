package attrvalue

import (
	"reflect"

	"github.com/zclconf/go-cty/cty"
)

// NewSimpleStringRule is syntactic sugar, returning a new SimpleRule with the type values set for string
// and the given resource type, attribute name, and expected values.
func NewSimpleStringRule(resourceType string, attributeName string, expectedValues []any) *SimpleRule {
	return &SimpleRule{
		resourceType:   resourceType,
		attributeName:  attributeName,
		expectedValues: expectedValues,
		ctyType:        cty.String,
		reflectType:    reflect.TypeOf(""),
	}
}

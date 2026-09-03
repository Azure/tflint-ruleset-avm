package interfaces

import (
	"github.com/Azure/tflint-ruleset-avm/internal/tfvarcheck/varcheck"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

const (
	resourceTypesRuleLink      = "https://azure.github.io/Azure-Verified-Modules/spec/TFFR6/"
	retryTimeoutsRuleLink      = "https://azure.github.io/Azure-Verified-Modules/spec/TFFR7/"
	ignoreBodyChangesRuleLink  = "https://azure.github.io/Azure-Verified-Modules/spec/TFFR8/"
	privateEndpointsSchemaLink = "https://azure.github.io/Azure-Verified-Modules/includes/interfaces/tf/int.pe.schema.tf"
	resourceTagsRuleLink       = "https://azure.github.io/Azure-Verified-Modules/specs/tf/interfaces/#tags"
)

const retryTypeString = `object({
  error_message_regex  = optional(list(string))
  interval_seconds     = optional(number)
  max_interval_seconds = optional(number)
})`

const timeoutsTypeString = `object({
  create = optional(string)
  read   = optional(string)
  update = optional(string)
  delete = optional(string)
})`

var requiredInterfaceRules = []tflint.Rule{
	newRequiredVariableRule(
		"avm_interface_resource_types",
		"resource_types",
		resourceTypesRuleLink,
		"when the module directly declares an AzAPI resource",
		appliesToDirectAzapiResource,
		validateVariableShape(variableShape{
			typeDescription:     "must recursively contain optional string resource leaves and optional object submodule slots with canonical defaults",
			defaultDescription:  "must be `{}`",
			nullableDescription: "must be `false`",
			validateType:        validateResourceTypesType,
			validateDefault:     emptyObjectDefault,
			nullable:            nullableMustBeFalse,
		}),
	),
	newRequiredVariableRule(
		"avm_interface_retry",
		"retry",
		retryTimeoutsRuleLink,
		"when the module directly declares an AzAPI resource",
		appliesToDirectAzapiResource,
		validateVariableShape(variableShape{
			typeDescription:     "must use the TFFR7 retry object shape with optional fields",
			defaultDescription:  "must be `null`, `{}`, or an object containing legal retry field defaults",
			nullableDescription: "must permit `null`",
			validateType: func(got varcheck.TypeConstraintWithDefaults) bool {
				return validateFlatOptionalObject(got, map[string]cty.Type{
					"error_message_regex":  cty.List(cty.String),
					"interval_seconds":     cty.Number,
					"max_interval_seconds": cty.Number,
				})
			},
			validateDefault: nullableObjectDefault,
			nullable:        nullableMayBeTrue,
		}),
	),
	newRequiredVariableRule(
		"avm_interface_timeouts",
		"timeouts",
		retryTimeoutsRuleLink,
		"when the module directly declares an AzAPI resource",
		appliesToDirectAzapiResource,
		validateVariableShape(variableShape{
			typeDescription:     "must use the TFFR7 timeouts object shape with optional fields",
			defaultDescription:  "must be `null`, `{}`, or an object containing legal timeout field defaults",
			nullableDescription: "must permit `null`",
			validateType: func(got varcheck.TypeConstraintWithDefaults) bool {
				return validateFlatOptionalObject(got, map[string]cty.Type{
					"create": cty.String,
					"read":   cty.String,
					"update": cty.String,
					"delete": cty.String,
				})
			},
			validateDefault: nullableObjectDefault,
			nullable:        nullableMayBeTrue,
		}),
	),
	newRequiredVariableRule(
		"avm_interface_ignore_body_changes",
		"ignore_body_changes",
		ignoreBodyChangesRuleLink,
		"when the module directly declares an AzAPI resource",
		appliesToDirectAzapiResource,
		validateVariableShape(variableShape{
			typeDescription:     "must recursively contain optional list(string) resource leaves with valid defaults and optional object submodule slots defaulting to `{}`",
			defaultDescription:  "must be `{}`",
			nullableDescription: "must be `false`",
			validateType:        validateIgnoreBodyChangesType,
			validateDefault:     emptyObjectDefault,
			nullable:            nullableMustBeFalse,
		}),
	),
	newRequiredVariableRule(
		"avm_interface_resource_tags",
		"resource_tags",
		resourceTagsRuleLink,
		"when per-resource or submodule tag overrides are exposed",
		appliesWhenVariableDeclared("resource_tags"),
		validateVariableShape(variableShape{
			typeDescription:     "must use optional `resources` and `modules` namespaces with optional `map(string)` resource leaves and recursive module objects",
			defaultDescription:  "must be `null`",
			nullableDescription: "must permit `null`",
			validateType:        validateResourceTagsType,
			validateDefault:     nullDefault,
			nullable:            nullableMayBeTrue,
		}),
	),
	newRequiredVariableRule(
		"avm_interface_private_endpoints_manage_dns_zone_group",
		"private_endpoints_manage_dns_zone_group",
		privateEndpointsSchemaLink,
		"when variable `private_endpoints` is declared",
		appliesWhenVariableDeclared("private_endpoints"),
		validateVariableShape(variableShape{
			typeDescription:     "must be `bool`",
			defaultDescription:  "must be `true`",
			nullableDescription: "must be `false`",
			validateType: func(got varcheck.TypeConstraintWithDefaults) bool {
				return got.Type.Equals(cty.Bool) && got.Default == nil
			},
			validateDefault: func(got cty.Value, _ varcheck.TypeConstraintWithDefaults) bool {
				return got.Type().Equals(cty.Bool) && !got.IsNull() && got.True()
			},
			nullable: nullableMustBeFalse,
		}),
	),
}

func validateResourceTypesType(got varcheck.TypeConstraintWithDefaults) bool {
	if !got.Type.IsObjectType() {
		return false
	}

	directLeaves, ok := validateRecursiveObject(got.Type, got.Default, 0, resourceTypeLeaf)
	return ok && directLeaves > 0
}

func validateIgnoreBodyChangesType(got varcheck.TypeConstraintWithDefaults) bool {
	if !got.Type.IsObjectType() {
		return false
	}

	directLeaves, ok := validateRecursiveObject(got.Type, got.Default, 0, ignoreBodyChangesLeaf)
	return ok && directLeaves > 0
}

func validateResourceTagsType(got varcheck.TypeConstraintWithDefaults) bool {
	resourceLeaves, ok := validateResourceTagsObject(got.Type, got.Default)
	return ok && resourceLeaves > 0
}

func validateResourceTagsObject(objectType cty.Type, defaults *typeexpr.Defaults) (int, bool) {
	if !objectType.IsObjectType() {
		return 0, false
	}

	attributes := objectType.AttributeTypes()
	if len(attributes) == 0 || len(attributes) > 2 {
		return 0, false
	}

	resourceLeaves := 0
	for name, attributeType := range attributes {
		if !objectType.AttributeOptional(name) {
			return 0, false
		}
		if _, hasDefault := defaultsValue(defaults, name); hasDefault {
			return 0, false
		}

		var (
			leaves int
			ok     bool
		)
		switch name {
		case "resources":
			leaves, ok = validateResourceTagsResources(attributeType, defaultsChild(defaults, name))
		case "modules":
			leaves, ok = validateResourceTagsModules(attributeType, defaultsChild(defaults, name))
		default:
			return 0, false
		}
		if !ok {
			return 0, false
		}
		resourceLeaves += leaves
	}

	return resourceLeaves, true
}

func validateResourceTagsResources(objectType cty.Type, defaults *typeexpr.Defaults) (int, bool) {
	if !objectType.IsObjectType() || len(objectType.AttributeTypes()) == 0 {
		return 0, false
	}

	for name, attributeType := range objectType.AttributeTypes() {
		if !objectType.AttributeOptional(name) || !attributeType.Equals(cty.Map(cty.String)) {
			return 0, false
		}
		if _, hasDefault := defaultsValue(defaults, name); hasDefault {
			return 0, false
		}
	}

	return len(objectType.AttributeTypes()), true
}

func validateResourceTagsModules(objectType cty.Type, defaults *typeexpr.Defaults) (int, bool) {
	if !objectType.IsObjectType() || len(objectType.AttributeTypes()) == 0 {
		return 0, false
	}

	resourceLeaves := 0
	for name, attributeType := range objectType.AttributeTypes() {
		if !objectType.AttributeOptional(name) || !attributeType.IsObjectType() {
			return 0, false
		}
		if _, hasDefault := defaultsValue(defaults, name); hasDefault {
			return 0, false
		}

		leaves, ok := validateResourceTagsObject(attributeType, defaultsChild(defaults, name))
		if !ok {
			return 0, false
		}
		resourceLeaves += leaves
	}

	return resourceLeaves, true
}

type recursiveLeafValidator func(cty.Type, *typeexpr.Defaults, string, int) bool

func validateRecursiveObject(
	objectType cty.Type,
	defaults *typeexpr.Defaults,
	depth int,
	validateLeaf recursiveLeafValidator,
) (int, bool) {
	if !objectType.IsObjectType() {
		return 0, false
	}

	directLeaves := 0
	for name, attributeType := range objectType.AttributeTypes() {
		if !objectType.AttributeOptional(name) {
			return 0, false
		}

		if attributeType.IsObjectType() {
			if !hasSemanticEmptyObjectDefault(defaults, name, attributeType) {
				return 0, false
			}

			childDefaults := defaultsChild(defaults, name)
			if _, ok := validateRecursiveObject(attributeType, childDefaults, depth+1, validateLeaf); !ok {
				return 0, false
			}
			continue
		}

		if !validateLeaf(attributeType, defaults, name, depth) {
			return 0, false
		}
		if depth == 0 {
			directLeaves++
		}
	}

	return directLeaves, true
}

func resourceTypeLeaf(attributeType cty.Type, defaults *typeexpr.Defaults, name string, depth int) bool {
	if !attributeType.Equals(cty.String) {
		return false
	}

	defaultValue, hasDefault := defaultsValue(defaults, name)
	if depth > 0 {
		return !hasDefault
	}
	if !hasDefault || !defaultValue.IsKnown() || defaultValue.IsNull() || !defaultValue.Type().Equals(cty.String) {
		return false
	}
	return defaultValue.AsString() != ""
}

func ignoreBodyChangesLeaf(attributeType cty.Type, defaults *typeexpr.Defaults, name string, _ int) bool {
	if !attributeType.Equals(cty.List(cty.String)) {
		return false
	}

	defaultValue, hasDefault := defaultsValue(defaults, name)
	if !hasDefault || !defaultValue.IsKnown() || defaultValue.IsNull() || !defaultValue.Type().Equals(cty.List(cty.String)) {
		return false
	}

	for iterator := defaultValue.ElementIterator(); iterator.Next(); {
		_, element := iterator.Element()
		if !element.IsKnown() || element.IsNull() || element.AsString() == "" {
			return false
		}
	}

	return true
}

func validateFlatOptionalObject(got varcheck.TypeConstraintWithDefaults, expected map[string]cty.Type) bool {
	if !got.Type.IsObjectType() || len(got.Type.AttributeTypes()) != len(expected) {
		return false
	}

	for name, expectedType := range expected {
		if !got.Type.HasAttribute(name) || !got.Type.AttributeOptional(name) {
			return false
		}
		if !got.Type.AttributeType(name).Equals(expectedType) {
			return false
		}
	}

	return true
}

func emptyObjectDefault(got cty.Value, _ varcheck.TypeConstraintWithDefaults) bool {
	return got.RawEquals(cty.EmptyObjectVal)
}

func nullDefault(got cty.Value, _ varcheck.TypeConstraintWithDefaults) bool {
	return got.IsNull()
}

func nullableObjectDefault(got cty.Value, typeConstraint varcheck.TypeConstraintWithDefaults) bool {
	if got.IsNull() {
		return true
	}
	if !got.Type().IsObjectType() {
		return false
	}

	for name, value := range got.AsValueMap() {
		if !typeConstraint.Type.HasAttribute(name) {
			return false
		}
		if value.IsNull() {
			continue
		}
		if _, err := convert.Convert(value, typeConstraint.Type.AttributeType(name)); err != nil {
			return false
		}
	}

	return true
}

func hasSemanticEmptyObjectDefault(defaults *typeexpr.Defaults, name string, objectType cty.Type) bool {
	actual, ok := defaultsValue(defaults, name)
	if !ok {
		return false
	}
	expected, err := convert.Convert(cty.EmptyObjectVal, objectType)
	return err == nil && actual.RawEquals(expected)
}

func defaultsValue(defaults *typeexpr.Defaults, name string) (cty.Value, bool) {
	if defaults == nil {
		return cty.NilVal, false
	}
	value, ok := defaults.DefaultValues[name]
	return value, ok
}

func defaultsChild(defaults *typeexpr.Defaults, name string) *typeexpr.Defaults {
	if defaults == nil {
		return nil
	}
	return defaults.Children[name]
}

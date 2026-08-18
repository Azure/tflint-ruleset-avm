package interfaces

import (
	"github.com/Azure/tflint-ruleset-avm/common"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

const deprecatedLockMessage = "lock uses deprecated interface variant 1; migrate to variant 2 by adding notes = optional(string, null). Support for variant 1 will be removed in the release following the v0.19.0 migration window."

const deprecatedRoleAssignmentsMessage = "role_assignments uses deprecated interface variant 1; migrate to variant 2 by adding name = optional(string, null). Support for variant 1 will be removed in the release following the v0.19.0 migration window."

const deprecatedPrivateEndpointsMessage = "private_endpoints uses a deprecated transitional schema; migrate nested lock, role_assignments, and ip_configurations to the canonical v0.19.0 schema. Transitional schemas will be removed in the release following the v0.19.0 migration window."

func newEitherInterfaceRule(name string, variants ...AvmInterface) tflint.Rule {
	rules := make([]tflint.Rule, len(variants))
	for i, variant := range variants {
		rules[i] = NewVarCheckRuleFromAvmInterface(variant)
	}
	return common.NewEitherCheckRule(name, true, tflint.ERROR, rules...)
}

type interfaceRuleRegistration struct {
	name     string
	variants []AvmInterface
}

var interfaceRuleRegistrations = []interfaceRuleRegistration{
	{name: "customer_managed_key", variants: []AvmInterface{CustomerManagedKeyV2, CustomerManagedKey}},
	{name: "location", variants: []AvmInterface{Location}},
	{name: "lock", variants: []AvmInterface{LockV2, LockV1}},
	{name: "managed_identities", variants: []AvmInterface{ManagedIdentities}},
	{name: "role_assignments", variants: []AvmInterface{RoleAssignmentsV2, RoleAssignmentsV1}},
	{name: "tags", variants: []AvmInterface{Tags}},
	{name: "private_endpoints", variants: privateEndpointVariants},
	{name: "diagnostic_settings", variants: []AvmInterface{DiagnosticSettings, DiagnosticSettingsV2}},
}

var Rules = func() []tflint.Rule {
	rules := make([]tflint.Rule, 0, len(interfaceRuleRegistrations)+3+len(requiredInterfaceRules))
	for _, registration := range interfaceRuleRegistrations {
		if len(registration.variants) == 1 {
			rules = append(rules, NewVarCheckRuleFromAvmInterface(registration.variants[0]))
			continue
		}
		rules = append(rules, newEitherInterfaceRule(registration.name, registration.variants...))
	}

	rules = append(rules,
		newDeprecatedInterfaceVariantRule(
			"deprecated_lock_interface",
			"lock",
			deprecatedLockMessage,
			LockV2.RuleLink,
			LockV1,
		),
		newDeprecatedInterfaceVariantRule(
			"deprecated_role_assignments_interface",
			"role_assignments",
			deprecatedRoleAssignmentsMessage,
			RoleAssignmentsV2.RuleLink,
			RoleAssignmentsV1,
		),
		newDeprecatedInterfaceVariantRule(
			"deprecated_private_endpoints_interface",
			"private_endpoints",
			deprecatedPrivateEndpointsMessage,
			PrivateEndpoints.RuleLink,
			privateEndpointVariants[1:]...,
		),
	)
	return append(rules, requiredInterfaceRules...)
}()

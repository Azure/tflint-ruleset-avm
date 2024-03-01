package interfaces_test

import (
	"fmt"
	"testing"

	"github.com/Azure/tflint-ruleset-avm/interfaces"
	"github.com/Azure/tflint-ruleset-avm/rules"
	"github.com/hashicorp/hcl/v2"
	"github.com/matt-FFFFFF/tfvarcheck/varcheck"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// complexVarTypeString is a string representation of a complex variable type used for testing.
// it is based on the private endpoints interface.
var complexVarTypeString = `map(object({
	name               = optional(string, null)
	role_assignments   = optional(map(object({
		role_definition_id_or_name             = string
		principal_id                           = string
		description                            = optional(string, null)
		skip_service_principal_aad_check       = optional(bool, false)
		condition                              = optional(string, null)
		condition_version                      = optional(string, null)
		delegated_managed_identity_resource_id = optional(string, null)
	})), {})
	lock = optional(object({
		kind = string
		name = optional(string, null)
	}), null)
	tags               = optional(map(string), null)
	subnet_resource_id = string
	private_dns_zone_group_name             = optional(string, "default")
	private_dns_zone_resource_ids           = optional(set(string), [])
	application_security_group_associations = optional(map(string), {})
	private_service_connection_name         = optional(string, null)
	network_interface_name                  = optional(string, null)
	location                                = optional(string, null)
	resource_group_name                     = optional(string, null)
	ip_configurations = optional(map(object({
		name               = string
		private_ip_address = string
	})), {})
}))`

var complexType = interfaces.StringToTypeConstraintWithDefaults(complexVarTypeString)

// ComplexVar represents a complex variable type used for testing.
var ComplexVar = interfaces.AvmInterface{
	VarCheck:      varcheck.NewVarCheck(complexType, cty.EmptyObjectVal, false),
	RuleName:      "complex",
	VarTypeString: complexVarTypeString,
	RuleEnabled:   true,
	RuleLink:      "https://complex",
	RuleSeverity:  tflint.ERROR,
}

func TestComplexInterface(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		JSON     bool
		Expected helper.Issues
	}{
		{
			Name: "not complex variable",
			Content: `
variable "not_complex" {
	default = "default"
}`,
			Expected: helper.Issues{},
		},
		{
			Name:     "correct",
			Content:  toTerraformVarType(ComplexVar),
			Expected: helper.Issues{},
		},
		{
			Name: "incorrect default on deep object",
			Content: `
variable "complex" {
	default  = {}
	nullable = false
	type     = map(object({
		name               = optional(string, null)
		role_assignments   = optional(map(object({
			role_definition_id_or_name             = string
			principal_id                           = string
			description                            = optional(string, null)
			skip_service_principal_aad_check       = optional(bool, false)
			condition                              = optional(string, null)
			condition_version                      = optional(string, null)
			delegated_managed_identity_resource_id = optional(string, null)
		})), {})
		lock = optional(object({
			kind = string
			name = optional(string, "foo")
		}), null)
		tags               = optional(map(string), null)
		subnet_resource_id = string
		private_dns_zone_group_name             = optional(string, "default")
		private_dns_zone_resource_ids           = optional(set(string), [])
		application_security_group_associations = optional(map(string), {})
		private_service_connection_name         = optional(string, null)
		network_interface_name                  = optional(string, null)
		location                                = optional(string, null)
		resource_group_name                     = optional(string, null)
		ip_configurations = optional(map(object({
			name               = string
			private_ip_address = string
		})), {})
	}))
}`,
			Expected: helper.Issues{
				&helper.Issue{
					Rule:    rules.NewVarCheckRuleFromAvmInterface(ComplexVar),
					Message: fmt.Sprintf("variable type does not comply with the interface specification:\n\n%s", complexVarTypeString),
					Range:   hcl.Range{Filename: "variables.tf", Start: hcl.Pos{Line: 5, Column: 2}, End: hcl.Pos{Line: 33, Column: 5}},
				},
			},
		},
	}

	rule := rules.NewVarCheckRuleFromAvmInterface(ComplexVar)

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			filename := "variables.tf"
			if tc.JSON {
				filename += ".json"
			}

			runner := helper.TestRunner(t, map[string]string{filename: tc.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}

package interfaces_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

// Canonical private endpoint schema snapshot:
// https://github.com/Azure/Azure-Verified-Modules/blob/005b4f94e9197d1e23ea4213c92b4263bf636a2e/docs/static/includes/interfaces/tf/int.pe.schema.tf
const validPrivateEndpointsManageDNSZoneGroup = `variable "private_endpoints" {
  type    = map(any)
  default = {}
}

variable "private_endpoints_manage_dns_zone_group" {
  type     = bool
  default  = true
  nullable = false
}`

func TestPrivateEndpointsManageDNSZoneGroupValid(t *testing.T) {
	rule := requiredInterfaceRule(t, "private_endpoints_manage_dns_zone_group")
	runner := helper.TestRunner(t, map[string]string{"variables.tf": validPrivateEndpointsManageDNSZoneGroup})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	helper.AssertIssues(t, helper.Issues{}, runner.Issues)
}

func TestPrivateEndpointsManageDNSZoneGroupRequiredConditionally(t *testing.T) {
	rule := requiredInterfaceRule(t, "private_endpoints_manage_dns_zone_group")
	content := `variable "private_endpoints" {
  type    = map(any)
  default = {}
}`
	runner := helper.TestRunner(t, map[string]string{"variables.tf": content})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	helper.AssertIssues(t, helper.Issues{
		{
			Rule:    rule,
			Message: "variable `private_endpoints_manage_dns_zone_group` must be declared when variable `private_endpoints` is declared; see: " + rule.Link(),
			Range: hcl.Range{
				Filename: "variables.tf",
				Start:    hcl.Pos{Line: 1, Column: 1},
				End:      hcl.Pos{Line: 1, Column: 29},
			},
		},
	}, runner.Issues)
}

func TestPrivateEndpointsManageDNSZoneGroupSkipsWhenPrivateEndpointsAbsent(t *testing.T) {
	rule := requiredInterfaceRule(t, "private_endpoints_manage_dns_zone_group")
	content := `variable "private_endpoints_manage_dns_zone_group" {
  type     = string
  default  = false
  nullable = true
}`
	runner := helper.TestRunner(t, map[string]string{"variables.tf": content})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	helper.AssertIssues(t, helper.Issues{}, runner.Issues)
}

func TestPrivateEndpointsManageDNSZoneGroupRejectsMalformedDeclaration(t *testing.T) {
	rule := requiredInterfaceRule(t, "private_endpoints_manage_dns_zone_group")
	cases := []struct {
		name        string
		content     string
		attribute   string
		expectation string
	}{
		{
			name:        "wrong type",
			content:     strings.Replace(validPrivateEndpointsManageDNSZoneGroup, "type     = bool", "type     = string", 1),
			attribute:   "type",
			expectation: "must be `bool`",
		},
		{
			name:        "false default",
			content:     strings.Replace(validPrivateEndpointsManageDNSZoneGroup, "default  = true", "default  = false", 1),
			attribute:   "default",
			expectation: "must be `true`",
		},
		{
			name:        "missing default",
			content:     strings.Replace(validPrivateEndpointsManageDNSZoneGroup, "  default  = true\n", "", 1),
			attribute:   "default",
			expectation: "must be `true`",
		},
		{
			name:        "nullable true",
			content:     strings.Replace(validPrivateEndpointsManageDNSZoneGroup, "nullable = false", "nullable = true", 1),
			attribute:   "nullable",
			expectation: "must be `false`",
		},
		{
			name:        "missing nullable",
			content:     strings.Replace(validPrivateEndpointsManageDNSZoneGroup, "  nullable = false\n", "", 1),
			attribute:   "nullable",
			expectation: "must be `false`",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"variables.tf": tc.content})
			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			helper.AssertIssuesWithoutRange(
				t,
				expectedVariableAttributeIssue(rule, tc.attribute, tc.expectation, hcl.Range{}),
				runner.Issues,
			)
		})
	}
}

func TestPrivateEndpointsManageDNSZoneGroupMalformedSourceRange(t *testing.T) {
	rule := requiredInterfaceRule(t, "private_endpoints_manage_dns_zone_group")
	content := strings.Replace(validPrivateEndpointsManageDNSZoneGroup, "type     = bool", "type     = string", 1)
	runner := helper.TestRunner(t, map[string]string{"variables.tf": content})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	helper.AssertIssues(t, expectedVariableAttributeIssue(
		rule,
		"type",
		"must be `bool`",
		hcl.Range{
			Filename: "variables.tf",
			Start:    hcl.Pos{Line: 7, Column: 3},
			End:      hcl.Pos{Line: 7, Column: 20},
		},
	), runner.Issues)
}

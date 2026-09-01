package interfaces_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Azure/tflint-ruleset-avm/interfaces"
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var azapiRequiredVariableNames = []string{
	"resource_types",
	"retry",
	"timeouts",
	"ignore_body_changes",
}

func TestAzapiRequiredInterfacesTriggerOnEveryResourceType(t *testing.T) {
	resourceTypes := []string{
		"azapi_resource",
		"azapi_data_plane_resource",
		"azapi_resource_action",
		"azapi_update_resource",
	}

	for _, resourceType := range resourceTypes {
		resourceType := resourceType
		t.Run(resourceType, func(t *testing.T) {
			content := fmt.Sprintf(`resource %q "example" {
  count    = 0
  for_each = {}
}`, resourceType)

			for _, variableName := range azapiRequiredVariableNames {
				rule := requiredInterfaceRule(t, variableName)
				runner := helper.TestRunner(t, map[string]string{"main.tf": content})

				if err := rule.Check(runner); err != nil {
					t.Fatalf("%s: unexpected error: %s", variableName, err)
				}
				helper.AssertIssuesWithoutRange(t, helper.Issues{
					{
						Rule: rule,
						Message: fmt.Sprintf(
							"variable `%s` must be declared when the module directly declares an AzAPI resource; see: %s",
							variableName,
							rule.Link(),
						),
					},
				}, runner.Issues)
			}
		})
	}
}

func TestAzapiRequiredInterfacesIgnoreNonTriggers(t *testing.T) {
	content := `provider "azapi" {}

data "azapi_resource" "example" {
  type      = "Microsoft.Example/widgets@2024-01-01"
  parent_id = "example"
  name      = "example"
}

resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "eastus"
}

module "child" {
  source = "./child"
}`

	for _, variableName := range azapiRequiredVariableNames {
		rule := requiredInterfaceRule(t, variableName)
		runner := helper.TestRunner(t, map[string]string{"main.tf": content})
		if err := rule.Check(runner); err != nil {
			t.Fatalf("%s: unexpected error: %s", variableName, err)
		}
		helper.AssertIssues(t, helper.Issues{}, runner.Issues)
	}
}

func TestAzapiRequiredInterfacesSkipValidationWhenNotApplicable(t *testing.T) {
	for _, variableName := range azapiRequiredVariableNames {
		variableName := variableName
		t.Run(variableName, func(t *testing.T) {
			content := fmt.Sprintf(`variable %q {
  type     = string
  default  = "invalid"
  nullable = true
}`, variableName)
			rule := requiredInterfaceRule(t, variableName)
			runner := helper.TestRunner(t, map[string]string{"variables.tf": content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			helper.AssertIssues(t, helper.Issues{}, runner.Issues)
		})
	}
}

func TestAzapiRequiredInterfacesEvaluateParentAndChildIndependently(t *testing.T) {
	rule := requiredInterfaceRule(t, "resource_types")
	parentRunner := helper.TestRunner(t, map[string]string{
		"main.tf": `module "child" {
  source = "./child"
}`,
	})
	if err := rule.Check(parentRunner); err != nil {
		t.Fatalf("parent: unexpected error: %s", err)
	}
	helper.AssertIssues(t, helper.Issues{}, parentRunner.Issues)

	childRunner := helper.TestRunner(t, map[string]string{
		"main.tf": `resource "azapi_resource" "child" {
  type = "Microsoft.Example/children@2024-01-01"
}`,
	})
	if err := rule.Check(childRunner); err != nil {
		t.Fatalf("child: unexpected error: %s", err)
	}
	helper.AssertIssuesWithoutRange(t, helper.Issues{
		{
			Rule:    rule,
			Message: "variable `resource_types` must be declared when the module directly declares an AzAPI resource; see: " + rule.Link(),
		},
	}, childRunner.Issues)
}

func TestAzapiRequiredInterfaceMissingVariableSourceRange(t *testing.T) {
	rule := requiredInterfaceRule(t, "resource_types")
	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `resource "azapi_resource" "example" {
  type = "Microsoft.Example/widgets@2024-01-01"
}`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	helper.AssertIssues(t, helper.Issues{
		{
			Rule:    rule,
			Message: "variable `resource_types` must be declared when the module directly declares an AzAPI resource; see: " + rule.Link(),
			Range: hcl.Range{
				Filename: "main.tf",
				Start:    hcl.Pos{Line: 1, Column: 1},
				End:      hcl.Pos{Line: 1, Column: 36},
			},
		},
	}, runner.Issues)
}

func requiredInterfaceRule(t *testing.T, name string) tflint.Rule {
	t.Helper()
	ruleName := "avm_interface_" + name
	for _, rule := range interfaces.Rules {
		if rule.Name() == ruleName {
			return rule
		}
	}
	t.Fatalf("required interface rule %q is not registered", name)
	return nil
}

func expectedVariableAttributeIssue(rule tflint.Rule, attribute, expectation string, sourceRange hcl.Range) helper.Issues {
	return helper.Issues{
		{
			Rule: rule,
			Message: fmt.Sprintf(
				"variable `%s` %s %s; see: %s",
				strings.TrimPrefix(rule.Name(), "avm_interface_"),
				attribute,
				expectation,
				rule.Link(),
			),
			Range: sourceRange,
		},
	}
}

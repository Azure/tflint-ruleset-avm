package interfaces_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Azure/tflint-ruleset-avm/common"
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
		"variables.tf": "",
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	helper.AssertIssues(t, helper.Issues{
		{
			Rule:    rule,
			Message: "variable `resource_types` must be declared when the module directly declares an AzAPI resource; see: " + rule.Link(),
			Range: hcl.Range{
				Filename: "variables.tf",
				Start:    hcl.Pos{Line: 1, Column: 1},
				End:      hcl.Pos{Line: 1, Column: 1},
			},
		},
	}, runner.Issues)
}

func TestAzapiRequiredInterfacesCheckEveryResourceWiring(t *testing.T) {
	tests := []struct {
		name     string
		variable string
		content  string
	}{
		{
			name:     "resource types",
			variable: "resource_types",
			content: validResourceTypes + `
resource "azapi_resource" "unwired" {
  type = "Microsoft.Example/widgets@2024-01-01"
}`,
		},
		{
			name:     "retry",
			variable: "retry",
			content: validRetry + `
resource "azapi_resource" "unwired" {
  type = "Microsoft.Example/widgets@2024-01-01"
}`,
		},
		{
			name:     "timeouts",
			variable: "timeouts",
			content: validTimeouts + `
resource "azapi_resource" "unwired" {
  type = "Microsoft.Example/widgets@2024-01-01"
}`,
		},
		{
			name:     "ignore body changes",
			variable: "ignore_body_changes",
			content: validIgnoreBodyChanges + `
resource "azapi_resource" "unwired" {
  type = "Microsoft.Example/widgets@2024-01-01"
}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rule := requiredInterfaceRule(t, test.variable)
			runner := helper.TestRunner(t, map[string]string{"main.tf": test.content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if len(runner.Issues) != 1 {
				t.Fatalf("expected one resource-local wiring issue, got %#v", runner.Issues)
			}
			if !strings.Contains(runner.Issues[0].Message, "using variable `"+test.variable+"`") {
				t.Fatalf("unexpected issue message %q", runner.Issues[0].Message)
			}
		})
	}
}

func TestAzapiRequiredInterfacesCheckMultipleUnwiredResources(t *testing.T) {
	rule := requiredInterfaceRule(t, "resource_types")
	runner := helper.TestRunner(t, map[string]string{"main.tf": validResourceTypes + `
resource "azapi_resource" "first_unwired" {
  type = "Microsoft.Example/widgets@2024-01-01"
}

resource "azapi_resource" "second_unwired" {
  type = "Microsoft.Example/widgets@2024-01-01"
}`})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(runner.Issues) != 2 {
		t.Fatalf("expected two resource-local wiring issues, got %#v", runner.Issues)
	}
}

func TestAzapiRequiredInterfacesAcceptLocalsAndConditionalForwarding(t *testing.T) {
	rule := requiredInterfaceRule(t, "retry")
	content := strings.Replace(validRetry, "retry = var.retry", "retry = local.retry", 1) + `
locals {
  retry = var.retry == null ? null : var.retry
}`
	runner := helper.TestRunner(t, map[string]string{"main.tf": content})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	helper.AssertIssues(t, helper.Issues{}, runner.Issues)
}

func TestAzapiRequiredInterfacesApplyIgnoredResourcesBeforeApplicability(t *testing.T) {
	rule := common.NewConfigurableRule(requiredInterfaceRule(t, "resource_types"))
	config := `rule "avm_interface_resource_types" {
  enabled           = true
  ignored_resources = ["azapi_resource.telemetry"]
}`
	telemetry := `resource "azapi_resource" "telemetry" {
  type = "Microsoft.Example/telemetry@2024-01-01"
}`

	for _, content := range []string{
		telemetry,
		telemetry + `
resource "azapi_resource" "business" {
  type = "Microsoft.Example/business@2024-01-01"
}`,
		`resource "azapi_resource" "business" {
  type = "Microsoft.Example/business@2024-01-01"
}` + "\n" + telemetry,
	} {
		runner := helper.TestRunner(t, map[string]string{
			".tflint.hcl":  config,
			"main.tf":      content,
			"variables.tf": "",
		})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		wantIssues := 0
		if strings.Contains(content, `"business"`) {
			wantIssues = 1
		}
		if len(runner.Issues) != wantIssues {
			t.Fatalf("expected %d issues for %q, got %#v", wantIssues, content, runner.Issues)
		}
	}
}

func TestAzapiRequiredInterfacesKeepExclusionsIsolatedPerRule(t *testing.T) {
	resourceTypesRule := common.NewConfigurableRule(requiredInterfaceRule(t, "resource_types"))
	retryRule := common.NewConfigurableRule(requiredInterfaceRule(t, "retry"))
	runner := helper.TestRunner(t, map[string]string{
		".tflint.hcl": `rule "avm_interface_resource_types" {
  enabled           = true
  ignored_resources = ["azapi_resource.telemetry"]
}

rule "avm_interface_retry" {
  enabled = true
}`,
		"main.tf": `resource "azapi_resource" "telemetry" {
  type = "Microsoft.Example/telemetry@2024-01-01"
}`,
		"variables.tf": "",
	})

	if err := resourceTypesRule.Check(runner); err != nil {
		t.Fatalf("resource types: unexpected error: %s", err)
	}
	if err := retryRule.Check(runner); err != nil {
		t.Fatalf("retry: unexpected error: %s", err)
	}
	if len(runner.Issues) != 1 || runner.Issues[0].Rule.Name() != "avm_interface_retry" {
		t.Fatalf("expected only the retry diagnostic, got %#v", runner.Issues)
	}
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

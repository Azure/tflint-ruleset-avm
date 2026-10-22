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
				if variableName == "ignore_body_changes" && resourceType != "azapi_resource" {
					continue
				}
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

func TestIgnoreBodyChangesOnlyAppliesToAzapiResource(t *testing.T) {
	rule := requiredInterfaceRule(t, "ignore_body_changes")
	for _, resourceType := range []string{"azapi_data_plane_resource", "azapi_resource_action", "azapi_update_resource"} {
		t.Run(resourceType, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{
				"main.tf": fmt.Sprintf(`resource %q "example" {
  type = "Microsoft.Example/widgets@2024-01-01"
}`, resourceType),
			})
			if err := rule.Check(runner); err != nil {
				t.Fatal(err)
			}
			helper.AssertIssues(t, helper.Issues{}, runner.Issues)
		})
	}
}

func TestTimeoutsBlockAppliesToEachSupportedAzapiResource(t *testing.T) {
	rule := requiredInterfaceRule(t, "timeouts")
	for _, resourceType := range []string{"azapi_resource", "azapi_data_plane_resource", "azapi_resource_action", "azapi_update_resource"} {
		t.Run(resourceType, func(t *testing.T) {
			content := strings.Replace(validTimeouts, `resource "azapi_resource"`, fmt.Sprintf(`resource %q`, resourceType), 1)
			runner := helper.TestRunner(t, map[string]string{"main.tf": content})
			if err := rule.Check(runner); err != nil {
				t.Fatal(err)
			}
			helper.AssertIssues(t, helper.Issues{}, runner.Issues)
		})
	}
}

func TestTimeoutsDynamicBlockForwarding(t *testing.T) {
	rule := requiredInterfaceRule(t, "timeouts")
	dynamic := `dynamic "timeouts" {
    for_each = var.timeouts == null ? [] : [var.timeouts]
    content {
      create = timeouts.value.create
      read   = timeouts.value.read
      update = timeouts.value.update
      delete = timeouts.value.delete
    }
  }`
	tests := []struct {
		name       string
		block      string
		wantIssues int
	}{
		{"nullable input", dynamic, 0},
		{"custom iterator", strings.NewReplacer("    content {", "    iterator = settings\n    content {", "timeouts.value", "settings.value").Replace(dynamic), 0},
		{"condition only collection", strings.Replace(dynamic, "[var.timeouts]", "[{ create = null, read = null, update = null, delete = null }]", 1), 1},
		{"unwired content", strings.NewReplacer("timeouts.value.create", "null", "timeouts.value.read", "null", "timeouts.value.update", "null", "timeouts.value.delete", "null").Replace(dynamic), 1},
		{"condition only content", strings.NewReplacer("timeouts.value.create", "timeouts.value.create != null ? null : null", "timeouts.value.read", "null", "timeouts.value.update", "null", "timeouts.value.delete", "null").Replace(dynamic), 1},
	}
	start := strings.Index(validTimeouts, "timeouts {")
	end := strings.Index(validTimeouts[start:], "\n  }") + start + len("\n  }")
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			content := validTimeouts[:start] + test.block + validTimeouts[end:]
			runner := helper.TestRunner(t, map[string]string{"main.tf": content})
			if err := rule.Check(runner); err != nil {
				t.Fatal(err)
			}
			if len(runner.Issues) != test.wantIssues {
				t.Fatalf("expected %d wiring issues, got %#v", test.wantIssues, runner.Issues)
			}
		})
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

func TestAzapiRequiredInterfacesAcceptSelectedLocalAliasMembers(t *testing.T) {
	tests := []struct {
		name     string
		variable string
		content  string
		old      string
		repl     string
	}{
		{
			name:     "resource types",
			variable: "resource_types",
			content:  validResourceTypes,
			old:      "type = var.resource_types.widget",
			repl:     `type = coalesce(local.overrides.widget, "Microsoft.Example/widgets@2024-01-01")`,
		},
		{
			name:     "ignore body changes",
			variable: "ignore_body_changes",
			content:  validIgnoreBodyChanges,
			old:      "ignore_body_changes = var.ignore_body_changes.widget",
			repl:     "ignore_body_changes = coalesce(local.overrides.widget, [])",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			content := strings.Replace(test.content, test.old, test.repl, 1) + `
locals {
  overrides = var.` + test.variable + `
}`
			rule := requiredInterfaceRule(t, test.variable)
			runner := helper.TestRunner(t, map[string]string{"main.tf": content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			helper.AssertIssues(t, helper.Issues{}, runner.Issues)
		})
	}
}

func TestAzapiRequiredInterfacesAcceptValueComposition(t *testing.T) {
	tests := []struct {
		name     string
		variable string
		content  string
		old      string
		repl     string
		locals   string
	}{
		{
			name:     "retry merge",
			variable: "retry",
			content:  validRetry,
			old:      "retry = var.retry",
			repl:     "retry = merge({ interval_seconds = 1 }, coalesce(var.retry, {}))",
		},
		{
			name:     "ignore body changes concat",
			variable: "ignore_body_changes",
			content:  validIgnoreBodyChanges,
			old:      "ignore_body_changes = var.ignore_body_changes.widget",
			repl:     `ignore_body_changes = concat(["properties.fixed"], coalesce(var.ignore_body_changes.widget, []))`,
		},
		{
			name:     "resource types lookup",
			variable: "resource_types",
			content:  validResourceTypes,
			old:      "type = var.resource_types.widget",
			repl:     `type = lookup(var.resource_types, "widget", "Microsoft.Example/widgets@2024-01-01")`,
		},
		{
			name:     "selected local object lookup",
			variable: "resource_types",
			content:  validResourceTypes,
			old:      "type = var.resource_types.widget",
			repl:     `type = lookup(local.overrides, "widget", "Microsoft.Example/widgets@2024-01-01")`,
			locals: `
locals {
  overrides = {
    widget = var.resource_types.widget
  }
}`,
		},
		{
			name:     "template wrap",
			variable: "retry",
			content:  validRetry,
			old:      "retry = var.retry",
			repl:     `retry = "${var.retry}"`,
		},
		{
			name:     "template parts",
			variable: "resource_types",
			content:  validResourceTypes,
			old:      "type = var.resource_types.widget",
			repl:     `type = "Microsoft.Example/${var.resource_types.widget}"`,
		},
		{
			name:     "composed retry object",
			variable: "retry",
			content:  validRetry,
			old:      "retry = var.retry",
			repl:     `retry = { interval_seconds = var.retry.interval_seconds }`,
		},
		{
			name:     "parenthesized local member",
			variable: "resource_types",
			content:  validResourceTypes,
			old:      "type = var.resource_types.widget",
			repl:     `type = (coalesce(local.overrides, {})).widget`,
			locals: `
locals {
  overrides = var.resource_types
}`,
		},
		{
			name:     "timeouts block",
			variable: "timeouts",
			content:  validTimeouts,
			old: `timeouts {
    create = try(var.timeouts.create, null)
    read   = try(var.timeouts.read, null)
    update = try(var.timeouts.update, null)
    delete = try(var.timeouts.delete, null)
  }`,
			repl: `timeouts {
  create = coalesce(var.timeouts.create, "30m")
  delete = coalesce(var.timeouts.delete, "30m")
}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rule := requiredInterfaceRule(t, test.variable)
			content := strings.Replace(test.content, test.old, test.repl, 1) + test.locals
			runner := helper.TestRunner(t, map[string]string{"main.tf": content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			helper.AssertIssues(t, helper.Issues{}, runner.Issues)
		})
	}
}

func TestAzapiRequiredInterfacesRejectConditionOnlyReferences(t *testing.T) {
	tests := []struct {
		name     string
		variable string
		content  string
		old      string
		repl     string
	}{
		{"resource types", "resource_types", validResourceTypes, "type = var.resource_types.widget", "type = var.resource_types != null ? null : null"},
		{"retry", "retry", validRetry, "retry = var.retry", "retry = var.retry != null ? null : null"},
		{"timeouts", "timeouts", validTimeouts, `timeouts {
    create = try(var.timeouts.create, null)
    read   = try(var.timeouts.read, null)
    update = try(var.timeouts.update, null)
    delete = try(var.timeouts.delete, null)
  }`, `timeouts {
    create = var.timeouts != null ? null : null
  }`},
		{"ignore body changes", "ignore_body_changes", validIgnoreBodyChanges, "ignore_body_changes = var.ignore_body_changes.widget", "ignore_body_changes = var.ignore_body_changes != null ? null : null"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rule := requiredInterfaceRule(t, test.variable)
			runner := helper.TestRunner(t, map[string]string{"main.tf": strings.Replace(test.content, test.old, test.repl, 1)})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if len(runner.Issues) != 1 {
				t.Fatalf("expected one wiring issue, got %#v", runner.Issues)
			}
		})
	}
}

func TestAzapiRequiredInterfacesRejectUnrelatedLocalObjectMember(t *testing.T) {
	rule := requiredInterfaceRule(t, "resource_types")
	content := strings.Replace(
		validResourceTypes,
		"type = var.resource_types.widget",
		"type = local.overrides.selected",
		1,
	) + `
locals {
  overrides = {
    selected  = "Microsoft.Example/widgets@2024-01-01"
    unrelated = var.resource_types.widget
  }
}`
	runner := helper.TestRunner(t, map[string]string{"main.tf": content})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(runner.Issues) != 1 {
		t.Fatalf("expected one wiring issue, got %#v", runner.Issues)
	}
}

func TestAzapiRequiredInterfacesAcceptJSONLocalAlias(t *testing.T) {
	rule := requiredInterfaceRule(t, "retry")
	tests := []struct {
		name       string
		expression string
		locals     string
		wantIssues int
	}{
		{"direct", "local.retry", `{"locals":{"retry":"${var.retry}"}}`, 0},
		{"composed", "local.retry", `{"locals":{"retry":"${merge({}, coalesce(var.retry, {}))}"}}`, 0},
		{"selected object", "local.overrides.selected", `{"locals":{"overrides":{"selected":"${var.retry}"}}}`, 0},
		{"condition only", "local.retry", `{"locals":{"retry":"${var.retry != null ? null : null}"}}`, 1},
		{"unselected object", "local.overrides.selected", `{"locals":{"overrides":{"selected":null,"unrelated":"${var.retry}"}}}`, 1},
		{"literal mention", "local.retry", `{"locals":{"retry":"var.retry"}}`, 1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{
				"main.tf":      strings.Replace(validRetry, "retry = var.retry", "retry = "+test.expression, 1),
				"main.tf.json": test.locals,
			})
			if err := rule.Check(runner); err != nil {
				t.Fatal(err)
			}
			if len(runner.Issues) != test.wantIssues {
				t.Fatalf("expected %d wiring issues, got %#v", test.wantIssues, runner.Issues)
			}
		})
	}
}

func TestAzapiRequiredInterfacesRejectLocalAliasCycles(t *testing.T) {
	rule := requiredInterfaceRule(t, "retry")
	content := strings.Replace(validRetry, "retry = var.retry", "retry = local.first", 1) + `
locals {
  first  = local.second
  second = local.first
}`
	runner := helper.TestRunner(t, map[string]string{"main.tf": content})
	if err := rule.Check(runner); err != nil {
		t.Fatal(err)
	}
	if len(runner.Issues) != 1 {
		t.Fatalf("expected one wiring issue, got %#v", runner.Issues)
	}
}

func TestAzapiRequiredInterfacesRejectUnselectedLookupMember(t *testing.T) {
	rule := requiredInterfaceRule(t, "retry")
	content := strings.Replace(validRetry, "retry = var.retry", `retry = lookup(local.overrides, "selected", null)`, 1) + `
locals {
  overrides = {
    selected  = null
    unrelated = var.retry
  }
}`
	runner := helper.TestRunner(t, map[string]string{"main.tf": content})
	if err := rule.Check(runner); err != nil {
		t.Fatal(err)
	}
	if len(runner.Issues) != 1 {
		t.Fatalf("expected one wiring issue, got %#v", runner.Issues)
	}
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

package rules_test

import (
	"testing"

	rules "github.com/Azure/tflint-ruleset-avm/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestDisallowedProviderRule_RequiredProviders(t *testing.T) {
	const content = `terraform {
    required_providers {
        azurerm = {
            source  = "hashicorp/azurerm"
            version = "~> 4.0"
        }
    }
}`
	rule := rules.NewDisallowedProviderRule("azurerm", "hashicorp/azurerm")
	runner := helper.TestRunner(t, map[string]string{"main.tf": content})
	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(runner.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(runner.Issues))
	}
	if runner.Issues[0].Rule.Name() != rule.Name() {
		t.Fatalf("unexpected rule name: %s", runner.Issues[0].Rule.Name())
	}
}

func TestDisallowedProviderRule_ResourceUsage(t *testing.T) {
	const content = `resource "random_string" "allowed" {
    length = 8
}
resource "azurerm_resource_group" "rg" {
    name     = "rg1"
    location = "uksouth"
}`
	rule := rules.NewDisallowedProviderRule("azurerm", "hashicorp/azurerm")
	runner := helper.TestRunner(t, map[string]string{"main.tf": content})
	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(runner.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(runner.Issues))
	}
	if runner.Issues[0].Rule.Name() != rule.Name() {
		t.Fatalf("unexpected rule name: %s", runner.Issues[0].Rule.Name())
	}
}

func TestDisallowedProviderRule_NoIssues(t *testing.T) {
	const content = `terraform {
    required_providers {
        random = {
            source  = "hashicorp/random"
            version = "~> 3.0"
        }
    }
}
resource "random_string" "s" {
    length = 4
}`
	rule := rules.NewDisallowedProviderRule("azurerm", "hashicorp/azurerm")
	runner := helper.TestRunner(t, map[string]string{"main.tf": content})
	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(runner.Issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(runner.Issues))
	}
}

package rules_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Azure/tflint-ruleset-avm/rules"
	"github.com/stretchr/testify/assert"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func TestAzurermProviderVersionRule(t *testing.T) {
	const upperBoundIssue = "provider `azurerm`'s version must exclude 5.0.0 and later, got %s. Recommended version constraint `~> 4.0`"
	const compatibilityIssue = "provider `azurerm`'s version should satisfy 4.999.0, got %s. Recommended version constraint `~> 4.0`"

	cases := []struct {
		desc       string
		constraint string
		wantIssue  string
	}{
		{desc: "next major allowed", constraint: ">= 4.2, < 6.0", wantIssue: upperBoundIssue},
		{desc: "open ended", constraint: ">= 4", wantIssue: upperBoundIssue},
		{desc: "excluded next major boundary", constraint: ">= 4, != 5.0.0", wantIssue: upperBoundIssue},
		{desc: "excluded future sentinel", constraint: ">= 4, != 5.999.0", wantIssue: upperBoundIssue},
		{desc: "multiple exclusions", constraint: ">= 4.2, < 6, != 5.0.0, != 5.999.0", wantIssue: upperBoundIssue},
		{desc: "part of next major allowed", constraint: ">= 4, < 5.1, != 5.0.0", wantIssue: upperBoundIssue},
		{desc: "inclusive next major bound", constraint: ">= 4, <= 5.0", wantIssue: upperBoundIssue},
		{desc: "inclusive bound without minimum", constraint: "<= 5.0.0", wantIssue: upperBoundIssue},
		{desc: "single segment pessimistic constraint", constraint: "~> 4", wantIssue: upperBoundIssue},
		{desc: "exclusion without upper bound", constraint: "!= 5", wantIssue: upperBoundIssue},
		{desc: "recommended constraint", constraint: "~> 4.0"},
		{desc: "minor minimum", constraint: "~> 4.2"},
		{desc: "explicit upper bound", constraint: ">= 4.2, < 5"},
		{desc: "upper bound without minimum", constraint: "< 5.0.0"},
		{desc: "earlier major remains accepted", constraint: ">= 3, < 5"},
		{desc: "exclusion with safe bound", constraint: ">= 4, != 4.2.0, < 5"},
		{desc: "excluded next major with safe bound", constraint: "!= 5.0.0, ~> 4.0"},
		{desc: "tightest bound wins", constraint: "< 6, >= 4.2, < 5"},
		{desc: "inclusive bound with excluded endpoint", constraint: "<= 5.0.0, != 5.0.0"},
		{desc: "normalized excluded endpoint", constraint: " != v5.0.0+build, <= 5.0 "},
		{desc: "inclusive current major bound", constraint: "<= 4.999.0"},
		{desc: "exact compatibility version", constraint: "= 4.999.0"},
		{desc: "implicit exact compatibility version", constraint: "4.999.0"},
		{desc: "narrow minor remains rejected", constraint: "~> 4.2.0", wantIssue: compatibilityIssue},
		{desc: "narrow range remains rejected", constraint: ">= 4.2, < 4.3", wantIssue: compatibilityIssue},
		{desc: "exact pin remains rejected", constraint: "4.2.0", wantIssue: compatibilityIssue},
		{desc: "excluded compatibility version remains rejected", constraint: "~> 4.0, != 4.999.0", wantIssue: compatibilityIssue},
		{desc: "next major only", constraint: "~> 5.0", wantIssue: compatibilityIssue},
	}

	rule := registeredProviderVersionRule(t, "azurerm")
	if !rule.Enabled() || rule.Severity() != tflint.ERROR {
		t.Fatal("AzureRM version rule must be enabled by default with error severity")
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"terraform.tf": fmt.Sprintf(`terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = %q
    }
  }
}`, c.constraint)})
			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			wantIssue := ""
			if c.wantIssue != "" {
				wantIssue = fmt.Sprintf(c.wantIssue, c.constraint)
			}
			assertProviderVersionIssue(t, rule, runner.Issues, wantIssue)
		})
	}
}

func TestAzurermProviderVersionRuleDeclarations(t *testing.T) {
	cases := []struct {
		desc      string
		config    string
		wantIssue string
		wantError string
	}{
		{desc: "no terraform block"},
		{desc: "no required providers", config: "terraform {}"},
		{desc: "empty required providers", config: `terraform {
  required_providers {}
}`},
		{desc: "AzAPI only module", config: `terraform {
  required_providers {
    azapi = {
      source  = "Azure/azapi"
      version = "~> 2.0"
    }
  }
}`},
		{desc: "case insensitive source", config: `terraform {
  required_providers {
    azurerm = {
      source  = "HashiCorp/AzureRM"
      version = "~> 4.0"
    }
  }
}`},
		{desc: "wrong source", config: `terraform {
  required_providers {
    azurerm = {
      source  = "other/azurerm"
      version = "~> 4.0"
    }
  }
}`, wantIssue: "provider `azurerm`'s source should be hashicorp/azurerm, got other/azurerm"},
		{desc: "malformed constraint", config: `terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">= 4 || < 5"
    }
  }
}`, wantError: "invalid version constraint:"},
	}

	rule := registeredProviderVersionRule(t, "azurerm")
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"terraform.tf": c.config})
			err := rule.Check(runner)
			if c.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), c.wantError) {
					t.Fatalf("expected error containing %q, got %v", c.wantError, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			assertProviderVersionIssue(t, rule, runner.Issues, c.wantIssue)
		})
	}
}

func TestOtherProviderVersionRulesKeepOpenEndedConstraints(t *testing.T) {
	cases := []struct {
		provider   string
		source     string
		constraint string
	}{
		{provider: "azapi", source: "Azure/azapi", constraint: ">= 2.0"},
		{provider: "modtm", source: "Azure/modtm", constraint: ">= 0.3"},
	}
	for _, c := range cases {
		t.Run(c.provider, func(t *testing.T) {
			rule := registeredProviderVersionRule(t, c.provider)
			runner := helper.TestRunner(t, map[string]string{"terraform.tf": fmt.Sprintf(`terraform {
  required_providers {
    %s = {
      source  = %q
      version = %q
    }
  }
}`, c.provider, c.source, c.constraint)})
			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			helper.AssertIssuesWithoutRange(t, helper.Issues{}, runner.Issues)
		})
	}
}

func registeredProviderVersionRule(t *testing.T, provider string) tflint.Rule {
	t.Helper()
	name := fmt.Sprintf("avm_provider_%s_version_constraint", provider)
	for _, rule := range rules.Rules {
		if rule.Name() == name {
			return rule
		}
	}
	t.Fatalf("provider version rule %q is not registered", name)
	return nil
}

func assertProviderVersionIssue(t *testing.T, rule tflint.Rule, issues helper.Issues, message string) {
	t.Helper()
	if message == "" {
		assert.Empty(t, issues)
		return
	}
	if assert.Len(t, issues, 1) {
		assert.Equal(t, rule.Name(), issues[0].Rule.Name())
		assert.Equal(t, rule.Severity(), issues[0].Rule.Severity())
		assert.Equal(t, message, issues[0].Message)
	}
}

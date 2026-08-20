package rules

import (
	"testing"

	basic "github.com/Azure/tflint-ruleset-avm/basic/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func TestConsolidatedRulesAreRegisteredDisabledByDefault(t *testing.T) {
	registered := make(map[string]tflint.Rule, len(Rules))
	for _, rule := range Rules {
		if _, exists := registered[rule.Name()]; exists {
			t.Fatalf("rule %q is registered more than once", rule.Name())
		}
		registered[rule.Name()] = rule
	}

	for _, rule := range basic.Rules {
		registeredRule, exists := registered[rule.Name()]
		if !exists {
			t.Errorf("consolidated rule %q is not registered", rule.Name())
			continue
		}
		if registeredRule.Enabled() {
			t.Errorf("consolidated rule %q must remain disabled by default", rule.Name())
		}
	}
}

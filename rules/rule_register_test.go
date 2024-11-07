package rules_test

import (
	"testing"

	"github.com/Azure/tflint-ruleset-avm/rules"
)

func TestDuplicateRuleNames(t *testing.T) {
	rules := rules.Rules

	names := make(map[string]bool)
	for _, rule := range rules {
		name := rule.Name()
		if names[name] {
			t.Errorf("duplicate rule name: %s", name)
		}
		names[name] = true
	}
}

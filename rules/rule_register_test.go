package rules_test

import (
	"regexp"
	"testing"

	"github.com/Azure/tflint-ruleset-avm/rules"
	"github.com/stretchr/testify/assert"
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

func TestRuleNamesUseCanonicalScheme(t *testing.T) {
	pattern := regexp.MustCompile(`^avm_[a-z0-9]+(?:_[a-z0-9]+)*$`)

	for _, rule := range rules.Rules {
		assert.Regexp(t, pattern, rule.Name())
	}
}

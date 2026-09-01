package rules_test

import (
	"regexp"
	"testing"

	"github.com/Azure/tflint-ruleset-avm/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	require.Len(t, rules.RuleNameMigrations, len(rules.Rules))
	for _, rule := range rules.Rules {
		assert.Regexp(t, pattern, rule.Name())
	}
}

package rules

import "testing"

func TestAllRulesAreEnabledByDefault(t *testing.T) {
	for _, rule := range Rules {
		if !rule.Enabled() {
			t.Errorf("rule %q must be enabled by default", rule.Name())
		}
	}
}

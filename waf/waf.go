// Package waf contains the rules for Well Architected Alignment.
// To add a new rule, create a new file and add a new function that returns a new rule.
// Then add the rule to the Rules slice.
package waf

import (
	"reflect"

	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// WafRules is a helper struct. Methods are created on this type that generate the rules for the WAF package.
// We then use reflection to automatically generate a slice of the rules to add the the ruleset.
type WafRules struct{}

// GetRules uses reflection to iterate over all the methods of the WafRules struct and add them to a slice of Rules to be included in the ruleset.
// See `GetRules()` implementation for more detail.
func GetRules() []tflint.Rule {
	// Create a slice to add the rules to
	rules := []tflint.Rule{}

	// Get an instance of the WefRules struct we can use to iterate and call functions on
	wafRulesInstance := reflect.ValueOf(WafRules{})

	// Iterate over all the functions of the WafRules struct
	for i := 0; i < wafRulesInstance.NumMethod(); i++ {
		// Get an instance of the function
		methodInstance := wafRulesInstance.Method(i)

		// Call the function with no parameters (WAF rules don't require them), this returns a slice of outputs
		ruleOutputs := methodInstance.Call([]reflect.Value{})

		// Cast the first (and only for WAF rules) output to its proper type
		rule := ruleOutputs[0].Interface().(tflint.Rule)

		// Apend the rule to slice
		rules = append(rules, rule)
	}

	return rules
}

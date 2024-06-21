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

// GetRules uses reflection to iterate over all the methods of the WafRules struct and add them to a slice of Rules to be included in the ruleset.
// See `GetRules()` for more detail.
type WafRules struct{}

func GetRules() []tflint.Rule {
	rules := []tflint.Rule{}

	wafRules := reflect.TypeOf(WafRules{})

	for i := 0; i < wafRules.NumMethod(); i++ {
		method := wafRules.Method(i)
		rule := reflect.ValueOf(WafRules{}).MethodByName(method.Name).Call([]reflect.Value{})
		rules = append(rules, rule[0].Interface().(tflint.Rule))
	}

	return rules
}

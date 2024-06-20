// Package waf contains the rules for Well Architected Alignment.
// To add a new rule, create a new file and add a new function that returns a new rule.
// Then add the rule to the Rules slice.
package waf

import (
	"reflect"

	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

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

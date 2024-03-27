package interfaces_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/Azure/tflint-ruleset-avm/interfaces"
	"github.com/hashicorp/hcl/v2"
	"github.com/matt-FFFFFF/tfvarcheck/varcheck"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

var simpleVarTypeString = `object({
  kind = string
  name = optional(string, null)
})`

var simpleType = interfaces.StringToTypeConstraintWithDefaults(simpleVarTypeString)

// SimpleVar represents a simple variable type used for testing.
var SimpleVar = interfaces.AvmInterface{
	VarCheck:      varcheck.NewVarCheck(simpleType, cty.NullVal(cty.DynamicPseudoType), true),
	RuleName:      "simple",
	VarTypeString: simpleVarTypeString,
	RuleEnabled:   true,
	RuleLink:      "https://simple",
	RuleSeverity:  tflint.ERROR,
}

func TestIncorrectInterfaceTypeStringShouldPanic(t *testing.T) {
	defer func() {
		err := recover()
		assert.NotNil(t, err)
	}()
	incorrectTypeString := `strin`
	interfaces.StringToTypeConstraintWithDefaults(incorrectTypeString)
	t.Fatal("incorrect type should panic")
}

func TestSimpleInterface(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		JSON     bool
		Expected helper.Issues
	}{
		{
			Name: "not simple variable",
			Content: `
variable "not_simple" {
	default = "default"
}`,
			Expected: helper.Issues{},
		},
		{
			Name:     "correct",
			Content:  toTerraformVarType(SimpleVar),
			Expected: helper.Issues{},
		},
		{
			Name: "incorrect nullable true",
			Content: fmt.Sprintf(`
variable "simple" {
	default = null
	type = %s
	nullable = true
}`, interfaces.LockTypeString),
			Expected: helper.Issues{
				&helper.Issue{
					Rule:    interfaces.NewVarCheckRuleFromAvmInterface(SimpleVar),
					Message: "nullable should not be set.",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 8, Column: 2},
						End:      hcl.Pos{Line: 8, Column: 17},
					},
				},
			},
		},
		{
			Name: "too many attributes in object",
			Content: `
variable "simple" {
	default = null
	type = object({
		kind     = string
		name     = optional(string, null)
		unwanted = string
	})
}`,
			Expected: helper.Issues{
				&helper.Issue{
					Rule:    interfaces.NewVarCheckRuleFromAvmInterface(SimpleVar),
					Message: fmt.Sprintf("variable type does not comply with the interface specification:\n\n%s", simpleVarTypeString),
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 4, Column: 2},
						End:      hcl.Pos{Line: 8, Column: 4},
					},
				},
			},
		},
		{
			Name: "missing attribute in object, but correct number of attributes",
			Content: `
variable "simple" {
	default = null
	type = object({
		# kind is missing
		name     = optional(string, null)
		unwanted = string
	})
}`,
			Expected: helper.Issues{
				&helper.Issue{
					Rule:    interfaces.NewVarCheckRuleFromAvmInterface(SimpleVar),
					Message: fmt.Sprintf("variable type does not comply with the interface specification:\n\n%s", simpleVarTypeString),
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 4, Column: 2},
						End:      hcl.Pos{Line: 8, Column: 4},
					},
				},
			},
		},
		{
			Name: "kind attribute incorrect type",
			Content: `
variable "simple" {
	default = null
	type = object({
		kind = number
		name = optional(string, null)
	})
}`,
			Expected: helper.Issues{
				&helper.Issue{
					Rule:    interfaces.NewVarCheckRuleFromAvmInterface(SimpleVar),
					Message: fmt.Sprintf("variable type does not comply with the interface specification:\n\n%s", simpleVarTypeString),
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 4, Column: 2},
						End:      hcl.Pos{Line: 7, Column: 4},
					},
				},
			},
		},
		{
			Name: "incorrect attribute type",
			Content: `
variable "simple" {
	default = null
	type = object({
		kind = string
		name = optional(number, null)
	})
}`,
			Expected: helper.Issues{
				&helper.Issue{
					Rule:    interfaces.NewVarCheckRuleFromAvmInterface(SimpleVar),
					Message: fmt.Sprintf("variable type does not comply with the interface specification:\n\n%s", simpleVarTypeString),
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 4, Column: 2},
						End:      hcl.Pos{Line: 7, Column: 4},
					},
				},
			},
		},
		{
			Name: "incorrect optional default",
			Content: `
variable "simple" {
	default = null
	type = object({
		kind = string
		name = optional(string, "")
	})
}`,
			Expected: helper.Issues{
				&helper.Issue{
					Rule:    interfaces.NewVarCheckRuleFromAvmInterface(SimpleVar),
					Message: fmt.Sprintf("variable type does not comply with the interface specification:\n\n%s", simpleVarTypeString),
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 4, Column: 2},
						End:      hcl.Pos{Line: 7, Column: 4},
					},
				},
			},
		},
		{
			Name: "no default",
			Content: `
variable "simple" {
	type = object({
		kind = string
		name = optional(string, null)
	})
}`,
			Expected: helper.Issues{
				&helper.Issue{
					Rule:    interfaces.NewVarCheckRuleFromAvmInterface(SimpleVar),
					Message: "`simple` default not declared",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 18},
					},
				},
			},
		},
		{
			Name: "incorrect default",
			Content: `
variable "simple" {
	default = {
		kind = "CanNotDelete"
	}
	type = object({
		kind = string
		name = optional(string, null)
	})
}`,
			Expected: helper.Issues{
				&helper.Issue{
					Rule:    interfaces.NewVarCheckRuleFromAvmInterface(SimpleVar),
					Message: fmt.Sprintf("default value is not correct, see: %s", SimpleVar.RuleLink),
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 18},
					},
				},
			},
		},
	}

	rule := interfaces.NewVarCheckRuleFromAvmInterface(SimpleVar)

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			filename := "variables.tf"
			if tc.JSON {
				filename += ".json"
			}

			runner := helper.TestRunner(t, map[string]string{filename: tc.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}

package rules

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"testing"
)

func Test_VariableNullableEqualTrueShouldRaiseIssue(t *testing.T) {
	code := `
variable "var" {
  type     = string
  nullable = true
}

variable "var2" {
  type     = string
  nullable = false
}

variable "var3" {
  type     = string
}
`
	rule := NewTerraformVariableNullableFalseRule()
	mockFileName := "test.tf"
	expectedIssue := helper.Issues{
		{
			Rule:    rule,
			Message: rule.errorMessage(),
			Range: struct {
				Filename   string
				Start, End hcl.Pos
			}{Filename: mockFileName, Start: hcl.Pos{
				Line:   4,
				Column: 3,
			}, End: hcl.Pos{
				Line:   4,
				Column: 18,
			}},
		},
	}
	runner := helper.TestRunner(t, map[string]string{mockFileName: code})
	if err := rule.Check(runner); err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}
	AssertIssues(t, expectedIssue, runner.Issues)
}

func Test_VariableNullableEqualFalseShouldNotRaiseIssue(t *testing.T) {
	code := `
variable "var" {
  type     = string
  nullable = false
}
`
	rule := NewTerraformVariableNullableFalseRule()
	mockFileName := "test.tf"
	runner := helper.TestRunner(t, map[string]string{mockFileName: code})
	if err := rule.Check(runner); err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}
	assert.Equal(t, 0, len(runner.Issues))
}

func Test_VariableWithoutNullableShouldNotRaiseIssue(t *testing.T) {
	code := `
variable "var" {
  type     = string
}
`
	rule := NewTerraformVariableNullableFalseRule()
	mockFileName := "test.tf"
	runner := helper.TestRunner(t, map[string]string{mockFileName: code})
	if err := rule.Check(runner); err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}
	assert.Equal(t, 0, len(runner.Issues))
}

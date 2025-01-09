package rules_test

import (
	"github.com/Azure/tflint-ruleset-avm/attrvalue"
	"github.com/Azure/tflint-ruleset-avm/rules"
	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"testing"
)

func TestPartialContentEvaluationFailureShouldFailValidTerraformConfigRule(t *testing.T) {
	sut := rules.NewValidTemplateInterpolationRule()
	content := `
    variable "null" {
	  type    = string
      default = null
    }
    locals {
      name = "${var.null}-rg"
    }
	resource "azurerm_resource_group" "example" {
	  location = "eastus"
	  name     = local.name
	}`

	filename := "main.tf"
	runner := helper.TestRunner(t, map[string]string{filename: content})
	stub := gostub.Stub(&attrvalue.AppFs, mockFs(content))
	defer stub.Reset()
	err := sut.Check(runner)
	require.NotNil(t, err)
	assert.Regexp(t, rules.TemplateInterpolationErrorRegex, err.Error())
}

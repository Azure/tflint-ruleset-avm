package waf_test

import (
	"testing"

	"github.com/Azure/tflint-ruleset-avm/attrvalue"
	"github.com/Azure/tflint-ruleset-avm/waf"
	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/assert"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestPartialContentEvaluationFailureShouldNotFailWafRule(t *testing.T) {
	var w waf.WafRules
	sut := w.AzurermLegacyVirtualMachineNotAllowed()
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
	if err := sut.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	assert.Empty(t, runner.Issues)
}

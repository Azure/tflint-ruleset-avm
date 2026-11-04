package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformPatternModuleSource(t *testing.T) {
	rule := NewTerraformPatternModuleSourceRule()
	runRuleTestCases(t, rule, []ruleTestCase{
		{
			Name: "resource, utility, and local modules are allowed",
			Files: map[string]string{"main.tf": `
module "avm_interfaces" {
  source  = "Azure/avm-utl-interfaces/azure"
  version = "0.5.0"
}

module "vnet" {
  source  = "Azure/avm-res-network-virtualnetwork/azurerm"
  version = "0.8.0"
}

module "secret" {
  source = "./modules/secret"
}`},
			Expected: helper.Issues{},
		},
		{
			Name: "pattern modules are flagged",
			Files: map[string]string{"main.tf": `
module "alz" {
  source  = "Azure/avm-ptn-alz/azurerm"
  version = "0.12.0"
}

module "hub" {
  source = "registry.terraform.io/azure/AVM-PTN-hubnetworking/azurerm//modules/hub"
}`},
			Expected: helper.Issues{
				{Rule: rule, Message: "Module `alz` must not call the AVM pattern module `Azure/avm-ptn-alz/azurerm`; resource modules may only call AVM resource or utility modules"},
				{Rule: rule, Message: "Module `hub` must not call the AVM pattern module `registry.terraform.io/azure/AVM-PTN-hubnetworking/azurerm//modules/hub`; resource modules may only call AVM resource or utility modules"},
			},
		},
		{
			Name: "non-literal sources are skipped",
			Files: map[string]string{"main.tf": `
module "dynamic" {
  source = local.source
}`},
			Expected: helper.Issues{},
		},
	})
}

func Test_isAVMPatternModuleSource(t *testing.T) {
	cases := map[string]bool{
		"Azure/avm-ptn-alz/azurerm":                      true,
		"azure/avm-ptn-alz/azurerm":                      true,
		"Azure/avm-ptn-alz/azurerm//modules/sub":         true,
		"registry.terraform.io/Azure/avm-ptn-alz/azurerm": true,
		"Azure/avm-res-keyvault-vault/azurerm":           false,
		"Azure/avm-utl-interfaces/azure":                 false,
		"contoso/avm-ptn-alz/azurerm":                    false,
		"./modules/avm-ptn-local":                        false,
		"../avm-ptn-alz":                                 false,
		"git::https://github.com/Azure/avm-ptn-alz.git":  false,
		"github.com/Azure/terraform-azurerm-avm-ptn-alz": false,
		"Azure/avm-ptn-alz":                              false,
	}
	for source, want := range cases {
		assert.Equal(t, want, isAVMPatternModuleSource(source), source)
	}
}

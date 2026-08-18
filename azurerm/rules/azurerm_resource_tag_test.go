package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"testing"
)

func Test_AzurermResourceTagRule(t *testing.T) {

	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "1. simple block",
			Content: `
resource "azurerm_container_group" "example" {
  location            = azurerm_resource_group.example.location
  name                = "example-continst"
  ip_address_type     = "Public"
  dns_name_label      = "aci-label"
  os_type             = "Linux"
  resource_group_name = azurerm_resource_group.example.name
  
  container {
    cpu    = "0.5"
    image  = "mcr.microsoft.com/azuredocs/aci-tutorial-sidecar"
    memory = "1.5"
    name   = "sidecar"
  }
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAzurermResourceTagRule(),
					Message: "`tags` argument is not set but supported in resource `azurerm_container_group`",
				},
			},
		},
		{
			Name: "2. tags are set anywhere if supported",
			Content: `
resource "azurerm_resource_group" "rg" {
  name     = "myTFResourceGroup"
  location = "westus2"
  tags = {
    Team = "DevOps"
  }
}

resource "azurerm_container_registry" "acr" {
  name                = "containerRegistry1"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku                 = "Premium"
  admin_enabled       = false
  tags = {
    environment = "testing"
  }

  georeplications {
    location                = "westeurope"
    zone_redundancy_enabled = true
    tags = {
      environment = "testing"
    }
  }
  dynamic "georeplications" {
    for_each = var.georeplications
    content {
      location                = georeplications.location
      zone_redundancy_enabled = georeplications.zone_redundancy_enabled
      tags = georeplications.tags
    }
  }
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAzurermResourceTagRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"config.tf": tc.Content})
		t.Run(tc.Name, func(t *testing.T) {
			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}
			AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}

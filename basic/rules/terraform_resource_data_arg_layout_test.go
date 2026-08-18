package rules

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformResourceDataArgLayout(t *testing.T) {

	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "1. separate attributes from blocks",
			Content: `
resource "azurerm_virtual_network" "vnet" {
  name                = "myTFVnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.rg.name
  container {
    name   = "sidecar"
    image  = "mcr.microsoft.com/azuredocs/aci-tutorial-sidecar"
    cpu    = "0.5"
    memory = "1.5"
  }
  dns_name_label      = "aci-label"
  os_type             = "Linux"
  container {
    name   = "sidecar"
    image  = "mcr.microsoft.com/azuredocs/aci-tutorial-sidecar"
    cpu    = "0.5"
    memory = "1.5"
  }
  tags = {
    Name = "VM network ${count.index}"
  }
}`,
			Expected: helper.Issues{
				{
					Rule: NewTerraformResourceDataArgLayoutRule(),
					Message: `Arguments are expected to be arranged in following Layout:
resource "azurerm_virtual_network" "vnet" {
  name                = "myTFVnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.rg.name
  dns_name_label      = "aci-label"
  os_type             = "Linux"
  tags = {
    Name = "VM network ${count.index}"
  }

  container {
    name   = "sidecar"
    image  = "mcr.microsoft.com/azuredocs/aci-tutorial-sidecar"
    cpu    = "0.5"
    memory = "1.5"
  }
  container {
    name   = "sidecar"
    image  = "mcr.microsoft.com/azuredocs/aci-tutorial-sidecar"
    cpu    = "0.5"
    memory = "1.5"
  }
}`,
				},
			},
		},
		{
			Name: "2. check for nested block",
			Content: `
resource "azurerm_virtual_network" "vnet" {
  name                = "myTFVnet"
  address_space       = ["10.0.0.0/16"]

  container {
    cpu    = "0.5"
    image  = "mcr.microsoft.com/azuredocs/aci-helloworld:latest"
    memory = "1.5"
    ports {
      port     = 443
      protocol = "TCP"
    }
    name   = "hello-world"
  }
}`,
			Expected: helper.Issues{
				{
					Rule: NewTerraformResourceDataArgLayoutRule(),
					Message: `Arguments are expected to be arranged in following Layout:
container {
  cpu    = "0.5"
  image  = "mcr.microsoft.com/azuredocs/aci-helloworld:latest"
  memory = "1.5"
  name   = "hello-world"

  ports {
    port     = 443
    protocol = "TCP"
  }
}`,
				},
			},
		},
		{
			Name: "3. Gap between different types of args",
			Content: `
resource "azurerm_virtual_network" "vnet" {
  count    = 4
  provider = azurerm.europe
  name                = "myTFVnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.rg.name
  dns_name_label      = "aci-label"
  os_type             = "Linux"
  tags = {
    Name = "VM network ${count.index}"
  }
  container {
    name   = "sidecar"
    image  = "mcr.microsoft.com/azuredocs/aci-tutorial-sidecar"
    cpu    = "0.5"
    memory = "1.5"
  }
  lifecycle {
    create_before_destroy = true
  }
  depends_on = [
    azurerm_resource_group.example
  ]
}`,
			Expected: helper.Issues{
				{
					Rule: NewTerraformResourceDataArgLayoutRule(),
					Message: `Arguments are expected to be arranged in following Layout:
resource "azurerm_virtual_network" "vnet" {
  provider = azurerm.europe
  count    = 4

  name                = "myTFVnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.rg.name
  dns_name_label      = "aci-label"
  os_type             = "Linux"
  tags = {
    Name = "VM network ${count.index}"
  }

  container {
    name   = "sidecar"
    image  = "mcr.microsoft.com/azuredocs/aci-tutorial-sidecar"
    cpu    = "0.5"
    memory = "1.5"
  }

  depends_on = [
    azurerm_resource_group.example
  ]

  lifecycle {
    create_before_destroy = true
  }
}`,
				},
			},
		},
		{
			Name: "4. Meta Arg",
			Content: `
resource "azurerm_virtual_network" "vnet1" {
  name          = "myTFVnet1"
  address_space = ["10.0.0.0/16"]
  provider      = azurerm.europe
  count         = 4
  depends_on = [
    azurerm_resource_group.example
  ]
  location = azurerm_resource_group.example.location
  lifecycle {
    create_before_destroy = true
  }
  resource_group_name = azurerm_resource_group.rg.name
  dns_name_label      = "aci-label"
  os_type             = "Linux"
  tags = {
    Name = "VM network ${count.index}"
  }

  container {
    name   = "sidecar"
    image  = "mcr.microsoft.com/azuredocs/aci-tutorial-sidecar"
    cpu    = "0.5"
    memory = "1.5"
  }
}

resource "azurerm_virtual_network" "vnet2" {
  name          = "myTFVnet2"
  address_space = ["10.0.0.0/16"]
  provider      = azurerm.europe
  depends_on = [
    azurerm_resource_group.example
  ]
  for_each = local.subnet_ids
  location = azurerm_resource_group.example.location
  lifecycle {
    create_before_destroy = true
  }
  resource_group_name = azurerm_resource_group.rg.name
  dns_name_label      = "aci-label"
  os_type             = "Linux"
  tags = {
    Name = "VM network ${each.key}"
  }

  container {
    name   = "sidecar"
    image  = "mcr.microsoft.com/azuredocs/aci-tutorial-sidecar"
    cpu    = "0.5"
    memory = "1.5"
  }
}`,
			Expected: helper.Issues{
				{
					Rule: NewTerraformResourceDataArgLayoutRule(),
					Message: `Arguments are expected to be arranged in following Layout:
resource "azurerm_virtual_network" "vnet1" {
  provider = azurerm.europe
  count    = 4

  name                = "myTFVnet1"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.rg.name
  dns_name_label      = "aci-label"
  os_type             = "Linux"
  tags = {
    Name = "VM network ${count.index}"
  }

  container {
    name   = "sidecar"
    image  = "mcr.microsoft.com/azuredocs/aci-tutorial-sidecar"
    cpu    = "0.5"
    memory = "1.5"
  }
  
  depends_on = [
    azurerm_resource_group.example
  ]

  lifecycle {
    create_before_destroy = true
  }
}`,
				},

				{
					Rule: NewTerraformResourceDataArgLayoutRule(),
					Message: `Arguments are expected to be arranged in following Layout:
resource "azurerm_virtual_network" "vnet2" {
  provider = azurerm.europe
  for_each = local.subnet_ids

  name                = "myTFVnet2"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.rg.name
  dns_name_label      = "aci-label"
  os_type             = "Linux"
  tags = {
    Name = "VM network ${each.key}"
  }

  container {
    name   = "sidecar"
    image  = "mcr.microsoft.com/azuredocs/aci-tutorial-sidecar"
    cpu    = "0.5"
    memory = "1.5"
  }

  depends_on = [
    azurerm_resource_group.example
  ]

  lifecycle {
    create_before_destroy = true
  }
}`,
				},
			},
		},

		{
			Name: "5. dynamic block",
			Content: `
resource "azurerm_container_group" "example" {
  location            = azurerm_resource_group.example.location
  name                = "example-continst"
  os_type             = "Linux"
  resource_group_name = azurerm_resource_group.example.name

  dynamic "container" {
    content {
	  name   = container.value["name"]
	  image  = container.value["image"]
	  cpu 	 = container.value["cpu"]
	  ports {
		port     = 443
		protocol = "TCP"
	  }
      memory = container.value["memory"]
    }
	for_each = var.containers
  }
}`,
			Expected: helper.Issues{
				{
					Rule: NewTerraformResourceDataArgLayoutRule(),
					Message: `Arguments are expected to be arranged in following Layout:
dynamic "container" {
  for_each = var.containers

  content {
    name   = container.value["name"]
    image  = container.value["image"]
    cpu 	 = container.value["cpu"]
    memory = container.value["memory"]

    ports {
	  port     = 443
	  protocol = "TCP"
    }
  }
}`,
				},
			},
		},

		{
			Name: "6. correct layout",
			Content: `
resource "azurerm_virtual_network" "vnet" {
  provider = azurerm.europe
  for_each = local.subnet_ids

  name                = "myTFVnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.rg.name
  dns_name_label      = "aci-label"
  os_type             = "Linux"
  tags = {
    Name = "VM network ${each.key}"
  }

  container {
    name   = "sidecar"
    image  = "mcr.microsoft.com/azuredocs/aci-tutorial-sidecar"
    cpu    = "0.5"
    memory = "1.5"
  }
  container {
    cpu    = "0.5"
    image  = "mcr.microsoft.com/azuredocs/aci-helloworld:latest"
    memory = "1.5"
    name   = "hello-world"

    ports {
      port     = 443
      protocol = "TCP"
    }
  }

  depends_on = [
    azurerm_resource_group.example
  ]

  lifecycle {
    create_before_destroy = true
  }
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "7. datasource",
			Content: `
data "azurerm_resources" "example" {
  testblock {
    environment = "production"
    role        = "webserver"
  }
  resource_group_name = "example-resources"
}`,
			Expected: helper.Issues{
				{
					Rule: NewTerraformResourceDataArgLayoutRule(),
					Message: `Arguments are expected to be arranged in following Layout:
data "azurerm_resources" "example" {
  resource_group_name = "example-resources"

  testblock {
    environment = "production"
    role        = "webserver"
  }
}`,
				},
			},
		},
		{
			Name: "correct arguments only resource",
			Content: `
resource "azurerm_virtual_network" "vnet" {
  name                = "myTFVnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.rg.name
  dns_name_label      = "aci-label"
  os_type             = "Linux"
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewTerraformResourceDataArgLayoutRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"config.tf": tc.Content})
		t.Run(tc.Name, func(t *testing.T) {
			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}
			AssertIssuesWithoutRange(t, tc.Expected, runner.Issues)
		})
	}
}

func Test_TerraformResourceDataArgLayout_DoNotCareNormalArgOrder(t *testing.T) {
	code := `
resource "azurerm_virtual_network" "vnet" {
  provider = azurerm.europe
  for_each = local.subnet_ids

  name                = "myTFVnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.rg.name
  dns_name_label      = "aci-label"
  os_type             = "Linux"
  tags = {
    Name = "VM network ${each.key}"
  }
}`

	rule := NewTerraformResourceDataArgLayoutRule()
	runner := helper.TestRunner(t, map[string]string{"config.tf": code})
	if err := rule.Check(runner); err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}
	assert.Empty(t, runner.Issues)
}

func Test_TerraformResourceDataArgLayout_DoNotCareNormalNestedBlockOrder(t *testing.T) {
	code := `
resource "azurerm_virtual_network" "vnet" {
  provider = azurerm.europe
  for_each = local.subnet_ids

  name                = "myTFVnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.rg.name
  dns_name_label      = "aci-label"
  os_type             = "Linux"
  tags = {
    Name = "VM network ${each.key}"
  }

  subnet {
	name 		   = "testsubnet"
    address_prefix = "192.168.0.0/24"
  }
  ddos_protection_plan {
	id      = azurerm_network_ddos_protection_plan.test.id
	enabled = true
  }
}`

	rule := NewTerraformResourceDataArgLayoutRule()
	runner := helper.TestRunner(t, map[string]string{"config.tf": code})
	if err := rule.Check(runner); err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}
	assert.Empty(t, runner.Issues)
}

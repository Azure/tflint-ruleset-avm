terraform {
  required_version = "~> 1.7.0"
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">= 3.97.0, < 3.99.0"
    }
  }
  module "other-module" {
    source  = "Azure/terraform-azurerm-avm-res-keyvault-vault"
    version = "0.5.3"
  }
}

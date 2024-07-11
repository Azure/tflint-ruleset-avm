terraform {
  required_version = "~> 1.7.0"
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">= 3.97.0, < 4.0.0"
    }
    modtm = {
      source  = "Azure/modtm"
      version = "~> 0.3.0"
    }

  }
}
module "other-module" {
  source  = "Azure/avm-res-keyvault-vault/azurerm"
  version = "0.5.3"
}

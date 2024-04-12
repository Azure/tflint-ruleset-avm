terraform {
  required_version = "~> 1.7.0"
  required_providers {
    azurerm = {
      version = ">= 3.98.0"
      source  = "hashicorp/azurerm"
    }
  }
}
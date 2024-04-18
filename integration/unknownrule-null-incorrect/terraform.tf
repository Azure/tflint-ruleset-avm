terraform {
  required_version = "~> 1.7.0"
  required_providers {
    azurerm = {
      version = ">= 3.97.0, < 3.99.0"
      source  = "hashicorp/azurerm"
    }
  }
}
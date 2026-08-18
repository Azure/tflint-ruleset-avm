terraform {
  required_version = "~> 1.11"

  required_providers {
    modtm = {
      source  = "Azure/modtm"
      version = "~> 0.3"
    }
  }
}

terraform {
  required_version = ">= 0.13"
  required_providers {
    azurecaf = {
      source = "aztfmod/azurecaf"
      version = "1.1.3"
    }
  }
}

provider "azurerm" {
  version = ">=2.31.0"
  features {}
}

resource "azurecaf_name" "rg_name" {
  name          = "policytest"
  random_length = 8
  resource_type = "azurerm_resource_group"
}

resource "azurerm_resource_group" "test" {
  name     = azurecaf_name.rg_name.result
  location = "westeurope"
}

output "test_resource_id" {
  value = azurerm_resource_group.test.id
}

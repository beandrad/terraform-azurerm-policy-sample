terraform {
  required_version = ">= 0.13"
  required_providers {
    azurecaf = {
      source  = "aztfmod/azurecaf"
      version = "1.1.3"
    }
  }
}

provider "azurerm" {
  version = ">=2.31.0"
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurecaf_name" "rg_name" {
  name          = "policytest"
  resource_type = "azurerm_resource_group"
  random_length = 10
}

resource "azurerm_resource_group" "test" {
  name     = azurecaf_name.rg_name.result
  location = "uksouth"
}

resource "azurecaf_name" "kv_name" {
  name          = "policytest"
  resource_type = "azurerm_key_vault"
  random_length = 10
}

resource "azurerm_key_vault" "test" {
  name                        = azurecaf_name.kv_name.result
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  enabled_for_disk_encryption = true
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  soft_delete_enabled         = false
  purge_protection_enabled    = false

  sku_name = "standard"
}

output "test_resource_id" {
  value = azurerm_key_vault.test.id
}

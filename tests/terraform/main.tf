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

module "policies" {
  source                      = "../../"
  definition_management_group = var.management_group_name
}

data "azurerm_subscription" "current" {
}

resource "azurerm_policy_assignment" "test_resource_location" {
  name                 = "test-resource-location"
  scope                = data.azurerm_subscription.current.id
  policy_definition_id = module.policies.resource_location_id
  display_name         = "Unit test resource location initiative assignment"
  parameters           = <<PARAMETERS
{
  "allowedLocations": {
    "value": [ "uksouth", "ukwest" ]
  }
}
PARAMETERS
}

resource "azurerm_policy_assignment" "test_key_vault_security" {
  name                 = "test-key-vault-security"
  scope                = data.azurerm_subscription.current.id
  policy_definition_id = module.policies.key_vault_security_id
  display_name         = "Unit test key vault security initiative assignment"
}

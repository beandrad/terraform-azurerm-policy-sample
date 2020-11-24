resource "azurerm_policy_definition" "key_vault_soft_delete" {
  name                  = "key-vault-soft-delete"
  policy_type           = "Custom"
  mode                  = "Indexed"
  display_name          = "Key Vault Soft Delete"
  management_group_name = var.definition_management_group
  policy_rule           = file("${path.module}/policies/key-vault-soft-delete/policy-rule.json")
}

data "azurerm_policy_definition" "key_vault_diagnostic_logs" {
  # display_name = "Diagnostic logs in Key Vault should be enabled"
  name = "cf820ca0-f99e-4f3e-84fb-66e913812d21"
}

resource "azurerm_policy_set_definition" "key_vault_security" {
  name                  = "key-vault-security"
  policy_type           = "Custom"
  display_name          = "Key Vault Security"
  management_group_name = var.definition_management_group

  policy_definition_reference {
    policy_definition_id = azurerm_policy_definition.key_vault_soft_delete.id
  }

  policy_definition_reference {
    policy_definition_id = data.azurerm_policy_definition.key_vault_diagnostic_logs.id
  }
}

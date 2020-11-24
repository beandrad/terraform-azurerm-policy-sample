resource "azurerm_policy_definition" "resource_location" {
  name                  = "resource-location"
  policy_type           = "Custom"
  mode                  = "All"
  display_name          = "Resource Location"
  management_group_name = var.definition_management_group
  policy_rule           = file("${path.module}/policies/resource-location/policy-rule.json")
  parameters            = file("${path.module}/policies/resource-location/policy-parameters.json")
}

resource "azurerm_policy_set_definition" "resource_location" {
  name                  = "resource-location"
  policy_type           = "Custom"
  display_name          = "Resource Location"
  management_group_name = var.definition_management_group

  parameters = file("${path.module}/initiative-parameters.json")

  policy_definition_reference {
    policy_definition_id = azurerm_policy_definition.resource_location.id
    parameter_values = jsonencode({
      allowedLocations = {
        value = "[parameters('allowedLocations')]"
      }
    })
  }
}

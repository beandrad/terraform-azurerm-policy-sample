# Sample AzureRM policy module

The aim of this project is to provide a baseline on how to set up a Terraform module and, in particular, one defining custom policies and initiatives.

## Terraform module conventions

Hashicorp published [naming conventions and general guidelines](https://www.terraform.io/docs/modules/index.html) on the structure of Terraform modules.

Module repository names should be `terraform-<PROVIDER>-<NAME>`, where `PROVIDER` in our case is AzureRM and `NAME` should be a label that describes the infrastructure provided. Please, note the use of hyphens to separate all fields.

The standard module structure looks as follows:

- `main.tf`, `variables.tf`, `outputs.tf`. Files configuring a minimal module. Apart from `main.tf`, more complex modules may have additional resource configuration files.
- `README.md`. It describes the module and its intent. It may also include some examples of use and a diagram of the infrastructure defined if it is complex.
- `LICENSE`. The license under which the module is available.
- `Examples` (Optional). The `examples/` directory should include examples on how to use the module.

In the case of public modules, those must be hosted as public repos in GitHub. The repository description is used as the module description and therefore, should be a simple sentence that conveys what the module provides. More details on how to publish modules to the public registry can be found [here](https://www.terraform.io/docs/registry/modules/publish.html#publishing-a-public-module).

Private modules can be either downloaded from the corresponding Git repo or from a [private Terraform module registry](https://www.terraform.io/docs/cloud/registry/publish.html). In our case, we decided to access the Git repo directly for simplicity.

## Custom policies and initiatives

This module defines custom policies and initiatives under a management group (`definition_management_group` in `variables.tf`). Those resources could then be assigned to that management group or any other management group or subscription under it.

Azure policies are defined as JSON in the `policies/` folder; each policy has its own folder, with the file `policy-rule.json` including the definition and `policy-parameters.json` defining the parameters if applicable.

Policies are grouped into initiatives based on the type of resources they monitor or by industry level standards (such as the CIS Hardening Guidelines). This module only exposes initiatives, as opposed to also exposing policy definitions, to keep the outputs and the general module structure simpler. Each initiative is defined on its own terraform file and, apart from custom policies, it may include build-in policies.

Custom Policy definitions are created using the `azurerm_policy_definition` resource and built-in policies are imported using the `azurerm_policy_definition` data resource. Both resources are included in the corresponding initiatives Terraform configuration file; unless they are shared across initiatives, in which case they are defined in the `main.tf` file.

It is important to note that policy data resource should be imported using its policy `name` (as opposed to they `displayName`), since the `displayName` is not unique and it may change, whereas the `name` is always unique and the same unless the policy is deleted. In the configuration, we kept the `displayName` commented out as it describes the policy definition being imported.

## Running Terraform configuration locally

The Terraform configuration defined in this module can be run locally inside the devcontainer. One way of doing so is by creating an `.env` file that defines the values of the variables in `.env.template`. The variables defined in the `.env` file are exported when creating the devcontainer.

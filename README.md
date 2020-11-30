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

## Dependency issues

In a Terraform configuration, when an `azurerm_policy_definition` resource is referenced from an `azurerm_policy_set_definition` resource if the policy definition and the corresponding reference from the policy set are removed in a later terraform configuration, the `terraform apply` fails. The reason is that Terraform tries to delete the definition before updating the policy set.

This is the expected Terraform behavior in the presence of dependencies: (1) delete old resource, (2) create new resource, (3) update resource dependencies (in the case above only (1) and (3) are executed). Terraform provides a way of altering this behavior: `lifecycle.create_before_destroy`, so that order in which the operations are executed is (1) create new resource, (2) update resource dependencies, (3) delete old resource.

Adding the `lifecycle.create_before_destroy` flag fixes the issue in the case above (where we just want to remove the policy definition and update the policy set) and it's a standard practice in resources from other providers. However, it creates another issue. The fact that the `create` operation is executed before `delete` poses a problem when trying to update any of the fields (different from the resource `name`) that force a resource recreation; the `terraform apply` will fail due to a conflict. The reason is that most Azure resources (such as `azurerm_policy_definition`) use the resource name to uniquely identify a resource.

Therefore, if no `lifecycle.create_before_destroy` is set in the policy definition, when a policy definition referenced by an initiative needs to be deleted, this change should be applied in two different terraform apply steps: (1) delete the policy definition reference from the initiative, (2) delete the policy definition.

## Tests

The acceptance tests in the `tests/` directory check whether the defined initiatives actually work: not compliant resources cannot be deployed whereas compliant ones are allowed to be deployed.

The following conventions were followed when testing the policy module:

- Tests check that the initiative is working as expected, as opposed to test individual policies. The reason is that this module only outputs initiatives, all the policies are linked from an initiative.
  
- Only the behavior of custom policies is tested; built-in policies are expected to work.

Tests are implemented using the [Go testing framework](https://golang.org/pkg/testing/) together with the [Terratest module](https://terratest.gruntwork.io/docs/). This configuration allows to call the Terraform configuration from the Go tests.

The way the test are design is as follows:

- Setup: load the policies module to define policies and initiatives and assign initiatives ([`tests/terraform/main.tf`](tests/terraform/main.tf)).
- Run: create compliant resources (for example, [`tests/terraform/resource-location-allow`](tests/terraform/resource-location-allow)) and non-compliant resources (for example, [`tests/terraform/resource-location-audit`](tests/terraform/resource-location-audit)).
- Assert: check whether the policy has been correctly applied using the returned error from the Terraform apply.
- Teardown: delete test resources.
  
Testing effects others than `deny` is particularly slow, since the policy evaluation is quite takes quite a long time, in the order of minutes. In order to speed up the tests, test cases are run in parallel using the [`Parallel()` function](https://golang.org/pkg/testing/#T.Parallel).

It is important to note that the teardown is done separately for the resources provisioned during the setup and for the resources created by each of the test cases. In the first case, the [`Cleanup` function](https://godoc.org/testing#T.Cleanup) is used; `defer` wouldn't work since deferred functions are run before parallel subtests are executed. On the other hand, resources created by the test cases are destroyed in a deferred function.

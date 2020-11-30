package test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

const AZURE_MANAGEMENT_URL = "https://management.azure.com"

func TestPolicies(t *testing.T) {
	policyAssignment := &terraform.Options{
		TerraformDir: "./terraform",
	}
	t.Cleanup(func() { terraform.Destroy(t, policyAssignment) })
	terraform.InitAndApply(t, policyAssignment)

	azurermClient := CreateAzureRMClient(os.Getenv("ARM_TENANT_ID"), os.Getenv("ARM_CLIENT_ID"), os.Getenv("ARM_CLIENT_SECRET"))

	testCases := []struct {
		Name   string
		Effect string
	}{
		{
			Name:   "resource-location",
			Effect: "allow",
		},
		{
			Name:   "resource-location",
			Effect: "audit",
		},
		{
			Name:   "key-vault-security",
			Effect: "allow",
		},
		{
			Name:   "key-vault-security",
			Effect: "deny",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s-%s", tc.Name, tc.Effect), func(t *testing.T) {
			t.Parallel()
			resourceCreation := &terraform.Options{
				TerraformDir: fmt.Sprintf("./terraform/%s-%s", tc.Name, tc.Effect),
			}
			defer terraform.Destroy(t, resourceCreation)
			terraform.Init(t, resourceCreation)

			_, err := terraform.ApplyE(t, resourceCreation)

			if tc.Effect == "allow" && err != nil {
				t.Fatalf("Compliant resources failed to deploy with error: %v", err)
			} else if tc.Effect == "deny" && err == nil {
				t.Fatalf("Policy breach was not prevented")
			} else if tc.Effect == "audit" && err != nil {
				t.Fatalf("Auditable resources failed to deploy with error: %v", err)
			} else if tc.Effect == "audit" && err == nil {
				assertComplianceState(t, resourceCreation, tc.Name, azurermClient, "NonCompliant")
			}
		})
	}
}

func assertComplianceState(t *testing.T, resourceCreation *terraform.Options, testCaseName string, azurermClient *AzureRMClient, expectedState string) {
	resourceId := terraform.Output(t, resourceCreation, "test_resource_id")
	resourceGroupId := strings.Join(strings.Split(resourceId, "/")[:5], "/")
	policyAssignmentName := fmt.Sprintf("test-%s", testCaseName)
	err := azurermClient.TriggerPolicyEvaluation(resourceGroupId)
	if err != nil {
		t.Fatalf("Trigger policy evaluation failed with error: %v", err)
	}
	actualState, err := azurermClient.GetComplianceState(resourceId, policyAssignmentName)
	if err != nil {
		t.Fatalf("Get compliance state failed with error: %v", err)
	}
	if actualState != expectedState {
		t.Fatalf("Expected compliance state %s but got %s", expectedState, actualState)
	}
}

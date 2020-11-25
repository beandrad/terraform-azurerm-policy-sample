package test

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestPolicies(t *testing.T) {
	policyAssignment := &terraform.Options{
		TerraformDir: "./terraform",
	}
	t.Cleanup(func() { terraform.Destroy(t, policyAssignment) })
	terraform.InitAndApply(t, policyAssignment)

	testCases := []struct {
		Name        string
		ShouldAllow bool
	}{
		{
			Name:        "resource-location-allow",
			ShouldAllow: true,
		},
		{
			Name:        "resource-location-deny",
			ShouldAllow: false,
		},
		{
			Name:        "key-vault-allow",
			ShouldAllow: true,
		},
		{
			Name:        "key-vault-deny",
			ShouldAllow: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			resourceCreation := &terraform.Options{
				TerraformDir: fmt.Sprintf("./terraform/%s", tc.Name),
			}
			defer terraform.Destroy(t, resourceCreation)
			terraform.Init(t, resourceCreation)

			_, err := terraform.ApplyE(t, resourceCreation)

			if (tc.ShouldAllow && err != nil) {
				t.Fatalf("Compliant resources failed to deploy with error: %v", err)
			} else if (!tc.ShouldAllow && err == nil) {
				t.Fatalf("Policy breach was not prevented")
			}
		})
	}
}

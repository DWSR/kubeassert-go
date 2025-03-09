package deployments_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/DWSR/kubeassert-go/internal/assertion"
	helpers "github.com/DWSR/kubeassert-go/internal/assertionhelpers"
	"github.com/DWSR/kubeassert-go/internal/deployments"
	"github.com/DWSR/kubeassert-go/internal/testhelpers"
)

const (
	badDeploymentPath = "./testdata/bad-deployment.yaml"
	deploymentPath    = "./testdata/deployment.yaml"
	configPath        = "./testdata/config.yaml"
)

func Test_1Deployment_Success(t *testing.T) {
	asserts := []testhelpers.SuccessfulAssert{
		{
			Name: "Exists_Labels",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				return deployments.NewDeploymentAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": "deployments_test"}),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(deploymentPath),
						helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath),
					),
				).Exists()
			},
		},
		{
			Name: "Exists_Name",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				return deployments.NewDeploymentAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceName("test-deployment"),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(deploymentPath),
						helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath),
					),
				).Exists()
			},
		},
		{
			Name: "IsAvailable_Labels",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				return deployments.NewDeploymentAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": "deployments_test"}),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(deploymentPath),
						helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath),
					),
				).Exists().IsAvailable()
			},
		},
		{
			Name: "IsAvailable_Name",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				return deployments.NewDeploymentAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceName("test-deployment"),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(deploymentPath),
						helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath),
					),
				).Exists().IsAvailable()
			},
		},
		{
			Name: "IsSystemClusterCritical_Labels",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				return deployments.NewDeploymentAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": "deployments_test"}),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(deploymentPath),
						helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath),
					),
				).Exists().IsSystemClusterCritical()
			},
		},
		{
			Name: "IsSystemClusterCritical_Name",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				return deployments.NewDeploymentAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceName("test-deployment"),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(deploymentPath),
						helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath),
					),
				).Exists().IsSystemClusterCritical()
			},
		},
		{
			Name: "HasNoCPULimits_Labels",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				return deployments.NewDeploymentAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": "deployments_test"}),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(deploymentPath),
						helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath),
					),
				).Exists().HasNoCPULimits()
			},
		},
		{
			Name: "HasNoCPULimits_Name",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				return deployments.NewDeploymentAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceName("test-deployment"),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(deploymentPath),
						helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath),
					),
				).Exists().HasNoCPULimits()
			},
		},
		{
			Name: "HasMemoryLimitsEqualToRequests_Labels",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				return deployments.NewDeploymentAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": "deployments_test"}),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(deploymentPath),
						helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath),
					),
				).Exists().HasMemoryLimitsEqualToRequests()
			},
		},
		{
			Name: "HasMemoryLimitsEqualToRequests_Name",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				return deployments.NewDeploymentAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceName("test-deployment"),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(deploymentPath),
						helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath),
					),
				).Exists().HasMemoryLimitsEqualToRequests()
			},
		},
		{
			Name: "HasMemoryLimits_Labels",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				return deployments.NewDeploymentAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": "deployments_test"}),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(deploymentPath),
						helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath),
					),
				).Exists().HasMemoryLimits()
			},
		},
		{
			Name: "HasMemoryLimits_Name",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				return deployments.NewDeploymentAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceName("test-deployment"),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(deploymentPath),
						helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath),
					),
				).Exists().HasMemoryLimits()
			},
		},
		{
			Name: "HasMemoryRequests_Labels",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				return deployments.NewDeploymentAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": "deployments_test"}),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(deploymentPath),
						helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath),
					),
				).Exists().HasMemoryRequests()
			},
		},
		{
			Name: "HasMemoryRequests_Name",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				return deployments.NewDeploymentAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceName("test-deployment"),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(deploymentPath),
						helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath),
					),
				).Exists().HasMemoryRequests()
			},
		},
		{
			Name: "HasCPURequests_Labels",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				return deployments.NewDeploymentAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": "deployments_test"}),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(deploymentPath),
						helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath),
					),
				).Exists().HasCPURequests()
			},
		},
		{
			Name: "HasCPURequests_Name",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				return deployments.NewDeploymentAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceName("test-deployment"),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(deploymentPath),
						helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath),
					),
				).Exists().HasCPURequests()
			},
		},
	}

	testhelpers.TestSuccessfulAsserts(t, testEnv, asserts...)
}

func Test_1Deployment_Fail(t *testing.T) {
	asserts := []testhelpers.FailingAssert{
		{
			Name: "Exists_Labels",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				return deployments.NewDeploymentAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": "deployments_test"}),
				).Exists()
			},
		},
		{
			Name: "Exists_Name",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				return deployments.NewDeploymentAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceName("test-deployment"),
				).Exists()
			},
		},
		{
			Name: "IsAvailable_Labels",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				return deployments.NewDeploymentAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": "deployments_test"}),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(deploymentPath),
					),
				).Exists().IsAvailable()
			},
		},
		{
			Name: "IsAvailable_Name",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				return deployments.NewDeploymentAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceName("test-deployment"),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(deploymentPath),
					),
				).Exists().IsAvailable()
			},
		},
		{
			Name: "IsSystemClusterCritical_Labels",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				return deployments.NewDeploymentAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": "deployments_test"}),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(badDeploymentPath),
					),
				).Exists().IsSystemClusterCritical()
			},
		},
		{
			Name: "IsSystemClusterCritical_Name",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				return deployments.NewDeploymentAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceName("test-deployment"),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(badDeploymentPath),
					),
				).Exists().IsSystemClusterCritical()
			},
		},
		{
			Name: "HasNoCPULimits_Labels",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				return deployments.NewDeploymentAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": "deployments_test"}),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(badDeploymentPath),
					),
				).Exists().HasNoCPULimits()
			},
		},
		{
			Name: "HasNoCPULimits_Name",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				return deployments.NewDeploymentAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceName("test-deployment"),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(badDeploymentPath),
					),
				).Exists().HasNoCPULimits()
			},
		},
		{
			Name: "HasMemoryLimitsEqualToRequests_Labels",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				return deployments.NewDeploymentAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": "deployments_test"}),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(badDeploymentPath),
					),
				).Exists().HasMemoryLimitsEqualToRequests()
			},
		},
		{
			Name: "HasMemoryLimitsEqualToRequests_Name",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				return deployments.NewDeploymentAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceName("test-deployment"),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(badDeploymentPath),
					),
				).Exists().HasMemoryLimitsEqualToRequests()
			},
		},
		{
			Name: "HasMemoryLimits_Labels",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				return deployments.NewDeploymentAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": "deployments_test"}),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(badDeploymentPath),
					),
				).Exists().HasMemoryLimits()
			},
		},
		{
			Name: "HasMemoryLimits_Name",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				return deployments.NewDeploymentAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceName("test-deployment"),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(badDeploymentPath),
					),
				).Exists().HasMemoryLimits()
			},
		},
		{
			Name: "HasCPURequests_Labels",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				return deployments.NewDeploymentAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": "deployments_test"}),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(badDeploymentPath),
					),
				).Exists().HasCPURequests()
			},
		},
		{
			Name: "HasCPURequests_Name",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				return deployments.NewDeploymentAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceName("test-deployment"),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(badDeploymentPath),
					),
				).Exists().HasCPURequests()
			},
		},
	}

	testhelpers.TestFailingAsserts(t, testEnv, asserts...)
}

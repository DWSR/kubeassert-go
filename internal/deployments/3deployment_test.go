package deployments_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	e2etypes "sigs.k8s.io/e2e-framework/pkg/types"

	"github.com/DWSR/kubeassert-go/internal/assertion"
	helpers "github.com/DWSR/kubeassert-go/internal/assertionhelpers"
	"github.com/DWSR/kubeassert-go/internal/deployments"
	"github.com/DWSR/kubeassert-go/internal/testhelpers"
)

func Test_3Deployments_Success(t *testing.T) {
	asserts := []testhelpers.SuccessfulAssert{
		{
			Name: "ExactlyNExist",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				deployNames := generateDeploymentNames()

				return deployments.NewDeploymentAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": deployNames[0]}),
					assertion.WithSetup(createGoodDeploys(deployNames)...),
				).ExactlyNExist(3)
			},
		},
		{
			Name: "AtLeastNExist",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				deployNames := generateDeploymentNames()

				return deployments.NewDeploymentAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": deployNames[0]}),
					assertion.WithSetup(createGoodDeploys(deployNames)...),
				).AtLeastNExist(2)
			},
		},
		{
			Name: "ExactlyNAreAvailable",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				deployNames := generateDeploymentNames()

				return deployments.NewDeploymentAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": deployNames[0]}),
					assertion.WithSetup(
						append([]e2etypes.StepFunc{helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath)},
							createGoodDeploys(deployNames)...)...,
					),
				).ExactlyNExist(3).ExactlyNAreAvailable(3)
			},
		},
		{
			Name: "AtLeastNAreAvailable",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				deployNames := generateDeploymentNames()

				return deployments.NewDeploymentAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": deployNames[0]}),
					assertion.WithSetup(
						append([]e2etypes.StepFunc{helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath)},
							createGoodDeploys(deployNames)...)...,
					),
				).ExactlyNExist(3).AtLeastNAreAvailable(2)
			},
		},
		{
			Name: "ExactlyNAreSystemClusterCritical",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				deployNames := generateDeploymentNames()

				return deployments.NewDeploymentAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": deployNames[0]}),
					assertion.WithSetup(
						append([]e2etypes.StepFunc{helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath)},
							createGoodDeploys(deployNames)...)...,
					),
				).ExactlyNExist(3).ExactlyNAreSystemClusterCritical(3)
			},
		},
		{
			Name: "AtLeastNAreSystemClusterCritical",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				deployNames := generateDeploymentNames()

				return deployments.NewDeploymentAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": deployNames[0]}),
					assertion.WithSetup(
						append([]e2etypes.StepFunc{helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath)},
							createGoodDeploys(deployNames)...)...,
					),
				).ExactlyNExist(3).AtLeastNAreSystemClusterCritical(2)
			},
		},
		{
			Name: "ExactlyNHaveNoCPULimits",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				deployNames := generateDeploymentNames()

				return deployments.NewDeploymentAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": deployNames[0]}),
					assertion.WithSetup(
						append([]e2etypes.StepFunc{helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath)},
							createGoodDeploys(deployNames)...)...,
					),
				).ExactlyNExist(3).ExactlyNHaveNoCPULimits(3)
			},
		},
		{
			Name: "AtLeastNHaveNoCPULimits",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				deployNames := generateDeploymentNames()

				return deployments.NewDeploymentAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": deployNames[0]}),
					assertion.WithSetup(
						append([]e2etypes.StepFunc{helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath)},
							createGoodDeploys(deployNames)...)...,
					),
				).ExactlyNExist(3).AtLeastNHaveNoCPULimits(2)
			},
		},
		{
			Name: "ExactlyNHaveMemoryLimitsEqualToRequests",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				deployNames := generateDeploymentNames()

				return deployments.NewDeploymentAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": deployNames[0]}),
					assertion.WithSetup(
						append([]e2etypes.StepFunc{helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath)},
							createGoodDeploys(deployNames)...)...,
					),
				).ExactlyNExist(3).ExactlyNHaveMemoryLimitsEqualToRequests(3)
			},
		},
		{
			Name: "AtLeastNHaveMemoryLimitsEqualToRequests",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				deployNames := generateDeploymentNames()

				return deployments.NewDeploymentAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": deployNames[0]}),
					assertion.WithSetup(
						append([]e2etypes.StepFunc{helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath)},
							createGoodDeploys(deployNames)...)...,
					),
				).ExactlyNExist(3).AtLeastNHaveMemoryLimitsEqualToRequests(2)
			},
		},
		{
			Name: "ExactlyNHaveMemoryLimits",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				deployNames := generateDeploymentNames()

				return deployments.NewDeploymentAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": deployNames[0]}),
					assertion.WithSetup(
						append([]e2etypes.StepFunc{helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath)},
							createGoodDeploys(deployNames)...)...,
					),
				).ExactlyNExist(3).ExactlyNHaveMemoryLimits(3)
			},
		},
		{
			Name: "AtLeastNHaveMemoryLimits",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				deployNames := generateDeploymentNames()

				return deployments.NewDeploymentAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": deployNames[0]}),
					assertion.WithSetup(
						append([]e2etypes.StepFunc{helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath)},
							createGoodDeploys(deployNames)...)...,
					),
				).ExactlyNExist(3).AtLeastNHaveMemoryLimits(2)
			},
		},
		{
			Name: "ExactlyNHaveMemoryRequests",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				deployNames := generateDeploymentNames()

				return deployments.NewDeploymentAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": deployNames[0]}),
					assertion.WithSetup(
						append([]e2etypes.StepFunc{helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath)},
							createGoodDeploys(deployNames)...)...,
					),
				).ExactlyNExist(3).ExactlyNHaveMemoryRequests(3)
			},
		},
		{
			Name: "AtLeastNHaveMemoryRequests",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				deployNames := generateDeploymentNames()

				return deployments.NewDeploymentAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": deployNames[0]}),
					assertion.WithSetup(
						append([]e2etypes.StepFunc{helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath)},
							createGoodDeploys(deployNames)...)...,
					),
				).ExactlyNExist(3).AtLeastNHaveMemoryRequests(2)
			},
		},
		{
			Name: "ExactlyNHaveCPURequests",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				deployNames := generateDeploymentNames()

				return deployments.NewDeploymentAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": deployNames[0]}),
					assertion.WithSetup(
						append([]e2etypes.StepFunc{helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath)},
							createGoodDeploys(deployNames)...)...,
					),
				).ExactlyNExist(3).ExactlyNHaveCPURequests(3)
			},
		},
		{
			Name: "AtLeastNHaveCPURequests",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				deployNames := generateDeploymentNames()

				return deployments.NewDeploymentAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": deployNames[0]}),
					assertion.WithSetup(
						append([]e2etypes.StepFunc{helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath)},
							createGoodDeploys(deployNames)...)...,
					),
				).ExactlyNExist(3).AtLeastNHaveCPURequests(2)
			},
		},
	}

	testhelpers.TestSuccessfulAsserts(t, testEnv, asserts...)
}

func Test_3Deployments_Fail(t *testing.T) {
	asserts := []testhelpers.FailingAssert{
		{
			Name: "ExactlyNExist",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				deployNames := generateDeploymentNames()

				return deployments.NewDeploymentAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": deployNames[0]}),
					assertion.WithSetup(createBadDeploys(deployNames)...),
				).ExactlyNExist(2)
			},
		},
		{
			Name: "AtLeastNExist",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				deployNames := generateDeploymentNames()

				return deployments.NewDeploymentAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": deployNames[0]}),
					assertion.WithSetup(createBadDeploys(deployNames)...),
				).AtLeastNExist(4)
			},
		},
		{
			Name: "ExactlyNAreAvailable",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				deployNames := generateDeploymentNames()

				return deployments.NewDeploymentAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(15*time.Second),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": deployNames[0]}),
					assertion.WithSetup(
						append([]e2etypes.StepFunc{helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath)},
							createGoodDeploys(deployNames)...)...,
					),
				).ExactlyNExist(3).ExactlyNAreAvailable(2)
			},
		},
		{
			Name: "AtLeastNAreAvailable",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				deployNames := generateDeploymentNames()

				return deployments.NewDeploymentAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(15*time.Second),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": deployNames[0]}),
					assertion.WithSetup(createGoodDeploys(deployNames)...),
				).ExactlyNExist(3).AtLeastNAreAvailable(4)
			},
		},
		{
			Name: "ExactlyNAreSystemClusterCritical",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				deployNames := generateDeploymentNames()

				return deployments.NewDeploymentAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": deployNames[0]}),
					assertion.WithSetup(createGoodDeploys(deployNames)...),
				).ExactlyNExist(3).ExactlyNAreSystemClusterCritical(2)
			},
		},
		{
			Name: "AtLeastNAreSystemClusterCritical",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				deployNames := generateDeploymentNames()

				return deployments.NewDeploymentAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": deployNames[0]}),
					assertion.WithSetup(createGoodDeploys(deployNames)...),
				).ExactlyNExist(3).AtLeastNAreSystemClusterCritical(4)
			},
		},
		{
			Name: "ExactlyNHaveNoCPULimits",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				deployNames := generateDeploymentNames()

				return deployments.NewDeploymentAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": deployNames[0]}),
					assertion.WithSetup(createGoodDeploys(deployNames)...),
				).ExactlyNExist(3).ExactlyNHaveNoCPULimits(2)
			},
		},
		{
			Name: "AtLeastNHaveNoCPULimits",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				deployNames := generateDeploymentNames()

				return deployments.NewDeploymentAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": deployNames[0]}),
					assertion.WithSetup(createGoodDeploys(deployNames)...),
				).ExactlyNExist(3).AtLeastNHaveNoCPULimits(4)
			},
		},
		{
			Name: "ExactlyNHaveMemoryLimitsEqualToRequests",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				deployNames := generateDeploymentNames()

				return deployments.NewDeploymentAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": deployNames[0]}),
					assertion.WithSetup(createGoodDeploys(deployNames)...),
				).ExactlyNExist(3).ExactlyNHaveMemoryLimitsEqualToRequests(2)
			},
		},
		{
			Name: "AtLeastNHaveMemoryLimitsEqualToRequests",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				deployNames := generateDeploymentNames()

				return deployments.NewDeploymentAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": deployNames[0]}),
					assertion.WithSetup(createGoodDeploys(deployNames)...),
				).ExactlyNExist(3).AtLeastNHaveMemoryLimitsEqualToRequests(4)
			},
		},
		{
			Name: "ExactlyNHaveMemoryLimits",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				deployNames := generateDeploymentNames()

				return deployments.NewDeploymentAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": deployNames[0]}),
					assertion.WithSetup(createGoodDeploys(deployNames)...),
				).ExactlyNExist(3).ExactlyNHaveMemoryLimits(2)
			},
		},
		{
			Name: "AtLeastNHaveMemoryLimits",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				deployNames := generateDeploymentNames()

				return deployments.NewDeploymentAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": deployNames[0]}),
					assertion.WithSetup(createGoodDeploys(deployNames)...),
				).ExactlyNExist(3).AtLeastNHaveMemoryLimits(4)
			},
		},
		{
			Name: "ExactlyNHaveMemoryRequests",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				deployNames := generateDeploymentNames()

				return deployments.NewDeploymentAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": deployNames[0]}),
					assertion.WithSetup(createGoodDeploys(deployNames)...),
				).ExactlyNExist(3).ExactlyNHaveMemoryRequests(2)
			},
		},
		{
			Name: "AtLeastNHaveMemoryRequests",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				deployNames := generateDeploymentNames()

				return deployments.NewDeploymentAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": deployNames[0]}),
					assertion.WithSetup(createGoodDeploys(deployNames)...),
				).ExactlyNExist(3).AtLeastNHaveMemoryRequests(4)
			},
		},
		{
			Name: "ExactlyNHaveCPURequests",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				deployNames := generateDeploymentNames()

				return deployments.NewDeploymentAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": deployNames[0]}),
					assertion.WithSetup(createGoodDeploys(deployNames)...),
				).ExactlyNExist(3).ExactlyNHaveCPURequests(2)
			},
		},
		{
			Name: "AtLeastNHaveCPURequests",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				deployNames := generateDeploymentNames()

				return deployments.NewDeploymentAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": deployNames[0]}),
					assertion.WithSetup(createGoodDeploys(deployNames)...),
				).ExactlyNExist(3).AtLeastNHaveCPURequests(4)
			},
		},
	}

	testhelpers.TestFailingAsserts(t, testEnv, asserts...)
}

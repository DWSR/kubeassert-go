package deployments

import (
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/DWSR/kubeassert-go/internal/assertion"
	helpers "github.com/DWSR/kubeassert-go/internal/assertionhelpers"
)

// DeploymentAssertion is a wrapper around assertion.Assertion that provides a set of assertion functions for
// Deployments.
type DeploymentAssertion struct {
	assertion.Assertion
}

func (da DeploymentAssertion) clone() DeploymentAssertion {
	return DeploymentAssertion{
		Assertion: assertion.Clone(da.Assertion),
	}
}

// Exists asserts that exactly one Deployment exists in the cluster that matches the provided options.
func (da DeploymentAssertion) Exists() DeploymentAssertion {
	return da.ExactlyNExist(1)
}

// ExactlyNExist asserts that exactly N Deployments exist in the cluster that match the provided options.
func (da DeploymentAssertion) ExactlyNExist(count int) DeploymentAssertion {
	stepFn := helpers.AsStepFunc(da, exist(), count, helpers.IntCompareFuncEqualTo, nil)

	res := da.clone()
	res.SetBuilder(res.GetBuilder().Assess("exactlyNExist", stepFn))

	return res
}

// AtLeastNExist asserts that at least N Deployments exist in the cluster that match the provided options.
func (da DeploymentAssertion) AtLeastNExist(count int) DeploymentAssertion {
	stepFn := helpers.AsStepFunc(da, exist(), count, helpers.IntCompareFuncGreaterThanOrEqualTo, nil)

	res := da.clone()
	res.SetBuilder(res.GetBuilder().Assess("atLeastNExist", stepFn))

	return res
}

// IsAvailable asserts that exactly one Deployment is available in the cluster that matches the provided options.
func (da DeploymentAssertion) IsAvailable() DeploymentAssertion {
	return da.ExactlyNAreAvailable(1)
}

// ExactlyNAreAvailable asserts that exactly N Deployments are available in the cluster that match the provided options.
func (da DeploymentAssertion) ExactlyNAreAvailable(count int) DeploymentAssertion {
	stepFn := helpers.AsStepFunc(
		da,
		areAvailable(),
		count,
		helpers.IntCompareFuncNotEqualTo,
		helpers.IntCompareFuncEqualTo,
	)

	res := da.clone()
	res.SetBuilder(res.GetBuilder().Assess("exactlyNAreAvailable", stepFn))

	return res
}

// AtLeastNAreAvailable asserts that at least N Deployments are available in the cluster that match the provided
// options.
func (da DeploymentAssertion) AtLeastNAreAvailable(count int) DeploymentAssertion {
	stepFn := helpers.AsStepFunc(
		da,
		areAvailable(),
		count,
		helpers.IntCompareFuncLessThan,
		helpers.IntCompareFuncGreaterThanOrEqualTo,
	)

	res := da.clone()
	res.SetBuilder(res.GetBuilder().Assess("atLeastNAreAvailable", stepFn))

	return res
}

// IsSystemClusterCritical asserts that exactly one Deployment is system cluster critical in the cluster that matches
// the provided options.
func (da DeploymentAssertion) IsSystemClusterCritical() DeploymentAssertion {
	return da.ExactlyNAreSystemClusterCritical(1)
}

// ExactlyNAreSystemClusterCritical asserts that exactly N Deployments are system cluster critical in the cluster that
// match the provided options.
func (da DeploymentAssertion) ExactlyNAreSystemClusterCritical(count int) DeploymentAssertion {
	stepFn := helpers.AsStepFunc(
		da,
		areSystemClusterCritical(),
		count,
		helpers.IntCompareFuncNotEqualTo,
		helpers.IntCompareFuncEqualTo,
	)

	res := da.clone()
	res.SetBuilder(res.GetBuilder().Assess("exactlyNAreSystemClusterCritical", stepFn))

	return res
}

// AtLeastNAreSystemClusterCritical asserts that at least N Deployments are system cluster critical in the cluster that
// match the provided options.
func (da DeploymentAssertion) AtLeastNAreSystemClusterCritical(count int) DeploymentAssertion {
	stepFn := helpers.AsStepFunc(
		da,
		areSystemClusterCritical(),
		count,
		helpers.IntCompareFuncLessThan,
		helpers.IntCompareFuncGreaterThanOrEqualTo,
	)

	res := da.clone()
	res.SetBuilder(res.GetBuilder().Assess("atLeastNAreSystemClusterCritical", stepFn))

	return res
}

// HasNoCPULimits asserts that exactly one Deployment has no CPU limits in the cluster that match the provided options.
func (da DeploymentAssertion) HasNoCPULimits() DeploymentAssertion {
	return da.ExactlyNHaveNoCPULimits(1)
}

// ExactlyNHaveNoCPULimits asserts that exactly N Deployments have no CPU limits in the cluster that match the provided
// options.
func (da DeploymentAssertion) ExactlyNHaveNoCPULimits(count int) DeploymentAssertion {
	stepFn := helpers.AsStepFunc(
		da,
		haveNoCPULimits(),
		count,
		helpers.IntCompareFuncNotEqualTo,
		helpers.IntCompareFuncEqualTo,
	)

	res := da.clone()
	res.SetBuilder(res.GetBuilder().Assess("exactlyNHaveNoCPULimits", stepFn))

	return res
}

// AtLeastNHaveNoCPULimits asserts that at least N Deployments have no CPU limits in the cluster that match the provided
// options.
func (da DeploymentAssertion) AtLeastNHaveNoCPULimits(count int) DeploymentAssertion {
	stepFn := helpers.AsStepFunc(
		da,
		haveNoCPULimits(),
		count,
		helpers.IntCompareFuncLessThan,
		helpers.IntCompareFuncGreaterThanOrEqualTo,
	)

	res := da.clone()
	res.SetBuilder(res.GetBuilder().Assess("atLeastNHaveNoCPULimits", stepFn))

	return res
}

// HasMemoryLimitsEqualToRequests asserts that exactly one Deployment has memory limits set equal to requests in the
// cluster that match the provided options.
func (da DeploymentAssertion) HasMemoryLimitsEqualToRequests() DeploymentAssertion {
	return da.ExactlyNHaveMemoryLimitsEqualToRequests(1)
}

// ExactlyNHaveMemoryLimitsEqualToRequests asserts that exactly N Deployments have memory limits set equal to requests
// in the cluster that match the provided options.
func (da DeploymentAssertion) ExactlyNHaveMemoryLimitsEqualToRequests(count int) DeploymentAssertion {
	stepFn := helpers.AsStepFunc(
		da,
		haveMemoryLimitsEqualToRequests(),
		count,
		helpers.IntCompareFuncNotEqualTo,
		helpers.IntCompareFuncEqualTo,
	)

	res := da.clone()
	res.SetBuilder(res.GetBuilder().Assess("exactlyNHaveMemoryLimitsEqualToRequests", stepFn))

	return res
}

// AtLeastNHaveMemoryLimitsEqualToRequests asserts that at least N Deployments have memory limits set equal to requests
// in the cluster that match the provided options.
func (da DeploymentAssertion) AtLeastNHaveMemoryLimitsEqualToRequests(count int) DeploymentAssertion {
	stepFn := helpers.AsStepFunc(
		da,
		haveMemoryLimitsEqualToRequests(),
		count,
		helpers.IntCompareFuncLessThan,
		helpers.IntCompareFuncGreaterThanOrEqualTo,
	)

	res := da.clone()
	res.SetBuilder(res.GetBuilder().Assess("atLeastNHaveMemoryLimitsEqualToRequests", stepFn))

	return res
}

// HasMemoryLimits asserts that exactly one Deployment has memory limits in the cluster that match the provided options.
func (da DeploymentAssertion) HasMemoryLimits() DeploymentAssertion {
	return da.ExactlyNHaveMemoryLimits(1)
}

// ExactlyNHaveMemoryLimits asserts that exactly N Deployments have memory limits in the cluster that match the provided
// options.
func (da DeploymentAssertion) ExactlyNHaveMemoryLimits(count int) DeploymentAssertion {
	stepFn := helpers.AsStepFunc(
		da,
		haveMemoryLimits(),
		count,
		helpers.IntCompareFuncNotEqualTo,
		helpers.IntCompareFuncEqualTo,
	)

	res := da.clone()
	res.SetBuilder(res.GetBuilder().Assess("exactlyNHaveMemoryLimits", stepFn))

	return res
}

// AtLeastNHaveMemoryLimits asserts that at least N Deployments have memory limits in the cluster that match the
// provided options.
func (da DeploymentAssertion) AtLeastNHaveMemoryLimits(count int) DeploymentAssertion {
	stepFn := helpers.AsStepFunc(
		da,
		haveMemoryLimits(),
		count,
		helpers.IntCompareFuncLessThan,
		helpers.IntCompareFuncGreaterThanOrEqualTo,
	)

	res := da.clone()
	res.SetBuilder(res.GetBuilder().Assess("atLeastNHaveMemoryLimits", stepFn))

	return res
}

// HasMemoryRequests asserts that exactly one Deployment has memory requests in the cluster that match the provided
// options.
func (da DeploymentAssertion) HasMemoryRequests() DeploymentAssertion {
	return da.ExactlyNHaveMemoryRequests(1)
}

// ExactlyNHaveMemoryRequests asserts that exactly N Deployments have memory requests in the cluster that match the
// provided options.
func (da DeploymentAssertion) ExactlyNHaveMemoryRequests(count int) DeploymentAssertion {
	stepFn := helpers.AsStepFunc(
		da,
		haveMemoryRequests(),
		count,
		helpers.IntCompareFuncNotEqualTo,
		helpers.IntCompareFuncEqualTo,
	)

	res := da.clone()
	res.SetBuilder(res.GetBuilder().Assess("exactlyNHaveMemoryRequests", stepFn))

	return res
}

// AtLeastNHaveMemoryRequests asserts that at least N Deployments have memory requests in the cluster that match the
// provided options.
func (da DeploymentAssertion) AtLeastNHaveMemoryRequests(count int) DeploymentAssertion {
	stepFn := helpers.AsStepFunc(
		da,
		haveMemoryRequests(),
		count,
		helpers.IntCompareFuncLessThan,
		helpers.IntCompareFuncGreaterThanOrEqualTo,
	)

	res := da.clone()
	res.SetBuilder(res.GetBuilder().Assess("atLeastNHaveMemoryRequests", stepFn))

	return res
}

// HasCPURequests asserts that exactly one Deployment has CPU requests in the cluster that match the provided options.
func (da DeploymentAssertion) HasCPURequests() DeploymentAssertion {
	return da.ExactlyNHaveCPURequests(1)
}

// ExactlyNHaveCPURequests asserts that exactly N Deployments have CPU requests in the cluster that match the provided
// options.
func (da DeploymentAssertion) ExactlyNHaveCPURequests(count int) DeploymentAssertion {
	stepFn := helpers.AsStepFunc(
		da,
		haveCPURequests(),
		count,
		helpers.IntCompareFuncNotEqualTo,
		helpers.IntCompareFuncEqualTo,
	)

	res := da.clone()
	res.SetBuilder(res.GetBuilder().Assess("exactlyNHaveCPURequests", stepFn))

	return res
}

// AtLeastNHaveCPURequests asserts that at least N Deployments have CPU requests in the cluster that match the provided
// options.
func (da DeploymentAssertion) AtLeastNHaveCPURequests(count int) DeploymentAssertion {
	stepFn := helpers.AsStepFunc(
		da,
		haveCPURequests(),
		count,
		helpers.IntCompareFuncLessThan,
		helpers.IntCompareFuncGreaterThanOrEqualTo,
	)

	res := da.clone()
	res.SetBuilder(res.GetBuilder().Assess("atLeastNHaveCPURequests", stepFn))

	return res
}

// NewDeploymentAssertion creates a new DeploymentAssertion with the provided options.
func NewDeploymentAssertion(opts ...assertion.Option) DeploymentAssertion {
	return DeploymentAssertion{
		Assertion: assertion.NewAssertion(
			append(
				[]assertion.Option{assertion.WithBuilder(features.New("Deployment").WithLabel("type", "deployment"))},
				opts...,
			)...,
		),
	}
}

package pods

import (
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/DWSR/kubeassert-go/internal/assertion"
	helpers "github.com/DWSR/kubeassert-go/internal/assertionhelpers"
)

// PodAssertion is a wrapper around the assertion.Assertion type and provides a set of assertions for Kubernetes Pods.
type PodAssertion struct {
	assertion.Assertion
}

func (pa PodAssertion) clone() PodAssertion {
	return PodAssertion{
		Assertion: assertion.Clone(pa.Assertion),
	}
}

// Exists asserts that exactly one Pod exists in the cluster that matches the provided options.
func (pa PodAssertion) Exists() PodAssertion {
	return pa.ExactlyNExist(1)
}

// ExactlyNExist asserts that exactly N Pods exist in the cluster that match the provided options.
func (pa PodAssertion) ExactlyNExist(count int) PodAssertion {
	stepFn := helpers.AsStepFunc(pa, exist(), count, helpers.IntCompareFuncEqualTo, nil)

	res := pa.clone()
	res.SetBuilder(res.GetBuilder().Assess("exactlyNExist", stepFn))

	return res
}

// AtLeastNExist asserts that at least N Pods exist in the cluster that match the provided options.
func (pa PodAssertion) AtLeastNExist(count int) PodAssertion {
	stepFn := helpers.AsStepFunc(pa, exist(), count, helpers.IntCompareFuncGreaterThanOrEqualTo, nil)

	res := pa.clone()
	res.SetBuilder(res.GetBuilder().Assess("atLeastNExist", stepFn))

	return res
}

// IsReady asserts that exactly one Pod is ready in the cluster that matches the provided options.
func (pa PodAssertion) IsReady() PodAssertion {
	return pa.ExactlyNAreReady(1)
}

// ExactlyNAreReady asserts that exactly N Pods are ready in the cluster that match the provided options.
func (pa PodAssertion) ExactlyNAreReady(count int) PodAssertion {
	stepFn := helpers.AsStepFunc(pa, areReady(), count, helpers.IntCompareFuncNotEqualTo, helpers.IntCompareFuncEqualTo)
	res := pa.clone()
	res.SetBuilder(pa.GetBuilder().Assess("exactlyNAreReady", stepFn))

	return res
}

// AtLeastNAreReady asserts that at least N Pods are ready in the cluster that match the provided options.
func (pa PodAssertion) AtLeastNAreReady(count int) PodAssertion {
	stepFn := helpers.AsStepFunc(
		pa,
		areReady(),
		count,
		helpers.IntCompareFuncLessThan,
		helpers.IntCompareFuncGreaterThanOrEqualTo,
	)
	res := pa.clone()
	res.SetBuilder(pa.GetBuilder().Assess("atLeastNAreReady", stepFn))

	return res
}

// NewPodAssertion creates a new PodAssertion with the provided options.
func NewPodAssertion(opts ...assertion.Option) PodAssertion {
	return PodAssertion{
		Assertion: assertion.NewAssertion(
			append(
				[]assertion.Option{assertion.WithBuilder(features.New("Pod").WithLabel("type", "pod"))},
				opts...,
			)...,
		),
	}
}

package secrets

import (
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/DWSR/kubeassert-go/internal/assertion"
	helpers "github.com/DWSR/kubeassert-go/internal/assertionhelpers"
)

// SecretAssertion is a wrapper around the assertion.Assertion type and provides a set of assertions for Kubernetes
// Secrets.
type SecretAssertion struct {
	assertion.Assertion
}

func (sa SecretAssertion) clone() SecretAssertion {
	return SecretAssertion{
		Assertion: assertion.Clone(sa.Assertion),
	}
}

// Exists asserts that exactly one Secret exists in the cluster that matches the provided options.
func (sa SecretAssertion) Exists() SecretAssertion {
	return sa.ExactlyNExist(1)
}

// ExactlyNExist asserts that exactly N Secrets exist in the cluster that match the provided options.
func (sa SecretAssertion) ExactlyNExist(count int) SecretAssertion {
	stepFn := helpers.AsStepFunc(sa, exist(), count, helpers.IntCompareFuncEqualTo, nil)

	res := sa.clone()
	res.SetBuilder(res.GetBuilder().Assess("exactlyNExist", stepFn))

	return res
}

// AtLeastNExist asserts that at least N Secrets exist in the cluster that match the provided options.
func (sa SecretAssertion) AtLeastNExist(count int) SecretAssertion {
	stepFn := helpers.AsStepFunc(sa, exist(), count, helpers.IntCompareFuncGreaterThanOrEqualTo, nil)

	res := sa.clone()
	res.SetBuilder(res.GetBuilder().Assess("atLeastNExist", stepFn))

	return res
}

// HasContent asserts that exactly one Secret in the cluster contains the provided content. This match is not exclusive
// meaning that the Secret can contain additional content.
func (sa SecretAssertion) HasContent(content map[string]string) SecretAssertion {
	return sa.ExactlyNHaveContent(1, content)
}

// ExactlyNHaveContent asserts that exactly N Secrets in the cluster contain the provided content. This match is not
// exclusive meaning that the Secrets can contain additional content.
func (sa SecretAssertion) ExactlyNHaveContent(count int, content map[string]string) SecretAssertion {
	stepFn := helpers.AsStepFunc(
		sa,
		haveContent(content),
		count,
		helpers.IntCompareFuncNotEqualTo,
		helpers.IntCompareFuncEqualTo,
	)
	res := sa.clone()
	res.SetBuilder(res.GetBuilder().Assess("exactlyNHaveContent", stepFn))

	return res
}

// AtLeastNHaveContent asserts that at least N Secrets in the cluster contain the provided content. This match is not
// exclusive meaning that the Secrets can contain additional content.
func (sa SecretAssertion) AtLeastNHaveContent(count int, content map[string]string) SecretAssertion {
	stepFn := helpers.AsStepFunc(
		sa,
		haveContent(content),
		count,
		helpers.IntCompareFuncLessThan,
		helpers.IntCompareFuncGreaterThanOrEqualTo,
	)

	res := sa.clone()
	res.SetBuilder(res.GetBuilder().Assess("atLeastNHaveContent", stepFn))

	return res
}

// NewSecretAssertion creates a new SecretAssertion with the provided options.
func NewSecretAssertion(opts ...assertion.Option) SecretAssertion {
	return SecretAssertion{
		Assertion: assertion.NewAssertion(
			append(
				[]assertion.Option{assertion.WithBuilder(features.New("Secret").WithLabel("type", "secret"))},
				opts...,
			)...,
		),
	}
}

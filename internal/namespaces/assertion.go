package namespaces

import (
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/DWSR/kubeassert-go/internal/assertion"
	helpers "github.com/DWSR/kubeassert-go/internal/assertionhelpers"
)

// NamespaceAssertion is an assertion for Kubernetes Namespaces.
type NamespaceAssertion struct {
	assertion.Assertion
}

const (
	podSecurityEnforceLabelKey = "pod-security.kubernetes.io/enforce"
)

func (na NamespaceAssertion) clone() NamespaceAssertion {
	return NamespaceAssertion{
		Assertion: assertion.Clone(na.Assertion),
	}
}

// Exists asserts that exactly one Namespace exists.
func (na NamespaceAssertion) Exists() NamespaceAssertion {
	return na.ExactlyNExist(1)
}

// ExactlyNExist asserts that exactly N Namespaces exist.
func (na NamespaceAssertion) ExactlyNExist(count int) NamespaceAssertion {
	stepFn := helpers.AsStepFunc(na, exist(), count, helpers.IntCompareFuncEqualTo, nil)

	res := na.clone()
	res.SetBuilder(res.GetBuilder().Assess("exactlyNExist", stepFn))

	return res
}

// AtLeastNExist asserts that at least N Namespaces exist.
func (na NamespaceAssertion) AtLeastNExist(count int) NamespaceAssertion {
	stepFn := helpers.AsStepFunc(na, exist(), count, helpers.IntCompareFuncGreaterThanOrEqualTo, nil)

	res := na.clone()
	res.SetBuilder(res.GetBuilder().Assess("atLeastNExist", stepFn))

	return res
}

// IsRestricted asserts that exactly one Namespace uses the default "restricted" pod security standard.
func (na NamespaceAssertion) IsRestricted() NamespaceAssertion {
	return na.ExactlyNAreRestricted(1)
}

// ExactlyNAreRestricted asserts that exactly N Namespaces use the default "restricted" pod security standard.
func (na NamespaceAssertion) ExactlyNAreRestricted(count int) NamespaceAssertion {
	stepFn := helpers.AsStepFunc(
		na,
		areRestricted(),
		count,
		helpers.IntCompareFuncNotEqualTo,
		helpers.IntCompareFuncEqualTo,
	)
	res := na.clone()
	res.SetBuilder(na.GetBuilder().Assess("exactlyNAreRestricted", stepFn))

	return res
}

// AtLeastNAreRestricted asserts that at least N Namespaces use the default "restricted" pod security standard.
func (na NamespaceAssertion) AtLeastNAreRestricted(count int) NamespaceAssertion {
	stepFn := helpers.AsStepFunc(
		na,
		areRestricted(),
		count,
		helpers.IntCompareFuncLessThan,
		helpers.IntCompareFuncGreaterThanOrEqualTo,
	)
	res := na.clone()
	res.SetBuilder(na.GetBuilder().Assess("atLeastNAreRestricted", stepFn))

	return res
}

// NewNamespaceAssertion creates a new NamespaceAssertion.
func NewNamespaceAssertion(opts ...assertion.Option) NamespaceAssertion {
	return NamespaceAssertion{
		Assertion: assertion.NewAssertion(
			append(
				[]assertion.Option{assertion.WithBuilder(features.New("Namespace").WithLabel("type", "namespace"))},
				opts...,
			)...,
		),
	}
}

package crds

import (
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/DWSR/kubeassert-go/internal/assertion"
	helpers "github.com/DWSR/kubeassert-go/internal/assertionhelpers"
)

// CRDAssertion is an assertion for CustomResourceDefinitions.
type CRDAssertion struct {
	assertion.Assertion
}

func (ca CRDAssertion) clone() CRDAssertion {
	return CRDAssertion{
		Assertion: assertion.Clone(ca.Assertion),
	}
}

// Exists asserts that exactly one CRD exists that matches the provided options.
func (ca CRDAssertion) Exists() CRDAssertion {
	stepFn := helpers.AsStepFunc(ca, exist(), 1, helpers.IntCompareFuncEqualTo, nil)

	res := ca.clone()
	res.SetBuilder(res.GetBuilder().Assess("exists", stepFn))

	return res
}

// HasVersion asserts that exactly one CRD that matches the supplied options has the supplied version.
func (ca CRDAssertion) HasVersion(crdVersion string) CRDAssertion {
	stepFn := helpers.AsStepFunc(
		ca,
		hasVersion(crdVersion),
		1,
		helpers.IntCompareFuncNotEqualTo,
		helpers.IntCompareFuncEqualTo,
	)
	res := ca.clone()
	res.SetBuilder(res.GetBuilder().Assess("hasVersion", stepFn))

	return res
}

// NewCRDAssertion creates a new CRDAssertion with the supplied options.
func NewCRDAssertion(opts ...assertion.Option) CRDAssertion {
	return CRDAssertion{
		Assertion: assertion.NewAssertion(
			append(
				[]assertion.Option{assertion.WithBuilder(features.New("CRD").WithLabel("type", "customresourcedefinition"))},
				opts...,
			)...,
		),
	}
}

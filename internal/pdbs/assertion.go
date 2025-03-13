package pdbs

import (
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/DWSR/kubeassert-go/internal/assertion"
	helpers "github.com/DWSR/kubeassert-go/internal/assertionhelpers"
)

// PDBAssertion is a wrapper around assertion.Assertion that provides additional functionality for PodDisruptionBudgets.
type PDBAssertion struct {
	assertion.Assertion
}

func (pa PDBAssertion) clone() PDBAssertion {
	return PDBAssertion{
		Assertion: assertion.Clone(pa.Assertion),
	}
}

// Exists asserts that exactly one PodDisruptionBudget exists in the cluster that matches the provided options.
func (pa PDBAssertion) Exists() PDBAssertion {
	stepFn := helpers.AsStepFunc(pa, exist(), 1, helpers.IntCompareFuncEqualTo, nil)

	res := pa.clone()
	res.SetBuilder(res.GetBuilder().Assess("exists", stepFn))

	return res
}

// NewPDBAssertion creates a new PDBAssertion with the provided options.
func NewPDBAssertion(opts ...assertion.Option) PDBAssertion {
	return PDBAssertion{
		Assertion: assertion.NewAssertion(
			append(
				[]assertion.Option{assertion.WithBuilder(features.New("PDB").WithLabel("type", "poddisruptionbudget"))},
				opts...,
			)...,
		),
	}
}

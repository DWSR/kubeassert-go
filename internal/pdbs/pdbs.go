package pdbs

import (
	"context"
	"testing"

	"github.com/DWSR/kubeassert-go/internal/assertion"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	policyv1 "k8s.io/api/policy/v1"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
	e2etypes "sigs.k8s.io/e2e-framework/pkg/types"
)

type PDBAssertion struct {
	assertion.Assertion
}

func (pa PDBAssertion) clone() PDBAssertion {
	return PDBAssertion{
		Assertion: assertion.CloneAssertion(pa.Assertion),
	}
}

func PodDisruptionBudgetExists(namespaceName, pdbName string) e2etypes.Feature {
	return features.New("PodDisruptionBudgetExists").
		WithLabel("type", "pdb").
		AssessWithDescription(
			"pdbExists",
			"PDB should exist",
			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				var pdb policyv1.PodDisruptionBudget

				err := cfg.Client().
					Resources("poddisruptionbudgets").
					WithNamespace(namespaceName).
					Get(ctx, pdbName, namespaceName, &pdb)
				require.NoError(t, err)

				return ctx
			}).
		Feature()
}

func PodDisruptionBudgetTargetsDeployment(namespaceName, pdbName, deployName string) e2etypes.Feature {
	return features.New("PodDisruptionBudgetTargetsDeployment").
		WithLabel("type", "pdb").
		AssessWithDescription(
			"pdbTargetsDeployment",
			"PDB should target deployment",
			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				var pdb policyv1.PodDisruptionBudget
				var deploy appsv1.Deployment

				err := cfg.Client().
					Resources("poddisruptionbudgets").
					WithNamespace(namespaceName).
					Get(ctx, pdbName, namespaceName, &pdb)
				require.NoError(t, err)

				err = cfg.Client().
					Resources("deployments").
					WithNamespace(namespaceName).
					Get(ctx, deployName, namespaceName, &deploy)
				require.NoError(t, err)

				for labelKey, labelValue := range pdb.Spec.Selector.MatchLabels {
					require.Equal(t, deploy.Spec.Selector.MatchLabels, labelKey)
					require.Equal(t, deploy.Spec.Selector.MatchLabels[labelKey], labelValue)
				}

				return ctx
			}).
		Feature()
}

func NewPDBAssertion(opts ...assertion.AssertionOption) PDBAssertion {
	return PDBAssertion{
		Assertion: assertion.NewAssertion(
			append(
				[]assertion.AssertionOption{assertion.WithBuilder(features.New("CRD").WithLabel("type", "poddisruptionbudget"))},
				opts...,
			)...,
		),
	}
}

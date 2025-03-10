package pdbs

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	policyv1 "k8s.io/api/policy/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/DWSR/kubeassert-go/internal/assertion"
	helpers "github.com/DWSR/kubeassert-go/internal/assertionhelpers"
)

type PDBAssertion struct {
	assertion.Assertion
}

func (pa PDBAssertion) clone() PDBAssertion {
	return PDBAssertion{
		Assertion: assertion.Clone(pa.Assertion),
	}
}

func (pa PDBAssertion) getPDBs(
	ctx context.Context,
	t require.TestingT,
	cfg *envconf.Config,
) (policyv1.PodDisruptionBudgetList, error) {
	client := helpers.DynamicClientFromEnvconf(t, cfg)

	var pdbList policyv1.PodDisruptionBudgetList

	list, err := client.Resource(policyv1.SchemeGroupVersion.WithResource("poddisruptionbudgets")).
		List(ctx, pa.ListOptions(cfg))
	if err != nil {
		return pdbList, err
	}

	err = runtime.DefaultUnstructuredConverter.FromUnstructured(list.UnstructuredContent(), &pdbList)
	if err != nil {
		return pdbList, err
	}

	return pdbList, nil
}

func (pa PDBAssertion) Exists() PDBAssertion {
	stepFn := func(ctx context.Context, testingT *testing.T, cfg *envconf.Config) context.Context {
		t := helpers.RequireTIfNotNil(testingT, pa.GetRequireT())

		conditionFunc := func(ctx context.Context) (bool, error) {
			pdbs, err := pa.getPDBs(ctx, t, cfg)
			require.NoError(t, err)

			return len(pdbs.Items) == 1, nil
		}

		require.NoError(t, helpers.WaitForCondition(ctx, pa, conditionFunc))

		return ctx
	}
	res := pa.clone()
	res.SetBuilder(res.GetBuilder().Assess("exists", stepFn))

	return res
}

func NewPDBAssertion(opts ...assertion.Option) PDBAssertion {
	return PDBAssertion{
		Assertion: assertion.NewAssertion(
			append(
				[]assertion.Option{assertion.WithBuilder(features.New("CRD").WithLabel("type", "poddisruptionbudget"))},
				opts...,
			)...,
		),
	}
}

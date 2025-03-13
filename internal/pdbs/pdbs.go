// pdbs contains assertions for Kubernetes PodDisruptionBudgets.
package pdbs

import (
	"context"

	"github.com/stretchr/testify/require"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/e2e-framework/pkg/envconf"

	"github.com/DWSR/kubeassert-go/internal/assertion"
	helpers "github.com/DWSR/kubeassert-go/internal/assertionhelpers"
)

func getPDBs(
	ctx context.Context,
	t require.TestingT,
	cfg *envconf.Config,
	listOpts metav1.ListOptions,
) (policyv1.PodDisruptionBudgetList, error) {
	client := helpers.DynamicClientFromEnvconf(t, cfg)

	var pdbList policyv1.PodDisruptionBudgetList

	list, err := client.Resource(policyv1.SchemeGroupVersion.WithResource("poddisruptionbudgets")).
		List(ctx, listOpts)
	if err != nil {
		return pdbList, err
	}

	err = runtime.DefaultUnstructuredConverter.FromUnstructured(list.UnstructuredContent(), &pdbList)
	if err != nil {
		return pdbList, err
	}

	return pdbList, nil
}

func exist() helpers.ConditionFuncFactory {
	return func(
		t require.TestingT,
		assert assertion.Assertion,
		cfg *envconf.Config,
		count int,
		itemCountFn, _ helpers.IntCompareFunc,
	) helpers.ConditionFunc {
		return func(ctx context.Context) (bool, error) {
			pdbs, err := getPDBs(ctx, t, cfg, assert.ListOptions(cfg))
			require.NoError(t, err)

			return itemCountFn(len(pdbs.Items), count), nil
		}
	}
}

// namespace contains assertions for Kubernetes Namespaces.
package namespaces

import (
	"context"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/e2e-framework/pkg/envconf"

	"github.com/DWSR/kubeassert-go/internal/assertion"
	helpers "github.com/DWSR/kubeassert-go/internal/assertionhelpers"
)

func getNamespaces(
	ctx context.Context,
	t require.TestingT,
	cfg *envconf.Config,
	listOpts metav1.ListOptions,
) (corev1.NamespaceList, error) {
	client := helpers.DynamicClientFromEnvconf(t, cfg)

	var nsList corev1.NamespaceList

	list, err := client.Resource(corev1.SchemeGroupVersion.WithResource("namespaces")).List(ctx, listOpts)
	if err != nil {
		return nsList, err
	}

	err = runtime.DefaultUnstructuredConverter.FromUnstructured(list.UnstructuredContent(), &nsList)
	if err != nil {
		return nsList, err
	}

	return nsList, nil
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
			secrets, err := getNamespaces(ctx, t, cfg, assert.ListOptions(cfg))
			require.NoError(t, err)

			return itemCountFn(len(secrets.Items), count), nil
		}
	}
}

func areRestricted() helpers.ConditionFuncFactory {
	return func(
		t require.TestingT,
		assert assertion.Assertion,
		cfg *envconf.Config,
		count int,
		itemCountFn, resultFn helpers.IntCompareFunc,
	) helpers.ConditionFunc {
		return func(ctx context.Context) (bool, error) {
			nsList, err := getNamespaces(ctx, t, cfg, assert.ListOptions(cfg))
			require.NoError(t, err)

			if itemCountFn(len(nsList.Items), count) {
				return false, nil
			}

			restrictedCount := 0

			for _, namespace := range nsList.Items {
				nsLabels := namespace.GetLabels()

				enforceLabel, ok := nsLabels[podSecurityEnforceLabelKey]
				if ok && enforceLabel == "restricted" {
					restrictedCount++
				}
			}

			return resultFn(restrictedCount, count), nil
		}
	}
}

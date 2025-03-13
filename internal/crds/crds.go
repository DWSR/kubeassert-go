package crds

import (
	"context"

	"github.com/stretchr/testify/require"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/e2e-framework/pkg/envconf"

	"github.com/DWSR/kubeassert-go/internal/assertion"
	helpers "github.com/DWSR/kubeassert-go/internal/assertionhelpers"
)

func getCRDs(
	ctx context.Context,
	t require.TestingT,
	cfg *envconf.Config,
	listOpts metav1.ListOptions,
) (extv1.CustomResourceDefinitionList, error) {
	client := helpers.DynamicClientFromEnvconf(t, cfg)

	var crdList extv1.CustomResourceDefinitionList

	list, err := client.
		Resource(extv1.SchemeGroupVersion.WithResource("customresourcedefinitions")).
		List(ctx, listOpts)
	if err != nil {
		return crdList, err
	}

	err = runtime.DefaultUnstructuredConverter.FromUnstructured(list.UnstructuredContent(), &crdList)
	if err != nil {
		return crdList, err
	}

	return crdList, nil
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
			crdList, err := getCRDs(ctx, t, cfg, assert.ListOptions(cfg))
			require.NoError(t, err)

			return itemCountFn(len(crdList.Items), count), nil
		}
	}
}

func hasVersion(crdVersion string) helpers.ConditionFuncFactory {
	return func(
		t require.TestingT,
		assert assertion.Assertion,
		cfg *envconf.Config,
		count int,
		itemCountFn, _ helpers.IntCompareFunc,
	) helpers.ConditionFunc {
		return func(ctx context.Context) (bool, error) {
			crdList, err := getCRDs(ctx, t, cfg, assert.ListOptions(cfg))
			require.NoError(t, err)

			if itemCountFn(len(crdList.Items), count) {
				return false, nil
			}

			foundVersion := false

			for _, version := range crdList.Items[0].Spec.Versions {
				if version.Name == crdVersion {
					foundVersion = true

					break
				}
			}

			return foundVersion, nil
		}
	}
}

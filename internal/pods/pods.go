package pods

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

// return default value instead of a nil pointer so that negative assertions (i.e. testing for false positives) can use
// a mock require.TestingT object.
func getPods(ctx context.Context, t require.TestingT, cfg *envconf.Config, listOpts metav1.ListOptions) (corev1.PodList, error) {
	client := helpers.DynamicClientFromEnvconf(t, cfg)

	var pods corev1.PodList

	list, err := client.
		Resource(corev1.SchemeGroupVersion.WithResource("pods")).
		List(ctx, listOpts)
	if err != nil {
		return pods, err
	}

	err = runtime.DefaultUnstructuredConverter.FromUnstructured(list.UnstructuredContent(), &pods)
	if err != nil {
		return pods, err
	}

	return pods, nil
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
			secrets, err := getPods(ctx, t, cfg, assert.ListOptions(cfg))
			require.NoError(t, err)

			return itemCountFn(len(secrets.Items), count), nil
		}
	}
}

func areReady() helpers.ConditionFuncFactory {
	return func(
		t require.TestingT,
		assert assertion.Assertion,
		cfg *envconf.Config,
		count int,
		itemCountFn, resultFn helpers.IntCompareFunc,
	) helpers.ConditionFunc {
		return func(ctx context.Context) (bool, error) {
			pods, err := getPods(ctx, t, cfg, assert.ListOptions(cfg))
			require.NoError(t, err)

			if itemCountFn(len(pods.Items), count) {
				return false, nil
			}

			readyCount := 0

			for _, pod := range pods.Items {
				for _, cond := range pod.Status.Conditions {
					if cond.Type == corev1.PodReady && cond.Status == corev1.ConditionTrue {
						readyCount++

						break
					}
				}
			}

			return resultFn(readyCount, count), nil
		}
	}
}

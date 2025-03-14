// secrets contains assertions for Kubernetes Secrets.
package secrets

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

func secretHasContent(secret corev1.Secret, content map[string]string) bool {
	hasContent := true

	for key, value := range content {
		secData, ok := secret.Data[key]
		if !ok || string(secData) != value {
			hasContent = false

			break
		}
	}

	return hasContent
}

func getSecrets(
	ctx context.Context,
	t require.TestingT,
	cfg *envconf.Config,
	listOpts metav1.ListOptions,
) (corev1.SecretList, error) {
	client := helpers.DynamicClientFromEnvconf(t, cfg)

	var secrets corev1.SecretList

	list, err := client.
		Resource(corev1.SchemeGroupVersion.WithResource("secrets")).
		List(ctx, listOpts)
	if err != nil {
		return secrets, err
	}

	err = runtime.DefaultUnstructuredConverter.FromUnstructured(list.UnstructuredContent(), &secrets)
	if err != nil {
		return secrets, err
	}

	return secrets, nil
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
			secrets, err := getSecrets(ctx, t, cfg, assert.ListOptions(cfg))
			require.NoError(t, err)

			return itemCountFn(len(secrets.Items), count), nil
		}
	}
}

func haveContent(content map[string]string) helpers.ConditionFuncFactory {
	return func(
		t require.TestingT,
		assert assertion.Assertion,
		cfg *envconf.Config,
		count int,
		itemCountFn, resultFn helpers.IntCompareFunc,
	) helpers.ConditionFunc {
		return func(ctx context.Context) (bool, error) {
			secrets, err := getSecrets(ctx, t, cfg, assert.ListOptions(cfg))
			require.NoError(t, err)

			if itemCountFn(len(secrets.Items), count) {
				return false, nil
			}

			haveContent := 0

			for _, secret := range secrets.Items {
				if secretHasContent(secret, content) {
					haveContent++
				}
			}

			return resultFn(haveContent, count), nil
		}
	}
}

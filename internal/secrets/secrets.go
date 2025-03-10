package secrets

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/DWSR/kubeassert-go/internal/assertion"
	helpers "github.com/DWSR/kubeassert-go/internal/assertionhelpers"
)

// SecretAssertion is a wrapper around the assertion.Assertion type and provides a set of assertions for Kubernetes
// Secrets.
type SecretAssertion struct {
	assertion.Assertion
}

func (sa SecretAssertion) clone() SecretAssertion {
	return SecretAssertion{
		Assertion: assertion.Clone(sa.Assertion),
	}
}

// Exists asserts that exactly one Secret exists in the cluster that matches the provided options.
func (sa SecretAssertion) Exists() SecretAssertion {
	return sa.ExactlyNExist(1)
}

// ExactlyNExist asserts that exactly N Secrets exist in the cluster that match the provided options.
func (sa SecretAssertion) ExactlyNExist(count int) SecretAssertion {
	stepFn := func(ctx context.Context, testingT *testing.T, cfg *envconf.Config) context.Context {
		t := helpers.RequireTIfNotNil(testingT, sa.GetRequireT())
		conditionFunc := func(ctx context.Context) (bool, error) {
			secrets, err := sa.getSecrets(ctx, t, cfg)
			require.NoError(t, err)

			return len(secrets.Items) == count, nil
		}

		require.NoError(t, helpers.WaitForCondition(ctx, sa, conditionFunc))

		return ctx
	}

	res := sa.clone()
	res.SetBuilder(res.GetBuilder().Assess("exactlyNExist", stepFn))

	return res
}

// AtLeastNExist asserts that at least N Secrets exist in the cluster that match the provided options.
func (sa SecretAssertion) AtLeastNExist(count int) SecretAssertion {
	stepFn := func(ctx context.Context, testingT *testing.T, cfg *envconf.Config) context.Context {
		t := helpers.RequireTIfNotNil(testingT, sa.GetRequireT())
		conditionFunc := func(ctx context.Context) (bool, error) {
			secrets, err := sa.getSecrets(ctx, t, cfg)
			require.NoError(t, err)

			return len(secrets.Items) >= count, nil
		}

		require.NoError(t, helpers.WaitForCondition(ctx, sa, conditionFunc))

		return ctx
	}

	res := sa.clone()
	res.SetBuilder(res.GetBuilder().Assess("atLeastNExist", stepFn))

	return res
}

func (sa SecretAssertion) getSecrets(
	ctx context.Context,
	t require.TestingT,
	cfg *envconf.Config,
) (corev1.SecretList, error) {
	client := helpers.DynamicClientFromEnvconf(t, cfg)

	var secrets corev1.SecretList

	list, err := client.
		Resource(corev1.SchemeGroupVersion.WithResource("secrets")).
		List(ctx, sa.ListOptions(cfg))
	if err != nil {
		return secrets, err
	}

	err = runtime.DefaultUnstructuredConverter.FromUnstructured(list.UnstructuredContent(), &secrets)
	if err != nil {
		return secrets, err
	}

	return secrets, nil
}

// HasContent asserts that exactly one Secret in the cluster contains the provided content. This match is not exclusive
// meaning that the Secret can contain additional content.
func (sa SecretAssertion) HasContent(content map[string]string) SecretAssertion {
	return sa.ExactlyNHaveContent(1, content)
}

// ExactlyNHaveContent asserts that exactly N Secrets in the cluster contain the provided content. This match is not
// exclusive meaning that the Secrets can contain additional content.
func (sa SecretAssertion) ExactlyNHaveContent(count int, content map[string]string) SecretAssertion {
	stepFn := func(ctx context.Context, testingT *testing.T, cfg *envconf.Config) context.Context {
		t := helpers.RequireTIfNotNil(testingT, sa.GetRequireT())
		conditionFunc := func(ctx context.Context) (bool, error) {
			secrets, err := sa.getSecrets(ctx, t, cfg)
			require.NoError(t, err)

			if len(secrets.Items) != count {
				return false, nil
			}

			haveContent := 0

			for _, secret := range secrets.Items {
				hasContent := true

				for key, value := range content {
					secData, ok := secret.Data[key]
					if !ok || string(secData) != value {
						hasContent = false

						break
					}
				}

				if hasContent {
					haveContent++
				}
			}

			return haveContent == count, nil
		}

		require.NoError(t, helpers.WaitForCondition(ctx, sa, conditionFunc))

		return ctx
	}

	res := sa.clone()
	res.SetBuilder(res.GetBuilder().Assess("exactlyNHaveContent", stepFn))

	return res
}

// AtLeastNHaveContent asserts that at least N Secrets in the cluster contain the provided content. This match is not
// exclusive meaning that the Secrets can contain additional content.
func (sa SecretAssertion) AtLeastNHaveContent(count int, content map[string]string) SecretAssertion {
	stepFn := func(ctx context.Context, testingT *testing.T, cfg *envconf.Config) context.Context {
		t := helpers.RequireTIfNotNil(testingT, sa.GetRequireT())
		conditionFunc := func(ctx context.Context) (bool, error) {
			secrets, err := sa.getSecrets(ctx, t, cfg)
			require.NoError(t, err)

			if len(secrets.Items) < count {
				return false, nil
			}

			haveContent := 0

			for _, secret := range secrets.Items {
				hasContent := true

				for key, value := range content {
					secData, ok := secret.Data[key]
					if !ok || string(secData) != value {
						hasContent = false

						break
					}
				}

				if hasContent {
					haveContent++
				}
			}

			return haveContent >= count, nil
		}

		require.NoError(t, helpers.WaitForCondition(ctx, sa, conditionFunc))

		return ctx
	}

	res := sa.clone()
	res.SetBuilder(res.GetBuilder().Assess("atLeastNHaveContent", stepFn))

	return res
}

// NewSecretAssertion creates a new SecretAssertion with the provided options.
func NewSecretAssertion(opts ...assertion.Option) SecretAssertion {
	return SecretAssertion{
		Assertion: assertion.NewAssertion(
			append(
				[]assertion.Option{assertion.WithBuilder(features.New("Secret").WithLabel("type", "secret"))},
				opts...,
			)...,
		),
	}
}

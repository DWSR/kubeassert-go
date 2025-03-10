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

type SecretAssertion struct {
	assertion.Assertion
}

func (sa SecretAssertion) clone() SecretAssertion {
	return SecretAssertion{
		Assertion: assertion.Clone(sa.Assertion),
	}
}

func (sa SecretAssertion) Exists() SecretAssertion {
	return sa.ExactlyNExist(1)
}

func (sa SecretAssertion) ExactlyNExist(count int) SecretAssertion {
	fn := func(ctx context.Context, testingT *testing.T, cfg *envconf.Config) context.Context {
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
	res.SetBuilder(res.GetBuilder().Assess("exactlyNExist", fn))

	return res
}

func (sa SecretAssertion) AtLeastNExist(count int) SecretAssertion {
	fn := func(ctx context.Context, testingT *testing.T, cfg *envconf.Config) context.Context {
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
	res.SetBuilder(res.GetBuilder().Assess("atLeastNExist", fn))

	return res
}

func (sa SecretAssertion) getSecrets(ctx context.Context, t require.TestingT, cfg *envconf.Config) (corev1.SecretList, error) {
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

func (sa SecretAssertion) HasContent(content map[string]string) SecretAssertion {
	return sa.ExactlyNHaveContent(1, content)
}

func (sa SecretAssertion) ExactlyNHaveContent(count int, content map[string]string) SecretAssertion {
	fn := func(ctx context.Context, testingT *testing.T, cfg *envconf.Config) context.Context {
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

			return haveContent == count, nil
		}

		require.NoError(t, helpers.WaitForCondition(ctx, sa, conditionFunc))

		return ctx
	}

	res := sa.clone()
	res.SetBuilder(res.GetBuilder().Assess("hasContent", fn))

	return res
}

func (sa SecretAssertion) AtLeastNHaveContent(count int, content map[string]string) SecretAssertion {
	fn := func(ctx context.Context, testingT *testing.T, cfg *envconf.Config) context.Context {
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
	res.SetBuilder(res.GetBuilder().Assess("hasContent", fn))

	return res
}

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

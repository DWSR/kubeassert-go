package pods

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

type PodAssertion struct {
	assertion.Assertion
}

func (pa PodAssertion) clone() PodAssertion {
	return PodAssertion{
		Assertion: assertion.CloneAssertion(pa.Assertion),
	}
}

func (pa PodAssertion) Exists() PodAssertion {
	return pa.ExactlyNExist(1)
}

func (pa PodAssertion) ExactlyNExist(count int) PodAssertion {
	fn := func(ctx context.Context, testingT *testing.T, cfg *envconf.Config) context.Context {
		t := helpers.RequireTIfNotNil(testingT, pa.GetRequireT())
		conditionFunc := func(ctx context.Context) (bool, error) {
			pods, err := pa.getPods(ctx, t, cfg)
			require.NoError(t, err)

			return len(pods.Items) == count, nil
		}

		require.NoError(t, pa.WaitForCondition(ctx, conditionFunc))

		return ctx
	}

	res := pa.clone()
	res.SetBuilder(res.GetBuilder().Assess("exists", fn))

	return res
}

func (pa PodAssertion) AtLeastNExist(count int) PodAssertion {
	fn := func(ctx context.Context, testingT *testing.T, cfg *envconf.Config) context.Context {
		t := helpers.RequireTIfNotNil(testingT, pa.GetRequireT())
		conditionFunc := func(ctx context.Context) (bool, error) {
			pods, err := pa.getPods(ctx, t, cfg)
			require.NoError(t, err)

			return len(pods.Items) >= count, nil
		}

		require.NoError(t, pa.WaitForCondition(ctx, conditionFunc))

		return ctx
	}

	res := pa.clone()
	res.SetBuilder(res.GetBuilder().Assess("exists", fn))

	return res
}

// return default value instead of a nil pointer so that negative assertions (i.e. testing for false positives) can use
// a mock require.TestingT object.
func (pa PodAssertion) getPods(ctx context.Context, t require.TestingT, cfg *envconf.Config) (corev1.PodList, error) {
	client := helpers.DynamicClientFromEnvconf(t, cfg)

	var pods corev1.PodList

	list, err := client.
		Resource(corev1.SchemeGroupVersion.WithResource("pods")).
		List(ctx, pa.ListOptions(cfg))
	if err != nil {
		return pods, err
	}

	err = runtime.DefaultUnstructuredConverter.FromUnstructured(list.UnstructuredContent(), &pods)
	if err != nil {
		return pods, err
	}

	return pods, nil
}

func (pa PodAssertion) IsReady() PodAssertion {
	return pa.ExactlyNAreReady(1)
}

func (pa PodAssertion) ExactlyNAreReady(count int) PodAssertion {
	fn := func(ctx context.Context, testingT *testing.T, cfg *envconf.Config) context.Context {
		t := helpers.RequireTIfNotNil(testingT, pa.GetRequireT())

		conditionFunc := func(ctx context.Context) (bool, error) {
			pods, err := pa.getPods(ctx, t, cfg)
			require.NoError(t, err)

			if len(pods.Items) < count {
				return false, nil
			}

			readyCount := 0

			for _, pod := range pods.Items {
				for _, cond := range pod.Status.Conditions {
					if cond.Type == corev1.PodReady && cond.Status == corev1.ConditionTrue {
						readyCount += 1
						break
					}
				}
			}

			return readyCount == count, nil
		}

		require.NoError(t, pa.WaitForCondition(ctx, conditionFunc))

		return ctx
	}

	res := pa.clone()
	res.SetBuilder(pa.GetBuilder().Assess("isReady", fn))

	return res
}

func (pa PodAssertion) AtLeastNAreReady(count int) PodAssertion {
	fn := func(ctx context.Context, testingT *testing.T, cfg *envconf.Config) context.Context {
		t := helpers.RequireTIfNotNil(testingT, pa.GetRequireT())

		conditionFunc := func(ctx context.Context) (bool, error) {
			pods, err := pa.getPods(ctx, t, cfg)
			require.NoError(t, err)

			if len(pods.Items) < count {
				return false, nil
			}

			readyCount := 0

			for _, pod := range pods.Items {
				for _, cond := range pod.Status.Conditions {
					if cond.Type == corev1.PodReady && cond.Status == corev1.ConditionTrue {
						readyCount += 1
						break
					}
				}
			}

			return readyCount >= count, nil
		}

		require.NoError(t, pa.WaitForCondition(ctx, conditionFunc))

		return ctx
	}

	res := pa.clone()
	res.SetBuilder(pa.GetBuilder().Assess("isReady", fn))

	return res
}

func NewPodAssertion(opts ...assertion.AssertionOption) PodAssertion {
	return PodAssertion{
		Assertion: assertion.NewAssertion(
			append(
				[]assertion.AssertionOption{assertion.WithBuilder(features.New("Pod").WithLabel("type", "pod"))},
				opts...,
			)...,
		),
	}
}

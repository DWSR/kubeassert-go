package pods

import (
	"context"

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
		Assertion: assertion.Clone(pa.Assertion),
	}
}

func (pa PodAssertion) Exists() PodAssertion {
	return pa.ExactlyNExist(1)
}

func (pa PodAssertion) ExactlyNExist(count int) PodAssertion {
	stepFn := helpers.AsStepFunc(pa, exist(), count, helpers.IntCompareFuncEqualTo, nil)

	res := pa.clone()
	res.SetBuilder(res.GetBuilder().Assess("exactlyNExist", stepFn))

	return res
}

func (pa PodAssertion) AtLeastNExist(count int) PodAssertion {
	stepFn := helpers.AsStepFunc(pa, exist(), count, helpers.IntCompareFuncGreaterThanOrEqualTo, nil)

	res := pa.clone()
	res.SetBuilder(res.GetBuilder().Assess("atLeastNExist", stepFn))

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
	stepFn := helpers.AsStepFunc(pa, areReady(), count, helpers.IntCompareFuncNotEqualTo, helpers.IntCompareFuncEqualTo)
	res := pa.clone()
	res.SetBuilder(pa.GetBuilder().Assess("exactlyNAreReady", stepFn))

	return res
}

func (pa PodAssertion) AtLeastNAreReady(count int) PodAssertion {
	stepFn := helpers.AsStepFunc(pa, areReady(), count, helpers.IntCompareFuncLessThan, helpers.IntCompareFuncGreaterThanOrEqualTo)
	res := pa.clone()
	res.SetBuilder(pa.GetBuilder().Assess("atLeastNAreReady", stepFn))

	return res
}

func NewPodAssertion(opts ...assertion.Option) PodAssertion {
	return PodAssertion{
		Assertion: assertion.NewAssertion(
			append(
				[]assertion.Option{assertion.WithBuilder(features.New("Pod").WithLabel("type", "pod"))},
				opts...,
			)...,
		),
	}
}

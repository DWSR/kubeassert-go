package namespaces

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

type NamespaceAssertion struct {
	assertion.Assertion
}

const (
	podSecurityEnforceLabelKey = "pod-security.kubernetes.io/enforce"
)

func (na NamespaceAssertion) clone() NamespaceAssertion {
	return NamespaceAssertion{
		Assertion: assertion.CloneAssertion(na.Assertion),
	}
}

func (na NamespaceAssertion) getNamespaces(ctx context.Context, t require.TestingT, cfg *envconf.Config) (corev1.NamespaceList, error) {
	client := helpers.DynamicClientFromEnvconf(t, cfg)

	var nsList corev1.NamespaceList

	list, err := client.Resource(corev1.SchemeGroupVersion.WithResource("namespaces")).List(ctx, na.ListOptions(cfg))
	if err != nil {
		return nsList, err
	}

	err = runtime.DefaultUnstructuredConverter.FromUnstructured(list.UnstructuredContent(), &nsList)
	if err != nil {
		return nsList, err
	}

	return nsList, nil
}

func (na NamespaceAssertion) Exists() NamespaceAssertion {
	return na.ExactlyNExist(1)
}

func (na NamespaceAssertion) ExactlyNExist(count int) NamespaceAssertion {
	fn := func(ctx context.Context, testingT *testing.T, cfg *envconf.Config) context.Context {
		t := helpers.RequireTIfNotNil(testingT, na.GetRequireT())
		conditionFunc := func(ctx context.Context) (bool, error) {
			nsList, err := na.getNamespaces(ctx, t, cfg)
			require.NoError(t, err)

			return len(nsList.Items) == count, nil
		}

		require.NoError(t, na.WaitForCondition(ctx, conditionFunc))

		return ctx
	}

	res := na.clone()
	res.SetBuilder(res.GetBuilder().Assess("exists", fn))

	return res
}

func (na NamespaceAssertion) AtLeastNExist(count int) NamespaceAssertion {
	fn := func(ctx context.Context, testingT *testing.T, cfg *envconf.Config) context.Context {
		t := helpers.RequireTIfNotNil(testingT, na.GetRequireT())
		conditionFunc := func(ctx context.Context) (bool, error) {
			nsList, err := na.getNamespaces(ctx, t, cfg)
			require.NoError(t, err)

			return len(nsList.Items) >= count, nil
		}

		require.NoError(t, na.WaitForCondition(ctx, conditionFunc))

		return ctx
	}

	res := na.clone()
	res.SetBuilder(res.GetBuilder().Assess("exists", fn))

	return res
}

func (na NamespaceAssertion) AtLeastNAreRestricted(count int) NamespaceAssertion {
	fn := func(ctx context.Context, testingT *testing.T, cfg *envconf.Config) context.Context {
		t := helpers.RequireTIfNotNil(testingT, na.GetRequireT())
		conditionFunc := func(ctx context.Context) (bool, error) {
			nsList, err := na.getNamespaces(ctx, t, cfg)
			require.NoError(t, err)

			if len(nsList.Items) < count {
				return false, nil
			}

			restrictedCount := 0

			for _, ns := range nsList.Items {
				nsLabels := ns.GetLabels()

				enforceLabel, ok := nsLabels[podSecurityEnforceLabelKey]
				if ok && enforceLabel == "restricted" {
					restrictedCount += 1
				}
			}

			return restrictedCount >= count, nil
		}

		require.NoError(t, na.WaitForCondition(ctx, conditionFunc))

		return ctx
	}

	res := na.clone()
	res.SetBuilder(na.GetBuilder().Assess("atLeastNAreRestricted", fn))

	return res
}

func (na NamespaceAssertion) IsRestricted() NamespaceAssertion {
	return na.ExactlyNAreRestricted(1)
}

func (na NamespaceAssertion) ExactlyNAreRestricted(count int) NamespaceAssertion {
	fn := func(ctx context.Context, testingT *testing.T, cfg *envconf.Config) context.Context {
		t := helpers.RequireTIfNotNil(testingT, na.GetRequireT())
		conditionFunc := func(ctx context.Context) (bool, error) {
			nsList, err := na.getNamespaces(ctx, t, cfg)
			require.NoError(t, err)

			if len(nsList.Items) < count {
				return false, nil
			}

			restrictedCount := 0

			for _, ns := range nsList.Items {
				nsLabels := ns.GetLabels()

				enforceLabel, ok := nsLabels[podSecurityEnforceLabelKey]
				if ok && enforceLabel == "restricted" {
					restrictedCount += 1
				}
			}

			return restrictedCount == count, nil
		}

		require.NoError(t, na.WaitForCondition(ctx, conditionFunc))

		return ctx
	}

	res := na.clone()
	res.SetBuilder(na.GetBuilder().Assess("atLeastNAreRestricted", fn))

	return res
}

func NewNamespaceAssertion(opts ...assertion.AssertionOption) NamespaceAssertion {
	return NamespaceAssertion{
		Assertion: assertion.NewAssertion(
			append(
				[]assertion.AssertionOption{assertion.WithBuilder(features.New("Namespace").WithLabel("type", "namespace"))},
				opts...,
			)...,
		),
	}
}

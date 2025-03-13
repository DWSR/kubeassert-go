package crds

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/DWSR/kubeassert-go/internal/assertion"
	helpers "github.com/DWSR/kubeassert-go/internal/assertionhelpers"
)

type CRDAssertion struct {
	assertion.Assertion
}

func (ca CRDAssertion) clone() CRDAssertion {
	return CRDAssertion{
		Assertion: assertion.Clone(ca.Assertion),
	}
}

func (ca CRDAssertion) Exists() CRDAssertion {
	stepFn := helpers.AsStepFunc(ca, exist(), 1, helpers.IntCompareFuncEqualTo, nil)

	res := ca.clone()
	res.SetBuilder(res.GetBuilder().Assess("exists", stepFn))

	return res
}

func (ca CRDAssertion) getCRDs(ctx context.Context, t require.TestingT, cfg *envconf.Config) (extv1.CustomResourceDefinitionList, error) {
	client := helpers.DynamicClientFromEnvconf(t, cfg)

	var crdList extv1.CustomResourceDefinitionList

	list, err := client.
		Resource(extv1.SchemeGroupVersion.WithResource("customresourcedefinitions")).
		List(ctx, ca.ListOptions(cfg))
	if err != nil {
		return crdList, err
	}

	err = runtime.DefaultUnstructuredConverter.FromUnstructured(list.UnstructuredContent(), &crdList)
	if err != nil {
		return crdList, err
	}

	return crdList, nil
}

func (ca CRDAssertion) HasVersion(crdVersion string) CRDAssertion {
	fn := func(ctx context.Context, testingT *testing.T, cfg *envconf.Config) context.Context {
		t := helpers.RequireTIfNotNil(testingT, ca.GetRequireT())
		conditionFunc := func(ctx context.Context) (bool, error) {
			crds, err := ca.getCRDs(ctx, t, cfg)
			require.NoError(t, err)

			if len(crds.Items) != 1 {
				return false, nil
			}

			foundVersion := false

			for _, version := range crds.Items[0].Spec.Versions {
				if version.Name == crdVersion {
					foundVersion = true

					break
				}
			}

			return foundVersion, nil
		}

		require.NoError(t, helpers.WaitForCondition(ctx, ca, conditionFunc))

		return ctx
	}

	res := ca.clone()
	res.SetBuilder(res.GetBuilder().Assess("hasVersion", fn))

	return res
}

func NewCRDAssertion(opts ...assertion.Option) CRDAssertion {
	return CRDAssertion{
		Assertion: assertion.NewAssertion(
			append(
				[]assertion.Option{assertion.WithBuilder(features.New("CRD").WithLabel("type", "customresourcedefinition"))},
				opts...,
			)...,
		),
	}
}

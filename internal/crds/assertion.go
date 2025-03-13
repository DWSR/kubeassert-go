package crds

import (
	"context"

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
	stepFn := helpers.AsStepFunc(ca, hasVersion(crdVersion), 1, helpers.IntCompareFuncNotEqualTo, helpers.IntCompareFuncEqualTo)
	res := ca.clone()
	res.SetBuilder(res.GetBuilder().Assess("hasVersion", stepFn))

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

package namespaces_test

import (
	"testing"
	"time"

	"github.com/DWSR/kubeassert-go/internal/assertion"
	helpers "github.com/DWSR/kubeassert-go/internal/assertionhelpers"
	"github.com/DWSR/kubeassert-go/internal/namespaces"
	"github.com/DWSR/kubeassert-go/internal/testhelpers"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/e2e-framework/klient/decoder"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

func Test_1Namespace_Success(t *testing.T) {
	asserts := []testhelpers.SuccessfulAssert{
		{
			Name: "Exists_Labels",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				namespaceName := envconf.RandomName("test", 20)

				return namespaces.NewNamespaceAssertion(
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": namespaceName}),
					assertion.WithSetup(
						helpers.CreateResourceFromPath(
							"./testdata/restricted-namespace.yaml",
							testhelpers.MutateResourceName(namespaceName),
							decoder.MutateLabels(map[string]string{"app.kubernetes.io/name": namespaceName}),
						),
					),
					assertion.WithTeardown(
						helpers.DeleteResourceFromPath(
							"./testdata/restricted-namespace.yaml",
							testhelpers.MutateResourceName(namespaceName),
						),
					),
				).Exists()
			},
		},
		{
			Name: "Exists_Name",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				namespaceName := envconf.RandomName("test", 20)

				return namespaces.NewNamespaceAssertion(
					assertion.WithResourceName(namespaceName),
					assertion.WithSetup(
						helpers.CreateResourceFromPath(
							"./testdata/restricted-namespace.yaml",
							testhelpers.MutateResourceName(namespaceName),
							decoder.MutateLabels(map[string]string{"app.kubernetes.io/name": namespaceName}),
						),
					),
					assertion.WithTeardown(
						helpers.DeleteResourceFromPath(
							"./testdata/restricted-namespace.yaml",
							testhelpers.MutateResourceName(namespaceName),
						),
					),
				).Exists()
			},
		},
		{
			Name: "IsRestricted_Labels",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				namespaceName := envconf.RandomName("test", 20)

				return namespaces.NewNamespaceAssertion(
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": namespaceName}),
					assertion.WithSetup(
						helpers.CreateResourceFromPath(
							"./testdata/restricted-namespace.yaml",
							testhelpers.MutateResourceName(namespaceName),
							decoder.MutateLabels(map[string]string{"app.kubernetes.io/name": namespaceName}),
						),
					),
					assertion.WithTeardown(
						helpers.DeleteResourceFromPath(
							"./testdata/restricted-namespace.yaml",
							testhelpers.MutateResourceName(namespaceName),
						),
					),
				).Exists().IsRestricted()
			},
		},
		{
			Name: "IsRestricted_Name",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				namespaceName := envconf.RandomName("test", 20)

				return namespaces.NewNamespaceAssertion(
					assertion.WithResourceName(namespaceName),
					assertion.WithSetup(
						helpers.CreateResourceFromPath(
							"./testdata/restricted-namespace.yaml",
							testhelpers.MutateResourceName(namespaceName),
							decoder.MutateLabels(map[string]string{"app.kubernetes.io/name": namespaceName}),
						),
					),
					assertion.WithTeardown(
						helpers.DeleteResourceFromPath(
							"./testdata/restricted-namespace.yaml",
							testhelpers.MutateResourceName(namespaceName),
						),
					),
				).Exists().IsRestricted()
			},
		},
	}

	testhelpers.TestSuccessfulAsserts(t, testEnv, asserts...)
}

func Test_1Namespace_Fail(t *testing.T) {
	asserts := []testhelpers.FailingAssert{
		{
			Name: "Exists_Labels",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				namespaceName := envconf.RandomName("test", 20)

				return namespaces.NewNamespaceAssertion(
					assertion.WithRequireT(t),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": namespaceName}),
				).Exists()
			},
		},
		{
			Name: "Exists_Name",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				namespaceName := envconf.RandomName("test", 20)

				return namespaces.NewNamespaceAssertion(
					assertion.WithRequireT(t),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithResourceName(namespaceName),
				).Exists()
			},
		},
		{
			Name: "IsRestricted_Labels",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				namespaceName := envconf.RandomName("test", 20)

				return namespaces.NewNamespaceAssertion(
					assertion.WithRequireT(t),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": namespaceName}),
					assertion.WithSetup(
						helpers.CreateResourceFromPath(
							"./testdata/unrestricted-namespace.yaml",
							testhelpers.MutateResourceName(namespaceName),
							decoder.MutateLabels(map[string]string{"app.kubernetes.io/name": namespaceName}),
						),
					),
					assertion.WithTeardown(
						helpers.DeleteResourceFromPath(
							"./testdata/unrestricted-namespace.yaml",
							testhelpers.MutateResourceName(namespaceName),
						),
					),
				).Exists().IsRestricted()
			},
		},
		{
			Name: "IsRestricted_Name",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				namespaceName := envconf.RandomName("test", 20)

				return namespaces.NewNamespaceAssertion(
					assertion.WithRequireT(t),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithResourceName(namespaceName),
					assertion.WithSetup(
						helpers.CreateResourceFromPath(
							"./testdata/unrestricted-namespace.yaml",
							testhelpers.MutateResourceName(namespaceName),
							decoder.MutateLabels(map[string]string{"app.kubernetes.io/name": namespaceName}),
						),
					),
					assertion.WithTeardown(
						helpers.DeleteResourceFromPath(
							"./testdata/unrestricted-namespace.yaml",
							testhelpers.MutateResourceName(namespaceName),
						),
					),
				).Exists().IsRestricted()
			},
		},
	}

	testhelpers.TestFailingAsserts(t, testEnv, asserts...)
}

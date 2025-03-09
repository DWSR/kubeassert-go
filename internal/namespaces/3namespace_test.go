package namespaces_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"sigs.k8s.io/e2e-framework/klient/decoder"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	e2etypes "sigs.k8s.io/e2e-framework/pkg/types"

	"github.com/DWSR/kubeassert-go/internal/assertion"
	helpers "github.com/DWSR/kubeassert-go/internal/assertionhelpers"
	"github.com/DWSR/kubeassert-go/internal/namespaces"
	"github.com/DWSR/kubeassert-go/internal/testhelpers"
)

func Test_3Namespace_Success(t *testing.T) {
	asserts := []testhelpers.SuccessfulAssert{
		{
			Name: "AtLeastNExist",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				namespaceNames := generateNamespaceNames()

				return namespaces.NewNamespaceAssertion(
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": namespaceNames[0]}),
					assertion.WithSetup(createNamespaces("./testdata/restricted-namespace.yaml", namespaceNames)...),
					assertion.WithTeardown(deleteNamespaces("./testdata/restricted-namespace.yaml", namespaceNames)...),
				).AtLeastNExist(2)
			},
		},
		{
			Name: "ExactlyNExist",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				namespaceNames := generateNamespaceNames()

				return namespaces.NewNamespaceAssertion(
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": namespaceNames[0]}),
					assertion.WithSetup(createNamespaces("./testdata/restricted-namespace.yaml", namespaceNames)...),
					assertion.WithTeardown(deleteNamespaces("./testdata/restricted-namespace.yaml", namespaceNames)...),
				).ExactlyNExist(3)
			},
		},
		{
			Name: "AtLeastNAreRestricted",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				namespaceNames := generateNamespaceNames()

				return namespaces.NewNamespaceAssertion(
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": namespaceNames[0]}),
					assertion.WithSetup(createNamespaces("./testdata/restricted-namespace.yaml", namespaceNames)...),
					assertion.WithTeardown(deleteNamespaces("./testdata/restricted-namespace.yaml", namespaceNames)...),
				).AtLeastNExist(2).AtLeastNAreRestricted(2)
			},
		},
		{
			Name: "ExactlyNAreRestricted",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				namespaceNames := generateNamespaceNames()

				return namespaces.NewNamespaceAssertion(
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": namespaceNames[0]}),
					assertion.WithSetup(createNamespaces("./testdata/restricted-namespace.yaml", namespaceNames)...),
					assertion.WithTeardown(deleteNamespaces("./testdata/restricted-namespace.yaml", namespaceNames)...),
				).ExactlyNExist(3).ExactlyNAreRestricted(3)
			},
		},
	}

	testhelpers.TestSuccessfulAsserts(t, testEnv, asserts...)
}

func Test_3Namespace_Fail(t *testing.T) {
	asserts := []testhelpers.FailingAssert{
		{
			Name: "AtLeastNExist",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				namespaceNames := generateNamespaceNames()

				return namespaces.NewNamespaceAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": namespaceNames[0]}),
					assertion.WithSetup(createNamespaces("./testdata/restricted-namespace.yaml", namespaceNames)...),
					assertion.WithTeardown(deleteNamespaces("./testdata/restricted-namespace.yaml", namespaceNames)...),
				).AtLeastNExist(4)
			},
		},
		{
			Name: "ExactlyNExist",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				namespaceNames := generateNamespaceNames()

				return namespaces.NewNamespaceAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": namespaceNames[0]}),
					assertion.WithSetup(createNamespaces("./testdata/restricted-namespace.yaml", namespaceNames)...),
					assertion.WithTeardown(deleteNamespaces("./testdata/restricted-namespace.yaml", namespaceNames)...),
				).ExactlyNExist(2)
			},
		},
		{
			Name: "AtLeastNAreRestricted",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				namespaceNames := generateNamespaceNames()

				return namespaces.NewNamespaceAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": namespaceNames[0]}),
					assertion.WithSetup(createNamespaces("./testdata/unrestricted-namespace.yaml", namespaceNames)...),
					assertion.WithTeardown(deleteNamespaces("./testdata/unrestricted-namespace.yaml", namespaceNames)...),
				).AtLeastNExist(3).AtLeastNAreRestricted(3)
			},
		},
		{
			Name: "ExactlyNAreRestricted",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				namespaceNames := generateNamespaceNames()

				return namespaces.NewNamespaceAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": namespaceNames[0]}),
					assertion.WithSetup(createNamespaces("./testdata/unrestricted-namespace.yaml", namespaceNames)...),
					assertion.WithTeardown(deleteNamespaces("./testdata/unrestricted-namespace.yaml", namespaceNames)...),
				).ExactlyNExist(3).ExactlyNAreRestricted(3)
			},
		},
	}

	testhelpers.TestFailingAsserts(t, testEnv, asserts...)
}

func generateNamespaceNames() []string {
	res := make([]string, 3)

	for i := range res {
		res[i] = envconf.RandomName("test", 20)
	}

	return res
}

func createNamespaces(resourcePath string, namespaceNames []string) []e2etypes.StepFunc {
	if len(namespaceNames) == 0 {
		panic("must supply namespace names")
	}

	res := make([]e2etypes.StepFunc, len(namespaceNames))
	labelValue := namespaceNames[0]

	for i, nsName := range namespaceNames {
		res[i] = helpers.CreateResourceFromPath(
			resourcePath,
			testhelpers.MutateResourceName(nsName),
			decoder.MutateLabels(map[string]string{"app.kubernetes.io/name": labelValue}),
		)
	}

	return res
}

func deleteNamespaces(resourcePath string, namespaceNames []string) []e2etypes.StepFunc {
	if len(namespaceNames) == 0 {
		panic("must supply namespace names")
	}

	res := make([]e2etypes.StepFunc, len(namespaceNames))

	for i, nsName := range namespaceNames {
		res[i] = helpers.DeleteResourceFromPath(
			resourcePath,
			testhelpers.MutateResourceName(nsName),
		)
	}

	return res
}

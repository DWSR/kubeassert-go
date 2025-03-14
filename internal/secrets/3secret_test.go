package secrets_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"sigs.k8s.io/e2e-framework/klient/decoder"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	e2etypes "sigs.k8s.io/e2e-framework/pkg/types"

	"github.com/DWSR/kubeassert-go/internal/assertion"
	helpers "github.com/DWSR/kubeassert-go/internal/assertionhelpers"
	"github.com/DWSR/kubeassert-go/internal/secrets"
	"github.com/DWSR/kubeassert-go/internal/testhelpers"
)

func Test_3Secret_Success(t *testing.T) {
	asserts := []testhelpers.SuccessfulAssert{
		{
			Name: "ExactlyNExist",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				secretNames := generateSecretNames()

				return secrets.NewSecretAssertion(
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": secretNames[0]}),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithSetup(createSecrets(secretNames)...),
				).ExactlyNExist(3)
			},
		},
		{
			Name: "ExactlyNHaveContent",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				secretNames := generateSecretNames()

				return secrets.NewSecretAssertion(
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": secretNames[0]}),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithSetup(createSecrets(secretNames)...),
				).ExactlyNExist(3).ExactlyNHaveContent(3, map[string]string{"foo": "bar"})
			},
		},
		{
			Name: "AtLeastNExist",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				secretNames := generateSecretNames()

				return secrets.NewSecretAssertion(
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": secretNames[0]}),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithSetup(createSecrets(secretNames)...),
				).AtLeastNExist(2)
			},
		},
		{
			Name: "AtLeastNHaveContent",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				secretNames := generateSecretNames()

				return secrets.NewSecretAssertion(
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": secretNames[0]}),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithSetup(createSecrets(secretNames)...),
				).AtLeastNExist(2).AtLeastNHaveContent(2, map[string]string{"foo": "bar"})
			},
		},
	}

	testhelpers.TestSuccessfulAsserts(t, testEnv, asserts...)
}

func Test_3Secret_Fail(t *testing.T) {
	asserts := []testhelpers.FailingAssert{
		{
			Name: "ExactlyNExist",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				secretNames := generateSecretNames()

				return secrets.NewSecretAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": secretNames[0]}),
					assertion.WithSetup(createSecrets(secretNames)...),
					assertion.WithResourceNamespaceFromTestEnv(),
				).ExactlyNExist(2)
			},
		},
		{
			Name: "AtLeastNExist",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				secretNames := generateSecretNames()

				return secrets.NewSecretAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": secretNames[0]}),
					assertion.WithSetup(createSecrets(secretNames)...),
					assertion.WithResourceNamespaceFromTestEnv(),
				).AtLeastNExist(4)
			},
		},
		{
			Name: "ExactlyNHaveContent",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				secretNames := generateSecretNames()

				return secrets.NewSecretAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": secretNames[0]}),
					assertion.WithSetup(createSecrets(secretNames)...),
					assertion.WithResourceNamespaceFromTestEnv(),
				).ExactlyNExist(3).ExactlyNHaveContent(2, map[string]string{"foo": "bar"})
			},
		},
		{
			Name: "AtLeastNHaveContent",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				secretNames := generateSecretNames()

				return secrets.NewSecretAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": secretNames[0]}),
					assertion.WithSetup(createSecrets(secretNames)...),
					assertion.WithResourceNamespaceFromTestEnv(),
				).ExactlyNExist(3).ExactlyNHaveContent(4, map[string]string{"foo": "bar"})
			},
		},
	}

	testhelpers.TestFailingAsserts(t, testEnv, asserts...)
}

func generateSecretNames() []string {
	res := make([]string, 3)

	for i := range res {
		res[i] = envconf.RandomName("test", 20)
	}

	return res
}

func createSecrets(secretNames []string) []e2etypes.StepFunc {
	if len(secretNames) == 0 {
		panic("must supply secret names")
	}

	res := make([]e2etypes.StepFunc, len(secretNames))
	labelValue := secretNames[0]

	for i, nsName := range secretNames {
		res[i] = helpers.CreateResourceFromPathWithNamespaceFromEnv(
			secretPath,
			testhelpers.MutateResourceName(nsName),
			decoder.MutateLabels(map[string]string{"app.kubernetes.io/name": labelValue}),
		)
	}

	return res
}

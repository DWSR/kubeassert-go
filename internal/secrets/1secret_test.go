package secrets_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/DWSR/kubeassert-go/internal/assertion"
	helpers "github.com/DWSR/kubeassert-go/internal/assertionhelpers"
	"github.com/DWSR/kubeassert-go/internal/secrets"
	"github.com/DWSR/kubeassert-go/internal/testhelpers"
)

func Test_1Secret_Success(t *testing.T) {
	asserts := []testhelpers.SuccessfulAssert{
		{
			Name: "Exists_Labels",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				return secrets.NewSecretAssertion(
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": "secrets_test"}),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithSetup(helpers.CreateResourceFromPathWithNamespaceFromEnv(secretPath)),
				).Exists()
			},
		},
		{
			Name: "Exists_Name",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				return secrets.NewSecretAssertion(
					assertion.WithResourceName("test-secret"),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithSetup(helpers.CreateResourceFromPathWithNamespaceFromEnv(secretPath)),
				).Exists()
			},
		},
		{
			Name: "HasContent_Labels",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				return secrets.NewSecretAssertion(
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": "secrets_test"}),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithSetup(helpers.CreateResourceFromPathWithNamespaceFromEnv(secretPath)),
				).Exists().HasContent(map[string]string{"foo": "bar"})
			},
		},
		{
			Name: "HasContent_Name",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				return secrets.NewSecretAssertion(
					assertion.WithResourceName("test-secret"),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithSetup(helpers.CreateResourceFromPathWithNamespaceFromEnv(secretPath)),
				).Exists().HasContent(map[string]string{"foo": "bar", "bar": "baz"})
			},
		},
	}

	testhelpers.TestSuccessfulAsserts(t, testEnv, asserts...)
}

func Test_1Secret_Fail(t *testing.T) {
	asserts := []testhelpers.FailingAssert{
		{
			Name: "Exists_Labels",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				return secrets.NewSecretAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": "secrets_test"}),
					assertion.WithResourceNamespaceFromTestEnv(),
				).Exists()
			},
		},
		{
			Name: "Exists_Name",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				return secrets.NewSecretAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceName("test-secret"),
					assertion.WithResourceNamespaceFromTestEnv(),
				).Exists()
			},
		},
		{
			Name: "HasContent_Labels",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				return secrets.NewSecretAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": "secrets_test"}),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithSetup(helpers.CreateResourceFromPathWithNamespaceFromEnv(secretPath)),
				).Exists().HasContent(map[string]string{"foo": "bar", "baz": "qux"})
			},
		},
		{
			Name: "HasContent_Name",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				return secrets.NewSecretAssertion(
					assertion.WithRequireT(t),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceName("test-secret"),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithSetup(helpers.CreateResourceFromPathWithNamespaceFromEnv(secretPath)),
				).Exists().HasContent(map[string]string{"baz": "qux"})
			},
		},
	}

	testhelpers.TestFailingAsserts(t, testEnv, asserts...)
}

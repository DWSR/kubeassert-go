package pods_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/DWSR/kubeassert-go/internal/assertion"
	helpers "github.com/DWSR/kubeassert-go/internal/assertionhelpers"
	"github.com/DWSR/kubeassert-go/internal/pods"
	"github.com/DWSR/kubeassert-go/internal/testhelpers"
)

func Test_1Pod_Success(t *testing.T) {
	asserts := []testhelpers.SuccessfulAssert{
		{
			Name: "Exists_Labels",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				return pods.NewPodAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": "pods_test"}),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(readyPodPath),
					),
				).Exists()
			},
		},
		{
			Name: "Exists_Name",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				return pods.NewPodAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceName("test-pod"),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(readyPodPath),
					),
				).Exists()
			},
		},
		{
			Name: "IsReady_Labels",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				return pods.NewPodAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": "pods_test"}),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath),
						helpers.CreateResourceFromPathWithNamespaceFromEnv(readyPodPath),
					),
				).Exists().IsReady()
			},
		},
		{
			Name: "IsReady_Name",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				return pods.NewPodAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithResourceName("test-pod"),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath),
						helpers.CreateResourceFromPathWithNamespaceFromEnv(readyPodPath),
					),
				).Exists().IsReady()
			},
		},
	}

	testhelpers.TestSuccessfulAsserts(t, testEnv, asserts...)
}

func Test_1Pod_Fail(t *testing.T) {
	asserts := []testhelpers.FailingAssert{
		{
			Name: "Exists_Labels",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				return pods.NewPodAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithRequireT(t),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": "pods_test"}),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
				).Exists()
			},
		},
		{
			Name: "Exists_Name",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				return pods.NewPodAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithRequireT(t),
					assertion.WithResourceName("test-pod"),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
				).Exists()
			},
		},
		{
			Name: "IsReady_Labels",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				return pods.NewPodAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithRequireT(t),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": "pods_test"}),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(readyPodPath),
					),
				).Exists().IsReady()
			},
		},
		{
			Name: "IsReady_Name",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				return pods.NewPodAssertion(
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithRequireT(t),
					assertion.WithResourceName("test-pod"),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(readyPodPath),
					),
				).Exists().IsReady()
			},
		},
	}

	testhelpers.TestFailingAsserts(t, testEnv, asserts...)
}

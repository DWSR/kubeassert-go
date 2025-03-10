package pods_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/DWSR/kubeassert-go/internal/assertion"
	helpers "github.com/DWSR/kubeassert-go/internal/assertionhelpers"
	"github.com/DWSR/kubeassert-go/internal/pods"
	"github.com/DWSR/kubeassert-go/internal/testhelpers"
)

func Test_3Pod_Success(t *testing.T) {
	testCases := []struct {
		name      string
		assertion pods.PodAssertion
	}{
		{
			name: "AtLeastNExist",
			assertion: pods.NewPodAssertion(
				assertion.WithResourceNamespaceFromTestEnv(),
				assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": "pods_test"}),
				assertion.WithSetup(
					helpers.CreateResourceFromPathWithNamespaceFromEnv(deployPath),
				),
			).AtLeastNExist(2),
		},
		{
			name: "ExactlyNExist",
			assertion: pods.NewPodAssertion(
				assertion.WithResourceNamespaceFromTestEnv(),
				assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": "pods_test"}),
				assertion.WithSetup(
					helpers.CreateResourceFromPathWithNamespaceFromEnv(deployPath),
				),
			).ExactlyNExist(3),
		},
		{
			name: "AtLeastNAreReady",
			assertion: pods.NewPodAssertion(
				assertion.WithResourceNamespaceFromTestEnv(),
				assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": "pods_test"}),
				assertion.WithSetup(
					helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath),
					helpers.CreateResourceFromPathWithNamespaceFromEnv(deployPath),
				),
			).AtLeastNExist(2).AtLeastNAreReady(2),
		},
		{
			name: "ExactlyNAreReady",
			assertion: pods.NewPodAssertion(
				assertion.WithResourceNamespaceFromTestEnv(),
				assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": "pods_test"}),
				assertion.WithSetup(
					helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath),
					helpers.CreateResourceFromPathWithNamespaceFromEnv(deployPath),
				),
			).ExactlyNExist(3).ExactlyNAreReady(3),
		},
	}

	features := make([]features.Feature, 0)

	for _, a := range testCases {
		features = append(features, assertion.AsFeature(a.assertion))
	}

	testEnv.TestInParallel(t, features...)
}

func Test_3Pod_Fail(t *testing.T) {
	testCases := []struct {
		name             string
		failingAssertion func(require.TestingT) pods.PodAssertion
	}{
		{
			name: "AtLeastNExist",
			failingAssertion: func(t require.TestingT) pods.PodAssertion {
				return pods.NewPodAssertion(
					assertion.WithRequireT(t),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": "pods_test"}),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(deployPath),
					),
				).AtLeastNExist(4)
			},
		},
		{
			name: "ExactlyNExist",
			failingAssertion: func(t require.TestingT) pods.PodAssertion {
				return pods.NewPodAssertion(
					assertion.WithRequireT(t),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": "pods_test"}),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(deployPath),
					),
				).ExactlyNExist(2)
			},
		},
		{
			name: "AtLeastNAreReady",
			failingAssertion: func(t require.TestingT) pods.PodAssertion {
				return pods.NewPodAssertion(
					assertion.WithRequireT(t),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithTimeout(10*time.Second),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": "pods_test"}),

					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(configPath),
						helpers.CreateResourceFromPathWithNamespaceFromEnv(deployPath),
					),
				).AtLeastNExist(3).AtLeastNAreReady(4)
			},
		},
		{
			name: "ExactlyNAreReady",
			failingAssertion: func(t require.TestingT) pods.PodAssertion {
				return pods.NewPodAssertion(
					assertion.WithRequireT(t),
					assertion.WithResourceNamespaceFromTestEnv(),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithResourceLabels(map[string]string{"app.kubernetes.io/name": "pods_test"}),
					assertion.WithSetup(
						helpers.CreateResourceFromPathWithNamespaceFromEnv(deployPath),
					),
				).ExactlyNExist(3).ExactlyNAreReady(3)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockT := &testhelpers.MockT{}
			testEnv.Test(t, assertion.AsFeature(tc.failingAssertion(mockT)))
			assert.True(t, mockT.Failed)
		})
	}
}

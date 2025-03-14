// testhelpers contains helper functions specifically for testing assertion functionality. Any code that could be
// useful to consumers of kubeassert should go in the assertionhelpers package.
package testhelpers

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/e2e-framework/klient/decoder"
	"sigs.k8s.io/e2e-framework/klient/k8s"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/envfuncs"
	e2etypes "sigs.k8s.io/e2e-framework/pkg/types"
)

type (
	// MockT implements the require.TestingT interface and enables the detection of failing assert/require statements.
	// This enables testing assertions for expected failures.
	MockT struct {
		Failed bool
	}
)

const randomNamespaceNameLength = 20

// Errorf is a no-op function that satisfies the require.TestingT interface.
func (*MockT) Errorf(_ string, _ ...interface{}) {}

// FailNow sets the Failed field to true, indicating that a failing assertion was detected.
func (t *MockT) FailNow() {
	t.Failed = true
}

// CreateNamespaceBeforeEachFeature is a FeatureEnvFunc that creates a namespace. This helps run
// Features in isolation without requiring that features handle their own setup or teardown.
func CreateNamespaceBeforeEachFeature(namespaceName string) e2etypes.FeatureEnvFunc {
	return func(ctx context.Context, cfg *envconf.Config, _ *testing.T, _ e2etypes.Feature) (context.Context, error) {
		namespace := corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespaceName}}

		client, err := cfg.NewClient()
		if err != nil {
			return ctx, err
		}

		if err := client.Resources().Create(ctx, &namespace); err != nil {
			return ctx, err
		}

		_ = cfg.WithNamespace(namespaceName)

		return context.WithValue(ctx, envfuncs.NamespaceContextKey(namespaceName), namespace), nil
	}
}

// DeleteNamespaceBeforeEachFeature is a FeatureEnvFunc that deletes a namespace. This helps run
// Features in isolation without requiring that features handle their own setup or teardown.
func DeleteNamespaceBeforeEachFeature(namespaceName string) e2etypes.FeatureEnvFunc {
	return func(ctx context.Context, cfg *envconf.Config, _ *testing.T, _ e2etypes.Feature) (context.Context, error) {
		var namespace corev1.Namespace

		nsFromCtx := ctx.Value(envfuncs.NamespaceContextKey(namespaceName))
		if nsFromCtx == nil {
			namespace = corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespaceName}}
		} else {
			var ok bool

			namespace, ok = nsFromCtx.(corev1.Namespace)

			if !ok {
				panic("namespace is not of type corev1.Namespace")
			}
		}

		client, err := cfg.NewClient()
		if err != nil {
			return ctx, err
		}

		if err := client.Resources().Delete(ctx, &namespace); err != nil {
			return ctx, err
		}

		cfg.WithNamespace("")

		return ctx, nil
	}
}

// CreateRandomNamespaceBeforeEachFeature is a FeatureEnvFunc that creates a namespace with a random name.
func CreateRandomNamespaceBeforeEachFeature() e2etypes.FeatureEnvFunc {
	return func(ctx context.Context, cfg *envconf.Config, t *testing.T, feat e2etypes.Feature) (context.Context, error) {
		return CreateNamespaceBeforeEachFeature(envconf.RandomName("test", randomNamespaceNameLength))(ctx, cfg, t, feat)
	}
}

// DeleteRandomNamespaceAfterEachFeature is a FeatureEnvFunc that deletes the namespace set in the test environment
// after each feature.
func DeleteRandomNamespaceAfterEachFeature() e2etypes.FeatureEnvFunc {
	return func(ctx context.Context, cfg *envconf.Config, t *testing.T, feat e2etypes.Feature) (context.Context, error) {
		nsName := cfg.Namespace()

		return DeleteNamespaceBeforeEachFeature(nsName)(ctx, cfg, t, feat)
	}
}

// MutateResourceName is a DecodeOption that mutates the name of a resource.
func MutateResourceName(resourceName string) decoder.DecodeOption {
	return decoder.MutateOption(func(obj k8s.Object) error {
		obj.SetName(resourceName)

		return nil
	})
}

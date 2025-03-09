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
	MockT struct {
		Failed bool
	}
)

func (t *MockT) Errorf(_ string, _ ...interface{}) {}

func (t *MockT) FailNow() {
	t.Failed = true
}

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

		cfg = cfg.WithNamespace(namespaceName)

		return context.WithValue(ctx, envfuncs.NamespaceContextKey(namespaceName), namespace), nil
	}
}

func DeleteNamespaceBeforeEachFeature(namespaceName string) e2etypes.FeatureEnvFunc {
	return func(ctx context.Context, cfg *envconf.Config, _ *testing.T, _ e2etypes.Feature) (context.Context, error) {
		var ns corev1.Namespace

		nsFromCtx := ctx.Value(envfuncs.NamespaceContextKey(namespaceName))
		if nsFromCtx == nil {
			ns = corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespaceName}}
		} else {
			ns = nsFromCtx.(corev1.Namespace)
		}

		client, err := cfg.NewClient()
		if err != nil {
			return ctx, err
		}

		if err := client.Resources().Delete(ctx, &ns); err != nil {
			return ctx, err
		}

		cfg.WithNamespace("")

		return ctx, nil
	}
}

func CreateRandomNamespaceBeforeEachFeature() e2etypes.FeatureEnvFunc {
	return func(ctx context.Context, cfg *envconf.Config, t *testing.T, feat e2etypes.Feature) (context.Context, error) {
		return CreateNamespaceBeforeEachFeature(envconf.RandomName("test", 20))(ctx, cfg, t, feat)
	}
}

func DeleteRandomNamespaceAfterEachFeature() e2etypes.FeatureEnvFunc {
	return func(ctx context.Context, cfg *envconf.Config, t *testing.T, feat e2etypes.Feature) (context.Context, error) {
		nsName := cfg.Namespace()

		return DeleteNamespaceBeforeEachFeature(nsName)(ctx, cfg, t, feat)
	}
}

func MutateResourceName(resourceName string) decoder.DecodeOption {
	return decoder.MutateOption(func(obj k8s.Object) error {
		obj.SetName(resourceName)

		return nil
	})
}

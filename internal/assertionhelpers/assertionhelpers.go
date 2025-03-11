package assertionhelpers

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	apimachinerywait "k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
	"sigs.k8s.io/e2e-framework/klient/decoder"
	"sigs.k8s.io/e2e-framework/klient/k8s/resources"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
	e2etypes "sigs.k8s.io/e2e-framework/pkg/types"
	"sigs.k8s.io/kustomize/api/krusty"
	kusttypes "sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/filesys"

	"github.com/DWSR/kubeassert-go/internal/assertion"
)

type (
	// IntCompareFunc is a function that compares two integers and returns a boolean.
	IntCompareFunc func(int, int) bool

	// ConditionFunc is a function that returns a boolean based on a condition being satisfied.
	ConditionFunc = apimachinerywait.ConditionWithContextFunc

	// ConditionFuncFactory is a function that returns a ConditionFunc.
	ConditionFuncFactory = func(
		require.TestingT, assertion.Assertion, *envconf.Config, int, IntCompareFunc, IntCompareFunc,
	) ConditionFunc

	// StepFunc is a function that performs a step in a test.
	StepFunc = features.Func
)

var (
	// IntCompareFuncLessThan is a function that compares two integers and returns true if the first integer is less
	// than the second.
	IntCompareFuncLessThan IntCompareFunc = func(a, b int) bool { return a < b }

	// IntCompareFuncLessThanOrEqualTo is a function that compares two integers and returns true if the first integer is
	// less than or equal to the second.
	IntCompareFuncLessThanOrEqualTo IntCompareFunc = func(a, b int) bool { return a <= b }

	// IntCompareFuncEqualTo is a function that compares two integers and returns true if the first integer is equal to
	// the second.
	IntCompareFuncEqualTo IntCompareFunc = func(a, b int) bool { return a == b }

	// IntCompareFuncGreaterThan is a function that compares two integers and returns true if the first integer is
	// greater than the second.
	IntCompareFuncGreaterThan IntCompareFunc = func(a, b int) bool { return a > b }

	// IntCompareFuncGreaterThanOrEqualTo is a function that compares two integers and returns true if the first integer
	// is greater than or equal to the second.
	IntCompareFuncGreaterThanOrEqualTo IntCompareFunc = func(a, b int) bool { return a >= b }

	// IntCompareFuncNotEqualTo is a function that compares two integers and returns true if the first integer is not
	// equal to the second.
	IntCompareFuncNotEqualTo IntCompareFunc = func(a, b int) bool { return a != b }
)

// CreateResourceFromPathWithNamespaceFromEnv creates a resource from a file at the provided path and sets the
// resource's namespace to the one provided in the environment configuration.
func CreateResourceFromPathWithNamespaceFromEnv(
	resourcePath string,
	decoderOpts ...decoder.DecodeOption,
) e2etypes.StepFunc {
	return func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		nsName := cfg.Namespace()

		res, err := resources.New(cfg.Client().RESTConfig())
		require.NoError(t, err)

		file, err := os.Open(filepath.Clean(resourcePath))
		require.NoError(t, err)

		defer func() { _ = file.Close() }()

		err = decoder.DecodeEach(
			ctx,
			file,
			decoder.CreateHandler(res),
			append(decoderOpts, decoder.MutateNamespace(nsName))...,
		)
		require.NoError(t, err)

		return ctx
	}
}

func CreateResourceFromPath(resourcePath string, decoderOpts ...decoder.DecodeOption) e2etypes.StepFunc {
	return func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		r, err := resources.New(cfg.Client().RESTConfig())
		require.NoError(t, err)

		file, err := os.Open(filepath.Clean(resourcePath))
		require.NoError(t, err)
		defer func() { _ = file.Close() }()

		err = decoder.DecodeEach(ctx, file, decoder.CreateHandler(r), decoderOpts...)
		require.NoError(t, err)

		return ctx
	}
}

func DeleteResourceFromPath(resourcePath string, decoderOpts ...decoder.DecodeOption) e2etypes.StepFunc {
	return func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		r, err := resources.New(cfg.Client().RESTConfig())
		require.NoError(t, err)

		file, err := os.Open(filepath.Clean(resourcePath))
		require.NoError(t, err)
		defer func() { _ = file.Close() }()

		err = decoder.DecodeEach(ctx, file, decoder.DeleteHandler(r), decoderOpts...)
		require.NoError(t, err)

		return ctx
	}
}

func Sleep(sleepTime time.Duration) e2etypes.StepFunc {
	return func(ctx context.Context, _ *testing.T, _ *envconf.Config) context.Context {
		time.Sleep(sleepTime)

		return ctx
	}
}

func DynamicClientFromEnvconf(t require.TestingT, cfg *envconf.Config) *dynamic.DynamicClient {
	klient, err := cfg.NewClient()
	require.NoError(t, err)

	client, err := dynamic.NewForConfig(klient.RESTConfig())
	require.NoError(t, err)

	return client
}

func RequireTIfNotNil(testingT *testing.T, requireT require.TestingT) require.TestingT {
	if requireT != nil {
		return requireT
	}

	return testingT
}

func TestAssertions(t *testing.T, testEnv env.Environment, assertions ...assertion.Assertion) {
	tests := make([]e2etypes.Feature, 0, len(assertions))

	for _, assert := range assertions {
		tests = append(tests, assertion.AsFeature(assert))
	}

	testEnv.Test(t, tests...)
}

func ApplyKustomization(kustDir string) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		diskFS := filesys.MakeFsOnDisk()
		opts := krusty.MakeDefaultOptions()
		opts.PluginConfig.HelmConfig = kusttypes.HelmConfig{
			Enabled: true,
			Command: "helm",
			Debug:   false,
		}
		opts.PluginConfig.FnpLoadingOptions.Network = true
		opts.LoadRestrictions = kusttypes.LoadRestrictionsNone
		opts.Reorder = krusty.ReorderOptionLegacy
		kust := krusty.MakeKustomizer(opts)

		slog.Debug("rendering kustomization")

		resMap, err := kust.Run(diskFS, kustDir)
		if err != nil {
			return ctx, err
		}

		slog.Debug("creating client")

		klient, err := cfg.NewClient()
		if err != nil {
			return ctx, err
		}

		client, err := dynamic.NewForConfig(klient.RESTConfig())
		if err != nil {
			return ctx, err
		}

		slog.Debug("applying kustomization")

		for _, res := range resMap.Resources() {
			// Do this inside the loop to account for new CRDs, etc. that get applied
			slog.Debug("creating resource mapper")

			discoveryClient, err := discovery.NewDiscoveryClientForConfig(klient.RESTConfig())
			if err != nil {
				return ctx, err
			}

			gr, err := restmapper.GetAPIGroupResources(discoveryClient)
			if err != nil {
				return ctx, err
			}

			restMapper := restmapper.NewDiscoveryRESTMapper(gr)

			slog.Debug("transmuting resMap resource to unstructured")

			yamlBytes, err := res.AsYAML()
			if err != nil {
				return ctx, err
			}

			obj := &unstructured.Unstructured{}

			decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(yamlBytes), len(yamlBytes))
			if err := decoder.Decode(obj); err != nil {
				return ctx, err
			}

			gvk := obj.GroupVersionKind()

			mapping, err := restMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
			if err != nil {
				return ctx, err
			}

			slog.Debug("applying resource", "kind", obj.GetKind(), "name", obj.GetName(), "gvr", mapping.Resource)

			var resourceClient dynamic.ResourceInterface

			switch mapping.Scope.Name() {
			case meta.RESTScopeNameNamespace:
				resourceClient = client.Resource(mapping.Resource).Namespace(obj.GetNamespace())
			case meta.RESTScopeNameRoot:
				resourceClient = client.Resource(mapping.Resource)
			}

			_, err = resourceClient.Apply(ctx, obj.GetName(), obj, metav1.ApplyOptions{
				Force:        true,
				FieldManager: "e2e-test",
			})
			if err != nil {
				return ctx, err
			}
		}

		return ctx, nil
	}
}

// WaitForCondition waits for a conditionFunc to be satisfied (i.e. return true) based on the timeout and interval set
// on the Assertion.
func WaitForCondition(
	ctx context.Context,
	assert assertion.Assertion,
	conditionFunc apimachinerywait.ConditionWithContextFunc,
) error {
	return wait.For(
		conditionFunc,
		wait.WithContext(ctx),
		wait.WithTimeout(assert.GetTimeout()),
		wait.WithInterval(assert.GetInterval()),
		wait.WithImmediate(),
	)
}

// AsStepFunc returns a StepFunc that waits for a condition to be satisfied based on the provided ConditionFuncFactory.
// The count parameter is used to determine the number of resources that should satisfy the condition. The itemCountFn
// parameter is used to evalute the number of resources under consideration to satisfy the condition. The resultFn
// is used to evaluate the number of items that satisfy the condition relative to the count.
func AsStepFunc(
	assert assertion.Assertion,
	conditionFactory ConditionFuncFactory,
	count int,
	itemCountFn, resultFn IntCompareFunc,
) StepFunc {
	return func(ctx context.Context, testingT *testing.T, cfg *envconf.Config) context.Context {
		t := RequireTIfNotNil(testingT, assert.GetRequireT())

		require.NoError(t, WaitForCondition(ctx, assert, conditionFactory(t, assert, cfg, count, itemCountFn, resultFn)))

		return ctx
	}
}

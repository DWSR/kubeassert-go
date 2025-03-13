package deployments

import (
	"context"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/e2e-framework/pkg/envconf"

	"github.com/DWSR/kubeassert-go/internal/assertion"
	helpers "github.com/DWSR/kubeassert-go/internal/assertionhelpers"
)

func getDeployments(
	ctx context.Context,
	t require.TestingT,
	cfg *envconf.Config,
	listOpts metav1.ListOptions,
) (appsv1.DeploymentList, error) {
	client := helpers.DynamicClientFromEnvconf(t, cfg)

	var deploys appsv1.DeploymentList

	list, err := client.
		Resource(appsv1.SchemeGroupVersion.WithResource("deployments")).
		List(ctx, listOpts)
	if err != nil {
		return deploys, err
	}

	err = runtime.DefaultUnstructuredConverter.FromUnstructured(list.UnstructuredContent(), &deploys)
	if err != nil {
		return deploys, err
	}

	return deploys, nil
}

func exist() helpers.ConditionFuncFactory {
	return func(
		t require.TestingT,
		assert assertion.Assertion,
		cfg *envconf.Config,
		count int,
		itemCountFn, _ helpers.IntCompareFunc,
	) helpers.ConditionFunc {
		return func(ctx context.Context) (bool, error) {
			deployments, err := getDeployments(ctx, t, cfg, assert.ListOptions(cfg))
			require.NoError(t, err)

			return itemCountFn(len(deployments.Items), count), nil
		}
	}
}

func areAvailable() helpers.ConditionFuncFactory {
	return func(
		t require.TestingT,
		assert assertion.Assertion,
		cfg *envconf.Config,
		count int,
		itemCountFn, resultFn helpers.IntCompareFunc,
	) helpers.ConditionFunc {
		return func(ctx context.Context) (bool, error) {
			deployments, err := getDeployments(ctx, t, cfg, assert.ListOptions(cfg))
			require.NoError(t, err)

			if itemCountFn(len(deployments.Items), count) {
				return false, nil
			}

			availableCount := 0

			for _, deploy := range deployments.Items {
				for _, condition := range deploy.Status.Conditions {
					if condition.Type == appsv1.DeploymentAvailable && condition.Status == corev1.ConditionTrue {
						availableCount++
					}
				}
			}

			return resultFn(availableCount, count), nil
		}
	}
}

func areSystemClusterCritical() helpers.ConditionFuncFactory {
	return func(
		t require.TestingT,
		assert assertion.Assertion,
		cfg *envconf.Config,
		count int,
		itemCountFn, resultFn helpers.IntCompareFunc,
	) helpers.ConditionFunc {
		return func(ctx context.Context) (bool, error) {
			deployments, err := getDeployments(ctx, t, cfg, assert.ListOptions(cfg))
			require.NoError(t, err)

			if itemCountFn(len(deployments.Items), count) {
				return false, nil
			}

			systemClusterCriticalCount := 0

			for _, deploy := range deployments.Items {
				if deploy.Spec.Template.Spec.PriorityClassName == "system-cluster-critical" {
					systemClusterCriticalCount++
				}
			}

			return resultFn(systemClusterCriticalCount, count), nil
		}
	}
}

func haveNoCPULimits() helpers.ConditionFuncFactory {
	return func(
		t require.TestingT,
		assert assertion.Assertion,
		cfg *envconf.Config,
		count int,
		itemCountFn, resultFn helpers.IntCompareFunc,
	) helpers.ConditionFunc {
		return func(ctx context.Context) (bool, error) {
			deployments, err := getDeployments(ctx, t, cfg, assert.ListOptions(cfg))
			require.NoError(t, err)

			if itemCountFn(len(deployments.Items), count) {
				return false, nil
			}

			hasNoCPULimits := 0

			for _, deploy := range deployments.Items {
				allContainersHaveNoCPULimits := true

				for _, container := range deploy.Spec.Template.Spec.Containers {
					if !container.Resources.Limits.Cpu().IsZero() {
						allContainersHaveNoCPULimits = false

						break
					}
				}

				if allContainersHaveNoCPULimits {
					hasNoCPULimits++
				}
			}

			return resultFn(hasNoCPULimits, count), nil
		}
	}
}

func haveMemoryLimitsEqualToRequests() helpers.ConditionFuncFactory {
	return func(
		t require.TestingT,
		assert assertion.Assertion,
		cfg *envconf.Config,
		count int,
		itemCountFn, resultFn helpers.IntCompareFunc,
	) helpers.ConditionFunc {
		return func(ctx context.Context) (bool, error) {
			deployments, err := getDeployments(ctx, t, cfg, assert.ListOptions(cfg))
			require.NoError(t, err)

			if itemCountFn(len(deployments.Items), count) {
				return false, nil
			}

			hasMemoryLimitsEqualToRequests := 0

			for _, deploy := range deployments.Items {
				allContainersHaveMemoryLimitsEqualToRequests := true

				for _, container := range deploy.Spec.Template.Spec.Containers {
					memoryRequests := container.Resources.Requests.Memory()
					memoryLimits := container.Resources.Limits.Memory()

					if !cmp.Equal(memoryLimits, memoryRequests) {
						allContainersHaveMemoryLimitsEqualToRequests = false

						break
					}
				}

				if allContainersHaveMemoryLimitsEqualToRequests {
					hasMemoryLimitsEqualToRequests++
				}
			}

			return resultFn(hasMemoryLimitsEqualToRequests, count), nil
		}
	}
}

func haveMemoryLimits() helpers.ConditionFuncFactory {
	return func(
		t require.TestingT,
		assert assertion.Assertion,
		cfg *envconf.Config,
		count int,
		itemCountFn, resultFn helpers.IntCompareFunc,
	) helpers.ConditionFunc {
		return func(ctx context.Context) (bool, error) {
			deployments, err := getDeployments(ctx, t, cfg, assert.ListOptions(cfg))
			require.NoError(t, err)

			if itemCountFn(len(deployments.Items), count) {
				return false, nil
			}

			hasMemoryLimits := 0

			for _, deploy := range deployments.Items {
				allContainersHaveMemoryLimits := true

				for _, container := range deploy.Spec.Template.Spec.Containers {
					if container.Resources.Limits.Memory().IsZero() {
						allContainersHaveMemoryLimits = false

						break
					}
				}

				if allContainersHaveMemoryLimits {
					hasMemoryLimits++
				}
			}

			return resultFn(hasMemoryLimits, count), nil
		}
	}
}

func haveMemoryRequests() helpers.ConditionFuncFactory {
	return func(
		t require.TestingT,
		assert assertion.Assertion,
		cfg *envconf.Config,
		count int,
		itemCountFn, resultFn helpers.IntCompareFunc,
	) helpers.ConditionFunc {
		return func(ctx context.Context) (bool, error) {
			deployments, err := getDeployments(ctx, t, cfg, assert.ListOptions(cfg))
			require.NoError(t, err)

			if itemCountFn(len(deployments.Items), count) {
				return false, nil
			}

			hasMemoryRequests := 0

			for _, deploy := range deployments.Items {
				allContainersHaveMemoryRequests := true

				for _, container := range deploy.Spec.Template.Spec.Containers {
					if container.Resources.Requests.Memory().IsZero() {
						allContainersHaveMemoryRequests = false

						break
					}
				}

				if allContainersHaveMemoryRequests {
					hasMemoryRequests++
				}
			}

			return resultFn(hasMemoryRequests, count), nil
		}
	}
}

func haveCPURequests() helpers.ConditionFuncFactory {
	return func(
		t require.TestingT,
		assert assertion.Assertion,
		cfg *envconf.Config,
		count int,
		itemCountFn, resultFn helpers.IntCompareFunc,
	) helpers.ConditionFunc {
		return func(ctx context.Context) (bool, error) {
			deployments, err := getDeployments(ctx, t, cfg, assert.ListOptions(cfg))
			require.NoError(t, err)

			if itemCountFn(len(deployments.Items), count) {
				return false, nil
			}

			hasCPURequests := 0

			for _, deploy := range deployments.Items {
				allContainersHaveCPURequests := true

				for _, container := range deploy.Spec.Template.Spec.Containers {
					if container.Resources.Requests.Cpu().IsZero() {
						allContainersHaveCPURequests = false

						break
					}
				}

				if allContainersHaveCPURequests {
					hasCPURequests++
				}
			}

			return resultFn(hasCPURequests, count), nil
		}
	}
}

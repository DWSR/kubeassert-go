package deployments

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/DWSR/kubeassert-go/internal/assertion"
	helpers "github.com/DWSR/kubeassert-go/internal/assertionhelpers"
)

type DeploymentAssertion struct {
	assertion.Assertion
}

func (da DeploymentAssertion) clone() DeploymentAssertion {
	return DeploymentAssertion{
		Assertion: assertion.CloneAssertion(da.Assertion),
	}
}

func (da DeploymentAssertion) Exists() DeploymentAssertion {
	return da.ExactlyNExist(1)
}

func (da DeploymentAssertion) ExactlyNExist(count int) DeploymentAssertion {
	fn := func(ctx context.Context, testingT *testing.T, cfg *envconf.Config) context.Context {
		t := helpers.RequireTIfNotNil(testingT, da.GetRequireT())
		conditionFunc := func(ctx context.Context) (bool, error) {
			deploys, err := da.getDeployments(ctx, t, cfg)
			require.NoError(t, err)

			return len(deploys.Items) == count, nil
		}

		require.NoError(t, da.WaitForCondition(ctx, conditionFunc))

		return ctx
	}

	res := da.clone()
	res.SetBuilder(res.GetBuilder().Assess("exactlyNExist", fn))

	return res
}

func (da DeploymentAssertion) AtLeastNExist(count int) DeploymentAssertion {
	fn := func(ctx context.Context, testingT *testing.T, cfg *envconf.Config) context.Context {
		t := helpers.RequireTIfNotNil(testingT, da.GetRequireT())
		conditionFunc := func(ctx context.Context) (bool, error) {
			deploys, err := da.getDeployments(ctx, t, cfg)
			require.NoError(t, err)

			return len(deploys.Items) >= count, nil
		}

		require.NoError(t, da.WaitForCondition(ctx, conditionFunc))

		return ctx
	}

	res := da.clone()
	res.SetBuilder(res.GetBuilder().Assess("atLeastNExist", fn))

	return res
}

func (da DeploymentAssertion) getDeployments(ctx context.Context, t require.TestingT, cfg *envconf.Config) (appsv1.DeploymentList, error) {
	client := helpers.DynamicClientFromEnvconf(t, cfg)

	var deploys appsv1.DeploymentList

	list, err := client.
		Resource(appsv1.SchemeGroupVersion.WithResource("deployments")).
		List(ctx, da.ListOptions(cfg))
	if err != nil {
		return deploys, err
	}

	err = runtime.DefaultUnstructuredConverter.FromUnstructured(list.UnstructuredContent(), &deploys)
	if err != nil {
		return deploys, err
	}

	return deploys, nil
}

func (da DeploymentAssertion) IsAvailable() DeploymentAssertion {
	return da.ExactlyNAreAvailable(1)
}

func (da DeploymentAssertion) ExactlyNAreAvailable(count int) DeploymentAssertion {
	stepFn := func(ctx context.Context, testingT *testing.T, cfg *envconf.Config) context.Context {
		t := helpers.RequireTIfNotNil(testingT, da.GetRequireT())

		conditionFunc := func(ctx context.Context) (bool, error) {
			deploys, err := da.getDeployments(ctx, t, cfg)
			require.NoError(t, err)

			if len(deploys.Items) != count {
				return false, nil
			}

			availableCount := 0

			for _, deploy := range deploys.Items {
				for _, condition := range deploy.Status.Conditions {
					if condition.Type == appsv1.DeploymentAvailable && condition.Status == corev1.ConditionTrue {
						availableCount += 1
					}
				}
			}

			return availableCount == count, nil
		}

		require.NoError(t, da.WaitForCondition(ctx, conditionFunc))

		return ctx
	}

	res := da.clone()
	res.SetBuilder(res.GetBuilder().Assess("exactlyNAreAvailable", stepFn))

	return res
}

func (da DeploymentAssertion) AtLeastNAreAvailable(count int) DeploymentAssertion {
	stepFn := func(ctx context.Context, testingT *testing.T, cfg *envconf.Config) context.Context {
		t := helpers.RequireTIfNotNil(testingT, da.GetRequireT())

		conditionFunc := func(ctx context.Context) (bool, error) {
			deploys, err := da.getDeployments(ctx, t, cfg)
			require.NoError(t, err)

			if len(deploys.Items) < count {
				return false, nil
			}

			availableCount := 0

			for _, deploy := range deploys.Items {
				for _, condition := range deploy.Status.Conditions {
					if condition.Type == appsv1.DeploymentAvailable && condition.Status == corev1.ConditionTrue {
						availableCount++
					}
				}
			}

			return availableCount >= count, nil
		}

		require.NoError(t, da.WaitForCondition(ctx, conditionFunc))

		return ctx
	}

	res := da.clone()
	res.SetBuilder(res.GetBuilder().Assess("atLeastNAreAvailable", stepFn))

	return res
}

func (da DeploymentAssertion) IsSystemClusterCritical() DeploymentAssertion {
	return da.ExactlyNAreSystemClusterCritical(1)
}

func (da DeploymentAssertion) ExactlyNAreSystemClusterCritical(count int) DeploymentAssertion {
	stepFn := func(ctx context.Context, testingT *testing.T, cfg *envconf.Config) context.Context {
		t := helpers.RequireTIfNotNil(testingT, da.GetRequireT())

		conditionFunc := func(ctx context.Context) (bool, error) {
			deploys, err := da.getDeployments(ctx, t, cfg)
			require.NoError(t, err)

			if len(deploys.Items) < count {
				return false, nil
			}

			systemClusterCriticalCount := 0

			for _, deploy := range deploys.Items {
				if deploy.Spec.Template.Spec.PriorityClassName == "system-cluster-critical" {
					systemClusterCriticalCount++
				}
			}

			return systemClusterCriticalCount == count, nil
		}

		require.NoError(t, da.WaitForCondition(ctx, conditionFunc))

		return ctx
	}

	res := da.clone()
	res.SetBuilder(res.GetBuilder().Assess("exactlyNAreSystemClusterCritical", stepFn))

	return res
}

func (da DeploymentAssertion) AtLeastNAreSystemClusterCritical(count int) DeploymentAssertion {
	stepFn := func(ctx context.Context, testingT *testing.T, cfg *envconf.Config) context.Context {
		t := helpers.RequireTIfNotNil(testingT, da.GetRequireT())

		conditionFunc := func(ctx context.Context) (bool, error) {
			deploys, err := da.getDeployments(ctx, t, cfg)
			require.NoError(t, err)

			if len(deploys.Items) < count {
				return false, nil
			}

			systemClusterCriticalCount := 0

			for _, deploy := range deploys.Items {
				if deploy.Spec.Template.Spec.PriorityClassName == "system-cluster-critical" {
					systemClusterCriticalCount++
				}
			}

			return systemClusterCriticalCount >= count, nil
		}

		require.NoError(t, da.WaitForCondition(ctx, conditionFunc))

		return ctx
	}

	res := da.clone()
	res.SetBuilder(res.GetBuilder().Assess("atLeastNAreSystemClusterCritical", stepFn))

	return res
}

func (da DeploymentAssertion) HasNoCPULimits() DeploymentAssertion {
	return da.ExactlyNHaveNoCPULimits(1)
}

func (da DeploymentAssertion) ExactlyNHaveNoCPULimits(count int) DeploymentAssertion {
	stepFn := func(ctx context.Context, testingT *testing.T, cfg *envconf.Config) context.Context {
		t := helpers.RequireTIfNotNil(testingT, da.GetRequireT())

		conditionFunc := func(ctx context.Context) (bool, error) {
			deploys, err := da.getDeployments(ctx, t, cfg)
			require.NoError(t, err)

			if len(deploys.Items) < count {
				return false, nil
			}

			hasNoCPULimits := 0

			for _, deploy := range deploys.Items {
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

			return hasNoCPULimits == count, nil
		}

		require.NoError(t, da.WaitForCondition(ctx, conditionFunc))

		return ctx
	}

	res := da.clone()
	res.SetBuilder(res.GetBuilder().Assess("exactlyNHaveNoCPULimits", stepFn))

	return res
}

func (da DeploymentAssertion) AtLeastNHaveNoCPULimits(count int) DeploymentAssertion {
	stepFn := func(ctx context.Context, testingT *testing.T, cfg *envconf.Config) context.Context {
		t := helpers.RequireTIfNotNil(testingT, da.GetRequireT())

		conditionFunc := func(ctx context.Context) (bool, error) {
			deploys, err := da.getDeployments(ctx, t, cfg)
			require.NoError(t, err)

			if len(deploys.Items) < count {
				return false, nil
			}

			hasNoCPULimits := 0

			for _, deploy := range deploys.Items {
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

			return hasNoCPULimits >= count, nil
		}

		require.NoError(t, da.WaitForCondition(ctx, conditionFunc))

		return ctx
	}

	res := da.clone()
	res.SetBuilder(res.GetBuilder().Assess("atLeastNHaveNoCPULimits", stepFn))

	return res
}

func (da DeploymentAssertion) HasMemoryLimitsEqualToRequests() DeploymentAssertion {
	return da.ExactlyNHaveMemoryLimitsEqualToRequests(1)
}

func (da DeploymentAssertion) ExactlyNHaveMemoryLimitsEqualToRequests(count int) DeploymentAssertion {
	stepFn := func(ctx context.Context, testingT *testing.T, cfg *envconf.Config) context.Context {
		t := helpers.RequireTIfNotNil(testingT, da.GetRequireT())

		conditionFunc := func(ctx context.Context) (bool, error) {
			deploys, err := da.getDeployments(ctx, t, cfg)
			require.NoError(t, err)

			if len(deploys.Items) < count {
				return false, nil
			}

			hasMemoryLimitsEqualToRequests := 0

			for _, deploy := range deploys.Items {
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

			return hasMemoryLimitsEqualToRequests == count, nil
		}

		require.NoError(t, da.WaitForCondition(ctx, conditionFunc))

		return ctx
	}

	res := da.clone()
	res.SetBuilder(res.GetBuilder().Assess("exactlyNHaveMemoryLimitsEqualToRequests", stepFn))

	return res
}

func (da DeploymentAssertion) AtLeastNHaveMemoryLimitsEqualToRequests(count int) DeploymentAssertion {
	stepFn := func(ctx context.Context, testingT *testing.T, cfg *envconf.Config) context.Context {
		t := helpers.RequireTIfNotNil(testingT, da.GetRequireT())

		conditionFunc := func(ctx context.Context) (bool, error) {
			deploys, err := da.getDeployments(ctx, t, cfg)
			require.NoError(t, err)

			if len(deploys.Items) < count {
				return false, nil
			}

			hasMemoryLimitsEqualToRequests := 0

			for _, deploy := range deploys.Items {
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

			return hasMemoryLimitsEqualToRequests >= count, nil
		}

		require.NoError(t, da.WaitForCondition(ctx, conditionFunc))

		return ctx
	}

	res := da.clone()
	res.SetBuilder(res.GetBuilder().Assess("atLeastNHaveMemoryLimitsEqualToRequests", stepFn))

	return res
}

func (da DeploymentAssertion) HasMemoryLimits() DeploymentAssertion {
	return da.ExactlyNHaveMemoryLimits(1)
}

func (da DeploymentAssertion) ExactlyNHaveMemoryLimits(count int) DeploymentAssertion {
	stepFn := func(ctx context.Context, testingT *testing.T, cfg *envconf.Config) context.Context {
		t := helpers.RequireTIfNotNil(testingT, da.GetRequireT())

		conditionFunc := func(ctx context.Context) (bool, error) {
			deploys, err := da.getDeployments(ctx, t, cfg)
			require.NoError(t, err)

			if len(deploys.Items) < count {
				return false, nil
			}

			hasMemoryLimits := 0

			for _, deploy := range deploys.Items {
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

			return hasMemoryLimits == count, nil
		}

		require.NoError(t, da.WaitForCondition(ctx, conditionFunc))

		return ctx
	}

	res := da.clone()
	res.SetBuilder(res.GetBuilder().Assess("exactlyNHaveMemoryLimits", stepFn))

	return res
}

func (da DeploymentAssertion) AtLeastNHaveMemoryLimits(count int) DeploymentAssertion {
	stepFn := func(ctx context.Context, testingT *testing.T, cfg *envconf.Config) context.Context {
		t := helpers.RequireTIfNotNil(testingT, da.GetRequireT())

		conditionFunc := func(ctx context.Context) (bool, error) {
			deploys, err := da.getDeployments(ctx, t, cfg)
			require.NoError(t, err)

			if len(deploys.Items) < count {
				return false, nil
			}

			hasMemoryLimits := 0

			for _, deploy := range deploys.Items {
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

			return hasMemoryLimits >= count, nil
		}

		require.NoError(t, da.WaitForCondition(ctx, conditionFunc))

		return ctx
	}

	res := da.clone()
	res.SetBuilder(res.GetBuilder().Assess("atLeastNHaveMemoryLimits", stepFn))

	return res
}

func (da DeploymentAssertion) HasMemoryRequests() DeploymentAssertion {
	return da.ExactlyNHaveMemoryRequests(1)
}

func (da DeploymentAssertion) ExactlyNHaveMemoryRequests(count int) DeploymentAssertion {
	stepFn := func(ctx context.Context, testingT *testing.T, cfg *envconf.Config) context.Context {
		t := helpers.RequireTIfNotNil(testingT, da.GetRequireT())

		conditionFunc := func(ctx context.Context) (bool, error) {
			deploys, err := da.getDeployments(ctx, t, cfg)
			require.NoError(t, err)

			if len(deploys.Items) < count {
				return false, nil
			}

			hasMemoryRequests := 0

			for _, deploy := range deploys.Items {
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

			return hasMemoryRequests == count, nil
		}

		require.NoError(t, da.WaitForCondition(ctx, conditionFunc))

		return ctx
	}

	res := da.clone()
	res.SetBuilder(res.GetBuilder().Assess("exactlyNHaveMemoryRequests", stepFn))

	return res
}

func (da DeploymentAssertion) AtLeastNHaveMemoryRequests(count int) DeploymentAssertion {
	stepFn := func(ctx context.Context, testingT *testing.T, cfg *envconf.Config) context.Context {
		t := helpers.RequireTIfNotNil(testingT, da.GetRequireT())

		conditionFunc := func(ctx context.Context) (bool, error) {
			deploys, err := da.getDeployments(ctx, t, cfg)
			require.NoError(t, err)

			if len(deploys.Items) < count {
				return false, nil
			}

			hasMemoryRequests := 0

			for _, deploy := range deploys.Items {
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

			return hasMemoryRequests >= count, nil
		}

		require.NoError(t, da.WaitForCondition(ctx, conditionFunc))

		return ctx
	}

	res := da.clone()
	res.SetBuilder(res.GetBuilder().Assess("atLeastNHaveMemoryRequests", stepFn))

	return res
}

func (da DeploymentAssertion) HasCPURequests() DeploymentAssertion {
	return da.ExactlyNHaveCPURequests(1)
}

func (da DeploymentAssertion) ExactlyNHaveCPURequests(count int) DeploymentAssertion {
	stepFn := func(ctx context.Context, testingT *testing.T, cfg *envconf.Config) context.Context {
		t := helpers.RequireTIfNotNil(testingT, da.GetRequireT())

		conditionFunc := func(ctx context.Context) (bool, error) {
			deploys, err := da.getDeployments(ctx, t, cfg)
			require.NoError(t, err)

			if len(deploys.Items) < count {
				return false, nil
			}

			hasCPURequests := 0

			for _, deploy := range deploys.Items {
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

			return hasCPURequests == count, nil
		}

		require.NoError(t, da.WaitForCondition(ctx, conditionFunc))

		return ctx
	}

	res := da.clone()
	res.SetBuilder(res.GetBuilder().Assess("exactlyNHaveCPURequests", stepFn))

	return res
}

func (da DeploymentAssertion) AtLeastNHaveCPURequests(count int) DeploymentAssertion {
	stepFn := func(ctx context.Context, testingT *testing.T, cfg *envconf.Config) context.Context {
		t := helpers.RequireTIfNotNil(testingT, da.GetRequireT())

		conditionFunc := func(ctx context.Context) (bool, error) {
			deploys, err := da.getDeployments(ctx, t, cfg)
			require.NoError(t, err)

			if len(deploys.Items) < count {
				return false, nil
			}

			hasCPURequests := 0

			for _, deploy := range deploys.Items {
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

			return hasCPURequests >= count, nil
		}

		require.NoError(t, da.WaitForCondition(ctx, conditionFunc))

		return ctx
	}

	res := da.clone()
	res.SetBuilder(res.GetBuilder().Assess("atLeastNHaveCPURequests", stepFn))

	return res
}

func NewDeploymentAssertion(opts ...assertion.AssertionOption) DeploymentAssertion {
	return DeploymentAssertion{
		Assertion: assertion.NewAssertion(
			append(
				[]assertion.AssertionOption{assertion.WithBuilder(features.New("Deployment").WithLabel("type", "deployment"))},
				opts...,
			)...,
		),
	}
}

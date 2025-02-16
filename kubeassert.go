package kubeassert

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
	e2etypes "sigs.k8s.io/e2e-framework/pkg/types"
)

func NamespaceIsRestricted(namespaceName string) e2etypes.Feature {
	return features.New("NamespaceIsRestricted").
		WithLabel("type", "namespace").
		AssessWithDescription(
			"restrictedNamespace",
			"Namespace should be restricted",
			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				var ns corev1.Namespace
				err := cfg.Client().Resources().Get(ctx, namespaceName, "", &ns)
				require.NoError(t, err)

				nsLabels := ns.GetLabels()

				assert.Contains(t, nsLabels, "pod-security.kubernetes.io/enforce")
				assert.Equal(t, "restricted", nsLabels["pod-security.kubernetes.io/enforce"])
				assert.Contains(t, nsLabels, "pod-security.kubernetes.io/audit")
				assert.Equal(t, "restricted", nsLabels["pod-security.kubernetes.io/audit"])

				return ctx
			}).
		Feature()
}

func NamespaceExists(namespaceName string) e2etypes.Feature {
	return features.New("NamespaceExists").
		WithLabel("type", "namespace").
		AssessWithDescription(
			"namespaceExists",
			"Namespace should exist",
			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				ns := &corev1.NamespaceList{
					Items: []corev1.Namespace{{ObjectMeta: metav1.ObjectMeta{Name: namespaceName}}},
				}

				err := wait.For(
					conditions.New(cfg.Client().Resources()).ResourcesFound(ns),
					wait.WithTimeout(3*time.Second),
					wait.WithImmediate(),
				)
				require.NoError(t, err)

				return ctx
			}).
		Feature()
}

func DeploymentExists(namespaceName, deploymentName string) e2etypes.Feature {
	return features.New("DeploymentExists").
		WithLabel("type", "deployment").
		AssessWithDescription(
			"deploymentExists",
			"Deployment should exist",
			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				var dep appsv1.Deployment

				err := cfg.Client().
					Resources("deployments").
					WithNamespace(namespaceName).
					Get(ctx, deploymentName, namespaceName, &dep)
				require.NoError(t, err)

				return ctx
			}).
		Feature()
}

func DeploymentAvailable(namespaceName, deploymentName string) e2etypes.Feature {
	return features.New("DeploymentAvailable").
		WithLabel("type", "deployment").
		AssessWithDescription(
			"deploymentAvailable",
			"Deployment should be available",
			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				err := wait.For(
					conditions.New(cfg.Client().Resources()).DeploymentAvailable(deploymentName, namespaceName),
					wait.WithTimeout(2*time.Minute),
					wait.WithImmediate(),
				)
				require.NoError(t, err)

				return ctx
			}).
		Feature()
}

func SecretExists(namespaceName, secretName string) e2etypes.Feature {
	return features.New("SecretExists").
		WithLabel("type", "secret").
		AssessWithDescription(
			"secretExists",
			"Secret should exist",
			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				var sec corev1.Secret

				err := cfg.Client().
					Resources("secrets").
					WithNamespace(namespaceName).
					Get(ctx, secretName, namespaceName, &sec)
				require.NoError(t, err)

				return ctx
			}).
		Feature()
}

func SecretHasContent(namespaceName, secretName string, content map[string]string) e2etypes.Feature {
	return features.New("SecretHasContent").
		WithLabel("type", "secret").
		AssessWithDescription(
			"secretHasContent",
			"Secret should have content",
			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				var sec corev1.Secret

				err := cfg.Client().
					Resources("secrets").
					WithNamespace(namespaceName).
					Get(ctx, secretName, namespaceName, &sec)
				require.NoError(t, err)

				for k, v := range content {
					secData, exists := sec.Data[k]

					require.True(t, exists)
					assert.Equal(t, v, string(secData))
				}

				return ctx
			}).
		Feature()
}

func DeploymentIsSystemClusterCritical(namespaceName, deploymentName string) e2etypes.Feature {
	return features.New("DeploymentIsSystemClusterCritical").
		WithLabel("type", "deployment").
		AssessWithDescription(
			"deploymentIsSystemClusterCritical",
			"Deployment should be system-cluster-critical priority",
			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				var deploy appsv1.Deployment

				err := cfg.Client().
					Resources("deployments").
					WithNamespace(namespaceName).
					Get(ctx, deploymentName, namespaceName, &deploy)
				require.NoError(t, err)

				assert.Equal(t, "system-cluster-critical", deploy.Spec.Template.Spec.PriorityClassName)

				return ctx
			}).
		Feature()
}

func DeploymentHasNoCPULimits(namespaceName, deploymentName string) e2etypes.Feature {
	return features.New("DeploymentHasNoCPULimits").
		WithLabel("type", "deployment").
		AssessWithDescription(
			"deploymentHasNoCPULimits",
			"Deployment should have no CPU limits",
			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				var deploy appsv1.Deployment

				err := cfg.Client().
					Resources("deployments").
					WithNamespace(namespaceName).
					Get(ctx, deploymentName, namespaceName, &deploy)
				require.NoError(t, err)

				for _, container := range deploy.Spec.Template.Spec.Containers {
					assert.True(t, container.Resources.Limits.Cpu().IsZero())
				}

				return ctx
			}).
		Feature()
}

func DeploymentHasMemoryLimitsEqualToRequests(namespaceName, deploymentName string) e2etypes.Feature {
	return features.New("DeploymentHasMemoryLimitsEqualToRequests").
		WithLabel("type", "deployment").
		AssessWithDescription(
			"deploymentHasMemoryLimitsEqualToRequests",
			"Deployment should have memory limits equal to requests",
			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				var deploy appsv1.Deployment

				err := cfg.Client().
					Resources("deployments").
					WithNamespace(namespaceName).
					Get(ctx, deploymentName, namespaceName, &deploy)
				require.NoError(t, err)

				for _, container := range deploy.Spec.Template.Spec.Containers {
					assert.NotNil(t, container.Resources.Limits.Memory())
					assert.Equal(t, container.Resources.Requests.Memory(), container.Resources.Limits.Memory())
				}

				return ctx
			}).
		Feature()
}

func DeploymentHasMemoryRequests(namespaceName, deploymentName string) e2etypes.Feature {
	return features.New("DeploymentHasMemoryRequests").
		WithLabel("type", "deployment").
		AssessWithDescription(
			"deploymentHasMemoryRequests",
			"Deployment should have memory requests",
			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				var deploy appsv1.Deployment

				err := cfg.Client().
					Resources("deployments").
					WithNamespace(namespaceName).
					Get(ctx, deploymentName, namespaceName, &deploy)
				require.NoError(t, err)

				for _, container := range deploy.Spec.Template.Spec.Containers {
					assert.False(t, container.Resources.Requests.Memory().IsZero())
				}

				return ctx
			}).
		Feature()
}

func DeploymentHasCPURequests(namespaceName, deploymentName string) e2etypes.Feature {
	return features.New("DeploymentHasCPURequests").
		WithLabel("type", "deployment").
		AssessWithDescription(
			"deploymentHasCPURequests",
			"Deployment should have CPU requests",
			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				var deploy appsv1.Deployment

				err := cfg.Client().
					Resources("deployments").
					WithNamespace(namespaceName).
					Get(ctx, deploymentName, namespaceName, &deploy)
				require.NoError(t, err)

				for _, container := range deploy.Spec.Template.Spec.Containers {
					assert.False(t, container.Resources.Requests.Cpu().IsZero())
				}

				return ctx
			}).
		Feature()
}

func CRDExists(crdName, crdVersion string) e2etypes.Feature {
	return features.New("CRDExists").
		WithLabel("type", "crd").
		AssessWithDescription(
			"crdExists",
			"CRD should exist",
			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				var crd extv1.CustomResourceDefinition

				klient, err := cfg.NewClient()
				require.NoError(t, err)

				client, err := dynamic.NewForConfig(klient.RESTConfig())
				require.NoError(t, err)

				unstructuredCRD, err := client.
					Resource(extv1.SchemeGroupVersion.WithResource("customresourcedefinitions")).
					Get(ctx, crdName, metav1.GetOptions{})
				require.NoError(t, err)

				err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredCRD.UnstructuredContent(), &crd)
				require.NoError(t, err)

				foundVersion := false
				for _, v := range crd.Spec.Versions {
					if crdVersion == v.Name {
						foundVersion = true
					}
				}

				assert.True(t, foundVersion)

				return ctx
			}).
		Feature()
}

func PodDisruptionBudgetExists(namespaceName, pdbName string) e2etypes.Feature {
	return features.New("PodDisruptionBudgetExists").
		WithLabel("type", "pdb").
		AssessWithDescription(
			"pdbExists",
			"PDB should exist",
			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				var pdb policyv1.PodDisruptionBudget

				err := cfg.Client().
					Resources("poddisruptionbudgets").
					WithNamespace(namespaceName).
					Get(ctx, pdbName, namespaceName, &pdb)
				require.NoError(t, err)

				return ctx
			}).
		Feature()
}

func PodDisruptionBudgetTargetsDeployment(namespaceName, pdbName, deployName string) e2etypes.Feature {
	return features.New("PodDisruptionBudgetTargetsDeployment").
		WithLabel("type", "pdb").
		AssessWithDescription(
			"pdbTargetsDeployment",
			"PDB should target deployment",
			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				var pdb policyv1.PodDisruptionBudget
				var deploy appsv1.Deployment

				err := cfg.Client().
					Resources("poddisruptionbudgets").
					WithNamespace(namespaceName).
					Get(ctx, pdbName, namespaceName, &pdb)
				require.NoError(t, err)

				err = cfg.Client().
					Resources("deployments").
					WithNamespace(namespaceName).
					Get(ctx, deployName, namespaceName, &deploy)
				require.NoError(t, err)

				for labelKey, labelValue := range pdb.Spec.Selector.MatchLabels {
					require.Equal(t, deploy.Spec.Selector.MatchLabels, labelKey)
					require.Equal(t, deploy.Spec.Selector.MatchLabels[labelKey], labelValue)
				}

				return ctx
			}).
		Feature()
}

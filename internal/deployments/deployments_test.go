package deployments_test

import (
	"log/slog"
	"os"
	"testing"

	"sigs.k8s.io/e2e-framework/klient/decoder"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/envfuncs"
	e2etypes "sigs.k8s.io/e2e-framework/pkg/types"
	"sigs.k8s.io/e2e-framework/support/kind"

	helpers "github.com/DWSR/kubeassert-go/internal/assertionhelpers"
	"github.com/DWSR/kubeassert-go/internal/testhelpers"
)

var testEnv env.Environment

func TestMain(m *testing.M) {
	kindClusterName := envconf.RandomName("kind", 16)

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	testEnv = env.New().
		Setup(
			envfuncs.CreateCluster(kind.NewProvider(), kindClusterName),
		).
		BeforeEachFeature(testhelpers.CreateRandomNamespaceBeforeEachFeature()).
		AfterEachFeature(testhelpers.DeleteRandomNamespaceAfterEachFeature()).
		Finish(
			envfuncs.DestroyCluster(kindClusterName),
		)

	os.Exit(testEnv.Run(m))
}

func generateDeploymentNames() []string {
	res := make([]string, 3)

	for i := range res {
		res[i] = envconf.RandomName("test", 20)
	}

	return res
}

func createGoodDeploys(deploymentNames []string) []e2etypes.StepFunc {
	if len(deploymentNames) == 0 {
		panic("must supply deployment names")
	}

	res := make([]e2etypes.StepFunc, len(deploymentNames))
	labelValue := deploymentNames[0]

	for i, deployName := range deploymentNames {
		res[i] = helpers.CreateResourceFromPathWithNamespaceFromEnv(
			deploymentPath,
			testhelpers.MutateResourceName(deployName),
			decoder.MutateLabels(map[string]string{"app.kubernetes.io/name": labelValue}),
		)
	}

	return res
}

func createBadDeploys(deploymentNames []string) []e2etypes.StepFunc {
	if len(deploymentNames) == 0 {
		panic("must supply deployment names")
	}

	res := make([]e2etypes.StepFunc, len(deploymentNames))
	labelValue := deploymentNames[0]

	for i, deployName := range deploymentNames {
		res[i] = helpers.CreateResourceFromPathWithNamespaceFromEnv(
			badDeploymentPath,
			testhelpers.MutateResourceName(deployName),
			decoder.MutateLabels(map[string]string{"app.kubernetes.io/name": labelValue}),
		)
	}

	return res
}

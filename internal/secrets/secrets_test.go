package secrets_test

import (
	"log/slog"
	"os"
	"testing"

	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/envfuncs"
	"sigs.k8s.io/e2e-framework/support/kind"

	"github.com/DWSR/kubeassert-go/internal/testhelpers"
)

//nolint:gosec
const secretPath = "./testdata/secret.yaml"

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

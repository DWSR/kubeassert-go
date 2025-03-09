package crds_test

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/envfuncs"
	"sigs.k8s.io/e2e-framework/support/kind"

	"github.com/DWSR/kubeassert-go/internal/assertion"
	helpers "github.com/DWSR/kubeassert-go/internal/assertionhelpers"
	"github.com/DWSR/kubeassert-go/internal/crds"
	"github.com/DWSR/kubeassert-go/internal/testhelpers"
)

const crdPath = "./testdata/crd.yaml"

var testEnv env.Environment

func TestMain(m *testing.M) {
	kindClusterName := envconf.RandomName("kind", 16)

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	testEnv = env.New().
		Setup(
			envfuncs.CreateCluster(kind.NewProvider(), kindClusterName),
		).
		Finish(
			envfuncs.DestroyCluster(kindClusterName),
		)

	os.Exit(testEnv.Run(m))
}

func Test_CRD_Success(t *testing.T) {
	asserts := []testhelpers.SuccessfulAssert{
		{
			Name: "Exists",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				return crds.NewCRDAssertion(
					assertion.WithResourceName("orders.acme.cert-manager.io"),
					assertion.WithSetup(helpers.CreateResourceFromPath(crdPath)),
					assertion.WithTeardown(
						helpers.DeleteResourceFromPath(crdPath),
						helpers.Sleep(5*time.Second),
					),
				).Exists()
			},
		},
		{
			Name: "HasVersion",
			SuccessfulAssert: func(_ require.TestingT) assertion.Assertion {
				return crds.NewCRDAssertion(
					assertion.WithResourceName("orders.acme.cert-manager.io"),
					assertion.WithSetup(helpers.CreateResourceFromPath(crdPath)),
					assertion.WithTeardown(
						helpers.DeleteResourceFromPath(crdPath),
						helpers.Sleep(5*time.Second),
					),
				).Exists().HasVersion("v1")
			},
		},
	}

	testhelpers.TestSuccessfulAsserts(t, testEnv, asserts...)
}

func Test_CRD_Fail(t *testing.T) {
	asserts := []testhelpers.FailingAssert{
		{
			Name: "Exists",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				return crds.NewCRDAssertion(
					assertion.WithRequireT(t),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithResourceName("orders.acme.cert-manager.io"),
				).Exists()
			},
		},
		{
			Name: "HasVersion",
			FailingAssert: func(t require.TestingT) assertion.Assertion {
				return crds.NewCRDAssertion(
					assertion.WithRequireT(t),
					assertion.WithInterval(100*time.Millisecond),
					assertion.WithTimeout(500*time.Millisecond),
					assertion.WithResourceName("orders.acme.cert-manager.io"),
					assertion.WithSetup(helpers.CreateResourceFromPath(crdPath)),
					assertion.WithTeardown(
						helpers.DeleteResourceFromPath(crdPath),
						helpers.Sleep(5*time.Second),
					),
				).Exists().HasVersion("v1alpha1")
			},
		},
	}
	testhelpers.TestFailingAsserts(t, testEnv, asserts...)
}

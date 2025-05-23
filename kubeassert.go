// kubeassert provides a set of assertion helpers for testing Kubernetes resources. It extends the e2e-framework
// package with a set of predefined assertions (e.g. a Deployment is ready) as well as additional helper functionality.
//
//revive:disable:exported
package kubeassert

import (
	"github.com/DWSR/kubeassert-go/internal/assertion"
	"github.com/DWSR/kubeassert-go/internal/assertionhelpers"
	"github.com/DWSR/kubeassert-go/internal/crds"
	"github.com/DWSR/kubeassert-go/internal/deployments"
	"github.com/DWSR/kubeassert-go/internal/namespaces"
	"github.com/DWSR/kubeassert-go/internal/pdbs"
	"github.com/DWSR/kubeassert-go/internal/pods"
	"github.com/DWSR/kubeassert-go/internal/secrets"
)

type (
	Assertion           = assertion.Assertion
	DeploymentAssertion = deployments.DeploymentAssertion
	NamespaceAssertion  = namespaces.NamespaceAssertion
	CRDAssertion        = crds.CRDAssertion
	PDBAssertion        = pdbs.PDBAssertion
	PodAssertion        = pods.PodAssertion
	SecretAssertion     = secrets.SecretAssertion
)

var (
	WithLabels           = assertion.WithResourceLabels
	WithFields           = assertion.WithResourceFields
	WithInterval         = assertion.WithInterval
	WithTimeout          = assertion.WithTimeout
	WithBuilder          = assertion.WithBuilder
	WithRequireT         = assertion.WithRequireT
	WithNamespace        = assertion.WithResourceNamespace
	WithNamespaceFromEnv = assertion.WithResourceNamespaceFromTestEnv
	WithResourceName     = assertion.WithResourceName
	WithSetup            = assertion.WithSetup
	WithTeardown         = assertion.WithTeardown

	NewDeploymentAssertion = deployments.NewDeploymentAssertion
	NewNamespaceAssertion  = namespaces.NewNamespaceAssertion
	NewCRDAssertion        = crds.NewCRDAssertion
	NewPDBAssertion        = pdbs.NewPDBAssertion
	NewPodAssertion        = pods.NewPodAssertion
	NewSecretAssertion     = secrets.NewSecretAssertion

	ApplyKustomization     = assertionhelpers.ApplyKustomization
	CreateResourceFromPath = assertionhelpers.CreateResourceFromPath
	DeleteResourceFromPath = assertionhelpers.DeleteResourceFromPath
	Sleep                  = assertionhelpers.Sleep
	TestAssertions         = assertionhelpers.TestAssertions
)

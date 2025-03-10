// The assertion package provides common functionality used to define assertions against a set of one or more
// Kubernetes resources.
package assertion

import (
	"time"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

type (
	// Assertion is an interface for a generic assertion (har) about the state of one or more resources in a
	// Kubernetes cluster. It is embedded into other assertion types to provide common functionality.
	Assertion interface {
		optionSetters

		clone() Assertion

		// GetLabels returns the labels (i.e. metadata.labels) used to select resources for the assertion.
		GetLabels() map[string]string

		// GetFields returns the fields used to select resources for the assertion.
		GetFields() map[string]string

		// GetListOptions returns the ListOptions used to when listing resources from the API server.
		ListOptions(cfg *envconf.Config) metav1.ListOptions

		// GetInterval returns the interval used when polling for the assertion to be true.
		GetInterval() time.Duration

		// GetTimeout returns the timeout used when polling for the assertion to be true.
		GetTimeout() time.Duration

		// GetBuilder returns the features.FeatureBuilder used to build the e2e-framework Feature.
		GetBuilder() *features.FeatureBuilder

		// SetBuilder sets the features.FeatureBuilder used to build the e2e-framework Feature.
		SetBuilder(builder *features.FeatureBuilder)

		// GetRequireT returns the require.TestingT used to report test failures. This is largely intended for testing
		// as it enables detection of failing require/assert statements.
		GetRequireT() require.TestingT
	}

	// private methods to enable the use of assertion options without accidentally exposing the functionality outside of
	// the package.
	optionSetters interface {
		setLabels(assertLabels map[string]string)
		setFields(assertFields map[string]string)
		setListOptionsFn(fn listOptionsFunc)
		setInterval(interval time.Duration)
		setTimeout(timeout time.Duration)
		setRequireT(t require.TestingT)
	}
)

// Clone clones an assertion. This is done this way instead of providing an exported Clone method in order
// to avoid having the Clone method exported on all assertion types.
//
//nolint:ireturn
func Clone(a Assertion) Assertion {
	return a.clone()
}

// AsFeature returns an e2e-framework Feature based on the supplied Assertion. This can be used to integrate Assertions
// into existing e2e-framework tests.
//
//nolint:ireturn
func AsFeature(assert Assertion) features.Feature {
	return assert.GetBuilder().Feature()
}

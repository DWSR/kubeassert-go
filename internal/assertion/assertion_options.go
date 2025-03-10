package assertion

import (
	"time"

	"github.com/stretchr/testify/require"
	"sigs.k8s.io/e2e-framework/pkg/features"
	e2etypes "sigs.k8s.io/e2e-framework/pkg/types"
)

// Option is a function that configures one or more facets of an Assertion.
type Option func(Assertion)

// WithResourceLabels sets the labels to be used when selecting resources for the assertion.
func WithResourceLabels(labels map[string]string) Option {
	return func(a Assertion) {
		a.setLabels(labels)
	}
}

// WithResourceFields sets the fields to be used when selecting resources for the assertion.
func WithResourceFields(fields map[string]string) Option {
	return func(a Assertion) {
		a.setFields(fields)
	}
}

// WithInterval sets the interval used when polling for the assertion to be true.
func WithInterval(interval time.Duration) Option {
	return func(a Assertion) {
		a.setInterval(interval)
	}
}

// WithTimeout sets the timeout used when polling for the assertion to be true.
func WithTimeout(timeout time.Duration) Option {
	return func(a Assertion) {
		a.setTimeout(timeout)
	}
}

// WithBuilder sets the features.FeatureBuilder used to build the e2e-framework Feature.
func WithBuilder(builder *features.FeatureBuilder) Option {
	return func(a Assertion) {
		a.SetBuilder(builder)
	}
}

// WithRequireT sets the require.TestingT used to report test failures.
func WithRequireT(requireT require.TestingT) Option {
	return func(a Assertion) {
		a.setRequireT(requireT)
	}
}

// WithResourceNamespace sets the namespace to be used when selecting resources for the assertion.
func WithResourceNamespace(namespaceName string) Option {
	return func(a Assertion) {
		newFields := a.GetFields()
		newFields["metadata.namespace"] = namespaceName
		a.setFields(newFields)
	}
}

// WithResourceNamespaceFromTestEnv sets the namespace to be used when selecting resources for the assertion to the
// namespace set in the test environment.
func WithResourceNamespaceFromTestEnv() Option {
	return func(a Assertion) {
		a.setListOptionsFn(listOptionsWithNamespaceFromEnv)
	}
}

// WithResourceName sets the name (i.e. metadata.name) to be used when selecting resources for the assertion.
func WithResourceName(name string) Option {
	return func(a Assertion) {
		newFields := a.GetFields()
		newFields["metadata.name"] = name
		a.setFields(newFields)
	}
}

// WithSetup adds setup steps to the assertion.
func WithSetup(steps ...e2etypes.StepFunc) Option {
	return func(assert Assertion) {
		builder := assert.GetBuilder()

		for _, s := range steps {
			builder = builder.Setup(s)
		}

		assert.SetBuilder(builder)
	}
}

// WithTeardown adds teardown steps to the assertion.
func WithTeardown(steps ...e2etypes.StepFunc) Option {
	return func(assert Assertion) {
		builder := assert.GetBuilder()

		for _, s := range steps {
			builder = builder.Teardown(s)
		}

		assert.SetBuilder(builder)
	}
}

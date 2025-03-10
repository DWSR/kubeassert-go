package assertion

import (
	"time"

	"github.com/stretchr/testify/require"
	"sigs.k8s.io/e2e-framework/pkg/features"
	e2etypes "sigs.k8s.io/e2e-framework/pkg/types"
)

// Option is a function that configures one or more facets of an Assertion.
type Option func(Assertion)

func WithResourceLabels(labels map[string]string) Option {
	return func(a Assertion) {
		a.setLabels(labels)
	}
}

func WithResourceFields(fields map[string]string) Option {
	return func(a Assertion) {
		a.setFields(fields)
	}
}

func WithInterval(interval time.Duration) Option {
	return func(a Assertion) {
		a.setInterval(interval)
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(a Assertion) {
		a.setTimeout(timeout)
	}
}

func WithBuilder(builder *features.FeatureBuilder) Option {
	return func(a Assertion) {
		a.SetBuilder(builder)
	}
}

func WithRequireT(requireT require.TestingT) Option {
	return func(a Assertion) {
		a.setRequireT(requireT)
	}
}

func WithResourceNamespace(namespaceName string) Option {
	return func(a Assertion) {
		newFields := a.GetFields()
		newFields["metadata.namespace"] = namespaceName
		a.setFields(newFields)
	}
}

func WithResourceNamespaceFromTestEnv() Option {
	return func(a Assertion) {
		a.setListOptionsFn(listOptionsWithNamespaceFromEnv)
	}
}

func WithResourceName(name string) Option {
	return func(a Assertion) {
		newFields := a.GetFields()
		newFields["metadata.name"] = name
		a.setFields(newFields)
	}
}

func WithSetup(steps ...e2etypes.StepFunc) Option {
	return func(a Assertion) {
		builder := a.GetBuilder()
		for _, s := range steps {
			builder = builder.Setup(s)
		}
		a.SetBuilder(builder)
	}
}

func WithTeardown(steps ...e2etypes.StepFunc) Option {
	return func(a Assertion) {
		builder := a.GetBuilder()
		for _, s := range steps {
			builder = builder.Teardown(s)
		}
		a.SetBuilder(builder)
	}
}

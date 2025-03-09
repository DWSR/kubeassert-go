package assertion

import (
	"time"

	"github.com/stretchr/testify/require"
	"sigs.k8s.io/e2e-framework/pkg/features"
	e2etypes "sigs.k8s.io/e2e-framework/pkg/types"
)

type AssertionOption func(Assertion)

func WithResourceLabels(labels map[string]string) AssertionOption {
	return func(a Assertion) {
		a.setLabels(labels)
	}
}

func WithResourceFields(fields map[string]string) AssertionOption {
	return func(a Assertion) {
		a.setFields(fields)
	}
}

func WithInterval(interval time.Duration) AssertionOption {
	return func(a Assertion) {
		a.setInterval(interval)
	}
}

func WithTimeout(timeout time.Duration) AssertionOption {
	return func(a Assertion) {
		a.setTimeout(timeout)
	}
}

func WithBuilder(builder *features.FeatureBuilder) AssertionOption {
	return func(a Assertion) {
		a.SetBuilder(builder)
	}
}

func WithRequireT(requireT require.TestingT) AssertionOption {
	return func(a Assertion) {
		a.setRequireT(requireT)
	}
}

func WithResourceNamespace(namespaceName string) AssertionOption {
	return func(a Assertion) {
		newFields := a.GetFields()
		newFields["metadata.namespace"] = namespaceName
		a.setFields(newFields)
	}
}

func WithResourceNamespaceFromTestEnv() AssertionOption {
	return func(a Assertion) {
		a.setListOptionsFn(listOptionsWithNamespaceFromEnv)
	}
}

func WithResourceName(name string) AssertionOption {
	return func(a Assertion) {
		newFields := a.GetFields()
		newFields["metadata.name"] = name
		a.setFields(newFields)
	}
}

func WithSetup(steps ...e2etypes.StepFunc) AssertionOption {
	return func(a Assertion) {
		builder := a.GetBuilder()
		for _, s := range steps {
			builder = builder.Setup(s)
		}
		a.SetBuilder(builder)
	}
}

func WithTeardown(steps ...e2etypes.StepFunc) AssertionOption {
	return func(a Assertion) {
		builder := a.GetBuilder()
		for _, s := range steps {
			builder = builder.Teardown(s)
		}
		a.SetBuilder(builder)
	}
}

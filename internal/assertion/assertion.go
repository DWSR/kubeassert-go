package assertion

import (
	"context"
	"time"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apimachinerywait "k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

type (
	Assertion interface {
		setLabels(assertLabels map[string]string)

		GetLabels() map[string]string

		setFields(assertFields map[string]string)

		GetFields() map[string]string

		setListOptionsFn(fn listOptionsFunc)

		ListOptions(cfg *envconf.Config) metav1.ListOptions

		setInterval(interval time.Duration)

		GetInterval() time.Duration

		setTimeout(timeout time.Duration)

		GetTimeout() time.Duration

		GetBuilder() *features.FeatureBuilder

		SetBuilder(builder *features.FeatureBuilder)

		GetRequireT() require.TestingT

		setRequireT(t require.TestingT)

		WaitForCondition(ctx context.Context, conditionFunc apimachinerywait.ConditionWithContextFunc) error

		clone() Assertion

		AsFeature() features.Feature
	}
)

func CloneAssertion(a Assertion) Assertion {
	return a.clone()
}

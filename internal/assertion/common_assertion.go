package assertion

import (
	"time"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

type (
	commonAssertion struct {
		builder           *features.FeatureBuilder
		interval          time.Duration
		assertFields      map[string]string
		assertLabels      map[string]string
		labelRequirements labels.Requirements
		timeout           time.Duration
		requireT          require.TestingT
		listOptionsFn     listOptionsFunc
	}
)

const (
	defaultTimeout  = 30 * time.Second
	defaultInterval = 1 * time.Second
)

func (ca *commonAssertion) setLabels(assertLabels map[string]string) {
	ca.assertLabels = assertLabels
}

func (ca *commonAssertion) GetLabels() map[string]string {
	return ca.assertLabels
}

func (ca *commonAssertion) setFields(assertFields map[string]string) {
	ca.assertFields = assertFields
}

func (ca *commonAssertion) GetFields() map[string]string {
	return ca.assertFields
}

func (ca *commonAssertion) setInterval(interval time.Duration) {
	ca.interval = interval
}

func (ca *commonAssertion) GetInterval() time.Duration {
	return ca.interval
}

func (ca *commonAssertion) setTimeout(timeout time.Duration) {
	ca.timeout = timeout
}

func (ca *commonAssertion) GetTimeout() time.Duration {
	return ca.timeout
}

func (ca *commonAssertion) SetBuilder(builder *features.FeatureBuilder) {
	ca.builder = builder
}

func (ca *commonAssertion) GetBuilder() *features.FeatureBuilder {
	return ca.builder
}

func (ca *commonAssertion) setRequireT(requireT require.TestingT) {
	ca.requireT = requireT
}

//nolint:ireturn
func (ca *commonAssertion) GetRequireT() require.TestingT {
	return ca.requireT
}

func (ca *commonAssertion) setListOptionsFn(fn listOptionsFunc) {
	ca.listOptionsFn = fn
}

func (ca *commonAssertion) ListOptions(cfg *envconf.Config) metav1.ListOptions {
	return ca.listOptionsFn(ca, cfg)
}

//nolint:ireturn
func (ca *commonAssertion) clone() Assertion {
	return &commonAssertion{
		builder:           ca.builder,
		interval:          ca.interval,
		assertFields:      ca.assertFields,
		assertLabels:      ca.assertLabels,
		labelRequirements: ca.labelRequirements,
		timeout:           ca.timeout,
		requireT:          ca.requireT,
		listOptionsFn:     ca.listOptionsFn,
	}
}

// NewAssertion creates a new Assertion with the provided options.
//
//nolint:ireturn
func NewAssertion(opts ...Option) Assertion {
	assertion := commonAssertion{
		builder:       features.New("default"),
		assertFields:  make(map[string]string),
		assertLabels:  make(map[string]string),
		timeout:       defaultTimeout,
		interval:      defaultInterval,
		requireT:      nil,
		listOptionsFn: defaultListOptions,
	}

	for _, opt := range opts {
		opt(&assertion)
	}

	return &assertion
}

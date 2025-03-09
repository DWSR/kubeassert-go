package testhelpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/DWSR/kubeassert-go/internal/assertion"
)

type (
	SuccessfulAssert struct {
		Name             string
		SuccessfulAssert func(t require.TestingT) assertion.Assertion
	}

	FailingAssert struct {
		Name          string
		FailingAssert func(t require.TestingT) assertion.Assertion
	}
)

func TestSuccessfulAsserts(t *testing.T, testEnv env.Environment, asserts ...SuccessfulAssert) {
	t.Helper()

	assertFeatures := make([]features.Feature, 0)

	for _, a := range asserts {
		assertFeatures = append(assertFeatures, a.SuccessfulAssert(t).AsFeature())
	}

	testEnv.Test(t, assertFeatures...)
}

func TestFailingAsserts(t *testing.T, testEnv env.Environment, asserts ...FailingAssert) {
	t.Helper()

	for _, tc := range asserts {
		t.Run(tc.Name, func(t *testing.T) {
			mockT := &MockT{}
			testEnv.Test(t, tc.FailingAssert(mockT).AsFeature())
			assert.True(t, mockT.Failed)
		})
	}
}

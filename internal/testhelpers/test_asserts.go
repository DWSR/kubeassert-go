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
	// SuccessfulAssert is a struct that contains the name of the test and a function that returns an assertion.Assertion
	// that should be true or pass.
	SuccessfulAssert struct {
		Name             string
		SuccessfulAssert func(t require.TestingT) assertion.Assertion
	}

	// FailingAssert is a struct that contains the name of the test and a function that returns an assertion.Assertion
	// that should be false or fail.
	FailingAssert struct {
		Name          string
		FailingAssert func(t require.TestingT) assertion.Assertion
	}
)

// TestSuccessfulAsserts is a helper function that runs a series of successful asserts.
func TestSuccessfulAsserts(t *testing.T, testEnv env.Environment, asserts ...SuccessfulAssert) {
	t.Helper()

	assertFeatures := make([]features.Feature, 0)

	for _, a := range asserts {
		assertFeatures = append(assertFeatures, assertion.AsFeature(a.SuccessfulAssert(t)))
	}

	testEnv.Test(t, assertFeatures...)
}

// TestFailingAsserts is a helper function that runs a series of failing asserts.
func TestFailingAsserts(t *testing.T, testEnv env.Environment, asserts ...FailingAssert) {
	t.Helper()

	for _, tc := range asserts {
		t.Run(tc.Name, func(t *testing.T) {
			mockT := &MockT{}
			testEnv.Test(t, assertion.AsFeature(tc.FailingAssert(mockT)))
			assert.True(t, mockT.Failed)
		})
	}
}

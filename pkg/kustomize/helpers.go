package kustomize

import (
	"testing"

	"github.com/deislabs/porter/pkg/context"
	"k8s.io/client-go/kubernetes"
	testclient "k8s.io/client-go/kubernetes/fake"
)

type TestMixin struct {
	*Mixin
	TestContext *context.TestContext
}

type testKubernetesFactory struct {
}

func (t *testKubernetesFactory) GetClient(configPath string) (kubernetes.Interface, error) {
	return testclient.NewSimpleClientset(), nil
}

// NewTestMixin initializes a kustomize mixin, with the output buffered, and an in-memory file system.
func NewTestMixin(t *testing.T) *TestMixin {
	c := context.NewTestContext(t)
	m := New()
	m.Context = c.Context
	return &TestMixin{
		Mixin:       m,
		TestContext: c,
	}
}

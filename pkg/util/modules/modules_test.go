package modules

import (
	"fmt"
	"testing"

	"github.com/cortexproject/cortex/pkg/util/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockInitFunc() (services.Service, error) { return nil, nil }

func TestDependencies(t *testing.T) {
	var testModules = map[string]module{
		"serviceA": {
			initFn: mockInitFunc,
		},

		"serviceB": {
			initFn: mockInitFunc,
		},

		"serviceC": {
			initFn: mockInitFunc,
		},
	}

	mm := NewManager()
	for name, mod := range testModules {
		mm.RegisterModule(name, mod.initFn)
	}
	assert.NoError(t, mm.AddDependency("serviceB", "serviceA"))
	assert.NoError(t, mm.AddDependency("serviceC", "serviceB"))
	assert.Equal(t, mm.modules["serviceB"].deps, []string{"serviceA"})

	invDeps := mm.findInverseDependencies("serviceA", []string{"serviceB", "serviceC"})
	require.Len(t, invDeps, 1)
	assert.Equal(t, invDeps[0], "serviceB")

	svcs, err := mm.InitModuleServices("serviceC")
	assert.Nil(t, svcs)
	assert.Error(t, err, fmt.Errorf("module serviceA returned nil service but has other modules dependent on it: [serviceB]"))
}

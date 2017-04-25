// This file provides orchestrator provider's registry related features.
//
// NOTE:
//    This is the new file w.r.t the deprecated plugins.go file
package orchprovider

import (
	"fmt"
	"sync"

	"github.com/golang/glog"
)

type OrchProviderFactory func(name string) (OrchestratorInterface, error)

// Registration is managed in a safe manner via these variables
var (
	orchProviderRegMutex sync.Mutex
	orchProviderRegistry = make(map[string]OrchProviderFactory)
)

// RegisterOrchestrator registers a orchestration provider by the provider's name.
// This registers the orchestrator provider name with the provider's instance
// creating function i.e. a Factory.
//
// NOTE:
//    Each implementation of orchestrator plugin need to call
// RegisterOrchestrator inside their init() function.
func RegisterOrchestrator(name string, oInstFactory OrchProviderFactory) {
	orchProviderRegMutex.Lock()
	defer orchProviderRegMutex.Unlock()

	if _, found := orchProviderRegistry[name]; found {
		glog.Fatalf("Orchestration provider %q was registered twice", name)
	}

	glog.V(1).Infof("Registered orchestration provider %q", name)
	orchProviderRegistry[name] = oInstFactory
}

// GetOrchestrator creates a new instance of the named orchestration provider,
// or nil if the name is unknown.
func GetOrchestrator(name string) (OrchestratorInterface, error) {
	orchProviderRegMutex.Lock()
	defer orchProviderRegMutex.Unlock()

	oInstFactory, found := orchProviderRegistry[name]
	if !found {
		return nil, fmt.Errorf("'%s' is not registered as an orchestration provider", name)
	}

	// Orchestration provider's instance creating function is invoked here
	return oInstFactory(name)
}

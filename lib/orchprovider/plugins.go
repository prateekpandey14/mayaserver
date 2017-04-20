// This file provides maya api server's plugin related utilities w.r.t
// orchestrator provider.
package orchprovider

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/golang/glog"
)

// OrchFactory is signature function that every orchestrator plugin implementor
// needs to implement. It should contain the initialization logic w.r.t a
// orchestrator plugin. This function signature i.e. functional type has been created
// to enable lazy initialization of orchestrator plugin. In other words, a orchestrator
// plugin can be initialized at runtime when the parameters are available or
// can be provided.
//
// `name` parameter signifies the name of the orchestrator plugin.
//
// `region` parameter signifies the name of the region where the orchestrator
// is currently operating.
//
// `config` parameter provides an io.Reader handler in order to load specific
// configurations. If no configuration is provided the parameter is nil.
//
// TODO
// name should be the only parameter that will be provided.
// Address, Region, dc & other config options should be provided/derived during
// volume provisioning.
//
// e.g. orchprofiles, volumeprofiles, etc packages may be required
// Aim will be to remove the config file completely. This will also remove the
// need to create multiple instances of specific orchestrators i.e. variants.
type OrchFactory func(name string, region string, config io.Reader) (OrchestratorInterface, error)

// All registered orchestration providers.
var (
	orchPluginsMutex   sync.Mutex
	orchPluginRegistry = make(map[string]OrchFactory)
)

// RegisterOrchProvider registers a orchprovider.OrchestratorInterface by name.
// This is just a registry entry.
//
// NOTE:
// Registration & Initialization are two different workflows.
//
// OrchFactory instance represents the initialization logic. This is
// executed lazily. The initialization logic accepts various parameters
// like:
//
//  1. orchestrator plugin name,
//  2. orchestrator plugin config file
//
// NOTE:
//    Each implementation of orchestrator plugin need to call
// RegisterOrchProvider inside their init() function.
func RegisterOrchProvider(name string, factory OrchFactory) {
	orchPluginsMutex.Lock()
	defer orchPluginsMutex.Unlock()

	if _, found := orchPluginRegistry[name]; found {
		glog.Fatalf("Orchestration provider %q was registered twice", name)
	}

	glog.V(1).Infof("Registered orchestration provider %q", name)
	orchPluginRegistry[name] = factory
}

// IsOrchProvider returns true if name corresponds to an already
// registered orchestration provider.
func IsOrchProvider(name string) bool {
	orchPluginsMutex.Lock()
	defer orchPluginsMutex.Unlock()

	_, found := orchPluginRegistry[name]
	return found
}

// OrchProviders returns the name of all registered orchestration
// providers in a string slice
func OrchProviders() []string {
	names := []string{}
	orchPluginsMutex.Lock()
	defer orchPluginsMutex.Unlock()

	for name := range orchPluginRegistry {
		names = append(names, name)
	}

	return names
}

// GetOrchProvider creates an instance of the named orchestration provider,
// or nil if the name is unknown.  The error return is only used if the named
// provider was known but failed to initialize. The config parameter specifies
// the io.Reader handler of the configuration file for the orchestration
// provider, or nil for no configuation.
func GetOrchProvider(name string, region string, config io.Reader) (OrchestratorInterface, error) {
	orchPluginsMutex.Lock()
	defer orchPluginsMutex.Unlock()

	oFactory, found := orchPluginRegistry[name]
	if !found {
		return nil, fmt.Errorf("Orchestrator plugin '%s' not registered", name)
	}

	return oFactory(name, region, config)
}

// InitOrchProvider creates an instance of the named orchestrator plugin.
//
// NOTE:
//    Who calls this ?
// This is triggered while starting the Mayaserver as a http service.
//
// Http service invokes this to initialize the default orchestrator plugin with the
// plugin's region.
//
// This can also be invoked at runtime depending on user initiated requests that
// demand a specific region based orchestrator plugin.
//
// NOTE:
//    However, the orchestrator plugin name should be registered before invoking this
// function.
func InitOrchProvider(name string, region string, configFilePath string) (OrchestratorInterface, error) {
	var orchestrator OrchestratorInterface
	var err error

	if name == "" {
		glog.Info("Orchestrator name not provided")
		return nil, nil
	}

	if region == "" {
		glog.Info("Orchestrator region not provided")
		return nil, nil
	}

	var config *os.File
	if configFilePath != "" {
		config, err = os.Open(configFilePath)
		if err != nil {
			glog.Warningf("%#s", err)
		}

		defer config.Close()
	}

	if config != nil {
		orchestrator, err = GetOrchProvider(name, region, config)
	} else {
		orchestrator, err = GetOrchProvider(name, region, nil)
	}

	if err != nil {
		return nil, err
	}

	if orchestrator == nil {
		return nil, fmt.Errorf("Could not create '%s' orchestration provider", name)
	}

	return orchestrator, nil
}

// This file provides persistent volume provisioner plugin's registration &/ registry
// features.
package volume

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/golang/glog"
	"github.com/openebs/mayaserver/lib/orchprovider"
)

// VolumeFactory is a signature function that each persistent volume provisioner
// plugin implementor needs to implement. It should contain the initialization
// logic specific to the persistent volume provisioner plugin. This function
// signature i.e. functional type enables lazy initialization of persistent volume
// provisioner plugin. In other words, a persistent volume provisioner plugin can
// be initialized at runtime when the parameters are available or can be provided.
//
// `name` parameter signifies the name of the volume plugin
//
// `config` parameter provides an io.Reader handler in order to load specific
// configurations. If no configuration is provided the parameter is nil.
//
// `aspect` parameter provides handles to various orthogonal aspects of the
// volume plugin.
//
// TODO
//    It might be good to remove the concept of a config & aspect from initialization.
// It has been found that these values are dynamic and differ from one request to
// another request.
//
// NOTE:
//    This also leads to removal of plugin registry concept. In other words, a
// request would create its own instance of persistent volume provisioner, set
// this instance with appropriate profiles & make use of it. The lifecycle of this
// instance will be tied to the lifecycle of the request.
//
// NOTE:
//    This also indicates the removal of registration & initialization logic.
//
// TODO
//    Simplify the registration logic to provide a persistent volume provisioner
// plugin instance based on the name of the volume provisioner plugin. In addition,
// each of the persistent volume provisioner plugin should fill this registry with
// its respective instance creation function. Each call to this registry will lead
// to creation & return of a newly created persistent volume provisioner plugin
// instance
type VolumeFactory func(name string, config io.Reader, aspect VolumePluginAspect) (VolumeInterface, error)

// TODO
// This may not be required once Persistent Volume Provisioner Profile is implemented.
//
// VolumePluginAspect interface abstracts persistent volume provisioner plugin's
// orthogonal requirements.
type VolumePluginAspect interface {

	// Get the suitable orchestration provider.
	// A persistent volume provisioner plugin may be linked with a orchestrator
	// e.g. K8s, Nomad, Mesos, Swarm, etc. It can be vanilla Docker engine as well.
	//
	// Note:
	//    OpenEBS believes in running storage software in containers & hence
	// above container specific orchestrators.
	GetOrchProvider() (orchprovider.OrchestratorInterface, error)

	// Get the datacenter typically within a region, that will be the target
	// for all requests.
	//
	// e.g. An orchestrator might be deployed in multiple datacenters
	// within a region. This will help in pointing the requests to a
	// particular datacenter.
	//
	// NOTE:
	//  This will be used only if user action request does not specify a datacenter
	//
	// TODO
	// This contract method will be removed once persistent volume provisioner &
	// orchestration provider specific profiles are introduced. These profiles &/
	// operator provided runtime values will be set in the requests.
	DefaultDatacenter() (string, error)
}

// All registered persistent volume provisioner plugins.
var (
	volumePluginsMutex sync.Mutex

	// Persistent volume provisioner plugin name mapped against the plugin's
	// initializer function.
	volumePluginRegistry = make(map[string]VolumeFactory)

	// Persistent volume provisioner plugin name with the actual plugin instance.
	//
	// Acts as a cache.
	volumePluginInstances = make(map[string]VolumeInterface)
)

// VolumePluginConfig is how volume plugins receive configuration.  An instance
// specific to the plugin will be passed to the plugin's
// ProbeVolumePlugins(config) func.  Reasonable defaults will be provided by
// the binary hosting the plugins while allowing override of those default
// values.  Those config values are then set to an instance of
// VolumePluginConfig and passed to the plugin.
//
// Values in VolumeConfig are intended to be relevant to several plugins, but
// not necessarily all plugins.  The preference is to leverage strong typing
// in this struct.  All config items must have a descriptive but non-specific
// name (i.e, RecyclerMinimumTimeout is OK but RecyclerMinimumTimeoutForNFS is
// !OK).  An instance of config will be given directly to the plugin, so
// config names specific to plugins are unneeded and wrongly expose plugins in
// this VolumeConfig struct.
//
// OtherAttributes is a map of string values intended for one-off
// configuration of a plugin or config that is only relevant to a single
// plugin.  All values are passed by string and require interpretation by the
// plugin. Passing config as strings is the least desirable option but can be
// used for truly one-off configuration. The binary should still use strong
// typing for this value when binding CLI values before they are passed as
// strings in OtherAttributes.
//
// TODO
//    This will not be required once the simplified registration logic is in
// place along with persistent volume provisioner profiles.
type VolumePluginConfig struct {

	// OtherAttributes stores config as strings.  These strings are opaque to
	// the system and only understood by the binary hosting the plugin and the
	// plugin itself.
	OtherAttributes map[string]string
}

// RegisterVolumePlugin registers a volume.VolumeInterface by name.
// This is just a registry entry.
//
// NOTE:
//    Each implementation of volume plugin need to call
// RegisterVolumePlugin inside their init() function.
func RegisterVolumePlugin(name string, factory VolumeFactory) {
	volumePluginsMutex.Lock()
	defer volumePluginsMutex.Unlock()

	if _, found := volumePluginRegistry[name]; found {
		glog.Fatalf("Volume plugin %q was registered twice", name)
	}

	glog.V(1).Infof("Registered volume plugin %q", name)
	volumePluginRegistry[name] = factory
}

// IsVolumePlugin returns true if name corresponds to an already
// registered volume plugin.
func IsVolumePlugin(name string) bool {
	volumePluginsMutex.Lock()
	defer volumePluginsMutex.Unlock()

	_, found := volumePluginRegistry[name]
	return found
}

// VolumePlugins returns the name of all registered volume
// plugins in a string slice
func VolumePlugins() []string {
	names := []string{}
	volumePluginsMutex.Lock()
	defer volumePluginsMutex.Unlock()

	for name := range volumePluginRegistry {
		names = append(names, name)
	}
	return names
}

// GetVolumePlugin creates an instance or returns previously created instance of
// the named volume plugin.
//
// NOTE:
//    This can be invoked just to get a cached instance by providing the name
// of the volume plugin only.
func GetVolumePlugin(name string, config io.Reader, aspect VolumePluginAspect) (VolumeInterface, error) {
	volumePluginsMutex.Lock()
	defer volumePluginsMutex.Unlock()

	factory, found := volumePluginRegistry[name]
	if !found {
		return nil, fmt.Errorf("Volume plugin '%s' not registered", name)
	}

	// Search from cache
	existingInstance, found := volumePluginInstances[name]
	if !found {
		// create the plugin instance
		newInstance, err := factory(name, config, aspect)
		if err != nil {
			return nil, err
		}

		// cache it
		volumePluginInstances[name] = newInstance
		return newInstance, nil
	}

	return existingInstance, nil
}

// InitVolumePlugin creates an instance of the named volume plugin.
//
// NOTE:
//    Who calls this ?
// This is triggered while starting the Mayaserver as a http service.
//
// Http service invokes this to initialize the default volume plugin with the
// plugin's default orchestrator pointing to the orchestrator's default region.
//
// This can also be invoked at runtime depending on user initiated requests that
// demand a specific volume plugin or a specific variant of volume plugin.
//
// NOTE:
//    However, the volume plugin name should be registered before invoking this
// function.
func InitVolumePlugin(name string, configFilePath string, aspect VolumePluginAspect) (VolumeInterface, error) {
	//var orchestrator Interface
	var volumeInterface VolumeInterface
	var err error

	if name == "" {
		glog.Info("No volume plugin specified.")
		return nil, nil
	}

	if configFilePath != "" {
		var config *os.File
		config, err = os.Open(configFilePath)
		if err != nil {
			glog.Fatalf("Couldn't open volume plugin configuration %s: %#v",
				configFilePath, err)
		}

		defer config.Close()
		volumeInterface, err = GetVolumePlugin(name, config, aspect)
	} else {
		// Pass explicit nil so plugins can actually check for nil. See
		// "Why is my nil error value not equal to nil?" in golang.org/doc/faq.
		volumeInterface, err = GetVolumePlugin(name, nil, aspect)
	}

	if err != nil {
		return nil, fmt.Errorf("could not init volume plugin %q: %v", name, err)
	}

	if volumeInterface == nil {
		return nil, fmt.Errorf("unknown volume plugin %q", name)
	}

	return volumeInterface, nil
}

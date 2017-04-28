// This file provides persistent volume provisioner's registry related features.
//
// NOTE:
//    This is the new file w.r.t the deprecated volume/plugins.go file
package volume

import (
	"fmt"
	"sync"

	"github.com/golang/glog"
	"github.com/openebs/mayaserver/lib/api/v1"
)

type VolumeProvisionerFactory func(label, name string) (VolumeInterface, error)

// Registration is managed in a safe manner via these variables
var (
	volProvisionerRegMutex sync.Mutex
	volProvisionerRegistry = make(map[string]VolumeProvisionerFactory)
)

// RegisterVolumeProvisioner registers a persistent volume provisioner by the
// provisioner's name. This registers the provisioner name with the provisioner's
// instance creating function i.e. a Factory.
//
// NOTE:
//    Each implementation of persistent volume provisioner plugin need to call
// RegisterVolumeProvisioner inside their init() function.
func RegisterVolumeProvisioner(name string, vpInstFactory VolumeProvisionerFactory) {
	volProvisionerRegMutex.Lock()
	defer volProvisionerRegMutex.Unlock()

	if _, found := volProvisionerRegistry[name]; found {
		glog.Fatalf("Persistent volume provisioner '%s' was registered twice", name)
	}

	glog.V(1).Infof("Registered '%s' as persistent volume provisioner", name)
	volProvisionerRegistry[name] = vpInstFactory
}

// GetVolumeProvisioner creates a new instance of the named persistent volume
// provisioner or nil if the name is unknown.
func GetVolumeProvisioner(name string) (VolumeInterface, error) {
	volProvisionerRegMutex.Lock()
	defer volProvisionerRegMutex.Unlock()

	vpInstFactory, found := volProvisionerRegistry[name]
	if !found {
		return nil, fmt.Errorf("'%s' is not registered as a persistent volume provisioner", name)
	}

	// Persistent volume provisioner's instance creating function is invoked here
	// The persistent volume provisioner label is decided here. This label is common
	// to all persistent volume provisioner implementors.
	return vpInstFactory(string(v1.VolumeProvisionerNameLbl), name)
}

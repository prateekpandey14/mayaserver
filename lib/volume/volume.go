// This file abstracts & exposes persistent volume provisioner features. All
// maya api server's persistent volume provisioners need to implement these
// contracts.
package volume

import (
	"github.com/openebs/mayaserver/lib/api/v1"
	"github.com/openebs/mayaserver/lib/orchprovider"
)

// VolumeInterface abstracts the persistent volume features of any persistent
// volume provisioner.
//
// NOTE:
//    maya api server can make use of any persistent volume provisioner & execute
// corresponding volume related operations.
type VolumeInterface interface {
	// Name of the persistent volume provisioner
	Name() string

	// Profile will set the persistent volume provisioner's profile
	//
	// Note:
	//    Will return false if profile is not supported by the persistent
	// volume provisioner. This is typically used to set the persistent volume
	// provisioner profile lazily i.e. much after the initialization of persistent
	// volume provisioner instance.
	Profile(VolumeProvisionerProfile) (bool, error)

	// TODO
	// Rename to Creator ??
	//
	// Provisioner gets the instance capable of provisioning volumes w.r.t this
	// persistent volume provisioner.
	//
	// Note:
	//    Will return false if provisioning of volumes is not supported by the
	// persistent volume provisioner.
	Provisioner() (Provisioner, bool)

	// Deleter gets the instance capable of deleting volumes w.r.t this
	// persistent volume provisioner.
	//
	// Note:
	//    Will return false if deletion of volumes is not supported by the
	// persistent volume provisioner.
	Deleter() (Deleter, bool)

	// TODO
	// Rename to Reader ??
	//
	// Informer gets the instance capable of providing volume information w.r.t this
	// persistent volume provisioner.
	//
	// Note:
	//    Will return false if providing volume information is not supported by
	// the persistent volume provisioner.
	Informer() (Informer, bool)
}

// TODO
// Rename to Reader ??
//
// Informer interface abstracts fetching of volume related information
// from a persistent volume provisioner.
type Informer interface {
	// Info tries to fetch the volume details from the persistent volume
	// provisioner.
	Info(*v1.PersistentVolumeClaim) (*v1.PersistentVolume, error)
}

// TODO
// Rename to Creator ??
//
// Provisioner interface abstracts creation of volume from a persistent volume
// provisioner.
type Provisioner interface {
	// Provision tries to create a volume of a persistent volume provisioner.
	Provision(*v1.PersistentVolumeClaim) (*v1.PersistentVolume, error)
}

// Deleter interface abstracts deletion of volume of a persistent volume
// provisioner.
type Deleter interface {
	// Delete tries to delete a volume of a persistent volume provisioner.
	Delete(*v1.PersistentVolume) (*v1.PersistentVolume, error)
}

// VolumeProvisionerProfile abstracts & exposes a persistent volume provisioner's
// runtime features.
//
// NOTE:
//    A persistent volume provisioner can align to a specific implementation of
// this profile & hence change its execution strategy at runtime.
type VolumeProvisionerProfile interface {
	// Profile name of the persistent volume provisioner.
	//
	// Note:
	//  An invalid or empty name will return error.
	Name() (string, error)

	// Get the suitable orchestration provider.
	// A persistent volume provisioner plugin may be linked with a orchestrator
	// e.g. K8s, Nomad, Mesos, Swarm, etc. It can be Docker engine as well.
	//
	// Note:
	//    It can return false in its second return argument if orchestrator support
	// is not applicable.
	//
	// Note:
	//    OpenEBS believes in running storage software in containers & hence
	// these container specific orchestrators.
	Orchestrator() (orchprovider.OrchestratorInterface, bool, error)
}

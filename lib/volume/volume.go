// This file abstracts & exposes persistent volume related features as required by
// maya api server.
package volume

import (
	"github.com/openebs/mayaserver/lib/api/v1"
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

	// Provisioner gets the instance capable of provisioning volumes w.r.t this
	// persistent volume provisioner. Will return false if provisioning of volumes
	// is not supported by this provisioner.
	Provisioner() (Provisioner, bool)

	// Deleter gets the instance capable of deleting volumes w.r.t this
	// persistent volume provisioner. Will return false if deletion of volumes is
	// not supported by this provisioner.
	Deleter() (Deleter, bool)

	// Informer gets the instance capable of providing volume information w.r.t this
	// persistent volume provisioner. Will return false if providing volume information
	// is not supported by this provisioner.
	Informer() (Informer, bool)
}

// Informer interface abstracts fetching of volume related information
// from a persistent volume provisioner.
type Informer interface {
	// Info tries to fetch the volume details from the persistent volume
	// provisioner.
	Info(*v1.PersistentVolumeClaim) (*v1.PersistentVolume, error)
}

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

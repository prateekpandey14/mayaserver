// This file defines interfaces that determines an orchestrator w.r.t maya api
// server. All the features that maya api server wants from an orchestrator is
// defined in these set of interfaces.
package orchprovider

import (
	"github.com/openebs/mayaserver/lib/api/v1"
	volProfile "github.com/openebs/mayaserver/lib/profile/volumeprovisioner"
)

// OrchestrationInterface is an interface abstraction of a real orchestrator.
// It represents an abstraction that maya api server expects from its
// orchestrator.
//
// NOTE:
//  OrchestratorInterface is an aggregator of specific interfaces.
type OrchestratorInterface interface {
	// Label assigned against the orchestration provider
	Label() string

	// Name of the orchestration provider
	Name() string

	// Region where this orchestration provider is running/deployed
	Region() string

	// TODO
	// Deprecate once orchestrator profiles are ready
	//
	// NetworkPlacements gets the NetworkPlacements related features. Will return
	// false if not supported.
	NetworkPlacements() (NetworkPlacements, bool)

	// TODO
	// Deprecate in favour of StorageOps
	//
	// StoragePlacements gets the StoragePlacements related features. Will return
	// false if not supported.
	StoragePlacements() (StoragePlacements, bool)

	// StorageOps gets the instance that deals with storage related operations.
	// Will return false if not supported.
	StorageOps() (StorageOps, bool)
}

// NetworkPlacements provides the interface abstraction for network related
// placements, scheduling, etc that are available at the orchestrator.
//
// TODO
// This interface will not be required once maya api server implements orchestrator
// provider specific profiles.
type NetworkPlacements interface {

	// NetworkPropsReq will try to fetch the networking details at the orchestrator
	// based on a particular datacenter
	//
	// NetworkPropsReq does not fall under CRUD operations. This is applicable
	// to fetching properties from a config, or database etc.
	//
	// NOTE:
	//    This interface will have no control over Create, Update, Delete operations
	// of network properties
	NetworkPropsReq(dc string) (map[v1.ContainerNetworkingLbl]string, error)
}

// TODO
// Deprecate in favour of StorageOps
//
// StoragePlacement provides the blueprint for storage related
// placements, scheduling, etc at the orchestrator end.
type StoragePlacements interface {

	// StoragePlacementReq will try to create storage resource(s) at the
	// infrastructure
	StoragePlacementReq(pvc *v1.PersistentVolumeClaim) (*v1.PersistentVolume, error)

	// StorageRemovalReq will try to delete the storage resource(s) at
	// the infrastructure
	StorageRemovalReq(pv *v1.PersistentVolume) (*v1.PersistentVolume, error)

	// StorageInfoReq will try to fetch the details of a particular storage
	// resource
	StorageInfoReq(pvc *v1.PersistentVolumeClaim) (*v1.PersistentVolume, error)

	// StoragePropsReq will try to fetch the storage details at the orchestrator
	// based on a particular datacenter
	//
	// StoragePropsReq does not fall under CRUD operations. This is applicable
	// to fetching properties from a config, or database etc.
	//
	// NOTE:
	//    This interface will have no control over Create, Update, Delete operations
	// of storage properties.
	//
	// NOTE:
	//    jiva requires these persistent storage properties to provision
	// its instances e.g. backing persistence location is required on which
	// a jiva replica can operate.
	//
	// TODO
	// This function will not be required once maya api server implements orchestrator
	// provider specific profiles.
	StoragePropsReq(dc string) (map[v1.ContainerStorageLbl]string, error)
}

// StorageOps exposes various storage related operations that deals with
// storage placements, scheduling, etc. The low level work is done by the
// orchestrator.
type StorageOps interface {

	// AddStorage will add persistent volume running as containers
	AddStorage(volProProfile volProfile.VolumeProvisionerProfile) (*v1.PersistentVolume, error)

	// DeleteStorage will remove the persistent volume
	DeleteStorage(volProProfile volProfile.VolumeProvisionerProfile) (*v1.PersistentVolume, error)

	// ReadStorage will fetch information about the persistent volume
	ReadStorage(volProProfile volProfile.VolumeProvisionerProfile) (*v1.PersistentVolume, error)
}

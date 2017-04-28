// This file registers Kubernetes as an orchestration provider plugin in maya
// api server. This orchestration is for persistent volume provisioners which
// also are registered in maya api server.
package k8s

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/openebs/mayaserver/lib/api/v1"
	"github.com/openebs/mayaserver/lib/orchprovider"
	volProfile "github.com/openebs/mayaserver/lib/profile/volumeprovisioner"
)

// The registration logic for the kubernetes as orchestrator provider plugin
//
// NOTE:
//    This function is executed once per application
func init() {
	orchprovider.RegisterOrchestrator(
		// Registration entry when Kubernetes is the orchestrator provider plugin
		// NOTE:
		//    This value remains same for all instances of Kubernetes
		v1.K8sOrchestrator,

		// Below is a callback function that creates a new instance of Kubernetes
		// orchestration provider
		func(label v1.NameLabel, name v1.OrchestratorRegistry) (orchprovider.OrchestratorInterface, error) {
			return NewK8sOrchestrator(label, name)
		})
}

// K8sOrchestrator is a concrete implementation of following
// interfaces:
//
//  1. orchprovider.OrchestratorInterface,
//  2. orchprovider.NetworkPlacements, &
//  3. orchprovider.StoragePlacements
type k8sOrchestrator struct {
	// label specified to this orchestrator
	label v1.NameLabel

	// Name of the orchestrator as registered in the registry
	name v1.OrchestratorRegistry
}

// NewK8sOrchestrator provides a new instance of K8sOrchestrator.
func NewK8sOrchestrator(label v1.NameLabel, name v1.OrchestratorRegistry) (orchprovider.OrchestratorInterface, error) {

	glog.Infof("Building '%s':'%s' orchestration provider", label, name)

	if string(label) == "" {
		return nil, fmt.Errorf("Label not found while building k8s orchestrator")
	}

	if string(name) == "" {
		return nil, fmt.Errorf("Name not found while building k8s orchestrator")
	}

	return &k8sOrchestrator{
		label: label,
		name:  name,
	}, nil
}

// Label provides the label assigned against this orchestrator.
// This is an implementation of the orchprovider.OrchestratorInterface interface.
func (k *k8sOrchestrator) Label() string {
	// TODO
	// Do not typecast. Make it typed
	// Ensure this for all orch provider implementors
	return string(k.label)
}

// Name provides the name of this orchestrator.
// This is an implementation of the orchprovider.OrchestratorInterface interface.
func (k *k8sOrchestrator) Name() string {
	// TODO
	// Do not typecast. Make it typed
	// Ensure this for all orch provider implementors
	return string(k.name)
}

// TODO
// Deprecate in favour of orchestrator profile
// Region is not supported by k8sOrchestrator.
// This is an implementation of the orchprovider.OrchestratorInterface interface.
func (k *k8sOrchestrator) Region() string {
	return ""
}

// StorageOps provides storage operations instance that deals with all storage
// related functionality by aligning with Kubernetes as the orchestration provider.
//
// NOTE:
//    This is an implementation of the orchprovider.OrchestratorInterface interface.
func (k *k8sOrchestrator) StorageOps() (orchprovider.StorageOps, bool) {
	return k, true
}

// AddStorage will add persistent volume running as containers
func (k *k8sOrchestrator) AddStorage(volProProfile volProfile.VolumeProvisionerProfile) (*v1.PersistentVolume, error) {
	return nil, nil
}

// DeleteStorage will remove the persistent volume
func (k *k8sOrchestrator) DeleteStorage(volProProfile volProfile.VolumeProvisionerProfile) (*v1.PersistentVolume, error) {
	return nil, nil
}

// ReadStorage will fetch information about the persistent volume
func (k *k8sOrchestrator) ReadStorage(volProProfile volProfile.VolumeProvisionerProfile) (*v1.PersistentVolume, error) {
	return nil, nil
}

// TODO
// Deprecate in favour of StorageOps
//
// StoragePlacements is not supported by k8sOrchestrator
// This is an implementation of the orchprovider.OrchestratorInterface interface.
func (k *k8sOrchestrator) StoragePlacements() (orchprovider.StoragePlacements, bool) {
	return nil, false
}

// TODO
// Deprecate
//
// NetworkPlacements is not supported by k8sOrchestrator.
// This is an implementation of the orchprovider.OrchestratorInterface interface.
func (k *k8sOrchestrator) NetworkPlacements() (orchprovider.NetworkPlacements, bool) {
	return nil, false
}

// StorageInfoReq is a contract method implementation of
// orchprovider.StoragePlacements interface. It will fetch the persistent volume
// details from its volume provisioner which is running in K8s setup.
//
// This is an implementation of the orchprovider.StoragePlacements interface.
//func (k *k8sOrchestrator) StorageInfoReq(pvc *v1.PersistentVolumeClaim) (*v1.PersistentVolume, error) {
//	return nil, nil
//}

// StoragePlacementReq is a contract method implementation of
// orchprovider.StoragePlacements interface. It will create a persistent volume
// at its volume provisioner which is running in K8s setup.
//
// This is an implementation of the orchprovider.StoragePlacements interface.
//func (k *k8sOrchestrator) StoragePlacementReq(pvc *v1.PersistentVolumeClaim) (*v1.PersistentVolume, error) {
//	return nil, nil
//}

// StorageRemovalReq is a contract method implementation of
// orchprovider.StoragePlacements interface. It will remove a persistent volume
// from its volume provisioner which is running in K8s setup.
//
// This is an implementation of the orchprovider.StoragePlacements interface.
//func (k *k8sOrchestrator) StorageRemovalReq(pv *v1.PersistentVolume) (*v1.PersistentVolume, error) {
//	return nil, nil
//}

// TODO
// Will be removed in future.
// This is an implementation of the orchprovider.StoragePlacements interface.
//func (k *k8sOrchestrator) StoragePropsReq(dc string) (map[v1.ContainerStorageLbl]string, error) {
//	return nil, fmt.Errorf("StoragePropsReq is not supported by '%s:%s'", k.Name())
//}

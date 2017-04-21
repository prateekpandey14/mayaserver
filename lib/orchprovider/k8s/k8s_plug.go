// This file plugs the following:
//
//    1. Generic orchprovider defined by maya api server &
//    2. Kubernetes as the orchestration provider
package k8s

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/openebs/mayaserver/lib/api/v1"
	v1k8s "github.com/openebs/mayaserver/lib/api/v1/k8s"
	"github.com/openebs/mayaserver/lib/orchprovider"
)

// The registration logic for the kubernetes as orchestrator prvoider plugin
//
// NOTE:
//    This function is executed once per application
func init() {
	orchprovider.RegisterOrchestrator(
		// Name when kubernetes is the orchestrator provider plugin
		v1k8s.K8sOrchProviderPluginName,

		// Below is a callback function that creates a new instance of k8s orchestrator
		// provider plugin
		func(name string) (orchprovider.OrchestratorInterface, error) {
			return NewK8sOrchestrator(name)
		})
}

// K8sOrchestrator is a concrete implementation of following
// interfaces:
//
//  1. orchprovider.OrchestratorInterface,
//  2. orchprovider.NetworkPlacements, &
//  3. orchprovider.StoragePlacements
type k8sOrchestrator struct {

	// Name of this orchestrator
	name string

	// kConfig represents an instance that provides the coordinates
	// of a K8s server / cluster deployment.
	//kConfig *K8sConfig
}

// NewK8sOrchestrator provides a new instance of K8sOrchestrator.
func NewK8sOrchestrator(name string) (orchprovider.OrchestratorInterface, error) {

	glog.Infof("Building k8s orchestration provider")

	if name == "" {
		return nil, fmt.Errorf("Name missing while building k8s orchestrator")
	}

	return &k8sOrchestrator{
		name: name,
	}, nil
}

// Name provides the name of this orchestrator.
// This is an implementation of the orchprovider.OrchestratorInterface interface.
func (k *k8sOrchestrator) Name() string {

	return k.name
}

// Region is not supported by k8sOrchestrator.
// This is an implementation of the orchprovider.OrchestratorInterface interface.
func (k *k8sOrchestrator) Region() string {

	return ""
}

// StoragePlacements returns self.
// This is an implementation of the orchprovider.OrchestratorInterface interface.
func (k *k8sOrchestrator) StoragePlacements() (orchprovider.StoragePlacements, bool) {

	return k, true
}

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
func (k *k8sOrchestrator) StorageInfoReq(pvc *v1.PersistentVolumeClaim) (*v1.PersistentVolume, error) {

	return nil, nil
}

// StoragePlacementReq is a contract method implementation of
// orchprovider.StoragePlacements interface. It will create a persistent volume
// at its volume provisioner which is running in K8s setup.
//
// This is an implementation of the orchprovider.StoragePlacements interface.
func (k *k8sOrchestrator) StoragePlacementReq(pvc *v1.PersistentVolumeClaim) (*v1.PersistentVolume, error) {

	return nil, nil
}

// StorageRemovalReq is a contract method implementation of
// orchprovider.StoragePlacements interface. It will remove a persistent volume
// from its volume provisioner which is running in K8s setup.
//
// This is an implementation of the orchprovider.StoragePlacements interface.
func (k *k8sOrchestrator) StorageRemovalReq(pv *v1.PersistentVolume) (*v1.PersistentVolume, error) {

	return nil, nil
}

// TODO
// Will be removed in future.
// This is an implementation of the orchprovider.StoragePlacements interface.
func (k *k8sOrchestrator) StoragePropsReq(dc string) (map[v1.ContainerStorageLbl]string, error) {
	return nil, fmt.Errorf("Not supported")
}

// This file registers Kubernetes as an orchestration provider plugin in maya
// api server. This orchestration is for persistent volume provisioners which
// also are registered in maya api server.
package k8s

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/openebs/mayaserver/lib/api/v1"
	"github.com/openebs/mayaserver/lib/orchprovider"
	orchProfile "github.com/openebs/mayaserver/lib/profile/orchprovider"
	volProfile "github.com/openebs/mayaserver/lib/profile/volumeprovisioner"
	k8sCv1 "k8s.io/client-go/pkg/api/v1"
)

// The registration logic for the kubernetes as orchestrator provider plugin
//
// NOTE:
//    This function is executed once per application
func init() {
	orchprovider.RegisterOrchestrator(
		// Registration entry when Kubernetes is the orchestrator provider plugin
		//
		// NOTE:
		//    This value remains same for all instances of Kubernetes
		v1.K8sOrchestrator,

		// Below is a callback function that creates a new instance of Kubernetes
		// orchestration provider
		func(label v1.NameLabel, name v1.OrchProviderRegistry) (orchprovider.OrchestratorInterface, error) {
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
	name v1.OrchProviderRegistry
}

// NewK8sOrchestrator provides a new instance of K8sOrchestrator.
func NewK8sOrchestrator(label v1.NameLabel, name v1.OrchProviderRegistry) (orchprovider.OrchestratorInterface, error) {

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
func (k *k8sOrchestrator) ReadStorage(volProProfile volProfile.VolumeProvisionerProfile) (*v1.PersistentVolumeList, error) {
	// Fetch pvc from volume provisioner profile
	pvc, err := volProProfile.PVC()
	if err != nil {
		return nil, err
	}

	// Get orchestrator provider profile from pvc
	oPrfle, err := orchProfile.GetOrchProviderProfileByPVC(pvc)
	if err != nil {
		return nil, err
	}

	// Which kind of request ? in-cluster or out-of-cluster ?
	isInCluster, err := oPrfle.InCluster()
	if err != nil {
		return nil, err
	}

	// Fetch appropriate kubernetes clientset
	cs, err := GetClusterCS(isInCluster)
	if err != nil {
		return nil, err
	}

	// Get the namespace which will be queried
	ns, err := oPrfle.NS()
	if err != nil {
		return nil, err
	}

	// NOTE:
	//    A VSM can be one or more k8s PODs
	//
	// NOTE:
	//    maya api service assigns the VSM name as one of the labels against all
	// the pods created during creation of persistent volume
	vsm, err := volProProfile.VSMName()
	if err != nil {
		return nil, err
	}

	lOpts := k8sCv1.ListOptions{
		LabelSelector: "vsm=" + vsm,
	}

	podList, err := cs.CoreV1().Pods(ns).List(lOpts)
	if err != nil {
		return nil, err
	}

	if podList == nil || len(podList.Items) == 0 {
		return nil, fmt.Errorf("Storage not found for VSM '%s' with '%s:%s' as orchestrator", vsm, k.Label(), k.Name())
	}

	pvl := &v1.PersistentVolumeList{
		Items: make([]v1.PersistentVolume, len(podList.Items)),
	}

	// TODO
	// Transform the POD type to persistent volume type
	// Do this in v1/k8s package ??
	for _, pod := range podList.Items {
		pv := v1.PersistentVolume{}
		pv.Name = pod.Name
		pvl.Items = append(pvl.Items, pv)
	}

	return pvl, nil
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

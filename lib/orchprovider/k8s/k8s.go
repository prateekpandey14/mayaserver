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
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sApiV1 "k8s.io/client-go/pkg/api/v1"
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
//  2. orchprovider.NetworkPlacements &
//  3. orchprovider.StoragePlacements
type k8sOrchestrator struct {
	// label specified to this orchestrator
	label v1.NameLabel

	// name of the orchestrator as registered in the registry
	name v1.OrchProviderRegistry

	// k8sUtlGtr provides the handle to fetch K8sUtilInterface
	// NOTE:
	//    This will be set at runtime.
	k8sUtlGtr K8sUtilGetter
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

// GetK8sUtil provides the k8sUtil instance that is capable of performing low
// level k8s operations
//
// NOTE:
//    This is an implementation of K8sUtilGetter interface
//
// NOTE:
//    This is meant to be used by k8sOrchestrator & is not a generic
// orchprovider.OrchestratorInterface contract
func (k *k8sOrchestrator) GetK8sUtil(volProfile volProfile.VolumeProvisionerProfile) K8sUtilInterface {

	// TODO validations
	// if volProfile == nil
	// if volProfile.PVC() == nil
	// if volProfile.PVC().Labels == nil

	return &k8sUtil{
		volProfile: volProfile,
	}
}

// k8sOrchUtil provides a common utility function for k8sOrchestrator to get an
// instance of k8sUtilInterface
func k8sOrchUtil(k *k8sOrchestrator, volProfile volProfile.VolumeProvisionerProfile) K8sUtilInterface {
	// k8sUtilGetter may or may not have been set earlier
	//
	// NOTE:
	//    If k8sUtilGetter was set earlier, it is known as dependency injection.
	// This means the dependency was injected at runtime. The flow of execution
	// will depend on the injected dependency
	//
	// NOTE:
	//    If k8sUtilGetter was not set, then use the default one
	if k.k8sUtlGtr == nil {
		// k8sOrchestrator is a k8sUtilGetter implementor
		k.k8sUtlGtr = k
	}

	return k.k8sUtlGtr.GetK8sUtil(volProfile)
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
func (k *k8sOrchestrator) AddStorage(volProProfile volProfile.VolumeProvisionerProfile) (*v1.PersistentVolumeList, error) {
	// TODO
	// This is jiva specific
	// Move this entire logic to a separate package that will couple jiva
	// provisioner with k8s orchestrator

	// create k8s pod of persistent volume controller
	_, err := k.createControllerPod(volProProfile)
	if err != nil {
		return nil, err
	}

	// create k8s service of persistent volume controller
	_, err = k.createControllerService(volProProfile)
	if err != nil {
		// TODO
		// Delete the persistent volume controller pod
		// Delegate to DeleteStorage which should handle stuff in a
		// robust way
		return nil, err
	}

	// TODO
	// Get the persistent volume controller service name & IP address
	_, ctrlIP, err := k.GetControllerService(volProProfile)
	if err != nil {
		// TODO
		// Delete the persistent volume controller pod
		// Delegate to DeleteStorage which should handle stuff in a
		// robust way
		return nil, err
	}

	// Create the k8s pod of persistent volume replica
	isRequested := volProProfile.ReqReplica()
	if !isRequested {
		return nil, nil
	}

	_, err = k.CreateReplicaPods(volProProfile, ctrlIP)
	if err != nil {
		// TODO
		// Delete the persistent volume controller pod
		// Delegate to DeleteStorage which should handle stuff in a
		// robust way
		return nil, err
	}

	return k.ReadStorage(volProProfile)
}

// DeleteStorage will remove the persistent volume
func (k *k8sOrchestrator) DeleteStorage(volProProfile volProfile.VolumeProvisionerProfile) (*v1.PersistentVolume, error) {
	return nil, nil
}

// ReadStorage will fetch information about the persistent volume
func (k *k8sOrchestrator) ReadStorage(volProProfile volProfile.VolumeProvisionerProfile) (*v1.PersistentVolumeList, error) {
	k8sUtl := k8sOrchUtil(k, volProProfile)

	kc, supported := k8sUtl.K8sClient()
	if !supported {
		return nil, fmt.Errorf("K8s client not supported by '%s'", k8sUtl.Name())
	}

	ns, err := kc.NS()
	if err != nil {
		return nil, err
	}

	pOps, err := kc.Pods()
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

	// This filtering logic does not work with client-go v2.0.0
	// Need to write code for filtering to work.
	// Need to upgrade client-go to stable version
	lOpts := k8sApiV1.ListOptions{
		LabelSelector: string(v1.VSMSelectorPrefix) + vsm,
	}

	podList, err := pOps.List(lOpts)
	if err != nil {
		return nil, err
	}

	if podList == nil || len(podList.Items) == 0 {
		return nil, fmt.Errorf("VSM '%s:%s' not found. Orchestrator '%s:%s'", ns, vsm, k.Label(), k.Name())
	}

	pvl := &v1.PersistentVolumeList{
		Items: make([]v1.PersistentVolume, len(podList.Items)),
	}

	// TODO
	// Transform the POD type to persistent volume type
	// Do this in v1/k8s package ??
	for i, pod := range podList.Items {

		if pod.Name != vsm || pod.Namespace != ns {
			return nil, fmt.Errorf("VSM mismatch. Expected: '%s:%s' Found: '%s:%s' Orchestrator: '%s:%s'", ns, vsm, pod.Namespace, pod.Name, k.Label(), k.Name())
		}

		//pv := v1.PersistentVolume{}
		//pv.Name = pod.Name
		pvl.Items[i].Name = pod.Name
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

// createControllerPod creates a persistent volume controller deployment in
// kubernetes
func (k *k8sOrchestrator) createControllerPod(volProProfile volProfile.VolumeProvisionerProfile) (*k8sApiV1.Pod, error) {
	// fetch VSM name
	vsm, err := volProProfile.VSMName()
	if err != nil {
		return nil, err
	}

	cImg, imgSupport, err := volProProfile.ControllerImage()
	if err != nil {
		return nil, err
	}

	if !imgSupport {
		return nil, fmt.Errorf("VSM '%s' requires a controller container image", vsm)
	}

	k8sUtl := k8sOrchUtil(k, volProProfile)

	kc, supported := k8sUtl.K8sClient()

	if !supported {
		return nil, fmt.Errorf("K8s client not supported by '%s'", k8sUtl.Name())
	}

	// fetch pod operator
	pOps, err := kc.Pods()
	if err != nil {
		return nil, err
	}

	// create persistent volume controller as a k8s pod
	ctrl := &k8sApiV1.Pod{}
	ctrl.Kind = string(v1.K8sKindDeployment)
	ctrl.APIVersion = string(v1.K8sPodVersion)
	ctrl.Name = vsm + string(v1.ControllerSuffix)
	// Labels will provide the VSM filtering options during GET/LIST calls
	ctrl.Labels = map[string]string{
		string(v1.VSMIdentifier): vsm,
	}

	// specify the controller pod's container properties
	ctrlCon := k8sApiV1.Container{}
	ctrlCon.Name = vsm + string(v1.ControllerSuffix) + string(v1.ContainerSuffix)
	ctrlCon.Image = cImg
	ctrlCon.Command = v1.JivaCtrlCmd
	ctrlCon.Args = v1.JivaCtrlArgs

	iscsiPort := k8sApiV1.ContainerPort{}
	iscsiPort.ContainerPort = v1.DefaultJivaISCSIPort()

	apiPort := k8sApiV1.ContainerPort{}
	apiPort.ContainerPort = v1.DefaultJivaAPIPort()

	// Set the ports
	ctrlCon.Ports = []k8sApiV1.ContainerPort{iscsiPort, apiPort}

	ctrlSpec := k8sApiV1.PodSpec{}
	// Set the container
	ctrlSpec.Containers = []k8sApiV1.Container{ctrlCon}
	// Set the pod spec
	ctrl.Spec = ctrlSpec

	// add persistent volume controller pod
	return pOps.Create(ctrl)
}

// CreateReplicaPods creates one or more persistent volume replica(s)
// deployment in Kubernetes
func (k *k8sOrchestrator) CreateReplicaPods(volProProfile volProfile.VolumeProvisionerProfile, ctrlIP string) (*k8sApiV1.Pod, error) {
	// fetch VSM name
	vsm, err := volProProfile.VSMName()
	if err != nil {
		return nil, err
	}

	rImg, imgSupport, err := volProProfile.ReplicaImage()
	if err != nil {
		return nil, err
	}

	if !imgSupport {
		return nil, fmt.Errorf("VSM '%s' requires a replica container image", vsm)
	}

	rCount, err := volProProfile.ReplicaCount()
	if err != nil {
		return nil, err
	}

	pCount, err := volProProfile.PersistentPathCount()
	if err != nil {
		return nil, err
	}

	if pCount != rCount {
		return nil, fmt.Errorf("VSM '%s' replica count '%d' does not match persistent path count '%d'", vsm, rCount, pCount)
	}

	for i := 1; i <= rCount; i++ {
		_, err := k.createReplicaPod(volProProfile, ctrlIP, rImg, vsm, i, rCount)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil

}

// CreateReplicaPod creates a persistent volume replica deployment in Kubernetes
func (k *k8sOrchestrator) createReplicaPod(volProProfile volProfile.VolumeProvisionerProfile, ctrlIP string, rImg string, vsm string, position int, rCount int) (*k8sApiV1.Pod, error) {
	k8sUtl := k8sOrchUtil(k, volProProfile)

	kc, supported := k8sUtl.K8sClient()
	if !supported {
		return nil, fmt.Errorf("K8s client not supported by '%s'", k8sUtl.Name())
	}

	// fetch k8s pod operator
	pOps, err := kc.Pods()
	if err != nil {
		return nil, err
	}

	// Create persistent volume replica as a k8s pod
	rep := &k8sApiV1.Pod{}
	rep.Kind = string(v1.K8sKindDeployment)
	rep.APIVersion = string(v1.K8sPodVersion)
	rep.Name = vsm + string(v1.JivaReplicaSuffix) + string(position)
	// Labels will provide the VSM filtering options during GET/LIST calls
	rep.Labels = map[string]string{
		string(v1.VSMIdentifier): vsm,
	}

	// Create the replica pod's container
	repCon := k8sApiV1.Container{}
	repCon.Name = vsm + string(v1.JivaReplicaSuffix) + string(v1.ContainerSuffix) + string(position)
	repCon.Image = rImg
	repCon.Command = v1.JivaReplicaCmd

	pvc, err := volProProfile.PVC()
	if err != nil {
		return nil, err
	}

	repCon.Args = v1.MakeOrDefJivaReplicaArgs(pvc.Labels, ctrlIP)

	// Create replica pod ports
	repPort1 := k8sApiV1.ContainerPort{}
	repPort1.ContainerPort = v1.DefaultJivaReplicaPort1()

	repPort2 := k8sApiV1.ContainerPort{}
	repPort2.ContainerPort = v1.DefaultJivaReplicaPort2()

	repPort3 := k8sApiV1.ContainerPort{}
	repPort3.ContainerPort = v1.DefaultJivaReplicaPort3()

	// Set the ports at container
	repCon.Ports = []k8sApiV1.ContainerPort{repPort1, repPort2, repPort3}

	// Create replica pod volume mounts
	repMount := k8sApiV1.VolumeMount{}
	repMount.Name = v1.DefaultJivaMountName()
	repMount.MountPath = v1.DefaultJivaMountPath()

	// Set the mount paths at container
	repCon.VolumeMounts = []k8sApiV1.VolumeMount{repMount}

	// Create the replica pod's backing volume
	repVol := k8sApiV1.Volume{}
	repVol.Name = v1.DefaultJivaMountName()

	// Create the replica pod's host path
	hostPath := &k8sApiV1.HostPathVolumeSource{}
	persistPath, err := volProProfile.PersistentPath(position, rCount)
	if err != nil {
		return nil, err
	}
	hostPath.Path = persistPath

	// Set the host path
	repVol.HostPath = hostPath

	// Create the replica pod spec
	repSpec := k8sApiV1.PodSpec{}
	// Set the container at pod spec
	repSpec.Containers = []k8sApiV1.Container{repCon}
	// Set the backing volume at pod spec
	repSpec.Volumes = []k8sApiV1.Volume{repVol}
	// Set the pod spec at pod
	rep.Spec = repSpec

	// add persistent volume replica pod
	return pOps.Create(rep)
}

// createControllerService creates a persistent volume controller service in
// kubernetes
func (k *k8sOrchestrator) createControllerService(volProProfile volProfile.VolumeProvisionerProfile) (*k8sApiV1.Service, error) {
	// fetch VSM name
	vsm, err := volProProfile.VSMName()
	if err != nil {
		return nil, err
	}

	k8sUtl := k8sOrchUtil(k, volProProfile)

	kc, supported := k8sUtl.K8sClient()
	if !supported {
		return nil, fmt.Errorf("K8s client not supported by '%s'", k8sUtl.Name())
	}

	// fetch k8s clientset & namespace
	sOps, err := kc.Services()
	if err != nil {
		return nil, err
	}

	// create persistent volume controller service
	svc := &k8sApiV1.Service{}
	svc.Kind = string(v1.K8sKindService)
	svc.APIVersion = string(v1.K8sServiceVersion)
	svc.Name = vsm + string(v1.ControllerSuffix) + string(v1.ServiceSuffix)
	svc.Labels = map[string]string{
		string(v1.VSMIdentifier): vsm,
	}

	iscsiPort := k8sApiV1.ServicePort{}
	iscsiPort.Name = string(v1.PortNameISCSI)
	iscsiPort.Port = v1.DefaultJivaISCSIPort()

	apiPort := k8sApiV1.ServicePort{}
	apiPort.Name = string(v1.PortNameAPI)
	apiPort.Port = v1.DefaultJivaAPIPort()

	svcSpec := k8sApiV1.ServiceSpec{}
	svcSpec.Ports = []k8sApiV1.ServicePort{iscsiPort, apiPort}
	// Set the selector that identifies the controller VSM
	svcSpec.Selector = map[string]string{
		string(v1.VSMIdentifier): vsm + string(v1.ControllerSuffix),
	}

	// Set the service spec
	svc.Spec = svcSpec

	// add controller service
	return sOps.Create(svc)
}

// GetControllerService fetches the service name & service IP address
// of the persistent volume controller
func (k *k8sOrchestrator) GetControllerService(volProProfile volProfile.VolumeProvisionerProfile) (string, string, error) {
	// fetch VSM name
	vsm, err := volProProfile.VSMName()
	if err != nil {
		return "", "", err
	}

	k8sUtl := k8sOrchUtil(k, volProProfile)

	kc, supported := k8sUtl.K8sClient()
	if !supported {
		return "", "", fmt.Errorf("K8s client not supported by '%s'", k8sUtl.Name())
	}

	// fetch k8s service operations
	sOps, err := kc.Services()
	if err != nil {
		return "", "", err
	}

	svc, err := sOps.Get(vsm + string(v1.ControllerSuffix) + string(v1.ServiceSuffix))
	if err != nil {
		return "", "", err
	}

	return svc.Name, svc.Spec.ClusterIP, nil
}

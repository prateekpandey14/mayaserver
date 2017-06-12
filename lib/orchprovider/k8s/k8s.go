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
	"k8s.io/apimachinery/pkg/labels"
	k8sUnversioned "k8s.io/client-go/pkg/api/unversioned"
	k8sApiV1 "k8s.io/client-go/pkg/api/v1"
	k8sApisExtnsBeta1 "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

// TODO
// NOTE:
//    This does not work. The bootstraping is done in server/server.go
// If self-registration is not feasible then remove this altogether.
//
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
//func (k *k8sOrchestrator) AddStorage(volProProfile volProfile.VolumeProvisionerProfile) (*v1.PersistentVolumeList, error) {
func (k *k8sOrchestrator) AddStorage(volProProfile volProfile.VolumeProvisionerProfile) (*v1.PersistentVolume, error) {

	// TODO
	// This is jiva specific
	// Move this entire logic to a separate package that will couple jiva
	// provisioner with k8s orchestrator

	// create k8s pod of persistent volume controller
	_, err := k.createControllerDeployment(volProProfile)
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
	_, ctrlIP, err := k.getControllerServiceDetails(volProProfile)
	if err != nil {
		// TODO
		// Delete the persistent volume controller pod
		// Delegate to DeleteStorage which should handle stuff in a
		// robust way
		return nil, err
	}

	// Create the k8s deployment of vsm replica
	isRequested := volProProfile.ReqReplica()
	if !isRequested {
		return nil, nil
	}

	_, err = k.createDeploymentReplicas(volProProfile, ctrlIP)
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
//func (k *k8sOrchestrator) ReadStorage(volProProfile volProfile.VolumeProvisionerProfile) (*v1.PersistentVolumeList, error) {
func (k *k8sOrchestrator) ReadStorage(volProProfile volProfile.VolumeProvisionerProfile) (*v1.PersistentVolume, error) {
	if volProProfile == nil {
		return nil, fmt.Errorf("Nil volume provisioner profile provided")
	}

	vsm, err := volProProfile.VSMName()
	if err != nil {
		return nil, err
	}

	annotations := map[string]string{}

	// deployment(s) related details
	err = k.readFromDeployments(vsm, volProProfile, annotations)
	if err != nil {
		return nil, err
	}

	// service related details
	err = k.readFromService(vsm, volProProfile, annotations)
	if err != nil {
		return nil, err
	}

	pv := &v1.PersistentVolume{}
	pv.Name = vsm
	pv.Annotations = annotations

	return pv, nil
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

// createControllerDeployment creates a persistent volume controller deployment in
// kubernetes
func (k *k8sOrchestrator) createControllerDeployment(volProProfile volProfile.VolumeProvisionerProfile) (*k8sApisExtnsBeta1.Deployment, error) {
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

	// fetch deployment operator
	dOps, err := kc.DeploymentOps()
	if err != nil {
		return nil, err
	}

	glog.Infof("Adding vsm controller for vsm 'name: %s'", vsm)

	deploy := &k8sApisExtnsBeta1.Deployment{
		ObjectMeta: k8sApiV1.ObjectMeta{
			Name: vsm + string(v1.ControllerSuffix),
			Labels: map[string]string{
				string(v1.VSMIdentifier): vsm,
			},
		},
		TypeMeta: k8sUnversioned.TypeMeta{
			Kind:       string(v1.K8sKindDeployment),
			APIVersion: string(v1.K8sDeploymentVersion),
		},
		Spec: k8sApisExtnsBeta1.DeploymentSpec{
			//Replicas: &int32(1),
			Template: k8sApiV1.PodTemplateSpec{
				ObjectMeta: k8sApiV1.ObjectMeta{
					Labels: map[string]string{
						string(v1.VSMIdentifier): vsm,
					},
				},
				Spec: k8sApiV1.PodSpec{
					Containers: []k8sApiV1.Container{
						k8sApiV1.Container{
							Name:    vsm + string(v1.ControllerSuffix) + string(v1.ContainerSuffix),
							Image:   cImg,
							Command: v1.JivaCtrlCmd,
							Args:    v1.MakeOrDefJivaControllerArgs(vsm),
							Ports: []k8sApiV1.ContainerPort{
								k8sApiV1.ContainerPort{
									ContainerPort: v1.DefaultJivaISCSIPort(),
								},
								k8sApiV1.ContainerPort{
									ContainerPort: v1.DefaultJivaAPIPort(),
								},
							},
						},
					},
				},
			},
		},
	}

	// add persistent volume controller deployment
	dd, err := dOps.Create(deploy)
	if err != nil {
		return nil, err
	}

	glog.Infof("Added vsm controller for vsm 'name: %s' as K8s 'kind: %s' 'apiversion: %s'", deploy.Name, deploy.Kind, deploy.APIVersion)

	return dd, nil
}

// createDeploymentReplicas creates one or more persistent volume deployment
// replica(s) in Kubernetes
func (k *k8sOrchestrator) createDeploymentReplicas(volProProfile volProfile.VolumeProvisionerProfile, ctrlIP string) (*k8sApisExtnsBeta1.Deployment, error) {
	// fetch VSM name
	vsm, err := volProProfile.VSMName()
	if err != nil {
		return nil, err
	}

	if ctrlIP == "" {
		return nil, fmt.Errorf("VSM controller IP is required to create replica(s) for vsm 'name: %s'", vsm)
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

	pvc, err := volProProfile.PVC()
	if err != nil {
		return nil, err
	}

	// TODO
	// The position is always send as 1
	// We might want to get the replica index & send it
	// However, this does not matter if replicas are placed on different hosts !!
	persistPath, err := volProProfile.PersistentPath(1, rCount)
	if err != nil {
		return nil, err
	}

	k8sUtl := k8sOrchUtil(k, volProProfile)

	kc, supported := k8sUtl.K8sClient()
	if !supported {
		return nil, fmt.Errorf("K8s client not supported by '%s'", k8sUtl.Name())
	}

	// fetch k8s deployment operator
	dOps, err := kc.DeploymentOps()
	if err != nil {
		return nil, err
	}

	glog.Infof("Adding vsm replica(s) for vsm 'name: %s'", vsm)

	deploy := &k8sApisExtnsBeta1.Deployment{
		ObjectMeta: k8sApiV1.ObjectMeta{
			Name: vsm + string(v1.JivaReplicaSuffix),
			Labels: map[string]string{
				string(v1.VSMIdentifier): vsm,
			},
		},
		TypeMeta: k8sUnversioned.TypeMeta{
			Kind:       string(v1.K8sKindDeployment),
			APIVersion: string(v1.K8sDeploymentVersion),
		},
		Spec: k8sApisExtnsBeta1.DeploymentSpec{
			Replicas: v1.Replicas(rCount),
			Template: k8sApiV1.PodTemplateSpec{
				ObjectMeta: k8sApiV1.ObjectMeta{
					Labels: map[string]string{
						string(v1.VSMIdentifier): vsm,
					},
				},
				Spec: k8sApiV1.PodSpec{
					Containers: []k8sApiV1.Container{
						k8sApiV1.Container{
							Name:    vsm + string(v1.JivaReplicaSuffix) + string(v1.ContainerSuffix),
							Image:   rImg,
							Command: v1.JivaReplicaCmd,
							Args:    v1.MakeOrDefJivaReplicaArgs(pvc.Labels, ctrlIP),
							Ports: []k8sApiV1.ContainerPort{
								k8sApiV1.ContainerPort{
									ContainerPort: v1.DefaultJivaReplicaPort1(),
								},
								k8sApiV1.ContainerPort{
									ContainerPort: v1.DefaultJivaReplicaPort2(),
								},
								k8sApiV1.ContainerPort{
									ContainerPort: v1.DefaultJivaReplicaPort3(),
								},
							},
							VolumeMounts: []k8sApiV1.VolumeMount{
								k8sApiV1.VolumeMount{
									Name:      v1.DefaultJivaMountName(),
									MountPath: v1.DefaultJivaMountPath(),
								},
							},
						},
					},
					Volumes: []k8sApiV1.Volume{
						k8sApiV1.Volume{
							Name: v1.DefaultJivaMountName(),
							VolumeSource: k8sApiV1.VolumeSource{
								HostPath: &k8sApiV1.HostPathVolumeSource{
									Path: persistPath,
								},
							},
						},
					},
				},
			},
		},
	}

	dd, err := dOps.Create(deploy)
	if err != nil {
		return nil, err
	}

	glog.Infof("Added vsm replica(s) 'count: %d' for vsm 'name: %s' as K8s 'kind: %s' 'apiversion: %s'", deploy.Spec.Replicas, deploy.Name, deploy.Kind, deploy.APIVersion)

	return dd, nil
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

	// TODO
	// log levels & logging context to be taken care of
	glog.Infof("Adding service 'vsm controller: %s'", vsm)

	// TODO
	// Code this like a golang struct template
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
	ssvc, err := sOps.Create(svc)

	// TODO
	// log levels & logging context to be taken care of
	if err == nil {
		glog.Infof("Service added 'vsm controller: %s' 'apiversion: %s'", vsm, svc.APIVersion)
	}

	return ssvc, err
}

// getControllerServiceDetails fetches the service name & service IP address
// associated with the VSM
func (k *k8sOrchestrator) getControllerServiceDetails(volProProfile volProfile.VolumeProvisionerProfile) (string, string, error) {
	svc, err := k.getControllerService(volProProfile)
	if err != nil {
		return "", "", err
	}

	return svc.Name, svc.Spec.ClusterIP, nil
}

// getControllerService fetches the service associated with the VSM
func (k *k8sOrchestrator) getControllerService(volProProfile volProfile.VolumeProvisionerProfile) (*k8sApiV1.Service, error) {
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

	// fetch k8s service operations
	sOps, err := kc.Services()
	if err != nil {
		return nil, err
	}

	return sOps.Get(vsm + string(v1.ControllerSuffix) + string(v1.ServiceSuffix))
}

// TODO
// To be deprecated to a standard VSM type !!
//
// Transform a VSM to a PersistentVolume
func (k *k8sOrchestrator) readFromDeployments(vsm string, volProProfile volProfile.VolumeProvisionerProfile, annotations map[string]string) error {

	dl, err := k.getDeploymentList(volProProfile)
	if err != nil {
		return err
	}

	if dl == nil {
		return fmt.Errorf("No deployments were found for vsm 'name: %s'", vsm)
	}

	for _, deploy := range dl.Items {
		// w.r.t controller
		if deploy.Name == vsm+string(v1.ControllerSuffix) {
			SetCtrlDeployConditions(deploy, annotations)
		}
		// w.r.t replica
		if deploy.Name == vsm+string(v1.JivaReplicaSuffix) {
			SetReplDeployConditions(deploy, annotations)
			SetReplIPs(deploy, annotations)
			SetReplCount(deploy, annotations)
			SetReplVolumeSize(deploy, annotations)
			SetIQN(vsm, deploy, annotations)
		}
	}

	return nil
}

func (k *k8sOrchestrator) readFromService(vsm string, volProProfile volProfile.VolumeProvisionerProfile, annotations map[string]string) error {
	// w.r.t service
	svc, err := k.getControllerService(volProProfile)
	if err != nil {
		return err
	}

	SetServiceStatus(svc, annotations)
	SetISCSITargetPortal(svc, annotations)

	return nil
}

// getDeploymentList fetches the deployments associated with the VSM
func (k *k8sOrchestrator) getDeploymentList(volProProfile volProfile.VolumeProvisionerProfile) (*k8sApisExtnsBeta1.DeploymentList, error) {
	// NOTE:
	//    A VSM can be one or more k8s deployments
	//
	// NOTE:
	//    maya api service assigns the VSM name as one of the labels against all
	// the pods created during creation of persistent volume
	vsm, err := volProProfile.VSMName()
	if err != nil {
		return nil, err
	}

	k8sUtl := k8sOrchUtil(k, volProProfile)

	kc, supported := k8sUtl.K8sClient()
	if !supported {
		return nil, fmt.Errorf("K8s client not supported by '%s'", k8sUtl.Name())
	}

	ns, err := kc.NS()
	if err != nil {
		return nil, err
	}

	dOps, err := kc.DeploymentOps()
	if err != nil {
		return nil, err
	}

	// TODO
	// This filtering logic DOES NOT WORK with client-go v2.0.0
	// Need to upgrade client-go to latest stable version for this to work
	lOpts := k8sApiV1.ListOptions{
		LabelSelector: string(v1.VSMSelectorPrefix) + vsm,
	}

	deployList, err := dOps.List(lOpts)
	if err != nil {
		return nil, err
	}

	if deployList == nil {
		return nil, fmt.Errorf("VSM(s) '%s:%s' not found at orchestrator '%s:%s'", ns, vsm, k.Label(), k.Name())
	}

	// NOTE:
	//    Workaround for above filtering logic that does not work
	// Get the filtered pod list based on expected label
	eLblStr := string(v1.VSMSelectorPrefix) + vsm
	eLbl, err := labels.Parse(eLblStr)
	if err != nil {
		return nil, err
	}

	// filtered deployment list
	// I believe above filtering does not work
	fdl := &k8sApisExtnsBeta1.DeploymentList{}

	for _, item := range deployList.Items {
		if eLbl.Matches(labels.Set(item.Labels)) {
			fdl.Items = append(fdl.Items, item)
		}
	}

	if fdl == nil || len(fdl.Items) == 0 {
		return nil, fmt.Errorf("VSM(s) '%s:%s' not found at orchestrator '%s:%s'", ns, vsm, k.Label(), k.Name())
	}

	return fdl, nil
}

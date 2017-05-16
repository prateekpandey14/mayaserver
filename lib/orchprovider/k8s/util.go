package k8s

import (
	"fmt"

	orchProfile "github.com/openebs/mayaserver/lib/profile/orchprovider"
	volProfile "github.com/openebs/mayaserver/lib/profile/volumeprovisioner"
	"k8s.io/client-go/kubernetes"
	k8sCoreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

// K8sUtilGetter is an abstraction to fetch instances of K8sUtilInterface
type K8sUtilGetter interface {
	K8sUtil(volProfile.VolumeProvisionerProfile) K8sUtilInterface
}

// K8sUtilInterface is an abstraction over communicating with K8s APIs
type K8sUtilInterface interface {
	// Name of K8s utility
	Name() string

	// K8sClient fetches an instance of K8sClients. Will return
	// false if the util does not support providing K8sClients instance.
	K8sClient() (K8sClient, bool)
}

// K8sClient is an abstraction to operate on various k8s entities.
//
// NOTE:
//    This abstraction makes use of K8s's client-go package.
type K8sClient interface {
	// InCluster indicates whether the operation is within cluster or in a
	// different cluster
	InCluster() (bool, error)

	// NS provides the namespace where operations will be executed
	NS() (string, error)

	// Pods provides all the CRUD operations associated w.r.t a POD
	Pods() (k8sCoreV1.PodInterface, error)

	// Services provides all the CRUD operations associated w.r.t a Service
	Services() (k8sCoreV1.ServiceInterface, error)
}

// k8sUtil provides the concrete implementation for below interfaces:
//
// 1. k8s.K8sUtilInterface interface
// 2. k8s.K8sClients interface
type k8sUtil struct {
	caCert     string
	caPath     string
	clientCert string
	clientKey  string
	insecure   bool

	// volProfile has context related information w.r.t k8s
	volProfile volProfile.VolumeProvisionerProfile
}

// This is a plain k8s utility & hence the name
func (k *k8sUtil) Name() string {
	ns, _ := k.NS()
	return fmt.Sprintf("k8sutil @ '%s'", ns)
}

// k8sUtil implements K8sClient interface. Hence it returns
// self
func (k *k8sUtil) K8sClient() (K8sClient, bool) {
	return k, true
}

// NS provides the namespace where operations will be executed
func (k *k8sUtil) NS() (string, error) {
	if nil == k.volProfile {
		return "", fmt.Errorf("Volume provisioner profile not initialized at '%s'", k.Name())
	}

	// Fetch pvc from volume provisioner profile
	pvc, err := k.volProfile.PVC()
	if err != nil {
		return "", err
	}

	// Get orchestrator provider profile from pvc
	oPrfle, err := orchProfile.GetOrchProviderProfileByPVC(pvc)
	if err != nil {
		return "", err
	}

	// Get the namespace which will be queried
	ns, err := oPrfle.NS()
	if err != nil {
		return "", err
	}

	return ns, nil
}

// InCluster indicates whether the operation is within cluster or in a
// different cluster
func (k *k8sUtil) InCluster() (bool, error) {
	if nil == k.volProfile {
		return false, fmt.Errorf("Volume provisioner profile not initialized at '%s'", k.Name())
	}

	// Fetch pvc from volume provisioner profile
	pvc, err := k.volProfile.PVC()
	if err != nil {
		return false, err
	}

	// Get orchestrator provider profile from pvc
	oPrfle, err := orchProfile.GetOrchProviderProfileByPVC(pvc)
	if err != nil {
		return false, err
	}

	// Which kind of request ? in-cluster or out-of-cluster ?
	isInCluster, err := oPrfle.InCluster()
	if err != nil {
		return false, err
	}

	return isInCluster, nil
}

// Pods is a utility function that provides a instance capable of
// executing various k8s pod related operations.
func (k *k8sUtil) Pods() (k8sCoreV1.PodInterface, error) {
	var cs *kubernetes.Clientset

	inC, err := k.InCluster()
	if err != nil {
		return nil, err
	}

	ns, err := k.NS()
	if err != nil {
		return nil, err
	}

	if inC {
		cs, err = k.inClusterCS()
	} else {
		cs, err = k.outClusterCS()
	}

	if err != nil {
		return nil, err
	}

	return cs.CoreV1().Pods(ns), nil
}

// Services is a utility function that provides a instance capable of
// executing various k8s service related operations.
func (k *k8sUtil) Services() (k8sCoreV1.ServiceInterface, error) {
	var cs *kubernetes.Clientset

	inC, err := k.InCluster()
	if err != nil {
		return nil, err
	}

	ns, err := k.NS()
	if err != nil {
		return nil, err
	}

	if inC {
		cs, err = k.inClusterCS()
	} else {
		cs, err = k.outClusterCS()
	}

	if err != nil {
		return nil, err
	}

	return cs.CoreV1().Services(ns), nil
}

// inClusterCS is used to initialize and return a new http client capable
// of invoking K8s APIs.
func (k *k8sUtil) inClusterCS() (*kubernetes.Clientset, error) {

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	// creates the in-cluster clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

// outClusterCS is used to initialize and return a new http client capable
// of invoking outside the cluster K8s APIs.
func (k *k8sUtil) outClusterCS() (*kubernetes.Clientset, error) {
	return nil, fmt.Errorf("outClusterCS not supported in '%s'", k.Name())
}

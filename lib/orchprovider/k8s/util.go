package k8s

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// K8sUtilInterface is an abstraction over communicating with K8s APIs
type K8sUtilInterface interface {
	// Name of K8s utility
	Name() string

	// K8sClients fetches an instance of K8sClients. Will return
	// false if the util does not support providing K8sClients instance.
	K8sClients() (K8sClients, bool)
}

// K8sClients is an abstraction over various connection modes (http, rpc)
// to K8s. Http client is currently supported.
//
// NOTE:
//    This abstraction makes use of K8s's client-go package.
type K8sClients interface {

	// ClientSet is capable to communicate with an in-cluster K8s
	// or outside the cluster depending on the passed flag.
	GetClusterCS(bool) (*kubernetes.Clientset, error)
}

// k8sUtil provides the concrete implementation for below interfaces:
//
// 1. k8s.K8sUtilInterface interface
// 2. k8s.K8sClients interface
type k8sUtil struct {
	// K8s server / cluster coordinates
	//k8sConf *K8sConfig

	caCert     string
	caPath     string
	clientCert string
	clientKey  string
	insecure   bool
}

// newK8sUtil provides a new instance of k8sUtil
func newK8sUtil() (K8sUtilInterface, error) {
	return &k8sUtil{}, nil
}

// This is a plain k8s utility & hence the name
func (k *k8sUtil) Name() string {
	return "k8sutil"
}

// k8sUtil implements K8sClients interface. Hence it returns
// self
func (k *k8sUtil) K8sClients() (K8sClients, bool) {
	return k, true
}

// GetClusterCS is a utility function to get clientset capable of communicating
// with k8s APIs.
func (k *k8sUtil) GetClusterCS(incluster bool) (*kubernetes.Clientset, error) {
	if incluster {
		return k.inClusterClientSet()
	} else {
		return k.outClusterClientSet()
	}
}

// ClientSet is used to initialize and return a new http client capable
// of invoking K8s APIs.
func (k *k8sUtil) inClusterClientSet() (*kubernetes.Clientset, error) {

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	// creates the in cluster clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

// OutClusterClientSet is used to initialize and return a new http client capable
// of invoking outside the cluster K8s APIs.
func (k *k8sUtil) outClusterClientSet() (*kubernetes.Clientset, error) {
	return nil, fmt.Errorf("OutClusterClientSet not supported in '%s'", k.Name())
}

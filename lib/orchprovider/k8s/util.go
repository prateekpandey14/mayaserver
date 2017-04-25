package k8s

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
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
	// Http returns the http client that is capable to communicate
	// with K8s
	// Http() (*a k8s Client, error)
	Http() error
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

// k8sUtil implements NomadClients interface. Hence it returns
// self
func (k *k8sUtil) K8sClients() (K8sClients, bool) {
	return k, true
}

// Http is used to initialize and return a new http client capable
// of invoking K8s APIs.
func (k *k8sUtil) Http() error {

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// This is a sample code to register the vendored dependencies due to k8s client-go
	pods, err := clientset.CoreV1().Pods("").List(v1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	return nil
}

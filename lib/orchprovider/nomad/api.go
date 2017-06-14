// This file transforms a Nomad scheduler as an orchestration
// platform for persistent volume placement. OpenEBS calls this as
// placement of storage pod.
package nomad

import (
	"fmt"

	"github.com/hashicorp/nomad/api"
	"github.com/openebs/mayaserver/lib/api/v1"
)

// NomadApiInterface provides a means to issue APIs against a Nomad cluster.
// These APIs are futher categorized into Networking & Storage specific APIs.
type NomadApiInterface interface {

	// Name of the Nomad API implementor
	Name() string

	// This returns a concrete implementation of NetworkApis
	NetworkApis() (NetworkApis, bool)

	// This returns a concrete implementation of StorageApis
	StorageApis() (StorageApis, bool)
}

// nomadApi is an implementation of
//
//  1. nomad.NomadApiInterface interface
//  2. nomad.NetworkApis interface
//  3. nomad.StorageApis interface
//
// It is composed of NomadUtilInterface
type nomadApi struct {
	nUtil NomadUtilInterface
}

// newNomadApi provides a new instance of nomadApi
func newNomadApi(nConfig *NomadConfig) (*nomadApi, error) {

	nUtil, err := newNomadUtil(nConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to create nomad api instance")
	}

	return &nomadApi{
		nUtil: nUtil,
	}, nil
}

// This is a plain nomad api implementor & hence the name
func (n *nomadApi) Name() string {
	return "nomadapi"
}

// nomadApi implements NetworkApis, hence it returns self.
func (n *nomadApi) NetworkApis() (NetworkApis, bool) {
	return n, true
}

// nomadApi implements StorageApis, hence it returns self.
func (n *nomadApi) StorageApis() (StorageApis, bool) {
	return n, true
}

// StorageApis provides a means to communicate with Nomad Apis
// w.r.t storage.
//
// NOTE:
//    A Nomad job spec is treated as a persistent volume storage
// spec & then submitted to a Nomad deployment.
//
// NOTE:
//    Nomad has no notion of Persistent Volume.
type StorageApis interface {
	// Create makes a request to Nomad to create a storage resource
	CreateStorage(job *api.Job) (*api.Evaluation, error)

	// Delete makes a request to Nomad to delete the storage resource
	DeleteStorage(job *api.Job) (*api.Evaluation, error)

	// Info provides the storage information w.r.t the provided job name
	StorageInfo(jobName string) (*api.Job, error)

	// StorageProps fetches appropriate persistent storage props that
	// is assumed to be supported at deployed Nomad environment.
	//StorageProps(dc string) (map[v1.ContainerStorageLbl]string, error)
	StorageProps(dc string) (map[v1.VolumeProvisionerProfileLabel]string, error)
}

// Fetch persistent storage properties that is supported at the deployed Nomad
// environment.
//func (n *nomadApi) StorageProps(dc string) (map[v1.ContainerStorageLbl]string, error) {
func (n *nomadApi) StorageProps(dc string) (map[v1.VolumeProvisionerProfileLabel]string, error) {

	nUtil := n.nUtil
	if nUtil == nil {
		return nil, fmt.Errorf("Nomad utility not initialized")
	}

	nStorages, ok := nUtil.NomadStorages()
	if !ok {
		return nil, fmt.Errorf("Nomad storages not supported by nomad utility '%s'", nUtil.Name())
	}

	ns, err := nStorages.CS(dc)
	if err != nil {
		return nil, err
	}

	return ns, nil
}

// Fetch info about a particular resource in Nomad cluster.
//
// NOTE:
//    Nomad does not have persistent volume as its first class citizen.
// Hence, this resource should exhibit storage characteristics. The validations
// for this should have been done at the volume plugin implementation.
func (n *nomadApi) StorageInfo(jobName string) (*api.Job, error) {

	nUtil := n.nUtil
	if nUtil == nil {
		return nil, fmt.Errorf("Nomad utility not initialized")
	}

	nClients, ok := nUtil.NomadClients()
	if !ok {
		return nil, fmt.Errorf("Nomad clients not supported by nomad utility '%s'", nUtil.Name())
	}

	nHttpClient, err := nClients.Http()
	if err != nil {
		return nil, err
	}

	// Fetch the job info
	job, _, err := nHttpClient.Jobs().Info(jobName, &api.QueryOptions{})

	if err != nil {
		return nil, err
	}

	return job, nil
}

// Creates a resource in Nomad cluster.
//
// NOTE:
//    Nomad does not have persistent volume as its first class citizen.
// Hence, this resource should exhibit storage characteristics. The validations
// for this should have been done at the volume plugin implementation.
func (n *nomadApi) CreateStorage(job *api.Job) (*api.Evaluation, error) {

	nUtil := n.nUtil
	if nUtil == nil {
		return nil, fmt.Errorf("Nomad utility not initialized")
	}

	nClients, ok := nUtil.NomadClients()
	if !ok {
		return nil, fmt.Errorf("Nomad clients not supported by nomad utility '%s'", nUtil.Name())
	}

	nHttpClient, err := nClients.Http()
	if err != nil {
		return nil, err
	}

	// Register a job & get its evaluation id
	evalID, _, err := nHttpClient.Jobs().Register(job, &api.WriteOptions{})

	if err != nil {
		return nil, err
	}

	// Get the evaluation details
	eval, _, err := nHttpClient.Evaluations().Info(evalID, &api.QueryOptions{})

	if err != nil {
		return nil, err
	}

	return eval, nil
}

// Remove a resource in Nomad cluster.
//
// NOTE:
//    Nomad does not have persistent volume as its first class citizen.
// Hence, this resource should exhibit storage characteristics. The validations
// for this should have been done at the volume plugin implementation.
func (n *nomadApi) DeleteStorage(job *api.Job) (*api.Evaluation, error) {

	nUtil := n.nUtil
	if nUtil == nil {
		return nil, fmt.Errorf("Nomad utility not initialized")
	}

	nClients, ok := nUtil.NomadClients()
	if !ok {
		return nil, fmt.Errorf("Nomad clients not supported by nomad utility '%s'", nUtil.Name())
	}

	nHttpClient, err := nClients.Http()
	if err != nil {
		return nil, err
	}

	evalID, _, err := nHttpClient.Jobs().Deregister(*job.Name, &api.WriteOptions{})

	if err != nil {
		return nil, err
	}

	eval, _, err := nHttpClient.Evaluations().Info(evalID, &api.QueryOptions{})
	if err != nil {
		return nil, err
	}

	return eval, nil
}

// NetworkApis provides a means to issue Nomad Apis
// w.r.t network.
//
// NOTE:
//    Nomad has no notion of Networking.
type NetworkApis interface {

	// NetworkInfo fetches appropriate container networking values that
	// is assumed to be supported at deployed Nomad environment.
	NetworkProps(dc string) (map[v1.ContainerNetworkingLbl]string, error)
}

// Fetch networking information that is supported at the deployed Nomad
// environment.
func (n *nomadApi) NetworkProps(dc string) (map[v1.ContainerNetworkingLbl]string, error) {

	nUtil := n.nUtil
	if nUtil == nil {
		return nil, fmt.Errorf("Nomad utility not initialized")
	}

	nNetworks, ok := nUtil.NomadNetworks()
	if !ok {
		return nil, fmt.Errorf("Nomad networks not supported by nomad utility '%s'", nUtil.Name())
	}

	nc, err := nNetworks.CN(dc)
	if err != nil {
		return nil, err
	}

	return nc, nil
}

// This file plugs the following:
//
//    1. Generic orchprovider &
//    2. Nomad orchprovider
package nomad

import (
	"fmt"
	"io"

	"github.com/golang/glog"
	"github.com/openebs/mayaserver/lib/api/v1"
	v1nomad "github.com/openebs/mayaserver/lib/api/v1/nomad"
	"github.com/openebs/mayaserver/lib/orchprovider"
)

// The registration logic for the default nomad orchestrator plugin
//
// NOTE:
//    This function is executed once per application
//
// NOTE:
//    Registration & Initialization are two different workflows.
//  Registration happens here while Initialization is a lazy logic that happens
//  at runtime with appropriate values.
func init() {
	orchprovider.RegisterOrchProvider(
		// A variant of nomad orchestrator plugin
		v1nomad.DefaultNomadPluginName,
		// Below is a functional implementation that holds the initialization
		// logic of nomad orchestrator plugin
		//
		// TODO
		// The region, config options should not be used as they are meant to be
		// dynamic. They should be available as profiles. Once that happens, any
		// instance of nomad orchestrator should do these jobs. Hence, variants will
		// no longer be required.
		func(name string, region string, config io.Reader) (orchprovider.OrchestratorInterface, error) {
			return NewNomadOrchestrator(name, region, config)
		})
}

// NomadOrchestrator is a concrete representation of following
// interfaces:
//
//  1. orchprovider.OrchestratorInterface,
//  2. orchprovider.NetworkPlacements, &
//  3. orchprovider.StoragePlacements
type NomadOrchestrator struct {

	// Name of this orchestrator
	name string

	// The region where this orchestrator is deployed
	// This is set during the initilization time.
	region string

	// nStorApis represents an instance capable of invoking
	// storage related APIs
	nStorApis StorageApis

	// nNetApis represents an instance capable of invoking
	// network related APIs
	nNetApis NetworkApis

	// nConfig represents an instance that provides the coordinates
	// of a Nomad server / cluster deployment.
	nConfig *NomadConfig
}

// NewNomadOrchestrator provides a new instance of NomadOrchestrator. This is
// invoked during binary startup.
func NewNomadOrchestrator(name string, region string, config io.Reader) (orchprovider.OrchestratorInterface, error) {

	glog.Infof("Building nomad orchestration provider")

	if name == "" {
		return nil, fmt.Errorf("Name missing while building nomad orchestrator")
	}

	if region == "" {
		return nil, fmt.Errorf("Region missing while building nomad orchestrator")
	}

	// Transform the Reader to a NomadConfig
	nCfg, err := readNomadConfig(config)
	if err != nil {
		return nil, fmt.Errorf("Unable to read Nomad orchestrator's config: %v", err)
	}

	// TODO
	// validations of the populated config structure

	// Get a new instance of Nomad API
	nApi, err := newNomadApi(nCfg)
	if err != nil {
		return nil, err
	}

	// Get Nomad's storage specific API provider
	nStorApis, ok := nApi.StorageApis()
	if !ok {
		return nil, fmt.Errorf("Storage APIs not supported in nomad api instance '%s'", nApi.Name())
	}

	nNetApis, ok := nApi.NetworkApis()
	if !ok {
		return nil, fmt.Errorf("Network APIs not supported in nomad api instance '%s'", nApi.Name())
	}

	// build the orchestrator instance
	nOrch := &NomadOrchestrator{
		nStorApis: nStorApis,
		nNetApis:  nNetApis,
		nConfig:   nCfg,
		region:    region,
		name:      name,
	}

	return nOrch, nil
}

// Name provides the name of this orchestrator.
// This is an implementation of the orchprovider.OrchestratorInterface interface.
func (n *NomadOrchestrator) Name() string {

	return n.name
}

// Region provides the region where this orchestrator is running.
// This is an implementation of the orchprovider.OrchestratorInterface interface.
func (n *NomadOrchestrator) Region() string {

	return n.region
}

// StoragePlacements is this orchestration provider's
// implementation of the orchprovider.OrchestratorInterface interface.
func (n *NomadOrchestrator) StoragePlacements() (orchprovider.StoragePlacements, bool) {

	return n, true
}

// NetworkPlacements is this orchestration provider's
// implementation of the orchprovider.OrchestratorInterface interface.
//
// TODO
// This implementation will not be required once maya api server implements
// orchestrator provider specific profiles.
func (n *NomadOrchestrator) NetworkPlacements() (orchprovider.NetworkPlacements, bool) {

	return n, true
}

// NetworkPropsReq is a contract method implementation of
// orchprovider.NetworkPlacements interface. In this implementation,
// network resource details will be fetched from a Nomad deployment.
func (n *NomadOrchestrator) NetworkPropsReq(dc string) (map[v1.ContainerNetworkingLbl]string, error) {

	return n.nNetApis.NetworkProps(dc)
}

// StoragePropsReq is a contract method implementation of
// orchprovider.StoragePlacements interface. In this implementation,
// persistent storage details will be fetched from a Nomad deployment.
//
// TODO
// This implementation will not be required once maya api server implements
// orchestrator provider specific profiles.
func (n *NomadOrchestrator) StoragePropsReq(dc string) (map[v1.ContainerStorageLbl]string, error) {

	return n.nStorApis.StorageProps(dc)
}

// StorageInfoReq is a contract method implementation of
// orchprovider.StoragePlacements interface. In this implementation,
// a resource details will be fetched from a Nomad deployment.
//
// NOTE:
//    Nomad does not have persistent volume as its first class citizen.
// Hence, this resource should exhibit storage characteristics. The validations
// for this should have been done at the volume plugin implementation.
func (n *NomadOrchestrator) StorageInfoReq(pvc *v1.PersistentVolumeClaim) (*v1.PersistentVolume, error) {

	jobName, err := PvcToJobName(pvc)
	if err != nil {
		return nil, err
	}

	job, err := n.nStorApis.StorageInfo(jobName)
	if err != nil {
		return nil, err
	}

	return JobToPv(job)
}

// StoragePlacementReq is a contract method implementation of
// orchprovider.StoragePlacements interface. In this implementation,
// a resource will be created at a Nomad deployment.
//
// NOTE:
//    Nomad does not have persistent volume as its first class citizen.
// Hence, this resource should exhibit storage characteristics. The validations
// for this should have been done at the volume plugin implementation.
func (n *NomadOrchestrator) StoragePlacementReq(pvc *v1.PersistentVolumeClaim) (*v1.PersistentVolume, error) {

	// TODO
	// Check for the presence of region
	// If region is present then:
	//    1. fetch the NomadConfig applicable for this region
	//    2. build the StorageApiClient
	//    3. add these to a map in NomadOrchestrator struct
	// If region is not present, then use the available region set against this
	// NomadOrchestrator instance
	// Set the pvc's Label with the region property

	// TODO
	// Check for the presence of datacenter
	// If datacenter is present then fetch its StorageApiClient
	// Set the pvc's Label with the datacenter property

	job, err := PvcToJob(pvc)
	if err != nil {
		return nil, err
	}

	eval, err := n.nStorApis.CreateStorage(job)
	if err != nil {
		return nil, err
	}

	glog.V(2).Infof("Volume '%s' was placed for provisioning with eval '%v'", *job.Name, eval)

	return JobEvalToPv(*job.Name, eval)
}

// StorageRemovalReq is a contract method implementation of
// orchprovider.StoragePlacements interface. In this implementation,
// the resource will be removed from the Nomad deployment.
//
// NOTE:
//    Nomad does not have persistent volume as its first class citizen.
// Hence, this resource should exhibit storage characteristics. The validations
// for this should have been done at the volume plugin implementation.
func (n *NomadOrchestrator) StorageRemovalReq(pv *v1.PersistentVolume) (*v1.PersistentVolume, error) {

	job, err := PvToJob(pv)
	if err != nil {
		return nil, err
	}

	eval, err := n.nStorApis.DeleteStorage(job)

	if err != nil {
		return nil, err
	}

	glog.V(2).Infof("Volume '%s' was placed for removal with eval '%v'", pv.Name, eval)

	return JobEvalToPv(*job.Name, eval)
}

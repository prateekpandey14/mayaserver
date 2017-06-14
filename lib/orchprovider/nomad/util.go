package nomad

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/golang/glog"
	"github.com/hashicorp/nomad/api"
	"github.com/openebs/mayaserver/lib/api/v1"
	v1nomad "github.com/openebs/mayaserver/lib/api/v1/nomad"
	gcfg "gopkg.in/gcfg.v1"
)

const (
	// Names of environment variables used to supply the coordinates
	// of a Nomad deployment
	EnvNomadAddress = "NOMAD_ADDR"
	EnvNomadRegion  = "NOMAD_REGION"
)

// NomadConfig provides the settings that has the coordinates of a
// Nomad server or a Nomad cluster deployment.
//
// In addition, it provides below:
//
// 1. networking options that can be consumed by the storage app container
// spawned inside a Nomad cluster.
//
// 2. storage options that can be consumed by the storage app container
// spawned inside a Nomad cluster.
//
// A NomadConfig file is a .INI extension file.
// NOTE:
//    This is as per gcfg lib's conventions
//
// Below is a sample:
//
// [datacenter "dc1"]
// ; Address of Nomad deployment
// address = http://20.0.0.2:4646
//
// ; Container Networking options
// cn-type = host
// cn-network-cidr = 172.28.128.1/24
// cn-interface = enp0s8
//
// ; Container Storage options
// cs-persistence-location = /tmp/
// cs-replica-count = 2
//
// TODO
// It is planned to remove this Config entirely. The values present here
// are mostly dynamic in nature & will differ from request to request.
// It might be well to label these under orch provider profile(s) stored in some
// DB &/ generated at runtime.
type NomadConfig struct {
	Datacenter map[string]*struct {
		// Address of Nomad cluster
		Address string

		// Whether it is a host based networking or something else
		// Required while placing a container inside Nomad
		CNType string `gcfg:"cn-type"`

		// The Network address in CIDR notation. Available IP addresses will
		// be considered from this network range.
		CNNetworkCIDR string `gcfg:"cn-network-cidr"`

		// Networking interface that is available in the Nomad cluster
		CNInterface string `gcfg:"cn-interface"`

		// The backing persistent storage location on which
		// containerized storage is supposed to operate
		CSPersistenceLocation string `gcfg:"cs-persistence-location"`

		// CSReplicaCount holds the default no. of storage replicas
		CSReplicaCount string `gcfg:"cs-replica-count"`
	}
}

// NomadUtilInterface is an abstraction over
//
// 1.   Hashicorp's Nomad properties & communication utilities.
// 2.   Networking options available at/derived from Nomad cluster.
// 3.   Storage options available at/derived from Nomad cluster.
type NomadUtilInterface interface {

	// Name of nomad utility
	Name() string

	// This is a builder for NomadClients interface. Will return
	// false if not supported.
	NomadClients() (NomadClients, bool)

	// This is a builder for NomadNetworks interface. Will return
	// false if not supported.
	//
	// TODO
	// This interface will not be required once maya api server implements
	// orchestrator provider specific profiles.
	NomadNetworks() (NomadNetworks, bool)

	// This is a builder for NomadStorages interface. Will return
	// false if not supported.
	//
	// TODO
	// This interface will not be required once maya api server implements
	// orchestrator provider specific profiles.
	NomadStorages() (NomadStorages, bool)
}

// NomadClients is an abstraction over various connection modes (http, rpc)
// to Nomad. Http client is currently supported.
//
// NOTE:
//    This abstraction makes use of Nomad's api package. Nomad's api
// package provides:
//
// 1. http client abstraction &
// 2. structures that can send http requests to Nomad's APIs.
type NomadClients interface {
	// Http returns the http client that is capable to communicate
	// with Nomad
	Http() (*api.Client, error)
}

// NomadNetworks is a blueprint to expose various networking options
// available in a Nomad cluster.
type NomadNetworks interface {
	// CN exposes various networking values that is supported at a
	// particular datacenter where Nomad is running
	CN(dc string) (map[v1.ContainerNetworkingLbl]string, error)
}

// NomadStorages is a blueprint to expose various persistence storage
// options available in a Nomad cluster.
type NomadStorages interface {
	// CS exposes various persistence storage options that is supported at a
	// particular datacenter where Nomad is running
	CS(dc string) (map[v1.VolumeProvisionerProfileLabel]string, error)
}

// nomadUtil is the concrete implementation for
//
// 1. nomad.NomadClients interface
// 2. nomad.NomadNetworks interface
type nomadUtil struct {

	// The region to send API requests to
	// TODO
	// This will be set during this instance creation time
	region string

	// The datacenter to send API requests to
	// TODO
	// This will be set during this instance creation time
	datacenter string

	// Nomad server / cluster coordinates
	// This will be set based on the region
	nomadConf *NomadConfig

	caCert     string
	caPath     string
	clientCert string
	clientKey  string
	insecure   bool
}

// newNomadUtil provides a new instance of nomadUtil
//
// TODO
// region may be passed as an argument
// & hence NomadConfig should be instantiated based on the region
// at this place
func newNomadUtil(nConfig *NomadConfig) (*nomadUtil, error) {
	return &nomadUtil{
		nomadConf: nConfig,
	}, nil
}

// This is a plain nomad utility & hence the name
func (m *nomadUtil) Name() string {
	return "nomadutil"
}

// nomadUtil implements NomadClients interface. Hence it returns
// self
func (m *nomadUtil) NomadClients() (NomadClients, bool) {
	return m, true
}

// nomadUtil implements NomadNetworks interface. Hence it returns
// self
func (m *nomadUtil) NomadNetworks() (NomadNetworks, bool) {
	return m, true
}

// nomadUtil implements NomadStorages interface. Hence it returns
// self
func (m *nomadUtil) NomadStorages() (NomadStorages, bool) {
	return m, true
}

// Client is used to initialize and return a new API client capable
// of calling Nomad APIs.
// TODO
// datacenter may be passed as a parameter
func (m *nomadUtil) Http() (*api.Client, error) {
	// Nomad API client config
	apiCConf := api.DefaultConfig()

	// Set from environment variable
	val, found := os.LookupEnv(EnvNomadAddress)

	if !found {
		glog.V(2).Infof("Env variable '%s' is not set", EnvNomadAddress)
	}

	if val != "" {
		glog.V(2).Infof("Nomad address is set to '%s' via env var", val)
		apiCConf.Address = val
	}

	// Override from conf structure
	if m.nomadConf != nil && m.nomadConf.Datacenter != nil {
		// TODO
		// Derive the datacenter at runtime
		// Remove the region & datacenter properties from Mayaconfig
		glog.V(2).Infof("Nomad address is set to: '%s' via conf", m.nomadConf.Datacenter["dc1"].Address)
		apiCConf.Address = m.nomadConf.Datacenter["dc1"].Address
	}

	if apiCConf.Address == "" {
		return nil, fmt.Errorf("Nomad address is not set")
	}

	glog.V(2).Infof("Nomad will be reached at: '%s'", apiCConf.Address)

	if v := os.Getenv(EnvNomadRegion); v != "" {
		apiCConf.Region = v
	}

	if m.region != "" {
		apiCConf.Region = m.region
	}

	// If we need custom TLS configuration, then set it
	if m.caCert != "" || m.caPath != "" || m.clientCert != "" || m.clientKey != "" || m.insecure {
		t := &api.TLSConfig{
			CACert:     m.caCert,
			CAPath:     m.caPath,
			ClientCert: m.clientCert,
			ClientKey:  m.clientKey,
			Insecure:   m.insecure,
		}
		apiCConf.TLSConfig = t
	}

	// This has the http address & authentication details
	// required to invoke Nomad APIs
	return api.NewClient(apiCConf)
}

// CN provides the container networking data in key-value pairs.
// These networking data are supposed to be available in the target Nomad
// cluster. These pairs are provided based on datacenter.
func (m *nomadUtil) CN(dcName string) (map[v1.ContainerNetworkingLbl]string, error) {

	err := m.validateConf(dcName)
	if err != nil {
		return nil, err
	}

	// build the cn map
	cn := map[v1.ContainerNetworkingLbl]string{
		// container networking properties
		v1.CNTypeLbl:            m.getCNType(dcName),
		v1.CNNetworkCIDRAddrLbl: m.getCNNetworkCIDR(dcName),
		v1.CNInterfaceLbl:       m.getCNInterface(dcName),
	}

	return cn, nil
}

// CS provides the container storage options in key-value pairs.
// These persistent storage specific properties are supposed to be specific to
// the target Nomad cluster. These pairs are provided based on datacenter.
func (m *nomadUtil) CS(dcName string) (map[v1.VolumeProvisionerProfileLabel]string, error) {

	err := m.validateConf(dcName)
	if err != nil {
		return nil, err
	}

	persistLoc := m.getCSPersistenceLocation(dcName)
	repCount, err := m.getCSReplicaCount(dcName)
	if err != nil {
		return nil, err
	}

	// build the cs map
	//cs := map[v1.ContainerStorageLbl]string{
	cs := map[v1.VolumeProvisionerProfileLabel]string{
		// container persistent storage properties
		//v1.CSPersistenceLocationLbl: persistLoc,
		v1.PVPPersistenceLocationLbl: persistLoc,
		//v1.CSReplicaCountLbl:        repCount,
		v1.PVPReplicaCountLbl: repCount,
	}

	return cs, nil
}

func (m *nomadUtil) validateConf(dcName string) error {

	if dcName == "" {
		return fmt.Errorf("DC name is empty")
	}

	if m.nomadConf == nil {
		return fmt.Errorf("Nil nomad config provided")
	}

	if m.nomadConf.Datacenter == nil {
		return fmt.Errorf("DC not available in nomad config")
	}

	if m.nomadConf.Datacenter[dcName] == nil {
		return fmt.Errorf("No details available for dc '%s'", dcName)
	}

	return nil
}

// getCNType extracts the network type from conf or returns the default value
func (m *nomadUtil) getCNType(dcName string) string {

	if m.nomadConf.Datacenter[dcName] != nil && m.nomadConf.Datacenter[dcName].CNType == "" {
		return v1nomad.DefaultNomadCNType
	}

	return m.nomadConf.Datacenter[dcName].CNType
}

// getCNNetworkCIDR extracts the network CIDR from conf or returns the default value
func (m *nomadUtil) getCNNetworkCIDR(dcName string) string {

	if m.nomadConf.Datacenter[dcName] != nil && m.nomadConf.Datacenter[dcName].CNNetworkCIDR == "" {
		return v1nomad.DefaultNomadCNNetworkCIDR
	}

	return m.nomadConf.Datacenter[dcName].CNNetworkCIDR
}

// getCNInterface extracts the interface from conf or returns the default value
func (m *nomadUtil) getCNInterface(dcName string) string {

	if m.nomadConf.Datacenter[dcName] != nil && m.nomadConf.Datacenter[dcName].CNInterface == "" {
		return v1nomad.DefaultNomadCNInterface
	}

	return m.nomadConf.Datacenter[dcName].CNInterface
}

// getCSPersistenceLocation extracts the backing persistence storage
// location from conf or returns the default value
func (m *nomadUtil) getCSPersistenceLocation(dcName string) string {

	if m.nomadConf.Datacenter[dcName] != nil && m.nomadConf.Datacenter[dcName].CSPersistenceLocation == "" {
		return v1nomad.DefaultNomadCSPersistenceLocation
	}

	return m.nomadConf.Datacenter[dcName].CSPersistenceLocation
}

// getCSReplicaCount returns the default no. of replicas
// as registered in the conf file (i.e. .INI file)
func (m *nomadUtil) getCSReplicaCount(dcName string) (string, error) {

	if m.nomadConf.Datacenter[dcName] != nil && m.nomadConf.Datacenter[dcName].CSReplicaCount == "" {
		return v1nomad.DefaultNomadCSReplicaCount, nil
	}

	repCount, err := strconv.Atoi(m.nomadConf.Datacenter[dcName].CSReplicaCount)
	if err != nil {
		return "", fmt.Errorf("Invalid replica count '%s' provided. '%v'", m.nomadConf.Datacenter[dcName].CSReplicaCount, err)
	}

	if repCount == 0 {
		return "", fmt.Errorf("Replica count can not be '0'")
	}

	return m.nomadConf.Datacenter[dcName].CSReplicaCount, nil
}

// readNomadConfig reads an instance of NomadConfig from config reader.
func readNomadConfig(config io.Reader) (*NomadConfig, error) {
	var nCfg NomadConfig
	var err error

	if config != nil {
		err = gcfg.ReadInto(&nCfg, config)
		if err != nil {
			return nil, err
		}
	}

	// TODO
	// validations w.r.t config

	return &nCfg, nil
}

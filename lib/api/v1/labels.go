package v1

// TODO
// Deprecate
// A type that acts on an infrastructure level i.e. orchestrator hosts
type ContainerNetworkingLbl string

// TODO
// Deprecate
const (
	CNTypeLbl            ContainerNetworkingLbl = "cn.openebs.io/type"
	CNNetworkCIDRAddrLbl ContainerNetworkingLbl = "cn.openebs.io/network-cidr-addr"
	CNSubnetLbl          ContainerNetworkingLbl = "cn.openebs.io/subnet"
	CNInterfaceLbl       ContainerNetworkingLbl = "cn.openebs.io/interface"
)

// TODO
// Deprecate
// A type that acts on an infrastructure level i.e. orchestrator hosts
type ContainerStorageLbl string

// TODO
// Deprecate
const (
	CSPersistenceLocationLbl ContainerStorageLbl = "cs.openebs.io/persistence-location"
	CSReplicaCountLbl        ContainerStorageLbl = "cs.openebs.io/replica-count"
)

// TODO
// Deprecate
const (
	VolumePluginNamePrefix string = "name.plugin.volume.openebs.io/"
)

// TODO
// Deprecate
const (
	DefaultOrchestratorConfigPath string = "/etc/mayaserver/orchprovider/"
)

const (
	JivaNomadPlacementSpecs string = "placement.specs.openebs.io/jnp-specs"
	JivaK8sPlacementSpecs   string = "placement.specs.openebs.io/jk8sp-specs"
)

// TODO
//    Need to standardize on the label's value.
type RequestsLbl string

const (
	// Old way to labels
	// TODO
	// Deprecate
	RegionLbl     RequestsLbl = "requests.openebs.io/region"
	DatacenterLbl RequestsLbl = "requests.openebs.io/dc"
)

// OrchProviderProfileLabel is a typed label to determine orchestration provider
// profile's values.
type OrchProviderProfileLabel string

const (
	// Label / Tag for an orchestrator profile name
	OrchProfileNameLbl OrchProviderProfileLabel = "orchprovider.mapi.openebs.io/profile-name"
	// Label / Tag for an orchestrator region
	OrchRegionLbl OrchProviderProfileLabel = "orchprovider.mapi.openebs.io/region"
	// Label / Tag for an orchestrator datacenter
	OrchDCLbl OrchProviderProfileLabel = "orchprovider.mapi.openebs.io/dc"
	// Label / Tag for an orchestrator namespace
	OrchNSLbl OrchProviderProfileLabel = "orchprovider.mapi.openebs.io/ns"
	// Label / Tag for an orchestrator network address in CIDR format
	OPNetworkAddrLbl OrchProviderProfileLabel = "orchprovider.mapi.openebs.io/network-cidr"
	// Label / Tag for an orchestrator's in-cluster flag
	OPInClusterLbl OrchProviderProfileLabel = "orchprovider.mapi.openebs.io/in-cluster"
)

// OrchProviderDefaults is a typed label to provide default values w.r.t
// orchestration provider properties.
type OrchProviderDefaults string

const (
	// Default value for orchestrator's network address
	// NOTE: Should be in valid CIDR notation
	OPNetworkAddrDef OrchProviderDefaults = "172.28.128.1/24"
	// Default value for orchestrator's network subnet
	OPNetworkSubnetDef OrchProviderDefaults = "24"
	// Default value for orchestrator's in-cluster flag
	OPInClusterDef OrchProviderDefaults = "true"
	// Default value for orchestrator namespace
	OrchNSDefLbl OrchProviderDefaults = "default"
)

// VolumeProvisionerProfileLabel is a typed label to determine volume provisioner
// profile values.
type VolumeProvisionerProfileLabel string

const (
	// Label / Tag for a persistent volume provisioner profile's name
	PVPProfileNameLbl VolumeProvisionerProfileLabel = "volumeprovisioner.mapi.openebs.io/profile-name"
	// Label / Tag for a persistent volume provisioner's VSM name
	PVPVSMNameLbl VolumeProvisionerProfileLabel = "volumeprovisioner.mapi.openebs.io/vsm-name"
	// Label / Tag for a persistent volume provisioner's persistence location
	PVPPersistenceLocationLbl VolumeProvisionerProfileLabel = "volumeprovisioner.mapi.openebs.io/persistence-location"
	// Label / Tag for a persistent volume provisioner's replica support
	PVPReqReplicaLbl VolumeProvisionerProfileLabel = "volumeprovisioner.mapi.openebs.io/req-replica"
	// Label / Tag for a persistent volume provisioner's networking support
	PVPReqNetworkingLbl VolumeProvisionerProfileLabel = "volumeprovisioner.mapi.openebs.io/req-networking"
	// Label / Tag for a persistent volume provisioner's replica count
	PVPReplicaCountLbl VolumeProvisionerProfileLabel = "volumeprovisioner.mapi.openebs.io/replica-count"
	// Label / Tag for a persistent volume provisioner's replica size
	PVPReplicaSizeLbl VolumeProvisionerProfileLabel = "volumeprovisioner.mapi.openebs.io/replica-size"
	// Label / Tag for a persistent volume provisioner's replica IPs
	PVPReplicaIPsLbl VolumeProvisionerProfileLabel = "volumeprovisioner.mapi.openebs.io/replica-ips"
	// Label / Tag for a persistent volume provisioner's replica image
	PVPReplicaImageLbl VolumeProvisionerProfileLabel = "volumeprovisioner.mapi.openebs.io/replica-image"
	// Label / Tag for a persistent volume provisioner's controller count
	PVPControllerCountLbl VolumeProvisionerProfileLabel = "volumeprovisioner.mapi.openebs.io/controller-count"
	// Label / Tag for a persistent volume provisioner's controller image
	PVPControllerImageLbl VolumeProvisionerProfileLabel = "volumeprovisioner.mapi.openebs.io/controller-image"
	// Label / Tag for a persistent volume provisioner's controller size
	PVPControllerSizeLbl VolumeProvisionerProfileLabel = "volumeprovisioner.mapi.openebs.io/controller-size"
	// Label / Tag for a persistent volume provisioner's controller IPs
	PVPControllerIPsLbl VolumeProvisionerProfileLabel = "volumeprovisioner.mapi.openebs.io/controller-ips"
)

// NameLabel type will be used to identify various maya api service components
// via this typed label
type NameLabel string

const (
	// Label / Tag for an orchestrator name
	OrchestratorNameLbl NameLabel = "orchprovider.mapi.openebs.io/name"
	// Label / Tag for a persistent volume provisioner name
	VolumeProvisionerNameLbl NameLabel = "volumeprovisioner.mapi.openebs.io/name"
)

// OrchestratorRegistry type will be used to register various maya api service
// orchestrators.
type OrchProviderRegistry string

const (
	// K8sOrchestrator states Kubernetes as orchestration provider plugin.
	// This is used for registering Kubernetes as an orchestration provider in maya
	// api server.
	K8sOrchestrator OrchProviderRegistry = "kubernetes"
	// NomadOrchestrator states Nomad as orchestration provider plugin.
	// This is used for registering Nomad as an orchestration provider in maya api
	// server.
	NomadOrchestrator OrchProviderRegistry = "nomad"
)

// VolumeProvisionerRegistry type will be used to register various maya api
// service volume provisioners.
type VolumeProvisionerRegistry string

const (
	// JivaVolumeProvisioner states Jiva as persistent volume provisioner plugin.
	// This is used for registering Jiva as a volume provisioner in maya api server.
	JivaVolumeProvisioner VolumeProvisionerRegistry = "jiva"
)

// OrchProviderProfileRegistry type will be used to register various maya api
// service orchestrator profiles
type OrchProviderProfileRegistry string

const (
	// This is the name of PVC as orchestration provider profile
	// This is used for labelling PVC as a orchestration provider profile
	PVCOrchestratorProfile OrchProviderProfileRegistry = "pvc"
)

// VolumeProvisionerProfileRegistry type will be used to register various maya api service
// persistent volume provisioner profiles
type VolumeProvisionerProfileRegistry string

const (
	// This is the name of PVC as persistent volume provisioner profile
	// This is used for labelling PVC as a persistent volume provisioner profile
	PVCProvisionerProfile VolumeProvisionerProfileRegistry = "pvc"
)

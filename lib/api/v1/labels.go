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
	//PVPVSMNameLbl VolumeProvisionerProfileLabel = "volumeprovisioner.mapi.openebs.io/vsm-name"
	// Label / Tag for a persistent volume provisioner's persistence location
	PVPPersistenceLocationLbl VolumeProvisionerProfileLabel = "volumeprovisioner.mapi.openebs.io/persistence-location"
	// Label / Tag for a persistent volume provisioner's replica support
	PVPReqReplicaLbl VolumeProvisionerProfileLabel = "volumeprovisioner.mapi.openebs.io/req-replica"
	// Label / Tag for a persistent volume provisioner's networking support
	PVPReqNetworkingLbl VolumeProvisionerProfileLabel = "volumeprovisioner.mapi.openebs.io/req-networking"
	// Label / Tag for a persistent volume provisioner's replica count
	PVPReplicaCountLbl VolumeProvisionerProfileLabel = "volumeprovisioner.mapi.openebs.io/replica-count"
	// Label / Tag for a persistent volume provisioner's persistent path count
	PVPPersistentPathCountLbl VolumeProvisionerProfileLabel = PVPReplicaCountLbl
	// Label / Tag for a persistent volume provisioner's storage size
	PVPStorageSizeLbl VolumeProvisionerProfileLabel = "volumeprovisioner.mapi.openebs.io/storage-size"
	// Label / Tag for a persistent volume provisioner's replica IPs
	PVPReplicaIPsLbl VolumeProvisionerProfileLabel = "volumeprovisioner.mapi.openebs.io/replica-ips"
	// Label / Tag for a persistent volume provisioner's replica image
	PVPReplicaImageLbl VolumeProvisionerProfileLabel = "volumeprovisioner.mapi.openebs.io/replica-image"
	// Label / Tag for a persistent volume provisioner's controller count
	PVPControllerCountLbl VolumeProvisionerProfileLabel = "volumeprovisioner.mapi.openebs.io/controller-count"
	// Label / Tag for a persistent volume provisioner's controller image
	PVPControllerImageLbl VolumeProvisionerProfileLabel = "volumeprovisioner.mapi.openebs.io/controller-image"
	// Label / Tag for a persistent volume provisioner's controller IPs
	PVPControllerIPsLbl VolumeProvisionerProfileLabel = "volumeprovisioner.mapi.openebs.io/controller-ips"
	// Label / Tag for a persistent volume provisioner's persistent path
	PVPPersistentPathLbl VolumeProvisionerProfileLabel = "volumeprovisioner.mapi.openebs.io/persistent-path"
)

// VolumeProvsionerDefaults is a typed label to provide default values w.r.t
// volume provisioner properties.
type VolumeProvisionerDefaults string

const (
	// Default value for persistent volume provisioner's controller count
	PVPControllerCountDef VolumeProvisionerDefaults = "1"
	// Default value for persistent volume provisioner's replica count
	PVPReplicaCountDef VolumeProvisionerDefaults = "2"
	// Default value for persistent volume provisioner's persistent path count
	// This should be equal to persistent volume provisioner's replica count
	PVPPersistentPathCountDef VolumeProvisionerDefaults = PVPReplicaCountDef
	// Default value for persistent volume provisioner's controller image
	PVPControllerImageDef VolumeProvisionerDefaults = "openebs/jiva:latest"
	// Default value for persistent volume provisioner's support for replica
	PVPReqReplicaDef VolumeProvisionerDefaults = "true"
	// Default value for persistent volume provisioner's replica image
	PVPReplicaImageDef VolumeProvisionerDefaults = "openebs/jiva:latest"
	// Default value for persistent volume provisioner's networking support
	PVPReqNetworkingDef VolumeProvisionerDefaults = "false"
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
	// DefaultOrchestrator provides the default orchestration provider
	DefaultOrchestrator = K8sOrchestrator
)

// VolumeProvisionerRegistry type will be used to register various maya api
// service volume provisioners.
type VolumeProvisionerRegistry string

const (
	// JivaVolumeProvisioner states Jiva as persistent volume provisioner plugin.
	// This is used for registering Jiva as a volume provisioner in maya api server.
	JivaVolumeProvisioner VolumeProvisionerRegistry = "jiva"
	// DefaultVolumeProvisioner provides the default persistent volume provisioner
	// plugin.
	DefaultVolumeProvisioner VolumeProvisionerRegistry = JivaVolumeProvisioner
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

// TODO
// Move these to jiva folder
//
// JivaAnnotations will be used to provide filtering options like
// named-labels, named-suffix, named-prefix, constants, etc.
//
// NOTE:
//    These value(s) are generally used / remembered by the consumers of
// maya api service
type JivaAnnotations string

// TODO
// Rename these const s.t. they start with Jiva as Key Word
const (
	// VSMIdentifier is used to filter vsm by name
	VSMIdentifier JivaAnnotations = "vsm"

	// VSMSelectorPrefix is used to filter vsm by name when
	// selector logic is used
	VSMSelectorPrefix JivaAnnotations = VSMIdentifier + "="

	// ControllerSuffix is used as a suffix for persistent volume controller
	// related names
	ControllerSuffix JivaAnnotations = "-ctrl"

	// JivaReplicaSuffix is used as a suffix for persistent volume replica
	// related names
	JivaReplicaSuffix JivaAnnotations = "-rep"

	// ServiceSuffix is used as a suffix for persistent volume controller
	// related names
	ServiceSuffix JivaAnnotations = "-svc"

	// ControllerSuffix is used as a suffix for persistent volume container
	// related names
	ContainerSuffix JivaAnnotations = "-con"

	// PortNameISCSI is the name given to iscsi ports
	PortNameISCSI JivaAnnotations = "iscsi"

	// PortNameAPI is the name given to api ports
	PortNameAPI JivaAnnotations = "api"

	// JivaCtrlIPHolder is used as a placeholder for persistent volume controller's
	// IP address
	//
	// NOTE:
	//    This is replaced at runtime
	JivaCtrlIPHolder JivaAnnotations = "__CTRL_IP__"

	// JivaStorageSizeHolder is used as a placeholder for persistent volume's
	// storage capacity
	//
	// NOTE:
	//    This is replaced at runtime
	JivaStorageSizeHolder JivaAnnotations = "__STOR_SIZE__"
)

// JivaDefaults is a typed label to provide default values to Jiva based
// persistent volume properties
type JivaDefaults string

const (
	// JivaControllerFrontendDef is used to provide default frontend for jiva
	// persistent volume controller
	JivaControllerFrontendDef JivaDefaults = "gotgt"

	// JivaVolumeNameDef is used to provide default name for jiva
	// based persistent volumes
	JivaVolumeNameDef JivaDefaults = "jvol"

	// JivaISCSIPortDef is used to provide default iscsi port value for jiva
	// based persistent volumes
	JivaISCSIPortDef JivaDefaults = "3260"

	// JivaPersistentMountPathDef is the default mount path used by jiva based
	// persistent volumes
	JivaPersistentMountPathDef JivaDefaults = "/openebs"

	// JivaPersistentMountNameDef is the default mount path name used by jiva based
	// persistent volumes
	JivaPersistentMountNameDef JivaDefaults = "openebs"

	// JivaAPIPortDef is used to provide management port for persistent volume
	// storage
	JivaAPIPortDef JivaDefaults = "9501"

	// JivaReplicaPortOneDef is used to provide port for jiva based persistent
	// volume replica
	JivaReplicaPortOneDef JivaDefaults = "9502"

	// JivaReplicaPortTwoDef is used to provide port for jiva based persistent
	// volume replica
	JivaReplicaPortTwoDef JivaDefaults = "9503"

	// JivaReplicaPortThreeDef is used to provide port for jiva based persistent
	// volume replica
	JivaReplicaPortThreeDef JivaDefaults = "9504"

	// JivaPersistentPathDef is used to set default value for
	// persistent volume's persist path
	JivaPersistentPathDef JivaDefaults = "/tmp"

	// JivaStorSizeDef is used to set default value for
	// persistent volume's persist path
	JivaStorSizeDef JivaDefaults = "1G"
)

// These will be used to provide array based constants that are
// related to jiva volume provisioner
var (
	// JivaCtrlCmd is the command used to start jiva controller
	JivaCtrlCmd = []string{"launch"}

	// JivaCtrlArgs is the set of arguments provided to JivaCtrlCmd
	JivaCtrlArgs = []string{"controller", "--frontend", string(JivaControllerFrontendDef), string(JivaVolumeNameDef)}

	// JivaReplicaCmd is the command used to start jiva replica
	JivaReplicaCmd = []string{"launch"}

	// JivaReplicaArgs is the set of arguments provided to JivaReplicaCmd
	JivaReplicaArgs = []string{"replica", "--frontendIP", string(JivaCtrlIPHolder), "--size", string(JivaStorageSizeHolder), string(JivaPersistentMountPathDef)}
)

// TODO
// Move these to k8s folder
//
// K8sAnnotations will be used to provide string based constants that are
// related to kubernetes as orchestration provider
type K8sAnnotations string

const (
	// K8sKindDeployment is used to state the k8s Deployment(s)
	K8sKindDeployment K8sAnnotations = "Deployment"
	// K8sKindService is used to state the k8s Service(s)
	K8sKindService K8sAnnotations = "Service"
	// K8sServiceVersion is used to state the k8s Service version
	K8sServiceVersion K8sAnnotations = "v1"
	// K8sPodVersion is used to state the k8s Pod version
	K8sPodVersion K8sAnnotations = "v1"
)

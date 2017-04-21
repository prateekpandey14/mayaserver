package v1

// A type that acts on an infrastructure level i.e. orchestrator hosts
type ContainerNetworkingLbl string

const (
	CNTypeLbl            ContainerNetworkingLbl = "cn.openebs.io/type"
	CNNetworkCIDRAddrLbl ContainerNetworkingLbl = "cn.openebs.io/network-cidr-addr"
	CNSubnetLbl          ContainerNetworkingLbl = "cn.openebs.io/subnet"
	CNInterfaceLbl       ContainerNetworkingLbl = "cn.openebs.io/interface"
)

// A type that acts on an infrastructure level i.e. orchestrator hosts
type ContainerStorageLbl string

const (
	CSPersistenceLocationLbl ContainerStorageLbl = "cs.openebs.io/persistence-location"
	CSReplicaCountLbl        ContainerStorageLbl = "cs.openebs.io/replica-count"
)

// TODO
//    Need to standardize on the label's value.
type RequestsLbl string

const (
	// Old way to labels
	RegionLbl     RequestsLbl = "requests.openebs.io/region"
	DatacenterLbl RequestsLbl = "requests.openebs.io/dc"

	// Current Label Standards
	// NOTE: Need not go with the versioning as we are at 0.2+ releases.
	//    Try to provide two Lbl variable names.
	OrchProviderNameLbl      RequestsLbl = "orchprovider.mapiserver.openebs.io/name"
	VolumeProvisionerNameLbl RequestsLbl = "volumeprovisioner.mapiserver.openebs.io/name"

	PVPersistenceLocationLbl ContainerStorageLbl = "persistentvolume.mapiserver.openebs.io/persistence-location"
	PVReplicaCountLbl        ContainerStorageLbl = "persistentvolume.mapiserver.openebs.io/replica-count"
)

const (
	VolumePluginNamePrefix string = "name.plugin.volume.openebs.io/"
)

const (
	DefaultOrchestratorConfigPath string = "/etc/mayaserver/orchprovider/"
)

const (
	JivaNomadPlacementSpecs string = "placement.specs.openebs.io/jnp-specs"
	JivaK8sPlacementSpecs   string = "placement.specs.openebs.io/jk8sp-specs"
)

package volumeprovisioner

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/openebs/mayaserver/lib/api/v1"
	"github.com/openebs/mayaserver/lib/util"
)

// VolumeProvisionerProfile abstracts & exposes a persistent volume provisioner's
// runtime features.
//
// NOTE:
//    A persistent volume provisioner can align to a specific implementation of
// this profile & hence change its execution strategy at runtime.
type VolumeProvisionerProfile interface {
	// Label assigned against the persistent volume provisioner profile.
	Label() v1.VolumeProvisionerProfileLabel

	// Registered volume provisioner profile name.
	Name() v1.VolumeProvisionerProfileRegistry

	// Get the persistent volume claim associated with this provisioner profile
	PVC() (*v1.PersistentVolumeClaim, error)

	// Gets the orchestration provider name.
	// A persistent volume provisioner plugin may be linked with a orchestrator
	// e.g. K8s, Nomad, Mesos, Swarm, etc. It can be Docker engine as well.
	//
	// Note:
	//    It can return false in its second return argument if orchestrator support
	// is not applicable.
	//
	// Note:
	//    OpenEBS believes in running storage software in containers & hence
	// these container specific orchestrators.
	Orchestrator() (v1.OrchProviderRegistry, bool, error)

	// Get the name of the VSM
	VSMName() (string, error)

	// Get the number of controllers
	ControllerCount() (int, error)

	// Gets the controller's image e.g. docker image version. The second return value
	// indicates if image based replica is supported or not.
	ControllerImage() (string, bool, error)

	// Get the IP addresses that needs to be assigned against the controller(s)
	ControllerIPs() ([]string, error)

	// If replica(s) are supported by the persistent volume provisioner
	ReqReplica() bool

	// Gets the replica's image e.g. docker image version. The second return value
	// indicates if image based replica is supported or not.
	ReplicaImage() (string, bool, error)

	// Get the storage size for each replica(s)
	StorageSize() (string, error)

	// Get the number of replicas
	ReplicaCount() (int, error)

	// Get the IP addresses that needs to be assigned against the replica(s)
	ReplicaIPs() ([]string, error)
}

// GetVolProProfileByPVC will return a specific persistent volume provisioner
// profile. It will decide first based on the provided specifications failing
// which will ensure a default profile is returned.
func GetVolProProfileByPVC(pvc *v1.PersistentVolumeClaim) (VolumeProvisionerProfile, error) {
	if pvc == nil {
		return GetDefaultVolProProfile()
	}

	// Extract the name of volume provisioner profile
	volProflName := v1.VolumeProvisionerProfileName(pvc.Labels)

	if volProflName == "" {
		return GetDefaultVolProProfile()
	}

	return GetVolProProfileByName(volProflName)
}

// GetDefaultVolProProfile will return the default volume provisioner
// profile.
//
// NOTE:
//    PVC based volume provisioner profile is considered as default
func GetDefaultVolProProfile() (VolumeProvisionerProfile, error) {
	return &pvcVolProProfile{}, nil
}

// TODO
//
// GetVolProProfileByName will return a volume provisioner profile by
// looking up from the provided profile name.
func GetVolProProfileByName(name string) (VolumeProvisionerProfile, error) {
	// TODO
	// Search from the in-memory registry

	// TODO
	// Alternatively, search from external discoverable DBs if any

	return nil, fmt.Errorf("GetVolProProfileByName is not yet implemented")
}

// pvcVolProProfile is a persistent volume provisioner profile that is based on
// persistent volume claim.
//
// NOTE:
//    This will use defaults in-case the values are not set in persistent volume
// claim.
//
// NOTE:
//    This is a concrete implementation of
// volumeprovisioner.VolumeProvisionerProfile
type pvcVolProProfile struct {
	pvc *v1.PersistentVolumeClaim
}

// newPvcVolProProfile provides a new instance of VolumeProvisionerProfile that is
// based on pvc (i.e. persistent volume claim).
func newPvcVolProProfile(pvc *v1.PersistentVolumeClaim) (VolumeProvisionerProfile, error) {
	// This does not care if pvc instance is nil
	return &pvcVolProProfile{
		pvc: pvc,
	}, nil
}

// Label provides the label assigned against the persistent volume provisioner
// profile.
//
// NOTE:
//    There can be many persistent volume provisioner profiles with this same label.
// This is used along with Name() method.
func (pp *pvcVolProProfile) Label() v1.VolumeProvisionerProfileLabel {
	return v1.PVPProfileNameLbl
}

// Name provides the name assigned to the persistent volume provisioner profile.
//
// NOTE:
//    Name provides the uniqueness among various variants of persistent volume
// provisioner profiles.
func (pp *pvcVolProProfile) Name() v1.VolumeProvisionerProfileRegistry {
	return v1.PVCProvisionerProfile
}

// PVC provides the persistent volume claim associated with this profile.
//
// NOTE:
//    This method provides a convinient way to access pvc. In other words
// volume provisioner profile acts as a wrapper over pvc.
func (pp *pvcVolProProfile) PVC() (*v1.PersistentVolumeClaim, error) {
	return pp.pvc, nil
}

// Orchestrator gets the suitable orchestration provider.
// A persistent volume provisioner plugin may be linked with a orchestrator
// e.g. K8s, Nomad, Mesos, Swarm, etc. It can be Docker engine as well.
func (pp *pvcVolProProfile) Orchestrator() (v1.OrchProviderRegistry, bool, error) {
	if pp.pvc == nil {
		return v1.DefaultOrchestratorName(), true, nil
	}

	// Extract the name of orchestration provider
	oName := v1.OrchestratorName(pp.pvc.Labels)

	if string(oName) == "" {
		return v1.DefaultOrchestratorName(), true, nil
	}

	// Get the orchestrator instance
	return oName, true, nil
}

// VSMName gets the name of the VSM
// Operator must provide this.
func (pp *pvcVolProProfile) VSMName() (string, error) {
	if pp.pvc == nil {
		return "", fmt.Errorf("Missing VSM name in '%s:%s'", pp.Label(), pp.Name())
	}

	// Extract the VSM name from pvc
	vsmName := v1.VSMName(pp.pvc.Labels)

	if vsmName == "" {
		return "", fmt.Errorf("Missing VSM name in '%s:%s'", pp.Label(), pp.Name())
	}

	return vsmName, nil
}

// ControllerCount gets the number of controllers
func (pp *pvcVolProProfile) ControllerCount() (int, error) {
	if pp.pvc == nil {
		return v1.DefaultControllerCount(), nil
	}

	// Extract the controller count from pvc
	cCount := v1.ControllerCount(pp.pvc.Labels)

	if cCount == "" {
		return v1.DefaultControllerCount(), nil
	}

	iCCount, err := strconv.Atoi(cCount)
	if err != nil {
		return 0, err
	}

	return iCCount, nil
}

// ControllerImage gets the controller's image currently its docker image label.
func (pp *pvcVolProProfile) ControllerImage() (string, bool, error) {
	if pp.pvc == nil {
		return v1.DefaultControllerImage(), true, nil
	}

	// Extract the controller image from pvc
	cImg := v1.ControllerImage(pp.pvc.Labels)

	if cImg == "" {
		return v1.DefaultControllerImage(), true, nil
	}

	return cImg, true, nil
}

// ReqReplica indicates if replica(s) are required by the persistent volume
// provisioner
func (pp *pvcVolProProfile) ReqReplica() bool {
	if pp.pvc == nil {
		return v1.DefaultReqReplica()
	}

	// Extract the replica truthiness (i.e. is replica required) from pvc
	reqReplica := v1.ReqReplica(pp.pvc.Labels)

	if reqReplica == "" {
		return v1.DefaultReqReplica()
	}

	return util.CheckTruthy(reqReplica)
}

// ReplicaImage gets the replica's image currently its docker image label.
func (pp *pvcVolProProfile) ReplicaImage() (string, bool, error) {
	if pp.pvc == nil {
		return v1.DefaultReplicaImage(), true, nil
	}

	// Extract the replica image from pvc
	rImg := v1.ReplicaImage(pp.pvc.Labels)

	if rImg == "" {
		return v1.DefaultReplicaImage(), true, nil
	}

	return rImg, true, nil
}

// StorageSize gets the storage size for each persistent volume replica(s)
func (pp *pvcVolProProfile) StorageSize() (string, error) {
	// Extract the storage size from pvc
	sSize := v1.StorageSize(pp.pvc.Labels)

	if sSize == "" {
		return "", fmt.Errorf("Missing storage size in '%s:%s'", pp.Label(), pp.Name())
	}

	return sSize, nil
}

// ReplicaCount get the number of replicas required
func (pp *pvcVolProProfile) ReplicaCount() (int, error) {
	if pp.pvc == nil {
		return v1.DefaultReplicaCount(), nil
	}

	// Extract the replica count from pvc
	rCount := v1.ReplicaCount(pp.pvc.Labels)

	if rCount == "" {
		return v1.DefaultReplicaCount(), nil
	}

	iRCount, err := strconv.Atoi(rCount)
	if err != nil {
		return 0, err
	}

	return iRCount, nil
}

// ReqNetworking indicates if any networking related operations are required to
// be executed by persistent volume provisioner
func (pp *pvcVolProProfile) ReqNetworking() bool {
	if pp.pvc == nil {
		return v1.DefaultReqNetworking()
	}

	// Extract the networking truthiness (i.e. is networking operations required)
	// from pvc
	reqNet := v1.ReqNetworking(pp.pvc.Labels)
	// Default is true
	if reqNet == "" {
		return v1.DefaultReqNetworking()
	}

	return util.CheckTruthy(reqNet)
}

// ControllerIPs gets the IP addresses that needs to be assigned against the
// controller(s)
//
// NOTE:
//    There is no default assignment of IPs
func (pp *pvcVolProProfile) ControllerIPs() ([]string, error) {
	if pp.pvc == nil {
		return nil, nil
	}

	// Extract the controller IPs from pvc
	cIPs := v1.ControllerIPs(pp.pvc.Labels)

	if cIPs == "" {
		return nil, nil
	}

	cIPsArr := strings.Split(cIPs, ",")

	if len(cIPsArr) == 0 {
		return nil, fmt.Errorf("Invalid controller IPs in '%s:%s'", pp.Label(), pp.Name())
	}

	return cIPsArr, nil
}

// ReplicaIPs gets the IP addresses that needs to be assigned against the
// replica(s)
//
// NOTE:
//    There is no default assignment of IPs
func (pp *pvcVolProProfile) ReplicaIPs() ([]string, error) {
	if pp.pvc == nil {
		return nil, nil
	}

	// Extract the controller IPs from pvc
	rIPs := v1.ReplicaIPs(pp.pvc.Labels)

	if rIPs == "" {
		return nil, nil
	}

	rIPsArr := strings.Split(rIPs, ",")

	if len(rIPsArr) == 0 {
		return nil, fmt.Errorf("Invalid replica IPs in '%s:%s'", pp.Label(), pp.Name())
	}

	return rIPsArr, nil
}

// srVolProProfile represents a single replica based persistent volume
// provisioner profile. It relies on persistent volume claim to source other
// volume provisioner properties.
//
// NOTE:
//    Replica count property specified in persistent volume claim will override
// the one specified in etcd
//
// NOTE:
//    This is a concrete implementation of volume.VolumeProvisionerProfile
type srVolProProfile struct {
	pvcVolProProfile
}

// ReqReplica indicates if replica(s) are required by the persistent volume
// provisioner
func (srp *srVolProProfile) ReqReplica() bool {
	return util.CheckTruthy("true")
}

// ReplicaCount gets the number of replicas required. In most of the cases, it
// returns 1. However, if replica count is specified in pvc then the specs value
// is considered.
func (srp *srVolProProfile) ReplicaCount() (int, error) {
	if srp.pvc == nil {
		return 1, nil
	}

	// Extract the replica count from pvc
	rCount := v1.ReplicaCount(srp.pvc.Labels)

	if rCount == "" {
		return 1, nil
	}

	iRCount, err := strconv.Atoi(rCount)
	if err != nil {
		return 0, err
	}

	return iRCount, nil
}

// etcdVolProProfile represents a generic volume provisioner profile whose
// properties are stored in etcd database.
//
// NOTE:
//    There can be multiple persistent volume provisioner profiles stored in
// etcd
//
// NOTE:
//    Properties specified in persistent volume claim will override the ones
// specified in etcd
//
// NOTE:
//    This is a concrete implementation of volume.VolumeProvisionerProfile
type etcdVolProProfile struct {
	pvc *v1.PersistentVolumeClaim
}

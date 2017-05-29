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

	// Get the count of persistent paths required for all the replicas
	PersistentPathCount() (int, error)

	// Get the persistent path based on the replica position.
	//
	// NOTE:
	//    `position` is just a parameter that determines a particular replica out
	// of the total replica count i.e. `rCount`.
	PersistentPath(position int, rCount int) (string, error)
}

// GetVolProProfileByPVC will return a specific persistent volume provisioner
// profile. It will decide first based on the provided specifications failing
// which will ensure a default profile is returned.
func GetVolProProfileByPVC(pvc *v1.PersistentVolumeClaim) (VolumeProvisionerProfile, error) {
	//if pvc == nil || pvc.Labels == nil {
	if pvc == nil {
		return nil, fmt.Errorf("PVC is required to create a volume provisioner profile")
	}

	// Extract the name of volume provisioner profile
	volProflName := v1.VolumeProvisionerProfileName(pvc.Labels)

	if volProflName == "" {
		return GetDefaultVolProProfile(pvc)
	}

	return GetVolProProfileByName(volProflName, pvc)
}

// GetDefaultVolProProfile will return the default volume provisioner
// profile.
//
// NOTE:
//    PVC based volume provisioner profile is considered as default
func GetDefaultVolProProfile(pvc *v1.PersistentVolumeClaim) (VolumeProvisionerProfile, error) {

	//if pvc == nil || pvc.Labels == nil {
	if pvc == nil {
		return nil, fmt.Errorf("PVC is required to create default volume provisioner profile")
	}

	return newPvcVolProProfile(pvc)
}

// TODO
//
// GetVolProProfileByName will return a volume provisioner profile by
// looking up from the provided profile name.
func GetVolProProfileByName(name string, pvc *v1.PersistentVolumeClaim) (VolumeProvisionerProfile, error) {
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
	// Extract the VSM name from PVC
	// Name of PVC is the name of VSM
	vsmName := v1.VSMName(pp.pvc.Name)

	if vsmName == "" {
		return "", fmt.Errorf("Missing VSM name in '%s:%s'", pp.Label(), pp.Name())
	}

	return vsmName, nil
}

// ControllerCount gets the number of controllers
func (pp *pvcVolProProfile) ControllerCount() (int, error) {
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
	// Extract the replica truthiness (i.e. is replica required) from pvc
	reqReplica := v1.ReqReplica(pp.pvc.Labels)

	if reqReplica == "" {
		return v1.DefaultReqReplica()
	}

	return util.CheckTruthy(reqReplica)
}

// ReplicaImage gets the replica's image currently its docker image label.
func (pp *pvcVolProProfile) ReplicaImage() (string, bool, error) {
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

// PersistentPathCount gets the count of persistent paths required for all the
// replicas.
//
// NOTE:
//    The count needs to be equal to no of replicas.
func (pp *pvcVolProProfile) PersistentPathCount() (int, error) {
	// Extract the persistent path count from pvc
	pCount := v1.PersistentPathCount(pp.pvc.Labels)

	if pCount == "" {
		return v1.DefaultPersistentPathCount(), nil
	}

	iPCount, err := strconv.Atoi(pCount)
	if err != nil {
		return 0, err
	}

	return iPCount, nil
}

// PersistentPath gets the persistent path based on the replica position.
//
// NOTE:
//    `position` is just a positional value that determines a particular replica
// out of the total replica count i.e. rCount.
func (pp *pvcVolProProfile) PersistentPath(position int, rCount int) (string, error) {
	if rCount <= 0 {
		return "", fmt.Errorf("Invalid replica count '%d' provided", rCount)
	}

	if position <= 0 {
		return "", fmt.Errorf("Invalid persistent path index '%d' provided", position)
	}

	vsm, err := pp.VSMName()
	if err != nil {
		return "", err
	}

	// Extract the persistent path from pvc
	pPath := v1.JivaPersistentPath(pp.pvc.Labels, vsm, position)

	if pPath == "" {
		return v1.DefaultJivaPersistentPath(vsm, position), nil
	}

	//pPathArr := strings.Split(pPath, ",")

	//if len(pPathArr) != rCount {
	//	return "", fmt.Errorf("VSM '%s' persistent paths '%d' and replicas '%d' mismatch", vsm, len(pPathArr), rCount)
	//}

	//iPPath := strings.TrimSpace(pPathArr[position-1])

	return pPath, nil
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

// ReplicaCount gets the number of replicas required. In this case it returns 1
// always.
func (srp *srVolProProfile) ReplicaCount() (int, error) {
	return 1, nil
}

// PersistentPathCount gets the count of persistent paths required for all the
// replicas. In this case it returns 1 always.
//
// NOTE:
//    The count needs to be equal to no of replicas.
func (srp *srVolProProfile) PersistentPathCount() (int, error) {
	return 1, nil
}

// PersistentPath gets the persistent path based on the position i.e. replica
// position.
//
// NOTE:
//    `position` is just a positional value that determines a particular replica
// out of the total replica count i.e. rCount.
func (srp *srVolProProfile) PersistentPath(position int, rCount int) (string, error) {
	if rCount != 1 {
		return "", fmt.Errorf("Invalid replica count. Expected '1' Provided '%d'", rCount)
	}

	if position != 1 {
		return "", fmt.Errorf("Invalid persistent path index. Expected '0' Provided '%d'", position)
	}

	vsm, err := srp.VSMName()
	if err != nil {
		return "", err
	}

	// Extract the persistent path from pvc
	pPath := v1.JivaPersistentPath(srp.pvc.Labels, vsm, position)

	if pPath == "" {
		return v1.DefaultJivaPersistentPath(vsm, position), nil
	}

	pPathArr := strings.Split(pPath, ",")

	if len(pPathArr) != rCount {
		return "", fmt.Errorf("VSM '%s' persistent paths '%d' and replicas '%d' mismatch", vsm, len(pPathArr), rCount)
	}

	iPPath := strings.TrimSpace(pPathArr[position-1])

	return iPPath, nil
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

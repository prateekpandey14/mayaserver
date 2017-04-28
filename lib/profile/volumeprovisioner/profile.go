package volumeprovisioner

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/openebs/mayaserver/lib/api/v1"
	"github.com/openebs/mayaserver/lib/nethelper"
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
	Orchestrator() (v1.OrchestratorRegistry, bool, error)

	// Get the name of the VSM
	VSMName() (string, error)

	// Get the number of controllers
	ControllerCount() (int, error)

	// Gets the controller's image e.g. docker image version. The second return value
	// indicates if image based replica is supported or not.
	ControllerImage() (string, bool, error)

	// Gets the controller storage size
	ControllerSize() (string, error)

	// Get the IP addresses that needs to be assigned against the controller(s)
	ControllerIPs() ([]string, error)

	// If replica(s) are supported by the persistent volume provisioner
	ReqReplica() bool

	// Gets the replica's image e.g. docker image version. The second return value
	// indicates if image based replica is supported or not.
	ReplicaImage() (string, bool, error)

	// Get the storage size for each replica(s)
	ReplicaSize() (string, error)

	// Get the number of replicas
	ReplicaCount() (int, error)

	// Get the IP addresses that needs to be assigned against the replica(s)
	ReplicaIPs() ([]string, error)

	// Get the network address in CIDR format
	NetworkAddr() (string, error)

	// Get the network subnet
	NetworkSubnet() (string, error)
}

// pvcVolProProfile is the concrete implementation of volume.VolumeProvisionerProfile
type pvcVolProProfile struct {
	pvc *v1.PersistentVolumeClaim
}

// NewPvcVolProProfile provides a new instance of VolumeProvisionerProfile that is
// based on pvc (i.e. persistent volume claim).
func NewPvcVolProProfile(pvc *v1.PersistentVolumeClaim) (VolumeProvisionerProfile, error) {
	if pvc == nil {
		return nil, fmt.Errorf("Nil pvc in pvc based persistent volume provisioner profile")
	}

	if pvc.Labels == nil {
		return nil, fmt.Errorf("Missing labels in pvc based persistent volume provisioner profile")
	}

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
func (pp *pvcVolProProfile) Orchestrator() (v1.OrchestratorRegistry, bool, error) {
	// Extract the orchestrator provider name from pvc
	orchestratorName := pp.pvc.Labels[string(v1.OrchestratorNameLbl)]
	if orchestratorName == "" {
		return "", true, fmt.Errorf("Missing orchestrator name in '%s:%s'", pp.Label(), pp.Name())
	}

	// Get the orchestrator instance
	return v1.OrchestratorRegistry(orchestratorName), true, nil
}

// VSMName gets the name of the VSM
func (pp *pvcVolProProfile) VSMName() (string, error) {
	// Extract the VSM name from pvc
	vsmName := pp.pvc.Labels[string(v1.PVPVSMNameLbl)]

	if vsmName == "" {
		return "", fmt.Errorf("Missing VSM name in '%s:%s'", pp.Label(), pp.Name())
	}

	return vsmName, nil
}

// ControllerCount gets the number of controllers
func (pp *pvcVolProProfile) ControllerCount() (int, error) {
	// Extract the controller count from pvc
	cCount := pp.pvc.Labels[string(v1.PVPControllerCountLbl)]

	iCCount, err := strconv.Atoi(cCount)
	if err != nil {
		return 0, err
	}

	return iCCount, nil
}

// ControllerImage gets the controller's image currently its docker image label.
func (pp *pvcVolProProfile) ControllerImage() (string, bool, error) {
	// Extract the controller image from pvc
	cImg := pp.pvc.Labels[string(v1.PVPControllerImageLbl)]

	if cImg == "" {
		return "", true, fmt.Errorf("Missing controller image in '%s:%s'", pp.Label(), pp.Name())
	}

	return cImg, true, nil
}

// ControllerSize gets the controller storage size
func (pp *pvcVolProProfile) ControllerSize() (string, error) {
	// Extract the controller size from pvc
	cSize := pp.pvc.Labels[string(v1.PVPControllerSizeLbl)]

	if cSize == "" {
		return "", fmt.Errorf("Missing controller size in '%s:%s'", pp.Label(), pp.Name())
	}

	return cSize, nil
}

// ReqReplica indicates if replica(s) are required by the persistent volume
// provisioner
func (pp *pvcVolProProfile) ReqReplica() bool {
	// Extract the replica truthiness (i.e. is replica required) from pvc
	reqReplica := pp.pvc.Labels[string(v1.PVPReqReplicaLbl)]
	// Default is true
	if reqReplica == "" {
		return true
	}

	return util.CheckTruthy(reqReplica)
}

// ReplicaImage gets the replica's image currently its docker image label.
func (pp *pvcVolProProfile) ReplicaImage() (string, bool, error) {
	// Extract the replica image from pvc
	rImg := pp.pvc.Labels[string(v1.PVPReplicaImageLbl)]

	if rImg == "" {
		return "", true, fmt.Errorf("Missing replica image in '%s:%s'", pp.Label(), pp.Name())
	}

	return rImg, true, nil
}

// ReplicaSize gets the storage size for each replica(s)
func (pp *pvcVolProProfile) ReplicaSize() (string, error) {
	// Extract the replica size from pvc
	rSize := pp.pvc.Labels[string(v1.PVPReplicaSizeLbl)]

	if rSize == "" {
		return "", fmt.Errorf("Missing replica size in '%s:%s'", pp.Label(), pp.Name())
	}

	return rSize, nil
}

// ReplicaCount get the number of replicas required
func (pp *pvcVolProProfile) ReplicaCount() (int, error) {
	// Extract the replica count from pvc
	rCount := pp.pvc.Labels[string(v1.PVPReplicaCountLbl)]

	if rCount == "" {
		return 0, fmt.Errorf("Missing replica count in '%s:%s'", pp.Label(), pp.Name())
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
	reqNet := pp.pvc.Labels[string(v1.PVPReqNetworkingLbl)]
	// Default is true
	if reqNet == "" {
		return true
	}

	return util.CheckTruthy(reqNet)
}

// ControllerIPs gets the IP addresses that needs to be assigned against the
// controller(s)
func (pp *pvcVolProProfile) ControllerIPs() ([]string, error) {
	// Extract the controller IPs from pvc
	cIPs := pp.pvc.Labels[string(v1.PVPControllerIPsLbl)]

	if cIPs == "" {
		return nil, fmt.Errorf("Missing controller IPs in '%s:%s'", pp.Label(), pp.Name())
	}

	cIPsArr := strings.Split(cIPs, ",")

	if len(cIPsArr) == 0 {
		return nil, fmt.Errorf("Invalid controller IPs in '%s:%s'", pp.Label(), pp.Name())
	}

	return cIPsArr, nil
}

// ReplicaIPs gets the IP addresses that needs to be assigned against the replica(s)
func (pp *pvcVolProProfile) ReplicaIPs() ([]string, error) {
	// Extract the controller IPs from pvc
	rIPs := pp.pvc.Labels[string(v1.PVPReplicaIPsLbl)]

	if rIPs == "" {
		return nil, fmt.Errorf("Missing replica IPs in '%s:%s'", pp.Label(), pp.Name())
	}

	rIPsArr := strings.Split(rIPs, ",")

	if len(rIPsArr) == 0 {
		return nil, fmt.Errorf("Invalid replica IPs in '%s:%s'", pp.Label(), pp.Name())
	}

	return rIPsArr, nil
}

// NetworkAddr gets the network address in CIDR format
func (pp *pvcVolProProfile) NetworkAddr() (string, error) {
	// Extract the network address from pvc
	nAddr := pp.pvc.Labels[string(v1.PVPNetworkAddrLbl)]

	if nAddr == "" {
		return "", fmt.Errorf("Missing network address in '%s:%s'", pp.Label(), pp.Name())
	}

	if !nethelper.IsCIDR(nAddr) {
		return "", fmt.Errorf("Network address not in CIDR format in '%s:%s'", pp.Label(), pp.Name())
	}

	return nAddr, nil
}

// NetworkSubnet gets the network's subnet in decimal format
func (pp *pvcVolProProfile) NetworkSubnet() (string, error) {
	// Extract the subnet from network address
	nAddr, err := pp.NetworkAddr()
	if err != nil {
		return "", err
	}

	subnet, err := nethelper.CIDRSubnet(nAddr)
	if err != nil {
		return "", err
	}

	return subnet, nil
}

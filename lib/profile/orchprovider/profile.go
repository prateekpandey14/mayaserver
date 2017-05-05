package orchprovider

import (
	"fmt"

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
type OrchProviderProfile interface {
	// Label assigned against the orchestration provider profile.
	Label() v1.OrchProviderProfileLabel

	// Registered orchestration provider profile name.
	Name() v1.OrchProviderProfileRegistry

	// Get the persistent volume claim associated with this orchestration provider
	PVC() (*v1.PersistentVolumeClaim, error)

	// Get the network address in CIDR format
	NetworkAddr() (string, error)

	// Get the network subnet
	NetworkSubnet() (string, error)

	// Get the namespace used at the orchestrator, where the request needs to be
	// operated on
	NS() (string, error)

	// InCluster indicates if the request to the orchestrator is scoped to the
	// cluster where this request originated
	//
	// TODO
	// Should this be termed as InDC ? Is a cluster same as a DataCenter ?
	// Cluster vs. DC vs. Region ?
	InCluster() (bool, error)
}

// GetOrchProviderProfileByPVC will return a specific orchestration provider profile.
// It will decide first based on the provided specifications failing which will
// ensure a default profile is returned.
func GetOrchProviderProfileByPVC(pvc *v1.PersistentVolumeClaim) (OrchProviderProfile, error) {
	if pvc == nil {
		return GetDefaultOrchProviderProfile()
	}

	// Extract the name of orchestrator profile
	orchProflName := v1.OrchProfileName(pvc.Labels)

	if orchProflName == "" {
		return GetDefaultOrchProviderProfile()
	}

	return GetOrchProviderProfileByName(orchProflName)
}

// GetDefaultOrchProviderProfile will return the default orchestration provider
// profile.
//
// NOTE:
//    PVC based orchestration provider profile is considered as default
func GetDefaultOrchProviderProfile() (OrchProviderProfile, error) {
	return &pvcOrchProviderProfile{}, nil
}

// TODO
//
// GetOrchProviderProfileByName will return a orchestration provider profile by
// looking up from the provided profile name.
func GetOrchProviderProfileByName(name string) (OrchProviderProfile, error) {
	// TODO
	// Search from the in-memory registry

	// TODO
	// Alternatively, search from external discoverable DBs if any

	return nil, fmt.Errorf("GetOrchProviderProfileByName is not yet implemented")
}

// pvcOrchProviderProfile is a orchestration provider profile that is based on
// persistent volume claim.
//
// NOTE:
//    This will use defaults in-case the values are not set in persistent volume
// claim.
//
// NOTE:
//    This is a concrete implementation of orchprovider.VolumeProvisionerProfile
type pvcOrchProviderProfile struct {
	pvc *v1.PersistentVolumeClaim
}

// newPvcOrchProviderProfile provides a new instance of OrchProviderProfile that
// is based on pvc (i.e. persistent volume claim).
func newPvcOrchProviderProfile(pvc *v1.PersistentVolumeClaim) (OrchProviderProfile, error) {
	// This does not care if pvc instance is nil
	return &pvcOrchProviderProfile{
		pvc: pvc,
	}, nil
}

// defaults signals if this profile should use defaults entirely
func (op *pvcOrchProviderProfile) defaults() bool {
	if op.pvc == nil || op.pvc.Labels == nil {
		return true
	}

	return false
}

// Label provides the label assigned against the persistent volume provisioner
// profile.
//
// NOTE:
//    There can be many persistent volume provisioner profiles with this same label.
// This is used along with Name() method.
func (op *pvcOrchProviderProfile) Label() v1.OrchProviderProfileLabel {
	return v1.OrchProfileNameLbl
}

// Name provides the name assigned to the orchestration provider profile.
//
// NOTE:
//    Name provides the uniqueness among various variants of orchestration
// provider profiles.
func (op *pvcOrchProviderProfile) Name() v1.OrchProviderProfileRegistry {
	return v1.PVCOrchestratorProfile
}

// PVC provides the persistent volume claim associated with this profile.
//
// NOTE:
//    This method provides a convinient way to access pvc. In other words
// orchestration provider profile acts as a wrapper over pvc.
func (op *pvcOrchProviderProfile) PVC() (*v1.PersistentVolumeClaim, error) {
	return op.pvc, nil
}

// NetworkAddr gets the network address in CIDR format
func (op *pvcOrchProviderProfile) NetworkAddr() (string, error) {
	var nAddr string

	if op.defaults() {
		return v1.NetworkAddrDef(), nil
	}

	if nAddr = v1.NetworkAddr(op.pvc.Labels); nAddr == "" {
		return v1.NetworkAddrDef(), nil
	}

	if !nethelper.IsCIDR(nAddr) {
		return "", fmt.Errorf("Network address not in CIDR format in '%s:%s'", op.Label(), op.Name())
	}

	return nAddr, nil
}

// NetworkSubnet gets the network's subnet in decimal format
func (op *pvcOrchProviderProfile) NetworkSubnet() (string, error) {
	var nAddr string

	if op.defaults() {
		return v1.NetworkSubnetDef(), nil
	}

	nAddr, err := op.NetworkAddr()
	if err != nil {
		return "", err
	}

	if nAddr == "" {
		return v1.NetworkSubnetDef(), nil
	}

	subnet, err := nethelper.CIDRSubnet(nAddr)
	if err != nil {
		return "", err
	}

	return subnet, nil
}

// Get the namespace used at the orchestrator, where the request needs to be
// operated on
func (op *pvcOrchProviderProfile) NS() (string, error) {
	var ns string

	if op.defaults() {
		return v1.NSDef(), nil
	}

	if ns = v1.NS(op.pvc.Labels); ns == "" {
		return v1.NSDef(), nil
	}

	return ns, nil
}

// InCluster indicates if the request to the orchestrator is scoped to the
// cluster where this request originated
func (op *pvcOrchProviderProfile) InCluster() (bool, error) {
	var inCluster string

	if op.defaults() {
		return util.CheckTruthy(v1.InClusterDef()), nil
	}

	if inCluster = v1.InCluster(op.pvc.Labels); inCluster == "" {
		return util.CheckTruthy(v1.InClusterDef()), nil
	}

	return util.CheckTruthy(inCluster), nil
}

// TODO
//
// etcdOrchProviderProfile represents a generic orchestration provider profile
// whose properties are stored in etcd database.
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
//    Properties missing in etcd & persistent volume claim will make use of the
// defaults provided by maya api service
//
// NOTE:
//    This is a concrete implementation of volume.VolumeProvisionerProfile
type etcdOrchProviderProfile struct {
}

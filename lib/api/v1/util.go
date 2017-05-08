package v1

import (
	"strconv"

	"github.com/openebs/mayaserver/lib/util"
)

// ReplicaCount will fetch the value specified against persistent volume replica
// count if available otherwise will return blank.
//
// NOTE:
//    This function need not bother about any validations
func ReplicaCount(profileMap map[string]string) string {
	if profileMap == nil {
		return ""
	}

	// Extract replica count
	return profileMap[string(PVPReplicaCountLbl)]
}

// DefaultReplicaCount will fetch the default value of persistent volume
// provisioner replica count
//
// NOTE:
//    This function need not bother about any validations
func DefaultReplicaCount() int {
	iRCount, _ := strconv.Atoi(string(PVPReplicaCountDef))
	return iRCount
}

// NetworkAddr will fetch the value specified against orchestration provider
// network address if available otherwise will return blank.
//
// NOTE:
//    This function need not bother about any validations
func NetworkAddr(profileMap map[string]string) string {
	if profileMap == nil {
		return ""
	}

	// Extract network addr
	return profileMap[string(OPNetworkAddrLbl)]
}

// NetworkAddrDef will fetch the default value of orchestration provider network
// address.
//
// NOTE:
//    This function need not bother about any validations
func NetworkAddrDef() string {
	return string(OPNetworkAddrDef)
}

// NetworkSubnetDef will fetch the default value of orchestration provider
// network subnet.
//
// NOTE:
//    This function need not bother about any validations
func NetworkSubnetDef() string {
	return string(OPNetworkSubnetDef)
}

// ReplicaIPs will fetch the value specified against persistent volume replica
// IPs if available otherwise will return blank.
//
// NOTE:
//    This function need not bother about any validations
func ReplicaIPs(profileMap map[string]string) string {
	if profileMap == nil {
		return ""
	}

	// Extract replica IPs
	return profileMap[string(PVPReplicaIPsLbl)]
}

// ControllerIPs will fetch the value specified against persistent volume controller
// IPs if available otherwise will return blank.
//
// NOTE:
//    This function need not bother about any validations
func ControllerIPs(profileMap map[string]string) string {
	if profileMap == nil {
		return ""
	}

	// Extract controller IPs
	return profileMap[string(PVPControllerIPsLbl)]
}

// ReqNetworking will fetch the value specified against persistent volume networking
// support if available otherwise will return blank.
//
// NOTE:
//    This function need not bother about any validations
func ReqNetworking(profileMap map[string]string) string {
	if profileMap == nil {
		return ""
	}

	// Extract networking support i.e. is networking required
	return profileMap[string(PVPReqNetworkingLbl)]
}

// DefaultReqNetworking will fetch the default flag w.r.t persistent volume
// networking support
//
// NOTE:
//    This function need not bother about any validations
func DefaultReqNetworking() bool {
	return util.CheckTruthy(string(PVPReqNetworkingDef))
}

// StorageSize will fetch the value specified against persistent volume storage
// size if available otherwise will return blank.
//
// NOTE:
//    This function need not bother about any validations
func StorageSize(profileMap map[string]string) string {
	if profileMap == nil {
		return ""
	}

	// Extract storage size
	return profileMap[string(PVPStorageSizeLbl)]
}

// ReplicaImage will fetch the value specified against persistent volume replica
// image if available otherwise will return blank.
//
// NOTE:
//    This function need not bother about any validations
func ReplicaImage(profileMap map[string]string) string {
	if profileMap == nil {
		return ""
	}

	// Extract replica image
	return profileMap[string(PVPReplicaImageLbl)]
}

// DefaultReplicaImage will fetch the default value for persistent
// volume replica image
//
// NOTE:
//    This function need not bother about any validations
func DefaultReplicaImage() string {
	return string(PVPReplicaImageDef)
}

// ReqReplica will fetch the value specified against persistent volume replica
// support if available otherwise will return blank.
//
// NOTE:
//    This function need not bother about any validations
func ReqReplica(profileMap map[string]string) string {
	if profileMap == nil {
		return ""
	}

	// Extract replica support i.e. is replica required
	return profileMap[string(PVPReqReplicaLbl)]
}

// DefaultReqReplica will fetch the default value for persistent volume
// replica support.
//
// NOTE:
//    This function need not bother about any validations
func DefaultReqReplica() bool {
	// Extract replica support i.e. is replica required
	return util.CheckTruthy(string(PVPReqReplicaDef))
}

// ControllerImage will fetch the value specified against persistent volume
// controller image if available otherwise will return blank.
//
// NOTE:
//    This function need not bother about any validations
func ControllerImage(profileMap map[string]string) string {
	if profileMap == nil {
		return ""
	}

	// Extract controller image
	return profileMap[string(PVPControllerImageLbl)]
}

// DefaultControllerImage will fetch the default value for persistent
// volume controller image
//
// NOTE:
//    This function need not bother about any validations
func DefaultControllerImage() string {
	return string(PVPControllerImageDef)
}

// ControllerCount will fetch the value specified against persistent volume
// controller count if available otherwise will return blank.
//
// NOTE:
//    This function need not bother about any validations
func ControllerCount(profileMap map[string]string) string {
	if profileMap == nil {
		return ""
	}

	// Extract controller count
	return profileMap[string(PVPControllerCountLbl)]
}

// DefaultControllerCount will fetch the default value for persistent
// volume controller count
//
// NOTE:
//    This function need not bother about any validations
func DefaultControllerCount() int {
	iCCount, _ := strconv.Atoi(string(PVPControllerCountDef))
	return iCCount
}

// VSMName will fetch the value specified against persistent volume
// VSM name if available otherwise will return blank.
//
// NOTE:
//    This function need not bother about any validations
func VSMName(profileMap map[string]string) string {
	if profileMap == nil {
		return ""
	}

	// Extract VSM name
	return profileMap[string(PVPVSMNameLbl)]
}

// OrchProfileName will fetch the value specified against persistent volume's
// orchestrator profile name if available otherwise will return blank.
//
// NOTE:
//    This function need not bother about any validations
func OrchProfileName(profileMap map[string]string) string {
	if profileMap == nil {
		return ""
	}

	// Extract orchestrator profile name
	return profileMap[string(OrchProfileNameLbl)]
}

// NS will fetch the value specified against orchestration provider
// namespace if available otherwise will return blank.
//
// NOTE:
//    This function need not bother about any validations
func NS(profileMap map[string]string) string {
	if profileMap == nil {
		return ""
	}

	// Extract orchestrator namespace
	return profileMap[string(OrchNSLbl)]
}

// NSDef will fetch the default value of orchestration provider namespace.
//
// NOTE:
//    This function need not bother about any validations
func NSDef() string {
	return string(OrchNSDefLbl)
}

// PodName will fetch the value specified against persistent volume
// VSM name if available otherwise will return blank.
//
// NOTE:
//    This function need not bother about any validations
func PodName(profileMap map[string]string) string {
	// Extract VSM name
	return VSMName(profileMap)
}

// InCluster will fetch the value specified against orchestration provider
// in-cluster flag if available otherwise will return blank.
//
// NOTE:
//    This function need not bother about any validations
func InCluster(profileMap map[string]string) string {
	if profileMap == nil {
		return ""
	}

	// Extract orchestrator in-cluster flag
	return profileMap[string(OPInClusterLbl)]
}

// InClusterDef will fetch the default value of orchestration provider
// in-cluster flag.
//
// NOTE:
//    This function need not bother about any validations
func InClusterDef() string {
	return string(OPInClusterDef)
}

// VolumeProvisionerProfileName will fetch the name of volume provisioner
// profile if available otherwise will return blank.
//
// NOTE:
//    This function need not bother about any validations
func VolumeProvisionerProfileName(profileMap map[string]string) string {
	if profileMap == nil {
		return ""
	}

	// Extract volume provisioner profile name
	return profileMap[string(PVPProfileNameLbl)]
}

// DefaultVolumeProvisionerName gets the default name of persistent volume
// provisioner plugin used to cater the provisioning requests to maya api
// service
//
// NOTE:
//    This function need not bother about any validations
func DefaultVolumeProvisionerName() VolumeProvisionerRegistry {
	return DefaultVolumeProvisioner
}

// DefaultOrchestratorName gets the default name of orchestration provider
// plugin used to cater the provisioning requests to maya api service
//
// NOTE:
//    This function need not bother about any validations
func DefaultOrchestratorName() OrchProviderRegistry {
	return DefaultOrchestrator
}

// OrchestratorName will fetch the value specified against persistent
// volume's orchestrator name if available otherwise will return blank.
//
// NOTE:
//    This function need not bother about any validations
func OrchestratorName(profileMap map[string]string) OrchProviderRegistry {
	if profileMap == nil {
		return OrchProviderRegistry("")
	}

	// Extract orchestrator name
	return OrchProviderRegistry(profileMap[string(OrchestratorNameLbl)])
}

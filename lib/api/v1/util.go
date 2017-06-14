package v1

import (
	"strconv"
	"strings"

	"github.com/openebs/mayaserver/lib/util"
)

// ReplicaCount will fetch the value specified against persistent volume replica
// count if available otherwise will return blank.
//
// NOTE:
//    This utility function does not validate & just returns if not capable of
// performing
func ReplicaCount(profileMap map[string]string) string {
	if profileMap == nil {
		return ""
	}

	// Extract replica count
	return profileMap[string(PVPReplicaCountLbl)]
}

// DefaultReplicaCount will fetch the default value of persistent volume
// provisioner replica count
func DefaultReplicaCount() int {
	iRCount, _ := strconv.Atoi(string(PVPReplicaCountDef))
	return iRCount
}

// NetworkAddr will fetch the value specified against orchestration provider
// network address if available otherwise will return blank.
//
// NOTE:
//    This utility function does not validate & just returns if not capable of
// performing
func NetworkAddr(profileMap map[string]string) string {
	if profileMap == nil {
		return ""
	}

	// Extract network addr
	return profileMap[string(OPNetworkAddrLbl)]
}

// NetworkAddrDef will fetch the default value of orchestration provider network
// address.
func NetworkAddrDef() string {
	return string(OPNetworkAddrDef)
}

// NetworkSubnetDef will fetch the default value of orchestration provider
// network subnet.
func NetworkSubnetDef() string {
	return string(OPNetworkSubnetDef)
}

// ReplicaIPs will fetch the value specified against persistent volume replica
// IPs if available otherwise will return blank.
//
// NOTE:
//    This utility function does not validate & just returns if not capable of
// performing
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
//    This utility function does not validate & just returns if not capable of
// performing
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
//    This utility function does not validate & just returns if not capable of
// performing
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
//    This utility function does not validate & just returns if not capable of
// performing
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
//    This utility function does not validate & just returns if not capable of
// performing
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
//    This utility function does not validate & just returns if not capable of
// performing
func ReqReplica(profileMap map[string]string) string {
	if profileMap == nil {
		return ""
	}

	// Extract replica support i.e. is replica required
	return profileMap[string(PVPReqReplicaLbl)]
}

// DefaultReqReplica will fetch the default value for persistent volume
// replica support.
func DefaultReqReplica() bool {
	// Extract replica support i.e. is replica required
	return util.CheckTruthy(string(PVPReqReplicaDef))
}

// ControllerImage will fetch the value specified against persistent volume
// controller image if available otherwise will return blank.
//
// NOTE:
//    This utility function does not validate & just returns if not capable of
// performing
func ControllerImage(profileMap map[string]string) string {
	if profileMap == nil {
		return ""
	}

	// Extract controller image
	return profileMap[string(PVPControllerImageLbl)]
}

// DefaultControllerImage will fetch the default value for persistent
// volume controller image
func DefaultControllerImage() string {
	return string(PVPControllerImageDef)
}

// ControllerCount will fetch the value specified against persistent volume
// controller count if available otherwise will return blank.
//
// NOTE:
//    This utility function does not validate & just returns if not capable of
// performing
func ControllerCount(profileMap map[string]string) string {
	if profileMap == nil {
		return ""
	}

	// Extract controller count
	return profileMap[string(PVPControllerCountLbl)]
}

// DefaultControllerCount will fetch the default value for persistent
// volume controller count
func DefaultControllerCount() int {
	iCCount, _ := strconv.Atoi(string(PVPControllerCountDef))
	return iCCount
}

// VSMName will fetch the value specified against persistent volume
// VSM name if available otherwise will return blank.
//
// NOTE:
//    This utility function does not validate & just returns if not capable of
// performing
func VSMName(pvcName string) string {
	// Name of PVC is the name of VSM
	return pvcName
}

// OrchProfileName will fetch the value specified against persistent volume's
// orchestrator profile name if available otherwise will return blank.
//
// NOTE:
//    This utility function does not validate & just returns if not capable of
// performing
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
//    This utility function does not validate & just returns if not capable of
// performing
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
//func PodName(profileMap map[string]string) string {
// Extract VSM name
//	return VSMName(profileMap)
//}

// InCluster will fetch the value specified against orchestration provider
// in-cluster flag if available otherwise will return blank.
//
// NOTE:
//    This utility function does not validate & just returns if not capable of
// performing
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
//    This utility function does not validate & just returns if not capable of
// performing
func InClusterDef() string {
	return string(OPInClusterDef)
}

// VolumeProvisionerProfileName will fetch the name of volume provisioner
// profile if available otherwise will return blank.
//
// NOTE:
//    This utility function does not validate & just returns if not capable of
// performing
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
//    This utility function does not validate & just returns if not capable of
// performing
func DefaultVolumeProvisionerName() VolumeProvisionerRegistry {
	return DefaultVolumeProvisioner
}

// DefaultOrchestratorName gets the default name of orchestration provider
// plugin used to cater the provisioning requests to maya api service
//
// NOTE:
//    This utility function does not validate & just returns if not capable of
// performing
func DefaultOrchestratorName() OrchProviderRegistry {
	return DefaultOrchestrator
}

// OrchestratorName will fetch the value specified against persistent
// volume's orchestrator name if available otherwise will return blank.
//
// NOTE:
//    This utility function does not validate & just returns if not capable of
// performing
func OrchestratorName(profileMap map[string]string) OrchProviderRegistry {
	if profileMap == nil {
		return OrchProviderRegistry("")
	}

	// Extract orchestrator name
	return OrchProviderRegistry(profileMap[string(OrchestratorNameLbl)])
}

// DefaultJivaISCSIPort will provide the port required to make ISCSI based
// connections
func DefaultJivaISCSIPort() int32 {
	iscsiPort, _ := strconv.Atoi(string(JivaISCSIPortDef))
	return int32(iscsiPort)
}

// DefaultJivaAPIPort will provide the port required for management of
// persistent volume
func DefaultJivaAPIPort() int32 {
	apiPort, _ := strconv.Atoi(string(JivaAPIPortDef))
	return int32(apiPort)
}

// DefaultPersistentPathCount will provide the default count of persistent
// paths required during provisioning.
func DefaultPersistentPathCount() int {
	pCount, _ := strconv.Atoi(string(PVPPersistentPathCountDef))
	return pCount
}

// PersistentPathCount will fetch the value specified against persistent volume
// persistent path count if available otherwise will return blank.
//
// NOTE:
//    This utility function does not validate & just returns if not capable of
// performing
func PersistentPathCount(profileMap map[string]string) string {
	if profileMap == nil {
		return ""
	}

	// Extract persistent path count
	return profileMap[string(PVPPersistentPathCountLbl)]
}

// JivaPersistentPath will fetch the value specified against persistent volume
// persistent host path if available otherwise will return blank.
//
// NOTE:
//    This utility function does not validate & just returns if not capable of
// performing
func JivaPersistentPath(profileMap map[string]string, vsm string, position int) string {
	if profileMap == nil || profileMap[string(PVPPersistentPathLbl)] == "" {
		return ""
	}

	// Extract persistent path count
	// We are not using position for the time-being
	// We may not require it probably.
	return profileMap[string(PVPPersistentPathLbl)] + "/" + vsm + string(JivaPersistentMountPathDef)
}

// Replicas returns a pointer to an int32 of a int value
func Replicas(rcount int) *int32 {
	o := int32(rcount)
	return &o
}

// TODO
// Move this to api/v1/jiva/util.go
//
// DefaultJivaPersistentPath provides the default persistent host path based on the
// name of the VSM & replica position
func DefaultJivaPersistentPath(vsm string, position int) string {
	// TODO
	// We are not using position for the time-being
	// We may not require it probably.
	return string(JivaPersistentPathDef) + "/" + vsm + string(JivaPersistentMountPathDef)
}

// MakeOrDefJivaReplicaArgs will set the placeholders in jiva replica args with
// their appropriate runtime values.
//
// NOTE:
//    The defaults will be set if the replica args are not available
//
// NOTE:
//    This utility function does not validate & just returns if not capable of
// performing
func MakeOrDefJivaReplicaArgs(profileMap map[string]string, ctrlIP string) []string {
	if strings.TrimSpace(ctrlIP) == "" {
		return nil
	}

	storSize := ""
	if profileMap == nil {
		storSize = string(JivaStorSizeDef)
	} else {
		// Extract the runtime values
		storSize = profileMap[string(PVPStorageSizeLbl)]
	}

	// If runtime was set with blank
	if storSize == "" {
		storSize = string(JivaStorSizeDef)
	}

	repArgs := make([]string, len(JivaReplicaArgs))

	for i, rArg := range JivaReplicaArgs {
		rArg = strings.Replace(rArg, string(JivaCtrlIPHolder), ctrlIP, 1)
		rArg = strings.Replace(rArg, string(JivaStorageSizeHolder), storSize, 1)
		repArgs[i] = rArg
	}

	return repArgs
}

//
func MakeOrDefJivaControllerArgs(vsm string) []string {
	if strings.TrimSpace(vsm) == "" {
		return nil
	}

	ctrlArgs := make([]string, len(JivaCtrlArgs))

	for i, cArg := range JivaCtrlArgs {
		cArg = strings.Replace(cArg, string(JivaVolumeNameHolder), vsm, 1)
		ctrlArgs[i] = cArg
	}

	return ctrlArgs
}

// DefaultJivaMountPath provides the default mount path for jiva based persistent
// volumes
func DefaultJivaMountPath() string {
	return string(JivaPersistentMountPathDef)
}

// DefaultJivaMountName provides the default mount path name for jiva based
// persistent volumes
func DefaultJivaMountName() string {
	return string(JivaPersistentMountNameDef)
}

// DefaultJivaReplicaPort1 provides the default port for jiva based
// persistent volume replicas
func DefaultJivaReplicaPort1() int32 {
	p, _ := strconv.Atoi(string(JivaReplicaPortOneDef))
	return int32(p)
}

// DefaultJivaReplicaPort2 provides the default port for jiva based
// persistent volume replicas
func DefaultJivaReplicaPort2() int32 {
	p, _ := strconv.Atoi(string(JivaReplicaPortTwoDef))
	return int32(p)
}

// DefaultJivaReplicaPort3 provides the default port for jiva based
// persistent volume replicas
func DefaultJivaReplicaPort3() int32 {
	p, _ := strconv.Atoi(string(JivaReplicaPortThreeDef))
	return int32(p)
}

//
func SanitiseVSMName(vsm string) string {
	// Trim the controller suffix if controller based
	v := strings.TrimSuffix(vsm, string(ControllerSuffix))
	// Or Trim the replica suffix if replica based
	v = strings.TrimSuffix(v, string(ReplicaSuffix))

	return v
}

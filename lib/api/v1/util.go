package v1

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

// ReplicaSize will fetch the value specified against persistent volume replica
// size if available otherwise will return blank.
//
// NOTE:
//    This function need not bother about any validations
func ReplicaSize(profileMap map[string]string) string {
	if profileMap == nil {
		return ""
	}

	// Extract replica size
	return profileMap[string(PVPReplicaSizeLbl)]
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

// ControllerSize will fetch the value specified against persistent volume
// controller size if available otherwise will return blank.
//
// NOTE:
//    This function need not bother about any validations
func ControllerSize(profileMap map[string]string) string {
	if profileMap == nil {
		return ""
	}

	// Extract controller size
	return profileMap[string(PVPControllerSizeLbl)]
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

// OrchestratorName will fetch the value specified against persistent volume's
// orchestrator name if available otherwise will return blank.
//
// NOTE:
//    This function need not bother about any validations
func OrchestratorName(profileMap map[string]string) string {
	if profileMap == nil {
		return ""
	}

	// Extract orchestrator name
	return profileMap[string(OrchestratorNameLbl)]
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

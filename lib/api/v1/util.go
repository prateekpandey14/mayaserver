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

// NetworkAddr will fetch the value specified against persistent volume network
// address if available otherwise will return blank.
//
// NOTE:
//    This function need not bother about any validations
func NetworkAddr(profileMap map[string]string) string {
	if profileMap == nil {
		return ""
	}

	// Extract network addr
	return profileMap[string(PVPNetworkAddrLbl)]
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

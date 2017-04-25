package jiva

import (
	"github.com/openebs/mayaserver/lib/api/v1"
	v1nomad "github.com/openebs/mayaserver/lib/api/v1/nomad"
)

// These label can be used in the persistent volume claims to specify
// exact requirements.
type JivaLbl string

const (
	JivaFrontEndImageLbl JivaLbl = "fe.jiva.volume.openebs.io/image-version"
	JivaBackEndImageLbl  JivaLbl = "be.jiva.volume.openebs.io/image-version"

	JivaFrontEndIPLbl JivaLbl = "fe.jiva.volume.openebs.io/ip"
	JivaBackEndIPLbl  JivaLbl = "be.jiva.volume.openebs.io/ip"

	// TODO remove this & use the generic v1.labels
	JivaFrontEndAllIPsLbl JivaLbl = "fe.jiva.volume.openebs.io/all-ips"
	JivaBackEndAllIPsLbl  JivaLbl = "be.jiva.volume.openebs.io/all-ips"

	// TODO remove this & use the generic v1.labels
	JivaFrontEndCountLbl JivaLbl = "fe.jiva.volume.openebs.io/count"
	JivaBackEndCountLbl  JivaLbl = "be.jiva.volume.openebs.io/count"

	// TODO
	// Will it be good to namespace these labels ?
	// JivaTargetPortalLbl is a label / tag that is used to provide a value for
	// jiva's frontend target portal.
	JivaTargetPortalLbl JivaLbl = "targetportal"
	// JivaIqnLbl is a label / tag that is used to provide a value for
	// jiva's iqn.
	JivaIqnLbl JivaLbl = "iqn"

	// TODO
	// This should be a mere constant !!
	JivaBackEndIPPrefixLbl JivaLbl = "JIVA_REP_IP_"
)

const (
	// JivaFrontEndVolSizeLbl is a label / tag that is used to provide a value for
	// jiva's frontend volume's size.
	JivaFrontEndVolSizeLbl v1.ResourceName = "fe.jiva.volume.openebs.io/vol-size"
	// JivaBackEndVolSizeLbl is a label / tag that is used to provide a value for
	// jiva's backend volume's size.
	JivaBackEndVolSizeLbl v1.ResourceName = "be.jiva.volume.openebs.io/vol-size"
)

const (
	// JivaVolumeProvisionerName is a value used for registering Jiva as a volume
	// provisioner in maya api server.
	//
	// NOTE:
	//    This label-value / tag-value can be overridden by user specified value when
	// provided with corresponding label / tag.
	JivaVolumeProvisionerName = "jiva"

	// Jiva's iSCSI Qualified port value.
	//
	// NOTE:
	//    This label-value / tag-value can be overridden by user specified value when
	// provided with corresponding label / tag.
	JivaIscsiTargetPortalPort string = "3260"

	// Jiva's iSCSI Qualified IQN value.
	//
	// NOTE:
	//    This label-value / tag-value can be overridden by user specified value when
	// provided with corresponding label / tag.
	JivaIqnFormatPrefix string = "iqn.2016-09.com.openebs.jiva"

	// TODO
	// Remove deprecated
	// Deprecated
	//
	// This naming is a result of considering Jiva volume plugin's default
	// orchestrator which is Nomad & this default orchestrator's default region
	// which is global.
	DefaultJivaVolumePluginName string = v1.VolumePluginNamePrefix + "jiva-nomad-" + v1nomad.DefaultNomadRegionName

	// TODO
	// Remove deprecate
	// Deprecate
	//
	// This just points to Nomad orchestrator's default dc.
	DefaultJivaDataCenter string = v1nomad.DefaultNomadDCName

	// TODO
	// Remove deprecate
	// Deprecate
	//
	// This naming signifies a prefix that can be used with variants of
	// jiva volume plugin.
	//
	// NOTE:
	// Some sample variants of jiva volume plugin:
	//
	//  Jiva + K8S + us-east-1
	//  Jiva + Nomad + global
	//  Jiva + Nomad + in-bang
	JivaVolumePluginNamePrefix string = v1.VolumePluginNamePrefix + "jiva-"
)

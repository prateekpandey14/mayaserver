package nomad

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/nomad/helper"
	"github.com/hashicorp/nomad/nomad/structs"
	"github.com/openebs/mayaserver/lib/api/v1"
	v1jiva "github.com/openebs/mayaserver/lib/api/v1/jiva"
)

// Get the job name from a persistent volume claim
func PvcToJobName(pvc *v1.PersistentVolumeClaim) (string, error) {

	if pvc == nil {
		return "", fmt.Errorf("Nil persistent volume claim provided")
	}

	if pvc.Name == "" {
		return "", fmt.Errorf("Missing VSM name in pvc")
	}

	return pvc.Name, nil
}

// Transform a PersistentVolumeClaim type to Nomad job type
//
// TODO
// There is redundancy in validation. These should be gone once the transformation
// is handled from jiva namespace.
//
// TODO
// It may be better to avoid using pvc as a direct argument
// It should be key:value pairs
// The keys can be of type SpecsLbl
// These keys can in turn be used in PVC as well
// However, the argument to this function should be a struct that consists of
// a map of SpecsLbl:Value
func PvcToJob(pvc *v1.PersistentVolumeClaim) (*api.Job, error) {

	if pvc == nil {
		return nil, fmt.Errorf("Nil persistent volume claim provided")
	}

	if pvc.Name == "" {
		return nil, fmt.Errorf("Missing VSM name in pvc")
	}

	if pvc.Labels == nil {
		return nil, fmt.Errorf("Missing labels in pvc")
	}

	if pvc.Labels[string(v1.OrchRegionLbl)] == "" {
		return nil, fmt.Errorf("Missing region in pvc")
	}

	if pvc.Labels[string(v1.OrchDCLbl)] == "" {
		return nil, fmt.Errorf("Missing datacenter in pvc")
	}

	//if pvc.Labels[string(v1jiva.JivaFrontEndImageLbl)] == "" {
	if pvc.Labels[string(v1.PVPControllerImageLbl)] == "" {
		return nil, fmt.Errorf("Missing controller image in pvc")
	}

	//if pvc.Labels[string(v1jiva.JivaFrontEndIPLbl)] == "" {
	if pvc.Labels[string(v1.PVPControllerIPsLbl)] == "" {
		return nil, fmt.Errorf("Missing controller IP(s) in pvc")
	}

	//if pvc.Labels[string(v1jiva.JivaBackEndAllIPsLbl)] == "" {
	if pvc.Labels[string(v1.PVPReplicaIPsLbl)] == "" {
		return nil, fmt.Errorf("Missing replica IPs in pvc")
	}

	// TODO These should be derived from:
	//
	// pvc.Spec.NetworkResources
	// pvc.Spec.StorageResources
	// These types are currently strict but less flexible

	if pvc.Labels[string(v1.CNTypeLbl)] == "" {
		return nil, fmt.Errorf("Missing CN type in pvc")
	}

	if pvc.Labels[string(v1.CNSubnetLbl)] == "" {
		return nil, fmt.Errorf("Missing CN subnet in pvc")
	}

	if pvc.Labels[string(v1.CNInterfaceLbl)] == "" {
		return nil, fmt.Errorf("Missing CN interface in pvc")
	}

	//if pvc.Labels[string(v1.CSPersistenceLocationLbl)] == "" {
	if pvc.Labels[string(v1.PVPPersistenceLocationLbl)] == "" {
		return nil, fmt.Errorf("Missing persistent location in pvc")
	}

	if pvc.Labels[string(v1.PVPStorageSizeLbl)] == "" {
		return nil, fmt.Errorf("Missing storage size in pvc")
	}

	jivaFEVolSize := pvc.Labels[string(v1.PVPStorageSizeLbl)]
	jivaBEVolSize := pvc.Labels[string(v1.PVPStorageSizeLbl)]

	// TODO
	// ID is same as Name currently
	// Do we need to think on it ?
	jobName := helper.StringToPtr(pvc.Name)
	region := helper.StringToPtr(pvc.Labels[string(v1.OrchRegionLbl)])
	dc := pvc.Labels[string(v1.OrchDCLbl)]

	jivaGroupName := "jiva-pod"
	jivaVolName := pvc.Name

	// Set storage size
	feTaskGroup := "fe" + "-" + jivaGroupName
	beTaskGroup := "be" + "-" + jivaGroupName

	// Default storage policy would required 1 FE & 2 BE
	feTaskName := "fe"
	beTaskName := "be"

	jivaFeVersion := pvc.Labels[string(v1.PVPControllerImageLbl)]
	jivaNetworkType := pvc.Labels[string(v1.CNTypeLbl)]
	jivaFeIP := pvc.Labels[string(v1.PVPControllerIPsLbl)]

	jivaBEPersistentStor := pvc.Labels[string(v1.PVPPersistenceLocationLbl)]
	jivaBECount := pvc.Labels[string(v1.PVPReplicaCountLbl)]
	iJivaBECount, err := strconv.Atoi(jivaBECount)
	if err != nil {
		return nil, err
	}

	jivaBeIPs := pvc.Labels[string(v1.PVPReplicaIPsLbl)]
	jivaBeIPArr := strings.Split(jivaBeIPs, ",")
	jivaFeSubnet := pvc.Labels[string(v1.CNSubnetLbl)]
	jivaFeInterface := pvc.Labels[string(v1.CNInterfaceLbl)]

	// Meta information will be used to:
	//    1. Persist metadata w.r.t this job
	//
	// NOTE:
	//    This enables to query various info w.r.t job much later.
	// In addition, job's ENV property can source these metadata by interpreting
	// them.
	jobMeta := map[string]string{
		string(v1jiva.JivaBackEndVolSizeLbl): jivaBEVolSize,
		string(v1jiva.JivaFrontEndIPLbl):     jivaFeIP,
		string(v1jiva.JivaTargetPortalLbl):   jivaFeIP + ":" + string(v1.JivaISCSIPortDef),
		string(v1jiva.JivaIqnLbl):            string(v1.JivaIqnFormatPrefix) + ":" + jivaVolName,
	}

	// Jiva FE's ENV among other things interpolates Nomad's built-in properties
	feEnv := map[string]string{
		"JIVA_CTL_NAME":    pvc.Name + "-" + feTaskName + "${NOMAD_ALLOC_INDEX}",
		"JIVA_CTL_VERSION": jivaFeVersion,
		"JIVA_CTL_VOLNAME": jivaVolName,
		"JIVA_CTL_VOLSIZE": jivaFEVolSize,
		"JIVA_CTL_IP":      jivaFeIP,
		"JIVA_CTL_SUBNET":  jivaFeSubnet,
		"JIVA_CTL_IFACE":   jivaFeInterface,
	}

	// Jiva BE's ENV among other things interpolates Nomad's built-in properties
	beEnv := map[string]string{
		"NOMAD_ALLOC_INDEX": "${NOMAD_ALLOC_INDEX}",
		"JIVA_REP_NAME":     pvc.Name + "-" + beTaskName + "${NOMAD_ALLOC_INDEX}",
		"JIVA_CTL_IP":       jivaFeIP,
		"JIVA_REP_VOLNAME":  jivaVolName,
		"JIVA_REP_VOLSIZE":  jivaBEVolSize,
		"JIVA_REP_VOLSTORE": jivaBEPersistentStor + pvc.Name + "/" + beTaskName + "${NOMAD_ALLOC_INDEX}",
		"JIVA_REP_VERSION":  jivaFeVersion,
		"JIVA_REP_NETWORK":  jivaNetworkType,
		"JIVA_REP_IFACE":    jivaFeInterface,
		"JIVA_REP_SUBNET":   jivaFeSubnet,
	}

	// This sets below variables with backend IP addresses:
	//
	//  1. job's backend's ENV pairs
	//  2. job's META pairs
	err = setBEIPs(beEnv, jobMeta, jivaBeIPArr, iJivaBECount)
	if err != nil {
		return nil, err
	}

	// TODO
	// Transformation from pvc or pv to nomad types & vice-versa:
	//
	//  1. Need an Interface or functional callback defined at
	// lib/api/v1/nomad/ &
	//  2. implemented by the volume plugins that want
	// to be orchestrated by Nomad
	//  3. This transformer instance needs to be injected from
	// volume plugin to orchestrator, in a generic way.

	// Hardcoded logic all the way
	// Nomad specific defaults, hardcoding is OK.
	// However, volume plugin specific stuff is BAD
	return &api.Job{
		Region:      region,
		Name:        jobName,
		ID:          jobName,
		Datacenters: []string{dc},
		Type:        helper.StringToPtr(api.JobTypeService),
		Priority:    helper.IntToPtr(50),
		Constraints: []*api.Constraint{
			api.NewConstraint("${attr.kernel.name}", "=", "linux"),
		},
		Meta: jobMeta,
		TaskGroups: []*api.TaskGroup{
			// jiva frontend
			&api.TaskGroup{
				Name:  helper.StringToPtr(feTaskGroup),
				Count: helper.IntToPtr(1),
				RestartPolicy: &api.RestartPolicy{
					Attempts: helper.IntToPtr(3),
					Interval: helper.TimeToPtr(5 * time.Minute),
					Delay:    helper.TimeToPtr(25 * time.Second),
					Mode:     helper.StringToPtr("delay"),
				},
				Tasks: []*api.Task{
					&api.Task{
						Name:   feTaskName,
						Driver: "raw_exec",
						Resources: &api.Resources{
							CPU:      helper.IntToPtr(50),
							MemoryMB: helper.IntToPtr(50),
							Networks: []*api.NetworkResource{
								&api.NetworkResource{
									MBits: helper.IntToPtr(50),
								},
							},
						},
						Env: feEnv,
						Artifacts: []*api.TaskArtifact{
							&api.TaskArtifact{
								GetterSource: helper.StringToPtr("https://raw.githubusercontent.com/openebs/jiva/master/scripts/launch-jiva-ctl-with-ip"),
								RelativeDest: helper.StringToPtr("local/"),
							},
						},
						Config: map[string]interface{}{
							"command": "launch-jiva-ctl-with-ip",
						},
						LogConfig: &api.LogConfig{
							MaxFiles:      helper.IntToPtr(3),
							MaxFileSizeMB: helper.IntToPtr(1),
						},
					},
				},
			},
			// jiva replica group
			&api.TaskGroup{
				Name: helper.StringToPtr(beTaskGroup),
				// Replica count
				Count: helper.IntToPtr(iJivaBECount),
				// We want the replicas to spread across hosts
				// This ensures high availability
				Constraints: []*api.Constraint{
					api.NewConstraint("", "distinct_hosts", "true"),
				},
				RestartPolicy: &api.RestartPolicy{
					Attempts: helper.IntToPtr(3),
					Interval: helper.TimeToPtr(5 * time.Minute),
					Delay:    helper.TimeToPtr(25 * time.Second),
					Mode:     helper.StringToPtr("delay"),
				},
				// This has multiple replicas as tasks
				Tasks: []*api.Task{
					&api.Task{
						Name:   beTaskName,
						Driver: "raw_exec",
						Resources: &api.Resources{
							CPU:      helper.IntToPtr(50),
							MemoryMB: helper.IntToPtr(50),
							Networks: []*api.NetworkResource{
								&api.NetworkResource{
									MBits: helper.IntToPtr(50),
								},
							},
						},
						Env: beEnv,
						Artifacts: []*api.TaskArtifact{
							&api.TaskArtifact{
								GetterSource: helper.StringToPtr("https://raw.githubusercontent.com/openebs/jiva/master/scripts/launch-jiva-rep-with-ip"),
								RelativeDest: helper.StringToPtr("local/"),
							},
						},
						Config: map[string]interface{}{
							"command": "launch-jiva-rep-with-ip",
						},
						LogConfig: &api.LogConfig{
							MaxFiles:      helper.IntToPtr(3),
							MaxFileSizeMB: helper.IntToPtr(1),
						},
					},
				},
			},
		},
	}, nil
}

// setBEIPs sets jiva backend environment with all backend IP addresses
func setBEIPs(beEnv, jobMeta map[string]string, jivaBeIPArr []string, iJivaBECount int) error {

	if iJivaBECount == 0 {
		return fmt.Errorf("Invalid replica count '0'")
	}

	if len(jivaBeIPArr) != iJivaBECount {
		return fmt.Errorf("Replica IP count '%d' does not match replica count '%d'", len(jivaBeIPArr), iJivaBECount)
	}

	jobMeta[string(v1.PVPReplicaCountLbl)] = strconv.Itoa(iJivaBECount)

	var k, v string

	for i := 0; i < iJivaBECount; i++ {
		k = string(v1jiva.JivaBackEndIPPrefixLbl) + strconv.Itoa(i)
		v = jivaBeIPArr[i]

		beEnv[k] = v
		jobMeta[k] = v
	}

	return nil
}

// TODO
// Transformation from JobSummary to pv
//
//  1. Need an Interface or functional callback defined at
// lib/api/v1/nomad.go &
//  2. implemented by the volume plugins that want
// to be orchestrated by Nomad
//  3. This transformer instance needs to be injected from
// volume plugin to orchestrator, in a generic way.
//func JobSummaryToPv(jobSummary *api.JobSummary) (*v1.PersistentVolume, error) {
//
//	if jobSummary == nil {
//		return nil, fmt.Errorf("Nil nomad job summary provided")
//	}
//
// TODO
// Needs to be filled up
//	return &v1.PersistentVolume{}, nil
//}

// TODO
// Transform the evaluation of a job to a PersistentVolume
func JobEvalToPv(jobName string, eval *api.Evaluation) (*v1.PersistentVolume, error) {

	if eval == nil {
		return nil, fmt.Errorf("Nil job evaluation provided")
	}

	pv := &v1.PersistentVolume{}
	pv.Name = jobName

	evalProps := map[string]string{
		"evalpriority":    strconv.Itoa(eval.Priority),
		"evaltype":        eval.Type,
		"evaltrigger":     eval.TriggeredBy,
		"evaljob":         eval.JobID,
		"evalstatus":      eval.Status,
		"evalstatusdesc":  eval.StatusDescription,
		"evalblockedeval": eval.BlockedEval,
	}
	pv.Annotations = evalProps

	pvs := v1.PersistentVolumeStatus{
		Message: eval.StatusDescription,
		Reason:  eval.Status,
	}
	pv.Status = pvs

	return pv, nil
}

// Transform a PersistentVolume type to Nomad Job type
func PvToJob(pv *v1.PersistentVolume) (*api.Job, error) {

	if pv == nil {
		return nil, fmt.Errorf("Nil persistent volume provided")
	}

	return &api.Job{
		Name: helper.StringToPtr(pv.Name),
		// TODO
		// ID is same as Name currently
		ID: helper.StringToPtr(pv.Name),
	}, nil
}

// Transform a Nomad Job to a PersistentVolume
func JobToPv(job *api.Job) (*v1.PersistentVolume, error) {
	if job == nil {
		return nil, fmt.Errorf("Nil job provided")
	}

	pv := &v1.PersistentVolume{}
	pv.Name = *job.Name

	pvs := v1.PersistentVolumeStatus{
		Message: *job.StatusDescription,
		Reason:  *job.Status,
	}
	pv.Status = pvs

	if *job.Status == structs.JobStatusRunning {
		// Remember we use the job's metadata to persist metadata w.r.t the job
		pv.Annotations = job.Meta
	}

	return pv, nil
}

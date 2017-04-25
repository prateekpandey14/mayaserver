// This file registers jiva as maya api server's persistent volume provisioner plugin.
package jiva

import (
	"fmt"
	"io"

	"github.com/golang/glog"
	"github.com/openebs/mayaserver/lib/api/v1"
	v1jiva "github.com/openebs/mayaserver/lib/api/v1/jiva"
	"github.com/openebs/mayaserver/lib/orchprovider"
	"github.com/openebs/mayaserver/lib/volume"
)

// The registration logic for jiva persistent volume provisioner plugin
//
// NOTE:
//    This function is executed once per application
//
// TODO
//  A simplified version of registration logic will be implemented. This in turn
// will enable the registry to create new instances of jiva persistent volume
// provisioner on each request.
func init() {
	// TODO
	// Remove the deprecated registration style
	// Deprecated registration style !!
	volume.RegisterVolumePlugin(
		// A variant of jiva volume plugin
		v1jiva.DefaultJivaVolumePluginName,
		// Below is a functional implementation that holds the initialization
		// logic of jiva volume plugin
		func(name string, config io.Reader, aspect volume.VolumePluginAspect) (volume.VolumeInterface, error) {
			return newJivaStor(name, config, aspect)
		})

	// Current/New registration style !!
	volume.RegisterVolumeProvisioner(
		// Name when jiva is the persistent volume provisioner plugin
		v1jiva.JivaVolumeProvisionerName,

		// Below is a callback function that creates a new instance of jiva as persistent
		// volume provisioner plugin
		func(name string) (volume.VolumeInterface, error) {
			return newJivaProvisioner(name)
		})
}

// TODO
// This will not be required once Persistent Volume Provisioner profile is
// implemented.
//
// JivaStorNomadAspect is a concrete implementation of following interface:
//
//  1. volume.VolumePluginAspect interface
type JivaStorNomadAspect struct {

	// The aspect that deals with orchestration needs for jiva
	// storage
	Nomad orchprovider.OrchestratorInterface

	// The datacenter which will be the target of API calls.
	// This is useful to set the default value of datacenter for
	// orchprovider.OrchestratorInterface instance.
	Datacenter string
}

func (jAspect *JivaStorNomadAspect) GetOrchProvider() (orchprovider.OrchestratorInterface, error) {

	if jAspect.Nomad == nil {
		return nil, fmt.Errorf("Nomad aspect is not set")
	}

	return jAspect.Nomad, nil
}

func (jAspect *JivaStorNomadAspect) DefaultDatacenter() (string, error) {
	return jAspect.Datacenter, nil
}

// TODO
// Rename to jivaProvisioner ??
//
// jivaStor is the concrete implementation that implements
// following interfaces:
//
//  1. volume.VolumeInterface interface
//  2. volume.Provisioner interface
//  3. volume.Deleter interface
type jivaStor struct {

	// name is the name of this jiva volume plugin.
	name string

	// jivaProUtil is the instance that does all the low level jiva persistent
	// volume provisioner works.
	jivaProUtil JivaInterface

	// TODO
	// Deprecate
	// Will be removed & jivaProUtil will be used instead.
	//
	// jStorOps abstracts the storage operations of this jivaStor
	// instance
	jStorOps StorageOps

	// TODO
	// jConfig provides a handle to tune the operations of
	// this jivaStor instance
	//jConfig *JivaConfig
}

// TODO
// Deprecate
// Remove this deprecated function in favour of newJivaProvisioner
//
// newJivaStor provides a new instance of jivaStor.
//
// This function aligns with VolumePluginFactory function type.
func newJivaStor(name string, config io.Reader, aspect volume.VolumePluginAspect) (*jivaStor, error) {

	glog.Infof("Building new instance of jiva storage '%s'", name)

	// TODO
	//jCfg, err := readJivaConfig(config)
	//if err != nil {
	//	return nil, fmt.Errorf("unable to read Jiva volume provisioner config file: %v", err)
	//}

	// TODO
	// validations of the populated config structure

	jivaUtil, err := newJivaUtil(aspect)
	if err != nil {
		return nil, err
	}

	jStorOps, ok := jivaUtil.StorageOps()
	if !ok {
		return nil, fmt.Errorf("Storage operations not supported by jiva util '%s'", jivaUtil.Name())
	}

	// build the provisioner instance
	jivaStor := &jivaStor{
		name: name,
		//aspect: aspect,
		jStorOps: jStorOps,
		//jConfig:    jCfg,
	}

	return jivaStor, nil
}

// newJivaProvisioner generates a new instance of jiva based persistent volume
// provisioner plugin.
//
// Note:
//    This function aligns with the callback function signature
func newJivaProvisioner(name string) (volume.VolumeInterface, error) {

	glog.Infof("Building new instance of jiva persistent volume provisioner '%s'", name)

	jUtil, err := newJivaProUtil()
	if err != nil {
		return nil, err
	}

	// build the provisioner instance
	jivaStor := &jivaStor{
		name:        name,
		jivaProUtil: jUtil,
	}

	return jivaStor, nil
}

// Name returns the namespaced name of this volume
//
// NOTE:
//    This is a contract implementation of volume.VolumeInterface
func (j *jivaStor) Name() string {
	return j.name
}

// Profile sets the persistent volume provisioner profile against this jiva volume
// provisioner.
func (j *jivaStor) Profile(volProProfile volume.VolumeProvisionerProfile) (bool, error) {
	return j.jivaProUtil.JivaProProfile(volProProfile)
}

// TODO
// Rename to Reader ??
//
// Informer provides a instance of volume.Informer interface.
// Since jivaStor implements volume.Informer, it returns self.
//
// NOTE:
//    This is a contract implementation of volume.VolumeInterface
func (j *jivaStor) Informer() (volume.Informer, bool) {
	return j, true
}

// TODO
// Rename to Creator ??
//
// Provisioner provides a instance of volume.Provisioner interace
// Since jivaStor implements volume.Provisioner, it returns self.
//
// NOTE:
//    This is a concrete implementation of volume.VolumeInterface
func (j *jivaStor) Provisioner() (volume.Provisioner, bool) {
	return j, true
}

// Deleter provides a instance of volume.Deleter interface
// Since jivaStor implements volume.Deleter, it returns self.
//
// NOTE:
//    This is a concrete implementation of volume.VolumeInterface
func (j *jivaStor) Deleter() (volume.Deleter, bool) {
	return j, true
}

// TODO
// Rename to Read ??
//
// Info provides information on a jiva volume
//
// NOTE:
//    This is a concrete implementation of volume.Informer interface
func (j *jivaStor) Info(pvc *v1.PersistentVolumeClaim) (*v1.PersistentVolume, error) {
	// TODO
	// Validations of input i.e. claim

	// Delegate to its provider
	return j.jStorOps.StorageInfo(pvc)
}

// TODO
// Rename to Create ??
//
// Provision provisions a jiva volume
//
// NOTE:
//    This is a concrete implementation of volume.Provisioner interface
func (j *jivaStor) Provision(pvc *v1.PersistentVolumeClaim) (*v1.PersistentVolume, error) {

	// TODO
	// Validations of input i.e. claim

	return j.jStorOps.ProvisionStorage(pvc)
}

// Delete removes a jiva volume
//
// NOTE:
//    This is a concrete implementation of volume.Deleter interface
func (j *jivaStor) Delete(pv *v1.PersistentVolume) (*v1.PersistentVolume, error) {

	// TODO
	// Validations if any

	return j.jStorOps.DeleteStorage(pv)
}

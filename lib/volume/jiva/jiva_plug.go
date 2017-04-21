// This file provides the necessary implementation to establish jiva
// as a mayaserver volume plugin.
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

// The registration logic for jiva storage volume plugin
//
// NOTE:
//    This is invoked at startup.
//
// NOTE:
//    Registration & Initialization are two different workflows. Both are
// mapped by volume plugin name.
//
// TODO
//  A simplified version of registration logic will be implemented. This in turn
// will enable the registry to create new instances of jiva persistent volume
// provisioner on each request.
func init() {
	volume.RegisterVolumePlugin(
		// A variant of jiva volume plugin
		v1jiva.DefaultJivaVolumePluginName,
		// Below is a functional implementation that holds the initialization
		// logic of jiva volume plugin
		func(name string, config io.Reader, aspect volume.VolumePluginAspect) (volume.VolumeInterface, error) {
			return newJivaStor(name, config, aspect)
		})
}

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

// jivaStor is the concrete implementation that implements
// following interfaces:
//
//  1. volume.VolumeInterface interface
//  2. volume.Provisioner interface
//  3. volume.Deleter interface
type jivaStor struct {

	// name is the name of this jiva volume plugin.
	name string

	// jStorOps abstracts the storage operations of this jivaStor
	// instance
	jStorOps StorageOps

	// TODO
	// jConfig provides a handle to tune the operations of
	// this jivaStor instance
	//jConfig *JivaConfig
}

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

// Name returns the namespaced name of this volume
//
// NOTE:
//    This is a contract implementation of volume.VolumeInterface
func (j *jivaStor) Name() string {
	return j.name
}

// Informer provides a instance of volume.Informer interface.
// Since jivaStor implements volume.Informer, it returns self.
//
// NOTE:
//    This is a contract implementation of volume.VolumeInterface
func (j *jivaStor) Informer() (volume.Informer, bool) {
	return j, true
}

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

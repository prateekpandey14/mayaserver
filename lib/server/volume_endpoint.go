package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/openebs/mayaserver/lib/api/v1"
	v1jiva "github.com/openebs/mayaserver/lib/api/v1/jiva"
	"github.com/openebs/mayaserver/lib/volumeprovisioner"
)

func (s *HTTPServer) VolumesRequest(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	switch req.Method {
	case "GET":
		return s.volumeListRequest(resp, req)
	case "PUT", "POST":
		return s.volumeProvision(resp, req, "")
	default:
		return nil, CodedError(405, ErrInvalidMethod)
	}
}

// TODO
// Not yet implemented
func (s *HTTPServer) volumeListRequest(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	return nil, CodedError(405, "Volume list not yet implemented")
}

// VolumeSpecificRequest is a http handler implementation.
// The URL path is parsed to match specific implementations.
//
// TODO
//    Should it return specific types than interface{} ?
func (s *HTTPServer) VolumeSpecificRequest(resp http.ResponseWriter, req *http.Request) (interface{}, error) {

	path := strings.TrimPrefix(req.URL.Path, "/latest/volume")

	// Is req valid ?
	if path == req.URL.Path {
		return nil, CodedError(405, ErrInvalidMethod)
	}

	switch {

	case strings.Contains(path, "/delete/"):
		volName := strings.TrimPrefix(path, "/delete/")
		return s.volumeDelete(resp, req, volName)
	case strings.Contains(path, "/info/"):
		volName := strings.TrimPrefix(path, "/info/")
		return s.volumeInfo(resp, req, volName)
	default:
		return nil, CodedError(405, ErrInvalidMethod)
	}
}

func (s *HTTPServer) volumeProvision(resp http.ResponseWriter, req *http.Request, volName string) (interface{}, error) {

	pvc := v1.PersistentVolumeClaim{}

	// The yaml/json spec is decoded to pvc struct
	if err := decodeBody(req, &pvc); err != nil {
		return nil, CodedError(400, err.Error())
	}

	// Name is expected to be available even in the minimalist specs
	if pvc.Name == "" {
		return nil, CodedError(400, fmt.Sprintf("Volume name hasn't been provided: '%v'", pvc))
	}

	// TODO
	// Get the variant of volume plugin as specified in:
	//
	//  1. http parameters
	//
	// If they have not been specified, then get the variant of volume plugin
	// from:
	//
	//  1. Mayaconfig & JivaConfig.

	// We shall hardcode the variant to jiva default type
	volPlugName := v1jiva.DefaultJivaVolumePluginName

	// Get jiva volume plugin instance which should have been initialized earlier
	jivaStor, err := volumeprovisioner.GetVolumePlugin(volPlugName, nil, nil)

	// Get jiva volume provisioner from the server
	jivaProv, ok := jivaStor.Provisioner()
	if !ok {
		return nil, fmt.Errorf("Volume provisioning not supported by '%s'", volPlugName)
	}

	pv, err := jivaProv.Provision(&pvc)

	if err != nil {
		return nil, err
	}

	return pv, nil
}

func (s *HTTPServer) volumeDelete(resp http.ResponseWriter, req *http.Request, volName string) (interface{}, error) {

	if volName == "" {
		return nil, fmt.Errorf("Volume name missing for deletion")
	}

	// TODO
	// Get the variant of volume plugin as specified in:
	//
	//  1. http parameters
	//
	// If they have not been specified, then get the variant of volume plugin
	// from:
	//
	//  1. Mayaconfig & JivaConfig.

	// We shall hardcode the variant to jiva default type
	volPlugName := v1jiva.DefaultJivaVolumePluginName

	// Get jiva volume plugin instance which should have been initialized earlier
	jivaStor, err := volumeprovisioner.GetVolumePlugin(volPlugName, nil, nil)

	// Get jiva volume deleter
	jivaDel, ok := jivaStor.Deleter()
	if !ok {
		return nil, fmt.Errorf("Deleting volume is not supported by '%s'", volPlugName)
	}

	// Delete a jiva volume
	pv := &v1.PersistentVolume{}
	pv.Name = volName

	dPV, err := jivaDel.Delete(pv)

	if err != nil {
		return nil, err
	}

	return dPV, nil
}

func (s *HTTPServer) volumeInfo(resp http.ResponseWriter, req *http.Request, volName string) (interface{}, error) {

	if volName == "" {
		return nil, fmt.Errorf("Volume name missing")
	}

	// TODO
	// Get the variant of volume plugin as specified in:
	//
	//  1. http parameters
	//
	// If they have not been specified, then get the variant of volume plugin
	// from:
	//
	//  1. Mayaconfig & JivaConfig.

	// We shall hardcode the variant to jiva default type
	volPlugName := v1jiva.DefaultJivaVolumePluginName

	// Get jiva volume plugin instance which should have been initialized earlier
	jivaStor, err := volumeprovisioner.GetVolumePlugin(volPlugName, nil, nil)

	jivaInfo, ok := jivaStor.Informer()
	if !ok {
		return nil, fmt.Errorf("Volume information is not supported by '%s'", volPlugName)
	}

	pvc := &v1.PersistentVolumeClaim{}
	pvc.Name = volName

	info, err := jivaInfo.Info(pvc)

	if err != nil {
		return nil, err
	}

	return info, nil
}

// VSMSpecificRequest is a http handler implementation. It deals with HTTP
// requests w.r.t a single VSM.
//
// TODO
//    Should it return specific types than interface{} ?
func (s *HTTPServer) VSMSpecificRequest(resp http.ResponseWriter, req *http.Request) (interface{}, error) {

	switch req.Method {
	case "PUT", "POST":
		return s.vsmAdd(resp, req)
	case "GET":
		return s.vsmSpecificGetRequest(resp, req)
	default:
		return nil, CodedError(405, ErrInvalidMethod)
	}
}

// vsmSpecificGetRequest deals with HTTP GET request w.r.t a single VSM
func (s *HTTPServer) vsmSpecificGetRequest(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	// Extract info from path after trimming
	path := strings.TrimPrefix(req.URL.Path, "/latest/volumes")

	// Is req valid ?
	if path == req.URL.Path {
		return nil, CodedError(405, ErrInvalidMethod)
	}

	switch {

	case strings.Contains(path, "/info/"):
		vsmName := strings.TrimPrefix(path, "/info/")
		return s.vsmRead(resp, req, vsmName)
	case path == "/":
		return s.vsmList(resp, req)
	default:
		return nil, CodedError(405, ErrInvalidMethod)
	}
}

// vsmList is the http handler that lists VSMs
func (s *HTTPServer) vsmList(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	// Get jiva persistent volume provisioner instance
	jiva, err := volumeprovisioner.GetVolumeProvisioner()
	if err != nil {
		return nil, err
	}

	// Create a PVC
	pvc := &v1.PersistentVolumeClaim{}

	// Set the volume provisioner profile to jiva provisioner
	_, err = jiva.Profile(pvc)
	if err != nil {
		return nil, err
	}

	lister, ok, err := jiva.Lister()
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, fmt.Errorf("VSM list is not supported by '%s:%s'", jiva.Label(), jiva.Name())
	}

	l, err := lister.List()
	if err != nil {
		return nil, err
	}

	return l, nil
}

// vsmRead is the http handler that fetches the details of a VSM
func (s *HTTPServer) vsmRead(resp http.ResponseWriter, req *http.Request, vsmName string) (interface{}, error) {

	if vsmName == "" {
		return nil, fmt.Errorf("VSM name is missing")
	}

	// Get jiva persistent volume provisioner instance
	jiva, err := volumeprovisioner.GetVolumeProvisioner()
	if err != nil {
		return nil, err
	}

	// Create a PVC
	pvc := &v1.PersistentVolumeClaim{}
	pvc.Name = vsmName

	// Set the volume provisioner profile to jiva provisioner
	_, err = jiva.Profile(pvc)
	if err != nil {
		return nil, err
	}

	reader, ok := jiva.Reader()
	if !ok {
		return nil, fmt.Errorf("VSM read is not supported by '%s:%s'", jiva.Label(), jiva.Name())
	}

	// TODO
	// pvc should not be passed again !!
	details, err := reader.Read(pvc)
	if err != nil {
		return nil, err
	}

	return details, nil
}

// vsmAdd is the http handler that fetches the details of a VSM
func (s *HTTPServer) vsmAdd(resp http.ResponseWriter, req *http.Request) (interface{}, error) {

	pvc := v1.PersistentVolumeClaim{}

	// The yaml/json spec is decoded to pvc struct
	if err := decodeBody(req, &pvc); err != nil {
		return nil, CodedError(400, err.Error())
	}

	// Name is expected to be available even in the minimalist specs
	if pvc.Name == "" {
		return nil, CodedError(400, fmt.Sprintf("VSM name missing in '%v'", pvc))
	}

	// Get jiva persistent volume provisioner instance
	jiva, err := volumeprovisioner.GetVolumeProvisioner()
	if err != nil {
		return nil, err
	}

	// Set the volume provisioner profile to jiva provisioner
	_, err = jiva.Profile(&pvc)
	if err != nil {
		return nil, err
	}

	adder, ok := jiva.Adder()
	if !ok {
		return nil, fmt.Errorf("VSM add is not supported by '%s:%s'", jiva.Label(), jiva.Name())
	}

	// TODO
	// pvc should not be passed again !!
	details, err := adder.Add(&pvc)
	if err != nil {
		return nil, err
	}

	return details, nil
}

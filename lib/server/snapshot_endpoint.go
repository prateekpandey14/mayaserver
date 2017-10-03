package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/openebs/maya/types/v1"
	"github.com/openebs/maya/volumes/provisioner/jiva"
)

/*func (s *HTTPServer) SnapshotSpecificRequest(resp http.ResponseWriter, req *http.Request) (interface{}, error) {

	fmt.Println("[DEBUG] Processing", req.Method, "request")

	switch req.Method {
	case "PUT", "POST":
		return s.SnapshotCreate(resp, req)
	case "GET":
		return s.SnapshotSpecificGetRequest(resp, req)
	default:
		return nil, CodedError(405, ErrInvalidMethod)
	}

}*/

// SnapshotSpecificGetRequest deals with HTTP GET request w.r.t a Volume Snapshot
func (s *HTTPServer) SnapshotSpecificRequest(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	// Extract info from path after trimming
	path := strings.TrimPrefix(req.URL.Path, "/latest/snapshot")

	// Is req valid ?
	if path == req.URL.Path {
		fmt.Println("Request coming", path)
		return nil, CodedError(405, ErrInvalidMethod)
	}

	switch {
	case strings.Contains(path, "/create/"):
		return s.SnapshotCreate(resp, req)
	case strings.Contains(path, "/revert/"):
		//	volName := strings.TrimPrefix(path, "/revert/")
		return s.SnapshotRevert(resp, req)
	case path == "/list":
		volName := strings.TrimPrefix(path, "/list/")
		return s.SnapshotList(resp, req, volName)
	default:
		return nil, CodedError(405, ErrInvalidMethod)
	}
}

func (s *HTTPServer) SnapshotCreate(resp http.ResponseWriter, req *http.Request) (interface{}, error) {

	if req.Method != "PUT" && req.Method != "POST" {
		return nil, CodedError(405, ErrInvalidMethod)
	}
	fmt.Println("Request Came")
	snap := v1.VolumeSnapshot{}

	// The yaml/json spec is decoded to pvc struct
	if err := decodeBody(req, &snap); err != nil {

		return nil, CodedError(400, err.Error())
	}

	// SnapshotName is expected to be available even in the minimalist specs
	if snap.Metadata.Name == "" {
		return nil, CodedError(400, fmt.Sprintf("Snapshot name missing in '%v'", snap.Metadata.Name))
	}

	// Name is expected to be available even in the minimalist specs
	if snap.Spec.PersistentVolumeClaimName == "" {

		return nil, CodedError(400, fmt.Sprintf("PVC Volume name missing in '%v'", snap.Spec.PersistentVolumeClaimName))
	}

	fmt.Println("Volume Name :", snap.Spec.PersistentVolumeClaimName)
	fmt.Println("[DEBUG] Processing snapshot-create request of volume:", snap.Spec.PersistentVolumeClaimName)

	//volSpec, err := s.vsmRead(resp, req, volName)
	//if err != nil {
	//	return nil, err
	//}

	//	fmt.Println("Details are :", details)
	//snapName := volName + "snap1"
	var labelMap map[string]string

	//snapSpec := v1.VolumeSnapshot{}

	id, err := jiva.Snapshot(snap.Spec.PersistentVolumeClaimName, snap.Metadata.Name, labelMap)
	if err != nil {
		log.Fatalf("Error running create snapshot command: %v", err)
		return nil, err
	}

	fmt.Println("Created Snapshot is:", id)

	return id, nil
}

func (s *HTTPServer) SnapshotRevert(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	if req.Method != "PUT" && req.Method != "POST" {
		return nil, CodedError(405, ErrInvalidMethod)
	}

	snap := v1.VolumeSnapshot{}

	// The yaml/json spec is decoded to pvc struct
	if err := decodeBody(req, &snap); err != nil {

		return nil, CodedError(400, err.Error())
	}

	// SnapshotName is expected to be available even in the minimalist specs
	if snap.Metadata.Name == "" {
		return nil, CodedError(400, fmt.Sprintf("Snapshot name missing in '%v'", snap.Metadata.Name))
	}

	// Name is expected to be available even in the minimalist specs
	if snap.Spec.PersistentVolumeClaimName == "" {

		return nil, CodedError(400, fmt.Sprintf("Volume name missing in '%v'", snap))
	}

	fmt.Println("Volume Name :", snap.Spec.PersistentVolumeClaimName)
	fmt.Println("[DEBUG] Processing snapshot-revert request of volume:", snap.Spec.PersistentVolumeClaimName)

	err := jiva.SnapshotRevert(snap.Spec.PersistentVolumeClaimName, snap.Metadata.Name)
	if err != nil {
		log.Fatalf("Error running revert snapshot command: %v", err)
		return nil, err
	}

	fmt.Println("Snapshot Reverted to:", snap.SnapshotName)
	return nil, nil

}

func (s *HTTPServer) SnapshotList(resp http.ResponseWriter, req *http.Request, volName string) (interface{}, error) {

	if req.Method != "GET" {
		return nil, CodedError(405, ErrInvalidMethod)
	}

	snap := v1.VolumeSnapshot{}
	snap.Spec.PersistentVolumeClaimName = volName

	// The yaml/json spec is decoded to pvc struct
	if err := decodeBody(req, &snap); err != nil {

		return nil, CodedError(400, err.Error())
	}

	// Name is expected to be available even in the minimalist specs
	if snap.Spec.PersistentVolumeClaimName == "" {

		return nil, CodedError(400, fmt.Sprintf("Volume name missing in '%v'", snap))
	}

	fmt.Println("Volume Name :", snap.Spec.PersistentVolumeClaimName)
	fmt.Println("[DEBUG] Processing snapshot-list request of volume:", snap.Spec.PersistentVolumeClaimName)

	err := jiva.SnapshotList(snap.Spec.PersistentVolumeClaimName)
	if err != nil {
		log.Fatalf("Error running list snapshot command: %v", err)
		return nil, err
	}

	fmt.Println("Snapshot :", snap.SnapshotName)
	return nil, nil

}

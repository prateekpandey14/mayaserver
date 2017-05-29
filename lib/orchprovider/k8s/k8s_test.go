package k8s

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/openebs/mayaserver/lib/api/v1"
	"github.com/openebs/mayaserver/lib/orchprovider"
	volProfile "github.com/openebs/mayaserver/lib/profile/volumeprovisioner"
	//k8sApiMetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	k8sCoreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
	k8sApi "k8s.io/client-go/pkg/api"
	k8sApiv1 "k8s.io/client-go/pkg/api/v1"
	policy "k8s.io/client-go/pkg/apis/policy/v1beta1"
	watch "k8s.io/client-go/pkg/watch"
	"k8s.io/client-go/rest"
)

// TestK8sOrchInterfaceCompliance verifies if k8sOrchestrator implements
// all the exposed methods of the desired interfaces.
//
// NOTE:
//    In case of non-compliance, this logic will error out during compile
// time itself.
func TestK8sOrchInterfaceCompliance(t *testing.T) {
	// k8sOrchestrator implements orchprovider.OrchestratorInterface
	var _ orchprovider.OrchestratorInterface = &k8sOrchestrator{}
	// k8sOrchestrator implements orchprovider.StorageOps
	var _ orchprovider.StorageOps = &k8sOrchestrator{}
	// k8sOrchestrator implements k8s.K8sUtilGetter
	var _ K8sUtilGetter = &k8sOrchestrator{}
}

// TestNewK8sOrchestrator verifies the function that creates a new instance of
// k8sOrchestrator. In addition, it verifies if the returned instance
// provides features it is expected of.
func TestNewK8sOrchestrator(t *testing.T) {
	cases := []struct {
		label string
		name  string
		err   string
	}{
		{"", "", "Label not found while building k8s orchestrator"},
		{"", "non-blank", "Label not found while building k8s orchestrator"},
		{"non-blank", "", "Name not found while building k8s orchestrator"},
		{"non-blank", "non-blank", ""},
		// These are real-world cases of using NewK8sOrchestrator(..) function
		{string(v1.OrchestratorNameLbl), string(v1.K8sOrchestrator), ""},
		{string(v1.OrchestratorNameLbl), string(v1.NomadOrchestrator), ""},
		{string(v1.OrchestratorNameLbl), string(v1.DefaultOrchestrator), ""},
	}

	for i, c := range cases {
		o, err := NewK8sOrchestrator(v1.NameLabel(c.label), v1.OrchProviderRegistry(c.name))

		if err != nil && c.err != err.Error() {
			t.Errorf("TestCase: '%d' ExpectedError: '%s' ActualError: '%s'", i, c.err, err.Error())
		}

		if err == nil && c.label != o.Label() {
			t.Errorf("TestCase: '%d' ExpectedLabel: '%s' ActualLabel: '%s'", i, c.label, o.Label())
		}

		if err == nil && c.name != o.Name() {
			t.Errorf("TestCase: '%d' ExpectedName: '%s' ActualName: '%s'", i, c.name, o.Name())
		}

		// Region is always blank currently in k8sOrchestrator
		if err == nil && "" != o.Region() {
			t.Errorf("TestCase: '%d' ExpectedRegion: '' ActualRegion: '%s'", i, o.Region())
		}

		// Storage Operations is always supported by k8sOrchestrator
		if err == nil {
			if _, supported := o.StorageOps(); !supported {
				t.Errorf("TestCase: '%d' ExpectedStorageOpsSupport: 'true' ActualStorageOpsSupport: '%t'", i, supported)
			}
		}
	}
}

// TestK8sStorageOps will verify the correctness of StorageOps() method of
// k8sOrchestrator
func TestK8sStorageOps(t *testing.T) {
	o, _ := NewK8sOrchestrator(v1.OrchestratorNameLbl, v1.DefaultOrchestrator)

	storOps, supported := o.StorageOps()
	if !supported {
		t.Errorf("ExpectedStorageOpsSupport: 'true' ActualStorageOpsSupport: 'false'")
	}

	if storOps == nil {
		t.Errorf("ExpectedStorageOps: 'non-nil' ActualStorageOps: 'nil'")
	}
}

// TestAddStorage will verify the correctness of AddStorage() method of
// k8sOrchestrator
//
// NOTE:
//    This test case expects the test run environment to NOT have k8s installed
// and hence fail with error.
func TestAddStorage(t *testing.T) {
	o, _ := NewK8sOrchestrator(v1.OrchestratorNameLbl, v1.K8sOrchestrator)

	cases := []struct {
		vsmname string
		err     string
	}{
		{"my-demo-vsm", "unable to load in-cluster configuration, KUBERNETES_SERVICE_HOST and KUBERNETES_SERVICE_PORT must be defined"},
	}

	for _, c := range cases {

		pvc := &v1.PersistentVolumeClaim{}
		pvc.Name = c.vsmname
		//pvc.Labels = map[string]string{
		//	string(v1.PVPVSMNameLbl): c.vsmname,
		//}

		volP, _ := volProfile.GetDefaultVolProProfile(pvc)

		sOps, _ := o.StorageOps()

		_, err := sOps.AddStorage(volP)
		if err != nil && c.err != err.Error() {
			t.Errorf("ExpectedAddStorageErr: '%s' ActualAddStorageErr: '%s'", c.err, err.Error())
		}
	}
}

// k8sVBTLbl represents those types that are used as KEYs for Value Based
// Testing
type k8sVBTLbl string

// These constants represent the Value Based Testing keys
const (
	testK8sUtlNameLbl            k8sVBTLbl = "k8s-utl-name"
	testK8sClientSupportLbl      k8sVBTLbl = "k8s-client-support"
	testK8sNSLbl                 k8sVBTLbl = "k8s-ns"
	testK8sInjectNSErrLbl        k8sVBTLbl = "k8s-inject-ns-err"
	testK8sInClusterLbl          k8sVBTLbl = "k8s-in-cluster"
	testK8sInjectInClusterErrLbl k8sVBTLbl = "k8s-inject-in-cluster-err"
	testK8sInjectPodErrLbl       k8sVBTLbl = "k8s-inject-pod-err"
	testK8sInjectSvcErrLbl       k8sVBTLbl = "k8s-inject-svc-err"
	testK8sInjectVSMLbl          k8sVBTLbl = "k8s-inject-vsm"
	testK8sErrorLbl              k8sVBTLbl = "k8s-err"
)

// mockK8sOrch represents the mock-ed struct of k8sOrchestrator.
//
// This embeds the original k8sOrchestrator to let the execution pass through
// the original code path (most of the times).
//
// NOTE:
//    mock instance(s) is/are injected into k8sOrchestrator's dependency when
// mock based code path is required to be executed.
//
// NOTE:
//    We require execution of mock code paths for unit testing purposes.
type mockK8sOrch struct {
	k8sOrchestrator
}

// StorageOps is the mocked version of the original's i.e. k8sOrchestrator.StorageOps()
func (m *mockK8sOrch) StorageOps() (orchprovider.StorageOps, bool) {
	return m, true
}

// K8sUtil is the mocked version of the original's i.e. k8sOrchestrator.K8sUtil()
func (m *mockK8sOrch) GetK8sUtil(volProfile volProfile.VolumeProvisionerProfile) K8sUtilInterface {

	pvc, _ := volProfile.PVC()

	// mockK8sUtil is instantiated based on a 'Value Based Test' record/row
	return &mockK8sUtil{
		name: pvc.Labels[string(testK8sUtlNameLbl)],
		//vsmName:            pvc.Labels[string(v1.PVPVSMNameLbl)],
		vsmName:            pvc.Name,
		kcSupport:          pvc.Labels[string(testK8sClientSupportLbl)],
		ns:                 pvc.Labels[string(v1.OrchNSLbl)],
		injectNSErr:        pvc.Labels[string(testK8sInjectNSErrLbl)],
		inCluster:          pvc.Labels[string(testK8sInClusterLbl)],
		injectInClusterErr: pvc.Labels[string(testK8sInjectInClusterErrLbl)],
		injectPodErr:       pvc.Labels[string(testK8sInjectPodErrLbl)],
		injectSvcErr:       pvc.Labels[string(testK8sInjectSvcErrLbl)],
		injectVsm:          pvc.Labels[string(testK8sInjectVSMLbl)],
		resultingErr:       pvc.Labels[string(testK8sErrorLbl)],
	}
}

// mockK8sUtil represents the mock-ed struct of k8sUtil & hence provides
// mocked code paths.
type mockK8sUtil struct {
	// name of this instance
	name string
	// name of the mocked VSM
	vsmName string
	// truthy value indicating support for k8s client
	kcSupport string
	// namespace
	ns string
	// injected error for NS() execution
	injectNSErr string
	// truthy value
	inCluster string
	// injected error for InCluster() execution
	injectInClusterErr string
	// injected error for Pods() execution
	injectPodErr string
	// injected error for Services() execution
	injectSvcErr string
	// truthy value
	injectVsm string
	// resultingErr is the error message that is returned
	resultingErr string
}

func (m *mockK8sUtil) Name() string {
	return m.name
}

func (m *mockK8sUtil) K8sClient() (K8sClient, bool) {
	if m.kcSupport == "true" {
		return m, true
	} else {
		return nil, false
	}
}

func (m *mockK8sUtil) InCluster() (bool, error) {
	if m.injectInClusterErr != "" {
		return false, errors.New(m.injectInClusterErr)
	}

	if m.inCluster == "true" {
		return true, nil
	} else {
		return false, nil
	}
}

func (m *mockK8sUtil) NS() (string, error) {
	if m.injectNSErr == "" {
		return m.ns, nil
	} else {
		return m.ns, errors.New(m.injectNSErr)
	}
}

func (m *mockK8sUtil) Pods() (k8sCoreV1.PodInterface, error) {
	if m.injectPodErr == "" {
		return &mockPodOps{
			ns:        m.ns,
			vsmName:   m.vsmName,
			injectVsm: m.injectVsm,
		}, nil
	} else {
		return nil, errors.New(m.injectPodErr)
	}
}

func (m *mockK8sUtil) Services() (k8sCoreV1.ServiceInterface, error) {
	if m.injectSvcErr == "" {
		return &mockSvcOps{}, nil
	} else {
		return nil, errors.New(m.injectSvcErr)
	}
}

// mockPodOps implements k8sCoreV1.PodInterface and hence provides
// necessary mock path
type mockPodOps struct {
	// namespace
	ns string
	// vsmName is the name of the mocked VSM
	vsmName string
	// truthy value
	injectVsm string
}

func (m *mockPodOps) Create(*k8sApiv1.Pod) (*k8sApiv1.Pod, error) {
	return &k8sApiv1.Pod{}, nil
}

func (m *mockPodOps) Update(*k8sApiv1.Pod) (*k8sApiv1.Pod, error) {
	return &k8sApiv1.Pod{}, nil
}

func (m *mockPodOps) UpdateStatus(*k8sApiv1.Pod) (*k8sApiv1.Pod, error) {
	return &k8sApiv1.Pod{}, nil
}

func (m *mockPodOps) Delete(name string, options *k8sApiv1.DeleteOptions) error {
	return nil
}

func (m *mockPodOps) DeleteCollection(options *k8sApiv1.DeleteOptions, listOptions k8sApiv1.ListOptions) error {
	return nil
}

func (m *mockPodOps) Get(name string) (*k8sApiv1.Pod, error) {
	return &k8sApiv1.Pod{}, nil
}

// List presents the mocked logic w.r.t pod list operation
func (m *mockPodOps) List(opts k8sApiv1.ListOptions) (*k8sApiv1.PodList, error) {

	if m.injectVsm == "true" {
		pod := k8sApiv1.Pod{
			ObjectMeta: k8sApiv1.ObjectMeta{
				Name:      m.vsmName,
				Namespace: m.ns,
				Labels: map[string]string{
					"vsm": m.vsmName,
				},
			},
		}

		return &k8sApiv1.PodList{
			Items: []k8sApiv1.Pod{pod},
		}, nil
	}

	return nil, nil
}

func (m *mockPodOps) Watch(opts k8sApiv1.ListOptions) (watch.Interface, error) {
	return nil, nil
}

func (m *mockPodOps) Patch(name string, pt k8sApi.PatchType, data []byte, subresources ...string) (result *k8sApiv1.Pod, err error) {
	return &k8sApiv1.Pod{}, nil
}

func (m *mockPodOps) Bind(binding *k8sApiv1.Binding) error {
	return nil
}

func (m *mockPodOps) Evict(eviction *policy.Eviction) error {
	return nil
}

func (m *mockPodOps) GetLogs(name string, opts *k8sApiv1.PodLogOptions) *rest.Request {
	return &rest.Request{}
}

// mockSvcOps implements k8sCoreV1.ServiceInterface and hence provides
// necessary mock path
type mockSvcOps struct{}

func (m *mockSvcOps) Create(*k8sApiv1.Service) (*k8sApiv1.Service, error) {
	return &k8sApiv1.Service{}, nil
}

func (m *mockSvcOps) Update(*k8sApiv1.Service) (*k8sApiv1.Service, error) {
	return &k8sApiv1.Service{}, nil
}

func (m *mockSvcOps) UpdateStatus(*k8sApiv1.Service) (*k8sApiv1.Service, error) {
	return &k8sApiv1.Service{}, nil
}

func (m *mockSvcOps) Delete(name string, options *k8sApiv1.DeleteOptions) error {
	return nil
}

func (m *mockSvcOps) DeleteCollection(options *k8sApiv1.DeleteOptions, listOptions k8sApiv1.ListOptions) error {
	return nil
}

func (m *mockSvcOps) Get(name string) (*k8sApiv1.Service, error) {
	return &k8sApiv1.Service{}, nil
}

func (m *mockSvcOps) List(opts k8sApiv1.ListOptions) (*k8sApiv1.ServiceList, error) {
	return &k8sApiv1.ServiceList{}, nil
}

func (m *mockSvcOps) Watch(opts k8sApiv1.ListOptions) (watch.Interface, error) {
	return nil, nil
}

func (m *mockSvcOps) Patch(name string, pt k8sApi.PatchType, data []byte, subresources ...string) (result *k8sApiv1.Service, err error) {
	return &k8sApiv1.Service{}, nil
}

func (m *mockSvcOps) ProxyGet(scheme, name, port, path string, params map[string]string) rest.ResponseWrapper {
	return nil
}

// TestAddStorageWithMocks will verify the correctness of AddStorage() method of
// k8sOrchestrator with the help of mock structures.
func TestAddStorageWithMocks(t *testing.T) {

	mockedO := &mockK8sOrch{
		// We are not going by the usual instantiation technique for k8sOrchestrator
		// Below style is to inject our mock
		k8sOrchestrator: k8sOrchestrator{
			label: v1.NameLabel("mock-orch-lbl"),
			name:  v1.OrchProviderRegistry("mock-k8s-orch"),
			// mockK8sOrch is also a k8sUtilInterface implementor
			k8sUtlGtr: &mockK8sOrch{},
		},
	}

	// NOTE:
	//    WATCH OUT: The order of entries are very important here
	cases := []struct {
		kUtlName           string
		vsmName            string
		kcSupport          string
		ns                 string
		injectNSErr        string
		inCluster          string
		injectInClusterErr string
		injectPodErr       string
		injectSvcErr       string
		injectVsm          string
		resultingErr       string
	}{
		{"mock-k8s-util", // kUtlName
			"",        // vsmName
			"true",    // kcSupport truthy
			"mock-ns", // ns
			"",        // injectNSErr Msg
			"",        // inCluster
			"",        // injectInClusterErr Msg
			"",        // injectPodErr Msg
			"",        // injectSvcErr Msg
			"false",   // injectVsm
			"Missing VSM name in 'volumeprovisioner.mapi.openebs.io/profile-name:pvc'", //resultingErr Msg
		},
		{"mock-k8s-util", // kUtlName
			"mock-vsm",     // vsmName
			"true",         // kcSupport truthy
			"mock-ns",      // ns
			"",             // injectNSErr Msg
			"true",         // inCluster truthy
			"",             // injectInClusterErr Msg
			"mock-pod-err", // injectPodErr Msg
			"",             // injectSvcErr Msg
			"false",        // injectVsm
			"mock-pod-err", // resultingErr Msg
		},
		{"mock-k8s-util", // kUtlName
			"mock-vsm",     // vsmName
			"true",         // kcSupport truthy
			"mock-ns",      //ns
			"",             // injectNSErr Msg
			"true",         // inCluster truthy
			"",             // injectInClusterErr Msg
			"",             // injectPodErr Msg
			"mock-svc-err", // injectSvcErr Msg
			"false",        // injectVsm
			"mock-svc-err", // resultingErr Msg
		},
		{"mock-k8s-util", // kUtlName
			"mock-vsm", // vsmName
			"true",     // kcSupport truthy
			"mock-ns",  // ns
			"",         // injectNSErr Msg
			"true",     // inCluster truthy
			"",         // injectInClusterErr Msg
			"",         // injectPodErr Msg
			"",         // injectSvcErr Msg
			"true",     // injectVsm
			"",         // resultingErr Msg
		},
	}

	for i, c := range cases {
		// We will use pvc to implement VALUE BASED TESTING ;)
		//
		// NOTE:
		//    This is just for testing purposes.
		// PVC is never meant to be used in this manner.
		pvc := &v1.PersistentVolumeClaim{}
		pvc.Name = c.vsmName
		pvc.Labels = map[string]string{
			string(testK8sUtlNameLbl): c.kUtlName,
			//string(v1.PVPVSMNameLbl):             c.vsmName,
			string(testK8sClientSupportLbl):      c.kcSupport,
			string(v1.OrchNSLbl):                 c.ns,
			string(testK8sInjectNSErrLbl):        c.injectNSErr,
			string(testK8sInClusterLbl):          c.inCluster,
			string(testK8sInjectInClusterErrLbl): c.injectInClusterErr,
			string(testK8sInjectPodErrLbl):       c.injectPodErr,
			string(testK8sInjectSvcErrLbl):       c.injectSvcErr,
			string(testK8sInjectVSMLbl):          c.injectVsm,
			string(testK8sErrorLbl):              c.resultingErr,
		}

		volP, _ := volProfile.GetDefaultVolProProfile(pvc)

		sOps, supported := mockedO.StorageOps()

		if !supported {
			t.Errorf("TestCase: #%d \n\tExpectedStorageOpsSupport: 'true' \n\tActualStorageOpsSupport: '%t'", i+1, supported)
			continue
		}

		pvList, err := sOps.AddStorage(volP)

		if err != nil && err.Error() != c.resultingErr {
			t.Errorf("TestCase: #%d \n\tExpectedAddStorageErr: '%s' \n\tActualAddStorageErr: '%s'", i+1, c.resultingErr, err.Error())

		} else if c.injectVsm == "true" && pvList == nil {
			t.Errorf("TestCase: #%d \n\tExpectedAddStoragePodList: 'non-nil' \n\tActualAddStoragePodList: 'nil'", i+1)

		} else if c.injectVsm == "true" && len(pvList.Items) == 0 {
			t.Errorf("TestCase: #%d \n\tExpectedAddStoragePodCount: 'non 0' \n\tActualAddStoragePodCount: '0'", i+1)

		} else if c.injectVsm == "true" {
			count := len(pvList.Items)

			for i := 0; i < count; i++ {
				pv := pvList.Items[i]
				if pv.Name != c.vsmName {
					t.Errorf("TestCase: #%d \n\tExpectedAddStoragePodName: '%s' \n\tActualAddStoragePodName: '%s' \n\tGotPod: '%v'", i+1, c.vsmName, pv.Name, pv)
				}
			}
		}
	}
}

// okVsmNameVolumeProfile focusses on NOT returning any error during invocation
// of VSMName() method
type okVsmNameVolumeProfile struct {
	volProfile.VolumeProvisionerProfile
}

// VSMName does not return any error
func (e *okVsmNameVolumeProfile) VSMName() (string, error) {
	return "ok-vsm-name", nil
}

// okCtrlImgVolumeProfile focusses on not returning any error during invocation
// of ControllerImage() method
type okCtrlImgVolumeProfile struct {
	okVsmNameVolumeProfile
}

// ControllerImage does not return any error
func (e *okCtrlImgVolumeProfile) ControllerImage() (string, bool, error) {
	return "ok-ctrl-img", true, nil
}

// errVsmNameVolumeProfile focusses on returning error during invocation of
// VSMName() method
type errVsmNameVolumeProfile struct {
	volProfile.VolumeProvisionerProfile
}

// VSMName returns error
func (e *errVsmNameVolumeProfile) VSMName() (string, error) {
	return "", fmt.Errorf("err-vsm-name")
}

// TestCreateControllerPodReturnsErrVsmName returns error while invoking
// createControllerPod(). This error is due to invocation of VSMName() within
// createControllerPod().
func TestCreateControllerPodReturnsErrVsmName(t *testing.T) {
	mockedO := &mockK8sOrch{
		k8sOrchestrator: k8sOrchestrator{},
	}

	_, err := mockedO.createControllerPod(&errVsmNameVolumeProfile{})
	if err == nil {
		t.Errorf("TestCase: Error Match \n\tExpectedErr: 'not-nil' \n\tActualErr: 'nil'")
	}

	if err != nil && err.Error() != "err-vsm-name" {
		t.Errorf("TestCase: Error Message Match \n\tExpectedErr: 'err-vsm-name' \n\tActualErr: '%s'", err.Error())
	}
}

// errCtrlImgVolumeProfile focusses on returning error during invocation of
// ControllerImage() method
type errCtrlImgVolumeProfile struct {
	volProfile.VolumeProvisionerProfile
}

// ControllerImage returns error
func (e *errCtrlImgVolumeProfile) ControllerImage() (string, bool, error) {
	return "", true, fmt.Errorf("err-ctrl-img")
}

// TestCreateControllerPodReturnsErrCtrlImg returns error while invoking
// createControllerPod(). This error is due to invocation of ControllerImage() within
// createControllerPod().
func TestCreateControllerPodReturnsErrCtrlImg(t *testing.T) {
	mockedO := &mockK8sOrch{
		k8sOrchestrator: k8sOrchestrator{},
	}

	volProfile := &errCtrlImgVolumeProfile{
		&okVsmNameVolumeProfile{},
	}

	_, err := mockedO.createControllerPod(volProfile)
	if err == nil {
		t.Errorf("TestCase: Error Match \n\tExpectedErr: 'not-nil' \n\tActualErr: 'nil'")
	}

	if err != nil && err.Error() != "err-ctrl-img" {
		t.Errorf("TestCase: Error Message Match \n\tExpectedErr: 'err-ctrl-img' \n\tActualErr: '%s'", err.Error())
	}
}

// noSupportCtrlImgVolumeProfile focusses on returning not supported during
// invocation of ControllerImage() method
type noSupportCtrlImgVolumeProfile struct {
	volProfile.VolumeProvisionerProfile
}

// ControllerImage returns not supported
func (e *noSupportCtrlImgVolumeProfile) ControllerImage() (string, bool, error) {
	return "", false, nil
}

// noSupportRepImgVolumeProfile focusses on returning not supported during
// invocation of ReplicaImage() method
type noSupportRepImgVolumeProfile struct {
	volProfile.VolumeProvisionerProfile
}

// ReplicaImage returns not supported
func (e *noSupportRepImgVolumeProfile) ReplicaImage() (string, bool, error) {
	return "", false, nil
}

// errRepImgVolumeProfile returns an error during invocation of ReplicaImage()
// method
type errRepImgVolumeProfile struct {
	volProfile.VolumeProvisionerProfile
}

// ReplicaImage returns an error
func (e *errRepImgVolumeProfile) ReplicaImage() (string, bool, error) {
	return "", false, fmt.Errorf("err-rep-image")
}

// errRepCountVolumeProfile returns an error during invocation of ReplicaCount()
// method
type errRepCountVolumeProfile struct {
	volProfile.VolumeProvisionerProfile
}

// ReplicaImage does not return any error
func (e *errRepCountVolumeProfile) ReplicaImage() (string, bool, error) {
	return "ok-rep-img", true, nil
}

// ReplicaCount returns an error
func (e *errRepCountVolumeProfile) ReplicaCount() (int, error) {
	return 0, fmt.Errorf("err-rep-count")
}

// errPersistentPathCountVolumeProfile returns an error during invocation of
// PersistentPathCount() method
type errPersistentPathCountVolumeProfile struct {
	volProfile.VolumeProvisionerProfile
}

// ReplicaImage does not return any error
func (e *errPersistentPathCountVolumeProfile) ReplicaImage() (string, bool, error) {
	return "ok-rep-img", true, nil
}

// ReplicaCount does not return any error
func (e *errPersistentPathCountVolumeProfile) ReplicaCount() (int, error) {
	return 0, nil
}

// PersistentPathCount returns an error
func (e *errPersistentPathCountVolumeProfile) PersistentPathCount() (int, error) {
	return 0, fmt.Errorf("err-persistent-path-count")
}

// errReplicaCountMatchVolumeProfile returns an error due to mismatch of
// replica count & persistent path count
type errReplicaCountMatchVolumeProfile struct {
	volProfile.VolumeProvisionerProfile
}

// ReplicaImage does not return any error
func (e *errReplicaCountMatchVolumeProfile) ReplicaImage() (string, bool, error) {
	return "ok-rep-img", true, nil
}

// ReplicaCount does not return any error
func (e *errReplicaCountMatchVolumeProfile) ReplicaCount() (int, error) {
	return 0, nil
}

// PersistentPathCount does not return any error
func (e *errReplicaCountMatchVolumeProfile) PersistentPathCount() (int, error) {
	return 1, nil
}

// okCreateReplicaPodVolumeProfile does not return any error
type okCreateReplicaPodVolumeProfile struct {
	volProfile.VolumeProvisionerProfile
}

// PVC does not return any error
func (e *okCreateReplicaPodVolumeProfile) PVC() (*v1.PersistentVolumeClaim, error) {
	pvc := &v1.PersistentVolumeClaim{}
	pvc.Labels = map[string]string{}
	return pvc, nil
}

// PersistentPath does not return any error
func (e *okCreateReplicaPodVolumeProfile) PersistentPath(position int, rCount int) (string, error) {
	return "/tmp/ok-vsm-name/openebs" + string(position), nil
}

// ReplicaImage does not return any error
func (e *okCreateReplicaPodVolumeProfile) ReplicaImage() (string, bool, error) {
	return "ok-rep-img", true, nil
}

// ReplicaCount does not return any error
func (e *okCreateReplicaPodVolumeProfile) ReplicaCount() (int, error) {
	return 2, nil
}

// PersistentPathCount does not return any error
func (e *okCreateReplicaPodVolumeProfile) PersistentPathCount() (int, error) {
	return 2, nil
}

// TestCreateControllerPodReturnsNoSupportCtrlImg returns not supported while
// invoking createControllerPod(). This error is due to invocation of
// ControllerImage() within createControllerPod().
func TestCreateControllerPodReturnsNoSupportCtrlImg(t *testing.T) {
	mockedO := &mockK8sOrch{
		k8sOrchestrator: k8sOrchestrator{},
	}

	volProfile := &noSupportCtrlImgVolumeProfile{
		&okVsmNameVolumeProfile{},
	}

	_, err := mockedO.createControllerPod(volProfile)
	if err == nil {
		t.Errorf("TestCase: Error Match \n\tExpectedErr: 'not-nil' \n\tActualErr: 'nil'")
	}

	n, _ := volProfile.VSMName()
	expErr := fmt.Sprintf("VSM '%s' requires a controller container image", n)

	if err != nil && err.Error() != expErr {
		t.Errorf("TestCase: Error Message Match \n\tExpectedErr: '%s' \n\tActualErr: '%s'", expErr, err.Error())
	}
}

// noK8sClientSupportK8sOrch is a k8s orchestrator that does not provide support
// to K8sClient
type noK8sClientSupportK8sOrch struct {
	k8sOrchestrator
}

// K8sUtil provides a k8sUtil instance that does not support K8sClient
func (m *noK8sClientSupportK8sOrch) GetK8sUtil(volProfile volProfile.VolumeProvisionerProfile) K8sUtilInterface {
	return &noK8sClientSupportK8sUtil{}
}

// noK8sClientSupportK8sUtil is a k8s util that does not provide support
// K8sClient
type noK8sClientSupportK8sUtil struct {
	K8sUtilInterface
}

// Name returns the name of noK8sClientSupportK8sUtil
func (m *noK8sClientSupportK8sUtil) Name() string {
	return "no-k8s-client-support-k8s-util"
}

// K8sClient does not support K8sClient
func (m *noK8sClientSupportK8sUtil) K8sClient() (K8sClient, bool) {
	return nil, false
}

// TestCreateControllerPodReturnsNoK8sClientSupport returns K8sClient not
// supported while invoking createControllerPod(). This error is due to
// invocation of K8sClient() within createControllerPod().
func TestCreateControllerPodReturnsNoK8sClientSupport(t *testing.T) {
	mockedO := &noK8sClientSupportK8sOrch{
		k8sOrchestrator: k8sOrchestrator{
			k8sUtlGtr: &noK8sClientSupportK8sOrch{},
		},
	}

	volProfile := &okCtrlImgVolumeProfile{}

	_, err := mockedO.createControllerPod(volProfile)
	if err == nil {
		t.Errorf("TestCase: Error Match \n\tExpectedErr: 'not-nil' \n\tActualErr: 'nil'")
	}

	expErr := fmt.Sprintf("K8s client not supported by '%s'", "no-k8s-client-support-k8s-util")

	if err != nil && err.Error() != expErr {
		t.Errorf("TestCase: Error Message Match \n\tExpectedErr: '%s' \n\tActualErr: '%s'", expErr, err.Error())
	}
}

// errNSK8sClientK8sOrch is a k8s orchestrator that returns
// errNSK8sClientK8sUtil
type errNSK8sClientK8sOrch struct {
	k8sOrchestrator
}

// K8sUtil returns errNSK8sClientK8sUtil
func (m *errNSK8sClientK8sOrch) GetK8sUtil(volProfile volProfile.VolumeProvisionerProfile) K8sUtilInterface {
	return &errNSK8sClientK8sUtil{}
}

// errNSK8sClientK8sUtil is a k8sUtil that provides errNSK8sClient
type errNSK8sClientK8sUtil struct {
	K8sUtilInterface
}

// Name returns the name of errNSK8sClientK8sUtil
func (m *errNSK8sClientK8sUtil) Name() string {
	return "err-ns-k8s-client-k8s-util"
}

// K8sClient returns an instance of errNSK8sClient
func (m *errNSK8sClientK8sUtil) K8sClient() (K8sClient, bool) {
	return &errNSK8sClient{}, true
}

// errNSK8sClient is a K8sClient that returns error during NS() invocation
type errNSK8sClient struct {
	K8sClient
}

// NS returns an error
func (e *errNSK8sClient) NS() (string, error) {
	return "", fmt.Errorf("err-ns")
}

// errPodOpsK8sClientK8sOrch is a k8s orchestrator that returns
// errPodOpsK8sClientK8sUtil
type errPodOpsK8sClientK8sOrch struct {
	k8sOrchestrator
}

// K8sUtil returns errPodOpsK8sClientK8sUtil
func (m *errPodOpsK8sClientK8sOrch) GetK8sUtil(volProfile volProfile.VolumeProvisionerProfile) K8sUtilInterface {
	return &errPodOpsK8sClientK8sUtil{}
}

// errPodOpsK8sClientK8sUtil is a k8sUtil that provides errPodOpsK8sClient
type errPodOpsK8sClientK8sUtil struct {
	K8sUtilInterface
}

// Name returns the name of errPodOpsK8sClientK8sUtil
func (m *errPodOpsK8sClientK8sUtil) Name() string {
	return "err-pod-ops-k8s-client-k8s-util"
}

// K8sClient does not support K8sClient
func (m *errPodOpsK8sClientK8sUtil) K8sClient() (K8sClient, bool) {
	return &errPodOpsK8sClient{}, true
}

// errPodOpsK8sClient is a K8sClient that returns error during Pods() invocation
type errPodOpsK8sClient struct {
	K8sClient
}

// NS will not return any error
func (e *errPodOpsK8sClient) NS() (string, error) {
	return "ok-ns", nil
}

// Pods returns an error
func (e *errPodOpsK8sClient) Pods() (k8sCoreV1.PodInterface, error) {
	return nil, fmt.Errorf("err-pod-ops")
}

// errSvcGetK8sOrch is a k8s orchestrator that returns
// errSvcGetK8sUtil
type errSvcGetK8sOrch struct {
	k8sOrchestrator
}

// K8sUtil returns errSvcGetK8sUtil
func (m *errSvcGetK8sOrch) GetK8sUtil(volProfile volProfile.VolumeProvisionerProfile) K8sUtilInterface {
	return &errSvcGetK8sUtil{}
}

// errSvcGetK8sUtil is a k8sUtil that provides errSvcGetK8sClient
type errSvcGetK8sUtil struct {
	K8sUtilInterface
}

// Name returns the name of errPodListPodOpsK8sUtil
func (m *errSvcGetK8sUtil) Name() string {
	return "err-svc-get-k8s-util"
}

// K8sClient returns an instance of errSvcGetK8sClient
func (m *errSvcGetK8sUtil) K8sClient() (K8sClient, bool) {
	return &errSvcGetK8sClient{}, true
}

// errSvcGetK8sClient is a K8sClient that returns errSvcGetPodOps
type errSvcGetK8sClient struct {
	K8sClient
}

// NS will not return any error
func (e *errSvcGetK8sClient) NS() (string, error) {
	return "ok-ns", nil
}

// Pods returns an instance of errSvcGetSvcOps
func (e *errSvcGetK8sClient) Services() (k8sCoreV1.ServiceInterface, error) {
	return &errSvcGetSvcOps{}, nil
}

// errSvcGetSvcOps is a k8sCoreV1.ServiceInterface that returns error during
// Get() invocation
type errSvcGetSvcOps struct {
	k8sCoreV1.ServiceInterface
}

// Get returns an error
func (m *errSvcGetSvcOps) Get(svc string) (*k8sApiv1.Service, error) {
	return nil, fmt.Errorf("err-svc-get")
}

// okGetServiceK8sOrch is a k8s orchestrator that returns
// okSvcGetK8sUtil
type okGetServiceK8sOrch struct {
	k8sOrchestrator
}

// K8sUtil returns an instance of okSvcGetK8sUtil
func (m *okGetServiceK8sOrch) GetK8sUtil(volProfile volProfile.VolumeProvisionerProfile) K8sUtilInterface {
	return &okSvcGetK8sUtil{}
}

// okSvcGetK8sUtil is a k8sUtil that provides errSvcGetK8sClient
type okSvcGetK8sUtil struct {
	K8sUtilInterface
}

// Name returns the name of okSvcGetK8sUtil
func (m *okSvcGetK8sUtil) Name() string {
	return "ok-svc-get-k8s-util"
}

// K8sClient returns an instance of okSvcGetK8sClient
func (m *okSvcGetK8sUtil) K8sClient() (K8sClient, bool) {
	return &okSvcGetK8sClient{}, true
}

// okSvcGetK8sClient is a K8sClient that returns okSvcGetPodOps
type okSvcGetK8sClient struct {
	K8sClient
}

// NS will not return any error
func (e *okSvcGetK8sClient) NS() (string, error) {
	return "ok-ns", nil
}

// Services returns an instance of okSvcGetSvcOps
func (e *okSvcGetK8sClient) Services() (k8sCoreV1.ServiceInterface, error) {
	return &okSvcGetSvcOps{}, nil
}

// okSvcGetSvcOps is a k8sCoreV1.ServiceInterface that does not
// return error during Get() invocation
type okSvcGetSvcOps struct {
	k8sCoreV1.ServiceInterface
}

// Get returns service that it receives without any error
func (m *okSvcGetSvcOps) Get(svc string) (*k8sApiv1.Service, error) {
	s := &k8sApiv1.Service{
		ObjectMeta: k8sApiv1.ObjectMeta{
			Name: "ok-svc-name",
		},
		Spec: k8sApiv1.ServiceSpec{
			ClusterIP: "1.1.1.1",
		},
	}

	return s, nil
}

// errCreateReplicaPodK8sOrch is a k8s orchestrator that returns
// errCreateReplicaPodK8sUtil
type errCreateReplicaPodK8sOrch struct {
	k8sOrchestrator
}

// GetK8sUtil returns an instance of errCreateReplicaPodK8sUtil
func (m *errCreateReplicaPodK8sOrch) GetK8sUtil(volProfile volProfile.VolumeProvisionerProfile) K8sUtilInterface {
	return &errCreateReplicaPodK8sUtil{}
}

// errCreateReplicaPodK8sUtil is a k8sUtil that provides errSvcGetK8sClient
type errCreateReplicaPodK8sUtil struct {
	K8sUtilInterface
}

// Name returns the name of errCreateReplicaPodK8sUtil
func (m *errCreateReplicaPodK8sUtil) Name() string {
	return "err-create-rep-pod-k8s-util"
}

// K8sClient does not support K8sClient
func (m *errCreateReplicaPodK8sUtil) K8sClient() (K8sClient, bool) {
	return nil, false
}

// errPodListPodOpsK8sOrch is a k8s orchestrator that returns
// errPodListPodOpsK8sUtil
type errPodListPodOpsK8sOrch struct {
	k8sOrchestrator
}

// K8sUtil returns errPodListPodOpsK8sUtil
func (m *errPodListPodOpsK8sOrch) GetK8sUtil(volProfile volProfile.VolumeProvisionerProfile) K8sUtilInterface {
	return &errPodListPodOpsK8sUtil{}
}

// errPodListPodOpsK8sUtil is a k8sUtil that provides errPodListPodOpsK8sClient
type errPodListPodOpsK8sUtil struct {
	K8sUtilInterface
}

// Name returns the name of errPodListPodOpsK8sUtil
func (m *errPodListPodOpsK8sUtil) Name() string {
	return "err-pod-list-pod-ops-k8s-util"
}

// K8sClient does not support K8sClient
func (m *errPodListPodOpsK8sUtil) K8sClient() (K8sClient, bool) {
	return &errPodListPodOpsK8sClient{}, true
}

// errPodListPodOpsK8sClient is a K8sClient that returns errPodListPodOps
type errPodListPodOpsK8sClient struct {
	K8sClient
}

// NS will not return any error
func (e *errPodListPodOpsK8sClient) NS() (string, error) {
	return "ok-ns", nil
}

// Pods returns an instance of errPodListPodOps
func (e *errPodListPodOpsK8sClient) Pods() (k8sCoreV1.PodInterface, error) {
	return &errPodListPodOps{}, nil
}

// errSvcOpsK8sOrch is a k8s orchestrator that returns
// errSvcOpsK8sUtil
type errSvcOpsK8sOrch struct {
	k8sOrchestrator
}

// K8sUtil returns an instance of errSvcOpsK8sUtil
func (m *errSvcOpsK8sOrch) GetK8sUtil(volProfile volProfile.VolumeProvisionerProfile) K8sUtilInterface {
	return &errSvcOpsK8sUtil{}
}

// errSvcOpsK8sUtil is a k8sUtil that provides errSvcOpsK8sClient
type errSvcOpsK8sUtil struct {
	K8sUtilInterface
}

// Name returns the name of errSvcOpsK8sUtil
func (m *errSvcOpsK8sUtil) Name() string {
	return "err-svc-ops-k8s-util"
}

// K8sClient returns an instance of errSvcOpsK8sClient
func (m *errSvcOpsK8sUtil) K8sClient() (K8sClient, bool) {
	return &errSvcOpsK8sClient{}, true
}

// errSvcOpsK8sClient is a K8sClient that returns error during Services()
// invocation
type errSvcOpsK8sClient struct {
	K8sClient
}

// NS will not return any error
func (e *errSvcOpsK8sClient) NS() (string, error) {
	return "ok-ns", nil
}

// Services returns error
func (e *errSvcOpsK8sClient) Services() (k8sCoreV1.ServiceInterface, error) {
	return nil, fmt.Errorf("err-svc-ops")
}

// errPodListPodOps is an instance of k8sCoreV1.PodInterface that returns error
// during List invocation
type errPodListPodOps struct {
	k8sCoreV1.PodInterface
}

// List retuns an error
func (e *errPodListPodOps) List(opts k8sApiv1.ListOptions) (*k8sApiv1.PodList, error) {
	return nil, fmt.Errorf("err-pod-list")
}

// errPodListMissK8sOrch is a k8s orchestrator that returns
// errPodListMissK8sUtil
type errPodListMissK8sOrch struct {
	k8sOrchestrator
}

// K8sUtil returns errPodListMissK8sUtil
func (m *errPodListMissK8sOrch) GetK8sUtil(volProfile volProfile.VolumeProvisionerProfile) K8sUtilInterface {
	return &errPodListMissK8sUtil{}
}

// errPodListPodOpsK8sUtil is a k8sUtil that provides errPodListMissK8sClient
type errPodListMissK8sUtil struct {
	K8sUtilInterface
}

// Name returns the name of errPodListMissK8sUtil
func (m *errPodListMissK8sUtil) Name() string {
	return "err-pod-list-miss-k8s-util"
}

// K8sClient returns an instance of errPodListMissK8sClient
func (m *errPodListMissK8sUtil) K8sClient() (K8sClient, bool) {
	return &errPodListMissK8sClient{}, true
}

// errPodListMissK8sClient is a K8sClient that returns errPodListMissPodOps
type errPodListMissK8sClient struct {
	K8sClient
}

// NS will not return any error
func (e *errPodListMissK8sClient) NS() (string, error) {
	return "ok-ns", nil
}

// Pods returns an instance of errPodListMissPodOps
func (e *errPodListMissK8sClient) Pods() (k8sCoreV1.PodInterface, error) {
	return &errPodListMissPodOps{}, nil
}

// errPodListMissPodOps is an instance of k8sCoreV1.PodInterface that returns
// a list of pods which are not expected during List invocation
type errPodListMissPodOps struct {
	k8sCoreV1.PodInterface
}

// List retuns a list of pods which are not expected
func (e *errPodListMissPodOps) List(opts k8sApiv1.ListOptions) (*k8sApiv1.PodList, error) {
	p := k8sApiv1.Pod{}
	p.Name = "err-pod-list"
	p.Labels = map[string]string{
		"err-key": "err-val",
	}

	l := &k8sApiv1.PodList{
		Items: []k8sApiv1.Pod{p},
	}

	return l, nil
}

// errPodListNilK8sOrch is a k8s orchestrator that returns
// errPodListNilK8sUtil
type errPodListNilK8sOrch struct {
	k8sOrchestrator
}

// K8sUtil returns errPodListNilK8sUtil
func (m *errPodListNilK8sOrch) GetK8sUtil(volProfile volProfile.VolumeProvisionerProfile) K8sUtilInterface {
	return &errPodListNilK8sUtil{}
}

// errPodListNilK8sUtil is a k8sUtil that provides errPodListNilK8sClient
type errPodListNilK8sUtil struct {
	K8sUtilInterface
}

// Name returns the name of errPodListNilK8sUtil
func (m *errPodListNilK8sUtil) Name() string {
	return "err-pod-list-nil-k8s-util"
}

// K8sClient returns an instance of errPodListNilK8sClient
func (m *errPodListNilK8sUtil) K8sClient() (K8sClient, bool) {
	return &errPodListNilK8sClient{}, true
}

// errPodListNilK8sClient is a K8sClient that returns errPodListNilPodOps
type errPodListNilK8sClient struct {
	K8sClient
}

// NS will not return any error
func (e *errPodListNilK8sClient) NS() (string, error) {
	return "ok-ns", nil
}

// Pods returns an instance of errPodListMissPodOps
func (e *errPodListNilK8sClient) Pods() (k8sCoreV1.PodInterface, error) {
	return &errPodListNilPodOps{}, nil
}

// errPodListNilPodOps is an instance of k8sCoreV1.PodInterface that returns
// nil during List invocation
type errPodListNilPodOps struct {
	k8sCoreV1.PodInterface
}

// List retuns nil
func (e *errPodListNilPodOps) List(opts k8sApiv1.ListOptions) (*k8sApiv1.PodList, error) {
	return nil, nil
}

// okReadStorageK8sOrch is a k8s orchestrator that returns
// okReadStorageK8sUtil
type okReadStorageK8sOrch struct {
	k8sOrchestrator
}

// K8sUtil returns okReadStorageK8sUtil
func (m *okReadStorageK8sOrch) GetK8sUtil(volProfile volProfile.VolumeProvisionerProfile) K8sUtilInterface {
	return &okReadStorageK8sUtil{}
}

// okReadStorageK8sUtil is a k8sUtil that provides okReadStorageK8sClient
type okReadStorageK8sUtil struct {
	K8sUtilInterface
}

// Name returns the name of okReadStorageK8sUtil
func (m *okReadStorageK8sUtil) Name() string {
	return "ok-read-storage-k8s-util"
}

// K8sClient returns an instance of okReadStorageK8sClient
func (m *okReadStorageK8sUtil) K8sClient() (K8sClient, bool) {
	return &okReadStorageK8sClient{}, true
}

// okReadStorageK8sClient is a K8sClient that returns okReadStoragePodOps
type okReadStorageK8sClient struct {
	K8sClient
}

// NS will not return any error
func (e *okReadStorageK8sClient) NS() (string, error) {
	return "ok-ns", nil
}

// Pods returns an instance of okReadStoragePodOps
func (e *okReadStorageK8sClient) Pods() (k8sCoreV1.PodInterface, error) {
	return &okReadStoragePodOps{}, nil
}

// okReadStoragePodOps is an instance of k8sCoreV1.PodInterface that returns
// a list of pods that are expected during List invocation
type okReadStoragePodOps struct {
	k8sCoreV1.PodInterface
}

// List retuns a list of expected pods
func (e *okReadStoragePodOps) List(opts k8sApiv1.ListOptions) (*k8sApiv1.PodList, error) {
	p := k8sApiv1.Pod{}
	p.Name = "ok-vsm-name"
	p.Labels = map[string]string{
		string(v1.VSMIdentifier): "ok-vsm-name",
	}

	l := &k8sApiv1.PodList{
		Items: []k8sApiv1.Pod{p},
	}

	return l, nil
}

// TestCreateControllerPodReturnsErrPodOps returns pod operator error
// while invoking createControllerPod(). This error is due to
// invocation of Pods() within createControllerPod().
func TestCreateControllerPodReturnsErrPodOps(t *testing.T) {
	mockedO := &errPodOpsK8sClientK8sOrch{
		k8sOrchestrator: k8sOrchestrator{
			k8sUtlGtr: &errPodOpsK8sClientK8sOrch{},
		},
	}

	volProfile := &okCtrlImgVolumeProfile{}

	_, err := mockedO.createControllerPod(volProfile)
	if err == nil {
		t.Errorf("TestCase: Error Match \n\tExpectedErr: 'not-nil' \n\tActualErr: 'nil'")
	}

	expErr := "err-pod-ops"

	if err != nil && err.Error() != expErr {
		t.Errorf("TestCase: Error Message Match \n\tExpectedErr: '%s' \n\tActualErr: '%s'", expErr, err.Error())
	}
}

// okCreatePodK8sOrch is a k8s orchestrator that returns an instance of
// okCreatePodK8sUtil
type okCreatePodK8sOrch struct {
	k8sOrchestrator
}

// GetK8sUtil returns an instance of okCreatePodK8sUtil
func (m *okCreatePodK8sOrch) GetK8sUtil(volProfile volProfile.VolumeProvisionerProfile) K8sUtilInterface {
	return &okCreatePodK8sUtil{}
}

// okCreatePodK8sUtil is a k8sUtil that returns an instance of
// okCreatePodK8sClient
type okCreatePodK8sUtil struct {
	K8sUtilInterface
}

// Name returns the name of okCreatePodK8sUtil
func (m *okCreatePodK8sUtil) Name() string {
	return "ok-create-pod-k8s-util"
}

// K8sClient returns an instance of okCreatePodK8sClient
func (m *okCreatePodK8sUtil) K8sClient() (K8sClient, bool) {
	return &okCreatePodK8sClient{}, true
}

// okCreatePodK8sClient is a K8sClient that returns okCreatePodPodOps
type okCreatePodK8sClient struct {
	K8sClient
}

// Pods returns an instance of okCreatePodPodOps
func (e *okCreatePodK8sClient) Pods() (k8sCoreV1.PodInterface, error) {
	return &okCreatePodPodOps{}, nil
}

// okCreatePodPodOps is a k8sCoreV1.PodInterface that does not return error
// during create()
type okCreatePodPodOps struct {
	k8sCoreV1.PodInterface
}

// Create returns the pod that it receives without any error
func (m *okCreatePodPodOps) Create(pod *k8sApiv1.Pod) (*k8sApiv1.Pod, error) {
	return pod, nil
}

// okCreateServiceK8sOrch is a k8s orchestrator that returns an
// instance of okCreateServiceK8sUtil
type okCreateServiceK8sOrch struct {
	k8sOrchestrator
}

// GetK8sUtil returns an instance of okCreateServiceK8sUtil
func (m *okCreateServiceK8sOrch) GetK8sUtil(volProfile volProfile.VolumeProvisionerProfile) K8sUtilInterface {
	return &okCreateServiceK8sUtil{}
}

// okCreateServiceK8sUtil is a k8sUtil that returns an instance of
// okCreateServiceK8sClient
type okCreateServiceK8sUtil struct {
	K8sUtilInterface
}

// Name returns the name of okCreateServiceK8sUtil
func (m *okCreateServiceK8sUtil) Name() string {
	return "ok-create-svc-k8s-util"
}

// K8sClient returns an instance of okCreateServiceK8sClient
func (m *okCreateServiceK8sUtil) K8sClient() (K8sClient, bool) {
	return &okCreateServiceK8sClient{}, true
}

// okCreateServiceK8sClient is a K8sClient that returns
// okCreateServiceSvcOps
type okCreateServiceK8sClient struct {
	K8sClient
}

// Services returns an instance of okCreateServiceSvcOps
func (e *okCreateServiceK8sClient) Services() (k8sCoreV1.ServiceInterface, error) {
	return &okCreateServiceSvcOps{}, nil
}

// okCreateServiceSvcOps is a k8sCoreV1.ServiceInterface that does not
// return error during Create() invocation
type okCreateServiceSvcOps struct {
	k8sCoreV1.ServiceInterface
}

// Create returns service that it receives without any error
func (m *okCreateServiceSvcOps) Create(svc *k8sApiv1.Service) (*k8sApiv1.Service, error) {
	return svc, nil
}

// TestCreateControllerPodReturnsOk returns a valid pod while invoking
// createControllerPod().
func TestCreateControllerPodReturnsOk(t *testing.T) {
	mockedO := &okCreatePodK8sOrch{
		k8sOrchestrator: k8sOrchestrator{
			k8sUtlGtr: &okCreatePodK8sOrch{},
		},
	}

	volProfile := &okCtrlImgVolumeProfile{}

	cp, err := mockedO.createControllerPod(volProfile)

	if err != nil {
		t.Errorf("TestCase: Nil Error Match \n\tExpectedErr: 'nil' \n\tActualErr: '%s'", err.Error())
	}

	// Verify the pod kind
	if cp.Kind != string(v1.K8sKindDeployment) {
		t.Errorf("TestCase: Kind Match \n\tExpectedKind: '%s' \n\tActualKind: '%s'", v1.K8sKindDeployment, cp.Kind)
	}

	// Verify the pod version
	if cp.APIVersion != string(v1.K8sPodVersion) {
		t.Errorf("TestCase: Pod Version Match \n\tExpectedAPIVersion: '%s' \n\tActualAPIVersion: '%s'", v1.K8sPodVersion, cp.APIVersion)
	}

	// Verify the pod name
	vsm, _ := volProfile.VSMName()
	ePodName := vsm + string(v1.ControllerSuffix)
	if cp.Name != ePodName {
		t.Errorf("TestCase: Pod Name Match \n\tExpectedName: '%s' \n\tActualName: '%s'", ePodName, cp.Name)
	}

	// Verify the pod labels
	eLblStr := string(v1.VSMSelectorPrefix) + vsm
	eLbl, _ := labels.Parse(eLblStr)
	if !eLbl.Matches(labels.Set(cp.Labels)) {
		t.Errorf("TestCase: Labels Match \n\tExpectedLabels: '%s' \n\tActualLabels: '%s'", eLbl, labels.Set(cp.Labels))
	}

	// Verify no. of containers within the pod
	if len(cp.Spec.Containers) != 1 {
		t.Errorf("TestCase: No. of Containers \n\tExpectedContainers: '1' \n\tActualContainers: '%d'", len(cp.Spec.Containers))
	}

	// Verify no. of ports within the container
	if len(cp.Spec.Containers[0].Ports) != 2 {
		t.Errorf("TestCase: No. of Ports \n\tExpectedPorts: '2' \n\tActualPorts: '%d'", len(cp.Spec.Containers[0].Ports))
	}

	// Verify container name
	if cp.Spec.Containers[0].Name != vsm+string(v1.ControllerSuffix)+string(v1.ContainerSuffix) {
		t.Errorf("TestCase: Container Name \n\tExpectedName: '%s' \n\tActualName: '%s'", vsm+string(v1.ControllerSuffix)+string(v1.ContainerSuffix), cp.Spec.Containers[0].Name)
	}

	// Verify controller launch command
	if !reflect.DeepEqual(cp.Spec.Containers[0].Command, v1.JivaCtrlCmd) {
		t.Errorf("TestCase: Controller Launch Cmd \n\tExpectedCmd: '%s' \n\tActualCmd: '%s'", v1.JivaCtrlCmd, cp.Spec.Containers[0].Command)
	}

	// Verify controller launch arguments
	if !reflect.DeepEqual(cp.Spec.Containers[0].Args, v1.JivaCtrlArgs) {
		t.Errorf("TestCase: Controller Launch Args \n\tExpectedArgs: '%s' \n\tActualArgs: '%s'", v1.JivaCtrlArgs, cp.Spec.Containers[0].Args)
	}
}

// TestReadStorageReturnsErrVsmName verifies the vsm name error
func TestReadStorageReturnsErrVsmName(t *testing.T) {
	mockedO := &mockK8sOrch{
		k8sOrchestrator: k8sOrchestrator{},
	}

	_, err := mockedO.ReadStorage(&errVsmNameVolumeProfile{})
	if err == nil {
		t.Errorf("TestCase: Error Match \n\tExpectedErr: 'not-nil' \n\tActualErr: 'nil'")
	}

	if err != nil && err.Error() != "err-vsm-name" {
		t.Errorf("TestCase: Error Message Match \n\tExpectedErr: 'err-vsm-name' \n\tActualErr: '%s'", err.Error())
	}
}

// TestReadStorageReturnsNoK8sClientSupport verifies no K8sClient support error
// error
func TestReadStorageReturnsNoK8sClientSupport(t *testing.T) {
	mockedO := &noK8sClientSupportK8sOrch{
		k8sOrchestrator: k8sOrchestrator{
			k8sUtlGtr: &noK8sClientSupportK8sOrch{},
		},
	}

	volProfile := &okVsmNameVolumeProfile{}

	_, err := mockedO.ReadStorage(volProfile)
	if err == nil {
		t.Errorf("TestCase: Error Match \n\tExpectedErr: 'not-nil' \n\tActualErr: 'nil'")
	}

	expErr := fmt.Sprintf("K8s client not supported by '%s'", "no-k8s-client-support-k8s-util")

	if err != nil && err.Error() != expErr {
		t.Errorf("TestCase: Error Message Match \n\tExpectedErr: '%s' \n\tActualErr: '%s'", expErr, err.Error())
	}
}

// TestReadStorageReturnsErrNS verifies the namespace error
func TestReadStorageReturnsErrNS(t *testing.T) {
	mockedO := &errNSK8sClientK8sOrch{
		k8sOrchestrator: k8sOrchestrator{
			k8sUtlGtr: &errNSK8sClientK8sOrch{},
		},
	}

	volProfile := &okCtrlImgVolumeProfile{}

	_, err := mockedO.ReadStorage(volProfile)
	if err == nil {
		t.Errorf("TestCase: Error Match \n\tExpectedErr: 'not-nil' \n\tActualErr: 'nil'")
	}

	expErr := "err-ns"

	if err != nil && err.Error() != expErr {
		t.Errorf("TestCase: Error Message Match \n\tExpectedErr: '%s' \n\tActualErr: '%s'", expErr, err.Error())
	}
}

// TestReadStorageReturnsErrPodOps verifies the pod operator error
func TestReadStorageReturnsErrPodOps(t *testing.T) {
	mockedO := &errPodOpsK8sClientK8sOrch{
		k8sOrchestrator: k8sOrchestrator{
			k8sUtlGtr: &errPodOpsK8sClientK8sOrch{},
		},
	}

	volProfile := &okCtrlImgVolumeProfile{}

	_, err := mockedO.ReadStorage(volProfile)
	if err == nil {
		t.Errorf("TestCase: Error Match \n\tExpectedErr: 'not-nil' \n\tActualErr: 'nil'")
	}

	expErr := "err-pod-ops"

	if err != nil && err.Error() != expErr {
		t.Errorf("TestCase: Error Message Match \n\tExpectedErr: '%s' \n\tActualErr: '%s'", expErr, err.Error())
	}
}

// TestReadStorageReturnsErrPodList verifies the error during listing of pod
func TestReadStorageReturnsErrPodList(t *testing.T) {
	mockedO := &errPodListPodOpsK8sOrch{
		k8sOrchestrator: k8sOrchestrator{
			k8sUtlGtr: &errPodListPodOpsK8sOrch{},
		},
	}

	volProfile := &okCtrlImgVolumeProfile{}

	_, err := mockedO.ReadStorage(volProfile)
	if err == nil {
		t.Errorf("TestCase: Error Match \n\tExpectedErr: 'not-nil' \n\tActualErr: 'nil'")
	}

	expErr := "err-pod-list"

	if err != nil && err.Error() != expErr {
		t.Errorf("TestCase: Error Message Match \n\tExpectedErr: '%s' \n\tActualErr: '%s'", expErr, err.Error())
	}
}

// TestReadStorageReturnsErrPodListNil verifies the nil returned during listing
// of pod
func TestReadStorageReturnsErrPodListNil(t *testing.T) {
	mockedO := &errPodListNilK8sOrch{
		k8sOrchestrator: k8sOrchestrator{
			label:     v1.NameLabel("err-pod-list-nil-k8s-orch-lbl"),
			name:      v1.OrchProviderRegistry("err-pod-list-nil-k8s-orch-name"),
			k8sUtlGtr: &errPodListNilK8sOrch{},
		},
	}

	volProfile := &okCtrlImgVolumeProfile{}

	_, err := mockedO.ReadStorage(volProfile)
	if err == nil {
		t.Errorf("TestCase: Error Match \n\tExpectedErr: 'not-nil' \n\tActualErr: 'nil'")
	}

	expErr := fmt.Sprintf("VSM(s) '%s:%s' not found at orchestrator '%s:%s'", "ok-ns", "ok-vsm-name", mockedO.Label(), mockedO.Name())

	if err != nil && err.Error() != expErr {
		t.Errorf("TestCase: Error Message Match \n\tExpectedErr: '%s' \n\tActualErr: '%s'", expErr, err.Error())
	}
}

// TestReadStorageReturnsErrPodMatch verifies the error due to mismatch of pods
func TestReadStorageReturnsErrPodMatch(t *testing.T) {
	mockedO := &errPodListMissK8sOrch{
		k8sOrchestrator: k8sOrchestrator{
			label:     v1.NameLabel("err-pod-list-miss-k8s-orch-lbl"),
			name:      v1.OrchProviderRegistry("err-pod-list-miss-k8s-orch-name"),
			k8sUtlGtr: &errPodListMissK8sOrch{},
		},
	}

	volProfile := &okCtrlImgVolumeProfile{}

	_, err := mockedO.ReadStorage(volProfile)
	if err == nil {
		t.Errorf("TestCase: Error Match \n\tExpectedErr: 'not-nil' \n\tActualErr: 'nil'")
	}

	expErr := fmt.Sprintf("VSM(s) '%s:%s' not found at orchestrator '%s:%s'", "ok-ns", "ok-vsm-name", mockedO.Label(), mockedO.Name())

	if err != nil && err.Error() != expErr {
		t.Errorf("TestCase: Error Message Match \n\tExpectedErr: '%s' \n\tActualErr: '%s'", expErr, err.Error())
	}
}

// TestReadStorageReturnsOk returns without error
func TestReadStorageReturnsOk(t *testing.T) {
	mockedO := &okReadStorageK8sOrch{
		k8sOrchestrator: k8sOrchestrator{
			k8sUtlGtr: &okReadStorageK8sOrch{},
		},
	}

	volProfile := &okCtrlImgVolumeProfile{}

	pvl, err := mockedO.ReadStorage(volProfile)

	if err != nil {
		t.Errorf("TestCase: Nil Error Match \n\tExpectedErr: 'nil' \n\tActualErr: '%s'", err.Error())
	}

	// Verify the pvl items
	if len(pvl.Items) != 1 {
		t.Errorf("TestCase: No of PVs \n\tExpectedNo: '1' \n\tActualNo: '%d'", len(pvl.Items))
	}

	// Verify the pv name
	expectedVSM := "ok-vsm-name"
	if pvl.Items[0].Name != expectedVSM {
		t.Errorf("TestCase: VSM Name Match \n\tExpectedName: '%s' \n\tActualName: '%s'", expectedVSM, pvl.Items[0].Name)
	}
}

// TestCreateControllerServiceReturnsErrVsmName verifies the vsm name error
func TestCreateControllerServiceReturnsErrVsmName(t *testing.T) {
	mockedO := &mockK8sOrch{
		k8sOrchestrator: k8sOrchestrator{},
	}

	_, err := mockedO.createControllerService(&errVsmNameVolumeProfile{})
	if err == nil {
		t.Errorf("TestCase: Error Match \n\tExpectedErr: 'not-nil' \n\tActualErr: 'nil'")
	}

	if err != nil && err.Error() != "err-vsm-name" {
		t.Errorf("TestCase: Error Message Match \n\tExpectedErr: 'err-vsm-name' \n\tActualErr: '%s'", err.Error())
	}
}

// TestCreateControllerServiceReturnsNoK8sClientSupport verifies no K8sClient
// support error
func TestCreateControllerServiceReturnsNoK8sClientSupport(t *testing.T) {
	mockedO := &noK8sClientSupportK8sOrch{
		k8sOrchestrator: k8sOrchestrator{
			k8sUtlGtr: &noK8sClientSupportK8sOrch{},
		},
	}

	volProfile := &okVsmNameVolumeProfile{}

	_, err := mockedO.createControllerService(volProfile)
	if err == nil {
		t.Errorf("TestCase: Error Match \n\tExpectedErr: 'not-nil' \n\tActualErr: 'nil'")
	}

	expErr := fmt.Sprintf("K8s client not supported by '%s'", "no-k8s-client-support-k8s-util")

	if err != nil && err.Error() != expErr {
		t.Errorf("TestCase: Error Message Match \n\tExpectedErr: '%s' \n\tActualErr: '%s'", expErr, err.Error())
	}
}

// TestCreateControllerServiceReturnsErrSvcOps verifies the services operator
// error
func TestCreateControllerServiceReturnsErrSvcOps(t *testing.T) {
	mockedO := &errSvcOpsK8sOrch{
		k8sOrchestrator: k8sOrchestrator{
			k8sUtlGtr: &errSvcOpsK8sOrch{},
		},
	}

	volProfile := &okCtrlImgVolumeProfile{}

	_, err := mockedO.createControllerService(volProfile)
	if err == nil {
		t.Errorf("TestCase: Error Match \n\tExpectedErr: 'not-nil' \n\tActualErr: 'nil'")
	}

	expErr := "err-svc-ops"

	if err != nil && err.Error() != expErr {
		t.Errorf("TestCase: Error Message Match \n\tExpectedErr: '%s' \n\tActualErr: '%s'", expErr, err.Error())
	}
}

// TestCreateControllerServiceReturnsOk verifies non error scenario
func TestCreateControllerServiceReturnsOk(t *testing.T) {
	mockedO := &okCreateServiceK8sOrch{
		k8sOrchestrator: k8sOrchestrator{
			k8sUtlGtr: &okCreateServiceK8sOrch{},
		},
	}

	volProfile := &okCtrlImgVolumeProfile{}

	svc, err := mockedO.createControllerService(volProfile)

	if err != nil {
		t.Errorf("TestCase: Nil Error Match \n\tExpectedErr: 'nil' \n\tActualErr: '%s'", err.Error())
	}

	// Verify the service kind
	if svc.Kind != string(v1.K8sKindService) {
		t.Errorf("TestCase: Kind Match \n\tExpectedKind: '%s' \n\tActualKind: '%s'", v1.K8sKindService, svc.Kind)
	}

	// Verify the service version
	if svc.APIVersion != string(v1.K8sServiceVersion) {
		t.Errorf("TestCase: Service Version Match \n\tExpectedAPIVersion: '%s' \n\tActualAPIVersion: '%s'", v1.K8sServiceVersion, svc.APIVersion)
	}

	// Verify the service name
	vsm, _ := volProfile.VSMName()
	eSvcName := vsm + string(v1.ControllerSuffix) + string(v1.ServiceSuffix)
	if svc.Name != eSvcName {
		t.Errorf("TestCase: Service Name Match \n\tExpectedName: '%s' \n\tActualName: '%s'", eSvcName, svc.Name)
	}

	// Verify the service labels
	eLblStr := string(v1.VSMSelectorPrefix) + vsm
	eLbl, _ := labels.Parse(eLblStr)
	if !eLbl.Matches(labels.Set(svc.Labels)) {
		t.Errorf("TestCase: Labels Match \n\tExpectedLabels: '%s' \n\tActualLabels: '%s'", eLbl, labels.Set(svc.Labels))
	}

	// Verify no. of ports within the service spec
	if len(svc.Spec.Ports) != 2 {
		t.Errorf("TestCase: No. of Ports \n\tExpectedPorts: '2' \n\tActualPorts: '%d'", len(svc.Spec.Ports))
	}

	// Verify the service spec labels
	eSelectorStr := string(v1.VSMSelectorPrefix) + vsm + string(v1.ControllerSuffix)
	eSelector, _ := labels.Parse(eSelectorStr)
	if !eSelector.Matches(labels.Set(svc.Spec.Selector)) {
		t.Errorf("TestCase: Selector Match \n\tExpectedSelector: '%s' \n\tActualSelector: '%s'", eSelector, labels.Set(svc.Spec.Selector))
	}
}

// TestGetControllerServiceReturnsErrVsmName verifies the vsm name error
func TestGetControllerServiceReturnsErrVsmName(t *testing.T) {
	mockedO := &mockK8sOrch{
		k8sOrchestrator: k8sOrchestrator{},
	}

	_, _, err := mockedO.getControllerService(&errVsmNameVolumeProfile{})
	if err == nil {
		t.Errorf("TestCase: Error Match \n\tExpectedErr: 'not-nil' \n\tActualErr: 'nil'")
	}

	if err != nil && err.Error() != "err-vsm-name" {
		t.Errorf("TestCase: Error Message Match \n\tExpectedErr: 'err-vsm-name' \n\tActualErr: '%s'", err.Error())
	}
}

// TestGetControllerServiceReturnsNoK8sClientSupport verifies no K8sClient
// support error
func TestGetControllerServiceReturnsNoK8sClientSupport(t *testing.T) {
	mockedO := &noK8sClientSupportK8sOrch{
		k8sOrchestrator: k8sOrchestrator{
			k8sUtlGtr: &noK8sClientSupportK8sOrch{},
		},
	}

	volProfile := &okVsmNameVolumeProfile{}

	_, _, err := mockedO.getControllerService(volProfile)
	if err == nil {
		t.Errorf("TestCase: Error Match \n\tExpectedErr: 'not-nil' \n\tActualErr: 'nil'")
	}

	expErr := fmt.Sprintf("K8s client not supported by '%s'", "no-k8s-client-support-k8s-util")

	if err != nil && err.Error() != expErr {
		t.Errorf("TestCase: Error Message Match \n\tExpectedErr: '%s' \n\tActualErr: '%s'", expErr, err.Error())
	}
}

// TestGetControllerServiceReturnsErrSvcOps verifies the services operator
// error
func TestGetControllerServiceReturnsErrSvcOps(t *testing.T) {
	mockedO := &errSvcOpsK8sOrch{
		k8sOrchestrator: k8sOrchestrator{
			k8sUtlGtr: &errSvcOpsK8sOrch{},
		},
	}

	volProfile := &okCtrlImgVolumeProfile{}

	_, _, err := mockedO.getControllerService(volProfile)
	if err == nil {
		t.Errorf("TestCase: Error Match \n\tExpectedErr: 'not-nil' \n\tActualErr: 'nil'")
	}

	expErr := "err-svc-ops"

	if err != nil && err.Error() != expErr {
		t.Errorf("TestCase: Error Message Match \n\tExpectedErr: '%s' \n\tActualErr: '%s'", expErr, err.Error())
	}
}

// TestGetControllerServiceReturnsErrSvcGet verifies the error received during
// a Get() call on services operator
func TestGetControllerServiceReturnsErrSvcGet(t *testing.T) {
	mockedO := &errSvcGetK8sOrch{
		k8sOrchestrator: k8sOrchestrator{
			k8sUtlGtr: &errSvcGetK8sOrch{},
		},
	}

	volProfile := &okCtrlImgVolumeProfile{}

	_, _, err := mockedO.getControllerService(volProfile)
	if err == nil {
		t.Errorf("TestCase: Error Match \n\tExpectedErr: 'not-nil' \n\tActualErr: 'nil'")
	}

	expErr := "err-svc-get"

	if err != nil && err.Error() != expErr {
		t.Errorf("TestCase: Error Message Match \n\tExpectedErr: '%s' \n\tActualErr: '%s'", expErr, err.Error())
	}
}

// TestGetControllerServiceReturnsOk verifies non error scenario
func TestGetControllerServiceReturnsOk(t *testing.T) {
	mockedO := &okGetServiceK8sOrch{
		k8sOrchestrator: k8sOrchestrator{
			k8sUtlGtr: &okGetServiceK8sOrch{},
		},
	}

	volProfile := &okCtrlImgVolumeProfile{}

	name, ip, err := mockedO.getControllerService(volProfile)

	if err != nil {
		t.Errorf("TestCase: Nil Error Match \n\tExpectedErr: 'nil' \n\tActualErr: '%s'", err.Error())
	}

	if name != "ok-svc-name" {
		t.Errorf("TestCase: Service Name Match \n\tExpectedName: '%s' \n\tActualName: '%s'", "ok-svc-name", name)
	}

	if ip != "1.1.1.1" {
		t.Errorf("TestCase: Service IP Match \n\tExpectedIP: '%s' \n\tActualIP: '%s'", "1.1.1.1", ip)
	}
}

// TestCreateReplicaPodsReturnsErrVsmName verifies the vsm name error
func TestCreateReplicaPodsReturnsErrVsmName(t *testing.T) {
	mockedO := &mockK8sOrch{
		k8sOrchestrator: k8sOrchestrator{},
	}

	_, err := mockedO.createReplicaPods(&errVsmNameVolumeProfile{}, "")
	if err == nil {
		t.Errorf("TestCase: Error Match \n\tExpectedErr: 'not-nil' \n\tActualErr: 'nil'")
	}

	if err != nil && err.Error() != "err-vsm-name" {
		t.Errorf("TestCase: Error Message Match \n\tExpectedErr: 'err-vsm-name' \n\tActualErr: '%s'", err.Error())
	}
}

// TestCreateReplicaPodsReturnsNoSupportReplicaImage verifies the no support
// during invocation of volProProfile.ReplicaImage()
func TestCreateReplicaPodsReturnsNoSupportReplicaImage(t *testing.T) {
	mockedO := &mockK8sOrch{
		k8sOrchestrator: k8sOrchestrator{},
	}

	volProfile := &noSupportRepImgVolumeProfile{
		&okVsmNameVolumeProfile{},
	}

	_, err := mockedO.createReplicaPods(volProfile, "")
	if err == nil {
		t.Errorf("TestCase: Error Match \n\tExpectedErr: 'not-nil' \n\tActualErr: 'nil'")
	}

	n, _ := volProfile.VSMName()
	expErr := fmt.Sprintf("VSM '%s' requires a replica container image", n)

	if err != nil && err.Error() != expErr {
		t.Errorf("TestCase: Error Message Match \n\tExpectedErr: '%s' \n\tActualErr: '%s'", expErr, err.Error())
	}
}

// TestCreateReplicaPodsReturnsErrReplicaImage verifies the error during
// invocation of volProProfile.ReplicaImage()
func TestCreateReplicaPodsReturnsErrReplicaImage(t *testing.T) {
	mockedO := &mockK8sOrch{
		k8sOrchestrator: k8sOrchestrator{},
	}

	volProfile := &errRepImgVolumeProfile{
		&okVsmNameVolumeProfile{},
	}

	_, err := mockedO.createReplicaPods(volProfile, "")
	if err == nil {
		t.Errorf("TestCase: Error Match \n\tExpectedErr: 'not-nil' \n\tActualErr: 'nil'")
	}

	if err != nil && err.Error() != "err-rep-image" {
		t.Errorf("TestCase: Error Message Match \n\tExpectedErr: '%s' \n\tActualErr: '%s'", "err-rep-image", err.Error())
	}
}

// TestCreateReplicaPodsReturnsErrReplicaCount verifies error
// during invocation of volProProfile.ReplicaCount()
func TestCreateReplicaPodsReturnsErrReplicaCount(t *testing.T) {
	mockedO := &mockK8sOrch{
		k8sOrchestrator: k8sOrchestrator{},
	}

	volProfile := &errRepCountVolumeProfile{
		&okVsmNameVolumeProfile{},
	}

	_, err := mockedO.createReplicaPods(volProfile, "")
	if err == nil {
		t.Errorf("TestCase: Error Match \n\tExpectedErr: 'not-nil' \n\tActualErr: 'nil'")
	}

	if err != nil && err.Error() != "err-rep-count" {
		t.Errorf("TestCase: Error Message Match \n\tExpectedErr: '%s' \n\tActualErr: '%s'", "err-rep-count", err.Error())
	}
}

// TestCreateReplicaPodsReturnsErrPersistentPathCount verifies error
// during invocation of volProProfile.PersistentPathCount()
func TestCreateReplicaPodsReturnsErrPersistentPathCount(t *testing.T) {
	mockedO := &mockK8sOrch{
		k8sOrchestrator: k8sOrchestrator{},
	}

	volProfile := &errPersistentPathCountVolumeProfile{
		&okVsmNameVolumeProfile{},
	}

	_, err := mockedO.createReplicaPods(volProfile, "")
	if err == nil {
		t.Errorf("TestCase: Error Match \n\tExpectedErr: 'not-nil' \n\tActualErr: 'nil'")
	}

	if err != nil && err.Error() != "err-persistent-path-count" {
		t.Errorf("TestCase: Error Message Match \n\tExpectedErr: '%s' \n\tActualErr: '%s'", "err-persistent-path-count", err.Error())
	}
}

// TestCreateReplicaPodsReturnsErrCountMatch verifies error
// during comparision of PersistentPathCount & ReplicaCount
func TestCreateReplicaPodsReturnsErrCountMatch(t *testing.T) {
	mockedO := &mockK8sOrch{
		k8sOrchestrator: k8sOrchestrator{},
	}

	volProfile := &errReplicaCountMatchVolumeProfile{
		&okVsmNameVolumeProfile{},
	}

	_, err := mockedO.createReplicaPods(volProfile, "")
	if err == nil {
		t.Errorf("TestCase: Error Match \n\tExpectedErr: 'not-nil' \n\tActualErr: 'nil'")
	}

	eError := "VSM 'ok-vsm-name' replica count '0' does not match persistent path count '1'"

	if err != nil && err.Error() != eError {
		t.Errorf("TestCase: Error Message Match \n\tExpectedErr: '%s' \n\tActualErr: '%s'", eError, err.Error())
	}
}

// TestCreateReplicaPodsReturnsErrCreateRepPod verifies error
// during invocation of createReplicaPod
func TestCreateReplicaPodsReturnsErrCreateRepPod(t *testing.T) {
	mockedO := &errCreateReplicaPodK8sOrch{
		// Most of the errCreateReplicaPodK8sOrch's methods are delegated to
		// k8sOrchestrator instance
		k8sOrchestrator: k8sOrchestrator{
			// We are only interested in modifying k8sUtlGtr property of k8sOrchestrator
			k8sUtlGtr: &errCreateReplicaPodK8sOrch{},
		},
	}

	volProfile := &okCreateReplicaPodVolumeProfile{
		&okVsmNameVolumeProfile{},
	}

	_, err := mockedO.createReplicaPods(volProfile, "")
	if err == nil {
		t.Errorf("TestCase: Error Match \n\tExpectedErr: 'not-nil' \n\tActualErr: 'nil'")
	}

	eError := "K8s client not supported by 'err-create-rep-pod-k8s-util'"
	if err != nil && err.Error() != eError {
		t.Errorf("TestCase: Error Message Match \n\tExpectedErr: '%s' \n\tActualErr: '%s'", eError, err.Error())
	}
}

// TestCreateReplicaPodsReturnsOk verifies non error scenario
func TestCreateReplicaPodsReturnsOk(t *testing.T) {
	mockedO := &okCreatePodK8sOrch{
		k8sOrchestrator: k8sOrchestrator{
			k8sUtlGtr: &okCreatePodK8sOrch{},
		},
	}

	volProfile := &okCreateReplicaPodVolumeProfile{
		&okVsmNameVolumeProfile{},
	}

	_, err := mockedO.createReplicaPods(volProfile, "")

	if err != nil {
		t.Errorf("TestCase: Nil Error Match \n\tExpectedErr: 'nil' \n\tActualErr: '%s'", err.Error())
	}
}

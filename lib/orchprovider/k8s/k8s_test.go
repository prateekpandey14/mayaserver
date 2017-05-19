package k8s

import (
	"errors"
	"testing"

	"github.com/openebs/mayaserver/lib/api/v1"
	"github.com/openebs/mayaserver/lib/orchprovider"
	volProfile "github.com/openebs/mayaserver/lib/profile/volumeprovisioner"
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
		pvc.Labels = map[string]string{
			string(v1.PVPVSMNameLbl): c.vsmname,
		}

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
		name:               pvc.Labels[string(testK8sUtlNameLbl)],
		vsmName:            pvc.Labels[string(v1.PVPVSMNameLbl)],
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
			"mock-ns", //ns
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
		pvc.Labels = map[string]string{
			string(testK8sUtlNameLbl):            c.kUtlName,
			string(v1.PVPVSMNameLbl):             c.vsmName,
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

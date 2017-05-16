package k8s

import (
	"testing"

	"github.com/openebs/mayaserver/lib/api/v1"
	"github.com/openebs/mayaserver/lib/orchprovider"
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

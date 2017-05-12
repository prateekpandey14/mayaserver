package k8s

import (
	"testing"
)

// TestK8sUtilInterfaceCompliance verifies if k8sUtil implements
// all the exposed methods of the desired interfaces.
//
// NOTE:
//    In case of non-compliance, this logic will error out during compile
// time itself.
func TestK8sUtilInterfaceCompliance(t *testing.T) {
	// k8sUtil implements K8sUtilInterface
	var _ K8sUtilInterface = &k8sUtil{}
	// k8sUtil implements K8sClients
	var _ K8sClients = &k8sUtil{}
}

// TestNewK8sUtil verifies the function that creates a new instance of
// k8sUtil. In addition, it verifies if the returned instance
// provides features it is expected of.
func TestNewK8sUtil(t *testing.T) {
	cases := []struct {
		incluster bool
		err       string
	}{
		{true, "unable to load in-cluster configuration, KUBERNETES_SERVICE_HOST and KUBERNETES_SERVICE_PORT must be defined"},
		{false, "OutClusterClientSet not supported in 'k8sutil'"},
	}

	for i, c := range cases {
		u, err := newK8sUtil()

		if err != nil {
			t.Errorf("TestCase: '%d' ExpectedError: 'nil' ActualError: '%s'", i, err.Error())
		}

		if "k8sutil" != u.Name() {
			t.Errorf("TestCase: '%d' ExpectedName: 'k8sutil' ActualName: '%s'", i, u.Name())
		}

		// K8s Clients is always supported by k8sUtil
		kc, supported := u.K8sClients()
		if !supported {
			t.Errorf("TestCase: '%d' ExpectedK8sClientsSupport: 'true' ActualK8sClientsSupport: '%t'", i, supported)
		}

		_, err = kc.GetClusterCS(c.incluster)

		if c.incluster && err != nil && c.err != err.Error() {
			t.Errorf("TestCase: '%d' ExpectedCSError: '%s' ActualCSError: '%s'", i, c.err, err.Error())
		}

		// out of cluster communication is not yet supported
		if !c.incluster && err == nil {
			t.Errorf("TestCase: '%d' ExpectedCSError: '%s' ActualCSError: 'nil'", i, c.err)
		}

		if !c.incluster && err != nil && c.err != err.Error() {
			t.Errorf("TestCase: '%d' ExpectedCSError: '%s' ActualCSError: '%s'", i, c.err, err.Error())
		}
	}
}

package server

import (
	//"fmt"
	"io"
	"log"
	"sync"

	"github.com/openebs/mayaserver/lib/api/v1"
	v1jiva "github.com/openebs/mayaserver/lib/api/v1/jiva"
	v1nomad "github.com/openebs/mayaserver/lib/api/v1/nomad"
	"github.com/openebs/mayaserver/lib/config"
	"github.com/openebs/mayaserver/lib/orchprovider"
	"github.com/openebs/mayaserver/lib/orchprovider/k8s"
	"github.com/openebs/mayaserver/lib/orchprovider/nomad"
	"github.com/openebs/mayaserver/lib/volumeprovisioner"
	"github.com/openebs/mayaserver/lib/volumeprovisioner/jiva"
)

// MayaApiServer is a long running stateless daemon that runs
// at openebs maya master(s)
type MayaApiServer struct {
	config    *config.MayaConfig
	logger    *log.Logger
	logOutput io.Writer

	shutdown     bool
	shutdownCh   chan struct{}
	shutdownLock sync.Mutex
}

// NewMayaApiServer is used to create a new maya api server
// with the given configuration
func NewMayaApiServer(config *config.MayaConfig, logOutput io.Writer) (*MayaApiServer, error) {

	ms := &MayaApiServer{
		config:     config,
		logger:     log.New(logOutput, "", log.LstdFlags|log.Lmicroseconds),
		logOutput:  logOutput,
		shutdownCh: make(chan struct{}),
	}

	err := ms.BootstrapPlugins()
	if err != nil {
		return nil, err
	}

	return ms, nil
}

// TODO
// Create a Bootstrap interface that facilitates initialization
// Create another Bootstraped interface that provides the initialized instances
// Perhaps at lib/bootstrap
// MayaServer struct will make use of above interfaces & hence specialized
// structs that cater to bootstraping & bootstraped features.
//
// NOTE:
//    The current implementation is tightly coupled & cannot be unit tested.
//
// NOTE:
//    Mayaserver should be entrusted to registering all possible variants of
// volume plugins.
//
// A volume plugin variant is composed of:
//    volume plugin + orchestrator of volume plugin + region of orchestrator
//
// In addition, Mayaserver should initialize the `default volume plugin`
// instance with its `default orchestrator` & `default region` of the
// orchestrator. User initiated requests requiring specific variants should be
// initialized at runtime.
func (ms *MayaApiServer) BootstrapPlugins() error {

	// TODO
	// Use MayaConfig
	// Fetch the names of volume plugins to be initialized
	// Iterate over the volumes:
	//  0. Fetch the config file location of orchestrator
	//  1. Initialize volume plugin's orchestrator
	//  2. Build an aspect that points to above orchestrator
	//  3. Fetch the config file location of volume plugin
	//  4. Initialize volume plugin

	// Typically the orchestrator should be initialized with its
	// region property. However, if initialization logic is not provided with
	// specific region then a default region specific to orchestrator is taken into
	// account.
	// TODO
	// Get this default value w.r.t orchestrator from mayaserver config file
	// e.g.
	//    NOMAD_REGION = global
	//    K8S_REGION = us-east-1
	// TODO
	// There may be cases where user initiated requests might specify region.
	// In those cases, a separate volumeplugin should be initialized with that
	// specific region.
	//
	// NOTE:
	//    In other words a particular volume plugin may have two
	// running instances pointing to different regions.

	// TODO
	// Deprecate
	// Old way to register
	found := orchprovider.IsOrchProvider(v1nomad.DefaultNomadPluginName)
	if !found {
		orchprovider.RegisterOrchProvider(
			// A variant of nomad orchestrator plugin
			v1nomad.DefaultNomadPluginName,
			// Below is a functional implementation that holds the initialization
			// logic of nomad orchestrator plugin
			func(name string, region string, config io.Reader) (orchprovider.OrchestratorInterface, error) {
				return nomad.NewNomadOrchestrator(name, region, config)
			})
	}

	orchestrator, err := orchprovider.InitOrchProvider(v1nomad.DefaultNomadPluginName, v1nomad.DefaultNomadRegionName, v1nomad.DefaultNomadConfigFile)
	if err != nil {
		return err
	}

	// Set the jiva aspects with its defaults
	jivaAspect := &jiva.JivaStorNomadAspect{
		Nomad: orchestrator,

		// The default datacenter. Typically user initiated actions will specify
		// a particular datacenter. This property is useful in cases where the actions
		// or requests do not specify a datacenter value.
		Datacenter: v1jiva.DefaultJivaDataCenter,
	}

	_, err = volumeprovisioner.InitVolumePlugin(v1jiva.DefaultJivaVolumePluginName, "", jivaAspect)
	if err != nil {
		return err
	}

	// New way to register orchestrator(s)
	isK8sOrchReg := orchprovider.HasOrchestrator(v1.K8sOrchestrator)
	if !isK8sOrchReg {
		orchprovider.RegisterOrchestrator(
			// Registration entry when Kubernetes is the orchestrator provider plugin
			v1.K8sOrchestrator,
			// Below is a callback function that creates a new instance of Kubernetes
			// orchestration provider
			func(label v1.NameLabel, name v1.OrchProviderRegistry) (orchprovider.OrchestratorInterface, error) {
				return k8s.NewK8sOrchestrator(label, name)
			})
	}

	return nil
}

// Shutdown is used to terminate MayaServer.
func (ms *MayaApiServer) Shutdown() error {

	ms.shutdownLock.Lock()
	defer ms.shutdownLock.Unlock()

	ms.logger.Println("[INFO] maya api server: requesting shutdown")

	if ms.shutdown {
		return nil
	}

	ms.logger.Println("[INFO] maya api server: shutdown complete")
	ms.shutdown = true

	close(ms.shutdownCh)

	return nil
}

// Leave is used gracefully exit.
func (ms *MayaApiServer) Leave() error {

	ms.logger.Println("[INFO] maya api server: exiting gracefully")

	// Nothing as of now
	return nil
}

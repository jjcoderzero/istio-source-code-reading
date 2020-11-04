package bootstrap

import (
	"encoding/json"

	"istio.io/istio/pkg/util/gogoprotomarshal"

	"istio.io/pkg/filewatcher"
	"istio.io/pkg/log"
	"istio.io/pkg/version"

	"istio.io/istio/pkg/config/mesh"
)

// initMeshConfiguration creates the mesh in the pilotConfig from the input arguments.
func (s *Server) initMeshConfiguration(args *PilotArgs, fileWatcher filewatcher.FileWatcher) {
	log.Info("initializing mesh configuration")
	defer func() {
		if s.environment.Watcher != nil {
			meshdump, _ := gogoprotomarshal.ToJSONWithIndent(s.environment.Mesh(), "    ")
			log.Infof("mesh configuration: %s", meshdump)
			log.Infof("version: %s", version.Info.String())
			argsdump, _ := json.MarshalIndent(args, "", "   ")
			log.Infof("flags: %s", argsdump)
		}
	}()

	var err error
	if args.MeshConfigFile != "" {
		s.environment.Watcher, err = mesh.NewWatcher(fileWatcher, args.MeshConfigFile)
		if err == nil {
			return
		}
		log.Warnf("Watching mesh config file %s failed: %v", args.MeshConfigFile, err)
	}

	// Config file either wasn't specified or failed to load - use a default mesh.
	meshConfig := mesh.DefaultMeshConfig()
	s.environment.Watcher = mesh.NewFixedWatcher(&meshConfig)
}

// initMeshNetworks loads the mesh networks configuration from the file provided
// in the args and add a watcher for changes in this file.
func (s *Server) initMeshNetworks(args *PilotArgs, fileWatcher filewatcher.FileWatcher) {
	log.Info("initializing mesh networks")
	if args.NetworksConfigFile != "" {
		var err error
		s.environment.NetworksWatcher, err = mesh.NewNetworksWatcher(fileWatcher, args.NetworksConfigFile)
		if err != nil {
			log.Infoa(err)
		}
	}

	if s.environment.NetworksWatcher == nil {
		log.Info("mesh networks configuration not provided")
		s.environment.NetworksWatcher = mesh.NewFixedNetworksWatcher(nil)
	}
}

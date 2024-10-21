package server

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/openrelayxyz/xplugeth"
	"github.com/openrelayxyz/xplugeth/types"
	"github.com/openrelayxyz/xplugeth/hooks/apis"
	"github.com/openrelayxyz/xplugeth/hooks/initialize"
	
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
)

func isValidPath(path string) bool {
    cleanPath := filepath.Clean(path)

    fileInfo, err := os.Stat(cleanPath)
    if err != nil {
        return false
    }
    return fileInfo.IsDir()
}

func disablePlugins() bool {
    val := os.Getenv("DISABLE_PLUGINS")
    if len(val) > 0 {
    	truthyValues := []string{"true", "1", "yes", "y"}
    	for _, truthy := range truthyValues {
        	if strings.EqualFold(val, truthy) {
            	return true
        	}
    	}
	}
    return false
}


func pluginsConfig() string {
	pluginsConfigEnv := os.Getenv("PLUGINS_CONFIG")

	if pluginsConfigEnv != "" && isValidPath(pluginsConfigEnv) {
		log.Info("plugins config path provided", "path", pluginsConfigEnv)
		return pluginsConfigEnv
	} 

	log.Info("plugins config path not provided or invalid")
	return ""
}

func pluginInitializeNode() {
	stack, ok := xplugeth.GetSingleton[*node.Node]()
	if !ok {
		panic("*node.Node singleton not set, xplugeth InitializeNode")
	}
	backend, ok := xplugeth.GetSingleton[types.Backend]()
	if !ok {
		panic("types.Backend singleton not set, xplugeth Initializenode")
	}
	for _, init := range xplugeth.GetModules[initialize.Initializer]() {
		init.InitializeNode(stack, backend)
	}
}

func pluginGetAPIs() []rpc.API {
	result := []rpc.API{}

	stack, ok := xplugeth.GetSingleton[*node.Node]()
	if !ok {
		panic("*node.Node singleton not set, xplugeth GetAPIs")
	}
	backend, ok := xplugeth.GetSingleton[types.Backend]()
	if !ok {
		panic("types.Backend singleton not set xplugeth GetAPIs")
	}

	for _, a := range xplugeth.GetModules[apis.GetAPIs]() {
		result = append(result, a.GetAPIs(stack, backend)...)
	}

	return result
}

func pluginOnShutdown() {
	for _, shutdown := range xplugeth.GetModules[initialize.Shutdown]() {
		shutdown.Shutdown()
	}
}

func pluginBlockchain() {
	for _, b := range xplugeth.GetModules[initialize.Blockchain]() {
		b.Blockchain()
	}
}

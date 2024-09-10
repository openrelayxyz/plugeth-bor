package server

import (
	"github.com/openrelayxyz/xplugeth"
	"github.com/openrelayxyz/xplugeth/types"
	
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
)

type initializer interface {
	InitializeNode(*node.Node, types.Backend)
}
type shutdown interface {
	Shutdown()
}
type blockchain interface {
	Blockchain()	
}
type getAPIs interface {
	GetAPIs(*node.Node, types.Backend) []rpc.API
}

func init() {
	xplugeth.RegisterHook[initializer]()
	xplugeth.RegisterHook[shutdown]()
	xplugeth.RegisterHook[blockchain]()
	xplugeth.RegisterHook[getAPIs]()
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
	for _, init := range xplugeth.GetModules[initializer]() {
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

	for _, a := range xplugeth.GetModules[getAPIs]() {
		result = append(result, a.GetAPIs(stack, backend)...)
	}

	return result
}

func pluginOnShutdown() {
	for _, shutdown := range xplugeth.GetModules[shutdown]() {
		shutdown.Shutdown()
	}
}

func pluginBlockchain() {
	for _, b := range xplugeth.GetModules[blockchain]() {
		b.Blockchain()
	}
}

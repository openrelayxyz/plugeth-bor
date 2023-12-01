package main

import (
	"context"
	// "time"
	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted"
)

var (
	log core.Logger
	backend restricted.Backend
)

type testService struct {
	backend restricted.Backend
}

func(*testService) Test(ctx context.Context) string {
	return "Returning from polygon plugeth light"
}

func Initialize(ctx core.Context, loader core.PluginLoader, logger core.Logger) {
	log = logger
	log.Info("Initialized plugin")
}

func InitializeNode(stack core.Node, b restricted.Backend) {
	backend = b
	log.Info("Initialized node block updater plugin")
}

func GetAPIs(stack core.Node, backend restricted.Backend) []core.API {
	log.Error("Returning APIs from plugin")
	return []core.API{
	 {
		 Namespace: "plugeth",
		 Version:	 "1.0",
		 Service:	 &testService{backend},
		 Public:		true,
	 },
 }
}
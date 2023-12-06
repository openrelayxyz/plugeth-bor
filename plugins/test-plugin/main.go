package main

import(
	"flag"

	"github.com/openrelayxyz/plugeth-utils/core"
)

var (
	pl      core.PluginLoader
	log     core.Logger
	Flags = *flag.NewFlagSet("test-plugin", flag.ContinueOnError)
	argsPresent = Flags.Bool("testFlag", false, "confirming the plugethArgs() injection is present")
)

var httpApiFlagName = "http.api"

func Initialize(ctx core.Context, loader core.PluginLoader, logger core.Logger) { 
	pl = loader
	log = logger
	log.Info("Loaded test plugin")
}
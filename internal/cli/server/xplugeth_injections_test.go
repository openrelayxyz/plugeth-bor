package server

import(
	"testing"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/openrelayxyz/xplugeth"
	"github.com/openrelayxyz/xplugeth/hooks/initialize"
	"github.com/openrelayxyz/xplugeth/hooks/apis"
	"github.com/openrelayxyz/xplugeth/types"
)

func init() {
	xplugeth.RegisterModule[gethTestModule]("gethTestModule")
}

func TestServerPkgInjections(t *testing.T) {
	cfg := DefaultConfig()
	
	h := &HeimdallConfig{
		Without:     true,
	}
	
	cfg.Heimdall = h
	cfg.Verbosity = 0

	srv, err := NewServer(cfg)
	if err != nil {
		t.Fatalf("error returned from NewServer, err %v", err)
	}

	time.Sleep(2 * time.Second)

	srv.Stop()
		
	time.Sleep(2 * time.Second)

	if len(injections) > 0 {
		var uncalledInjections []string
		for k, _ := range injections {
			uncalledInjections = append(uncalledInjections, k)
		}
		t.Fatalf("test failed, injections not called %v", uncalledInjections)
	} 
}

type gethTestModule struct {}

func (*gethTestModule) InitializeNode(s *node.Node, b types.Backend) {
	if s != nil && b != nil {
		delete(injections, "InitializeNode")
	} else {
		fmt.Printf("get singleton error stack %v, backend %v", s, b)
	}
}

func (*gethTestModule) GetAPIs(*node.Node, types.Backend) []rpc.API {
	delete(injections, "GetAPIs")
	return nil
}

func (*gethTestModule) Blockchain() {
	delete(injections, "BlockChain")
}

func (*gethTestModule) Shutdown() {
	delete(injections, "Shutdown")
}

var injections map[string]struct{} = map[string]struct{}{
	"InitializeNode":struct{}{},
	"GetAPIs":struct{}{},
	"BlockChain":struct{}{},
	"Shutdown":struct{}{},
}


var (
	_ initialize.Initializer = (*gethTestModule)(nil)
	_ initialize.Blockchain = (*gethTestModule)(nil)
	_ initialize.Shutdown = (*gethTestModule)(nil)
	_ apis.GetAPIs = (*gethTestModule)(nil)
)
package fetcher

import(
	// "testing"
	"fmt"
	// "syscall"
	// "flag"
	// "time"

	// "github.com/urfave/cli/v2"

	gtypes "github.com/ethereum/go-ethereum/core/types"
	// "github.com/ethereum/go-ethereum/rpc"

	"github.com/openrelayxyz/xplugeth"
	"github.com/openrelayxyz/xplugeth/hooks/blockchain"
	// "github.com/openrelayxyz/xplugeth/hooks/apis"
	// "github.com/openrelayxyz/xplugeth/types"
)

func init() {
	xplugeth.RegisterModule[gethTestModule]("gethTestModule")
}

// func TestGethPkgInjections(t *testing.T) {

// 	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	
// 	flagSet.String("gcmode", "full", "Blockchain garbage collection mode")
// 	flagSet.String("crypto.kzg", "gokzg", "KZG library implementation to use; gokzg (recommended) or ckzg")
	
// 	if err := flagSet.Parse([]string{}); err != nil {
// 		t.Fatalf("Failed to parse flag set: %v", err)
// 	}

// 	app := &cli.App{
// 		Flags: nodeFlags,
// 	}
	
// 	ctx := cli.NewContext(app, flagSet, nil)

// 	go func() {
// 		if err := geth(ctx); err != nil {
// 			t.Fatalf("err calling geth, err %v", err)
// 		}
// 	}()

// 	time.Sleep(2 * time.Second)

// 	if err := syscall.Kill(syscall.Getpid(), syscall.SIGINT); err != nil {
// 		t.Fatalf("failed to send SIGINT: %v", err)
// 	}

// 	time.Sleep(2 * time.Second)
	
// 	if len(injections) > 0 {
// 		var uncalledInjections []string
// 		for k, _ := range injections {
// 			uncalledInjections = append(uncalledInjections, k)
// 		}
// 		t.Fatalf("test failed, injections not called %v", uncalledInjections)
// 	} 
// }

type gethTestModule struct {}

func (p *gethTestModule) PeerEval(id string, headers []*gtypes.Header) {
	fmt.Println("PeerEval plugin engaged")
	// delete(injections, "PeerEval")
}

var injections map[string]struct{} = map[string]struct{}{
	"PeerEval":struct{}{},
	
}


var (
	_ blockchain.PeerEvalPlugin = (*gethTestModule)(nil)
)
package state

import(
	"testing"
	
	"github.com/ethereum/go-ethereum/common"

	"github.com/openrelayxyz/xplugeth"
	"github.com/openrelayxyz/xplugeth/hooks/stateupdates"
)

func init() {
	xplugeth.RegisterModule[stateTestModule]("stateTestModule")
}

type stateTestModule struct {}

func TestStateInjections(t *testing.T) {
	xplugeth.Initialize("")

	TestStateChanges(t)

	if len(injections) > 0 {
		var uncalledInjections []string
		for k, _ := range injections {
			uncalledInjections = append(uncalledInjections, k)
		}
		t.Fatalf("injections not called %v", uncalledInjections)
	} 
}

func (*stateTestModule) StateUpdate(blockRoot, parentRoot common.Hash, destructs map[common.Hash]struct{}, accounts map[common.Hash][]byte, storage map[common.Hash]map[common.Hash][]byte, codeUpdates map[common.Hash][]byte) {
	delete(injections, "StateUpdate")
}

var injections map[string]struct{} = map[string]struct{}{
	"StateUpdate":struct{}{},
}

var (
	_ stateupdates.StateUpdatePlugin = (*stateTestModule)(nil)
)
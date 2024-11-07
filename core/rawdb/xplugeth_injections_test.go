package rawdb

import(
	"testing"

	
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/openrelayxyz/xplugeth"
	"github.com/openrelayxyz/xplugeth/hooks/modifyancients"
)

func init() {
	xplugeth.RegisterModule[rawdbTestModule]("rawdbTestModule")
}

type rawdbTestModule struct {}

func TestRawDBInjections(t *testing.T) {
	xplugeth.Initialize("")

	TestAncientStorage(t)

	if len(injections) > 0 {
		var uncalledInjections []string
		for k, _ := range injections {
			uncalledInjections = append(uncalledInjections, k)
		}
		t.Fatalf("injections not called %v", uncalledInjections)
	} 
}

func (*rawdbTestModule) ModifyAncients(uint64, *types.Header) {
	delete(injections, "ModifyAncients")
}

var injections map[string]struct{} = map[string]struct{}{
	"ModifyAncients":struct{}{},
}

var (
	_ modifyancients.ModifyAncientsPlugin = (*rawdbTestModule)(nil)
)
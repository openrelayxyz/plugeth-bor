package core

import(
	"math/big"
	"testing"
	"time"

	
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/openrelayxyz/xplugeth"
	"github.com/openrelayxyz/xplugeth/hooks/blockchain"
)

func init() {
	xplugeth.RegisterModule[coreTestModule]("coreTestModule")
}

type coreTestModule struct {}

func TestCoreInjections(t *testing.T) {
	xplugeth.Initialize("")

	testBlockchainHeaderchainReorgConsistency(t, rawdb.HashScheme)
	testBlockchainHeaderchainReorgConsistency(t, rawdb.PathScheme)

	if len(injections) > 0 {
		var uncalledInjections []string
		for k, _ := range injections {
			uncalledInjections = append(uncalledInjections, k)
		}
		t.Fatalf("injections not called %v", uncalledInjections)
	} 
}

func (*coreTestModule) NewHead(block *types.Block, hash common.Hash, logs []*types.Log, td *big.Int) {
	delete(injections, "NewHead")
}

func (*coreTestModule) Reorg(common common.Hash, oldChain []common.Hash, newChain []common.Hash) {
	delete(injections, "Reorg")
}

func (*coreTestModule) NewSideBlock(block *types.Block, hash common.Hash, logs []*types.Log) {
	delete(injections, "SideBlock")
}

func (*coreTestModule) SetTrieFlushIntervalClone(flushInterval time.Duration) time.Duration {
	delete(injections, "SetTrieFlush")
	return flushInterval
}

var injections map[string]struct{} = map[string]struct{}{
	"NewHead":struct{}{},
	"Reorg":struct{}{},
	"SideBlock":struct{}{},
	"SetTrieFlush":struct{}{},
}

var (
	_ blockchain.NewHeadPlugin = (*coreTestModule)(nil)
	_ blockchain.ReorgPlugin = (*coreTestModule)(nil)
	_ blockchain.NewSideBlockPlugin = (*coreTestModule)(nil)
	_ blockchain.SetTrieFlushIntervalClonePlugin = (*coreTestModule)(nil)
)
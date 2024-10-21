package core

import (
	"math/big"
	"sync"
	"time"

	"github.com/openrelayxyz/xplugeth"
	"github.com/openrelayxyz/xplugeth/hooks/blockchain"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

func pluginNewHead(block *types.Block, hash common.Hash, logs []*types.Log, td *big.Int) {
	for _, m := range xplugeth.GetModules[blockchain.NewHeadPlugin]() {
		m.NewHead(block, hash, logs, td)
	}
}

func pluginNewSideBlock(block *types.Block, hash common.Hash, logs []*types.Log) {
	for _, m := range xplugeth.GetModules[blockchain.NewSideBlockPlugin]() {
		m.NewSideBlock(block, hash, logs)
	}
}

func pluginReorg(commonBlock *types.Block, oldChain, newChain types.Blocks) {
	oldChainHashes := make([]common.Hash, len(oldChain))
	for i, block := range oldChain {
		oldChainHashes[i] = block.Hash()
	}
	newChainHashes := make([]common.Hash, len(newChain))
	for i, block := range newChain {
		newChainHashes[i] = block.Hash()
	}
	for _, m := range xplugeth.GetModules[blockchain.ReorgPlugin]() {
		m.Reorg(commonBlock.Hash(), oldChainHashes, newChainHashes)
	}
}

func pluginSetTrieFlushIntervalClone(flushInterval time.Duration) time.Duration {
	m := xplugeth.GetModules[blockchain.SetTrieFlushIntervalClonePlugin]()
	var snc sync.Once
	if len(m) > 1 {
		snc.Do(func() { log.Warn("The blockChain flushInterval value is being accessed by multiple plugins") })
	}
	for _, m := range m {
		flushInterval = m.SetTrieFlushIntervalClone(flushInterval)
	}
	return flushInterval
}

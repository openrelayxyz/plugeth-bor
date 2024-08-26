package core

import (
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/openrelayxyz/xplugeth"
)

type newHeadPlugin interface {
	NewHead(block *types.Block, hash common.Hash, logs []*types.Log, td *big.Int)
}

type newSideBlockPlugin interface {
	NewSideBlock(block *types.Block, hash common.Hash, logs []*types.Log)
}

type reorgPlugin interface {
	Reorg(commonBlock common.Hash, oldChain, newChain []common.Hash)
}

type setTrieFlushIntervalClonePlugin interface {
	SetTrieFlushIntervalClone(flushInterval time.Duration) time.Duration
}

func init() {
	xplugeth.RegisterHook[newHeadPlugin]()
	xplugeth.RegisterHook[newSideBlockPlugin]()
	xplugeth.RegisterHook[reorgPlugin]()
	xplugeth.RegisterHook[setTrieFlushIntervalClonePlugin]()
}

func pluginNewHead(block *types.Block, hash common.Hash, logs []*types.Log, td *big.Int) {
	for _, m := range xplugeth.GetModules[newHeadPlugin]() {
		m.NewHead(block, hash, logs, td)
	}
}

func pluginNewSideBlock(block *types.Block, hash common.Hash, logs []*types.Log) {
	for _, m := range xplugeth.GetModules[newSideBlockPlugin]() {
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
	for _, m := range xplugeth.GetModules[reorgPlugin]() {
		m.Reorg(commonBlock.Hash(), oldChainHashes, newChainHashes)
	}
}

func pluginSetTrieFlushIntervalClone(flushInterval time.Duration) time.Duration {
	m := xplugeth.GetModules[setTrieFlushIntervalClonePlugin]()
	var snc sync.Once
	if len(m) > 1 {
		snc.Do(func() { log.Warn("The blockChain flushInterval value is being accessed by multiple plugins") })
	}
	for _, m := range m {
		flushInterval = m.SetTrieFlushIntervalClone(flushInterval)
	}
	return flushInterval
}

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
	PluginNewHead(block *types.Block, hash common.Hash, logs []*types.Log, td *big.Int)
}

type newSideBlockPlugin interface {
	PluginNewSideBlock(block *types.Block, hash common.Hash, logs []*types.Log)
}

type reorgPlugin interface {
	PluginReorg(commonBlock *types.Block, oldChain, newChain types.Blocks)
}

type setTrieFlushIntervalClonePlugin interface {
	PluginSetTrieFlushIntervalClone(flushInterval time.Duration) time.Duration
}

func init() {
	xplugeth.RegisterHook[newHeadPlugin]()
	xplugeth.RegisterHook[newSideBlockPlugin]()
	xplugeth.RegisterHook[reorgPlugin]()
	xplugeth.RegisterHook[setTrieFlushIntervalClonePlugin]()
}

func PluginNewHead(block *types.Block, hash common.Hash, logs []*types.Log, td *big.Int) {
	for _, m := range xplugeth.GetModules[newHeadPlugin]() {
		m.PluginNewHead(block, hash, logs, td)
	}
}

func PluginNewSideBlock(block *types.Block, hash common.Hash, logs []*types.Log) {
	for _, m := range xplugeth.GetModules[newSideBlockPlugin]() {
		m.PluginNewSideBlock(block, hash, logs)
	}
}

func PluginReorg(commonBlock *types.Block, oldChain, newChain types.Blocks) {
	for _, m := range xplugeth.GetModules[reorgPlugin]() {
		m.PluginReorg(commonBlock, oldChain, newChain)
	}
}

func PluginSetTrieFlushIntervalClone(flushInterval time.Duration) time.Duration {
	m := xplugeth.GetModules[setTrieFlushIntervalClonePlugin]()
	var snc sync.Once
	if len(m) > 1 {
		snc.Do(func() { log.Warn("The blockChain flushInterval value is being accessed by multiple plugins") })
	}
	for _, m := range m {
		flushInterval = m.PluginSetTrieFlushIntervalClone(flushInterval)
	}
	return flushInterval
}

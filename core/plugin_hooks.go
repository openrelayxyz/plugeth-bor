package core

import (
	"encoding/json"
	"math/big"
	"time"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/plugins"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/openrelayxyz/plugeth-utils/core"
)

var (
	preprocessBlockTimer   = metrics.NewRegisteredTimer("plugins/core/block/preprocess", nil)
	postprocessBlockTimer   = metrics.NewRegisteredTimer("plugins/core/block/postprocess", nil)
	sideblockTimer   = metrics.NewRegisteredTimer("plugins/core/block/side", nil)
	preprocessTxTimer   = metrics.NewRegisteredTimer("plugins/core/tx/preprocess", nil)
	postprocessTxTimer   = metrics.NewRegisteredTimer("plugins/core/tx/postprocess", nil)
	processBlockErrorTimer   = metrics.NewRegisteredTimer("plugins/core/block/error", nil)
	newheadTimer   = metrics.NewRegisteredTimer("plugins/core/newhead", nil)
	reorgTimer   = metrics.NewRegisteredTimer("plugins/core/reorg", nil)
)

func PluginPreProcessBlock(pl *plugins.PluginLoader, block *types.Block) {
	start := time.Now()
	defer preprocessBlockTimer.UpdateSince(start)
	fnList := pl.Lookup("PreProcessBlock", func(item interface{}) bool {
		_, ok := item.(func(core.Hash, uint64, []byte))
		return ok
	})
	encoded, _ := rlp.EncodeToBytes(block)
	for _, fni := range fnList {
		if fn, ok := fni.(func(core.Hash, uint64, []byte)); ok {
			fn(core.Hash(block.Hash()), block.NumberU64(), encoded)
		}
	}
}
func pluginPreProcessBlock(block *types.Block) {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting PreProcessBlock, but default PluginLoader has not been initialized")
		return
	}
	PluginPreProcessBlock(plugins.DefaultPluginLoader, block) // TODO
}
func PluginPreProcessTransaction(pl *plugins.PluginLoader, tx *types.Transaction, block *types.Block, i int) {
	start := time.Now()
	defer preprocessTxTimer.UpdateSince(start)
	fnList := pl.Lookup("PreProcessTransaction", func(item interface{}) bool {
		_, ok := item.(func([]byte, core.Hash, core.Hash, int))
		return ok
	})
	txBytes, _ := tx.MarshalBinary()
	for _, fni := range fnList {
		if fn, ok := fni.(func([]byte, core.Hash, core.Hash, int)); ok {
			fn(txBytes, core.Hash(tx.Hash()), core.Hash(block.Hash()), i)
		}
	}
}
func pluginPreProcessTransaction(tx *types.Transaction, block *types.Block, i int) {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting PreProcessTransaction, but default PluginLoader has not been initialized")
		return
	}
	PluginPreProcessTransaction(plugins.DefaultPluginLoader, tx, block, i)
}
func PluginBlockProcessingError(pl *plugins.PluginLoader, tx *types.Transaction, block *types.Block, err error) {
	start := time.Now()
	defer processBlockErrorTimer.UpdateSince(start)
	fnList := pl.Lookup("BlockProcessingError", func(item interface{}) bool {
		_, ok := item.(func(core.Hash, core.Hash, error))
		return ok
	})
	for _, fni := range fnList {
		if fn, ok := fni.(func(core.Hash, core.Hash, error)); ok {
			fn(core.Hash(tx.Hash()), core.Hash(block.Hash()), err)
		}
	}
}
func pluginBlockProcessingError(tx *types.Transaction, block *types.Block, err error) {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting BlockProcessingError, but default PluginLoader has not been initialized")
		return
	}
	PluginBlockProcessingError(plugins.DefaultPluginLoader, tx, block, err)
}
func PluginPostProcessTransaction(pl *plugins.PluginLoader, tx *types.Transaction, block *types.Block, i int, receipt *types.Receipt) {
	start := time.Now()
	defer postprocessTxTimer.UpdateSince(start)
	fnList := pl.Lookup("PostProcessTransaction", func(item interface{}) bool {
		_, ok := item.(func(core.Hash, core.Hash, int, []byte))
		return ok
	})
	receiptBytes, _ := json.Marshal(receipt)
	for _, fni := range fnList {
		if fn, ok := fni.(func(core.Hash, core.Hash, int, []byte)); ok {
			fn(core.Hash(tx.Hash()), core.Hash(block.Hash()), i, receiptBytes)
		}
	}
}
func pluginPostProcessTransaction(tx *types.Transaction, block *types.Block, i int, receipt *types.Receipt) {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting PostProcessTransaction, but default PluginLoader has not been initialized")
		return
	}
	PluginPostProcessTransaction(plugins.DefaultPluginLoader, tx, block, i, receipt)
}
func PluginPostProcessBlock(pl *plugins.PluginLoader, block *types.Block) {
	start := time.Now()
	defer postprocessBlockTimer.UpdateSince(start)
	fnList := pl.Lookup("PostProcessBlock", func(item interface{}) bool {
		_, ok := item.(func(core.Hash))
		return ok
	})
	for _, fni := range fnList {
		if fn, ok := fni.(func(core.Hash)); ok {
			fn(core.Hash(block.Hash()))
		}
	}
}
func pluginPostProcessBlock(block *types.Block) {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting PostProcessBlock, but default PluginLoader has not been initialized")
		return
	}
	PluginPostProcessBlock(plugins.DefaultPluginLoader, block)
}

func PluginNewHead(pl *plugins.PluginLoader, block *types.Block, hash common.Hash, logs []*types.Log, td *big.Int) {
	start := time.Now()
	defer newheadTimer.UpdateSince(start)
	fnList := pl.Lookup("NewHead", func(item interface{}) bool {
		_, ok := item.(func([]byte, core.Hash, [][]byte, *big.Int))
		return ok
	})
	blockBytes, _ := rlp.EncodeToBytes(block)
	logBytes := make([][]byte, len(logs))
	for i, l := range logs {
		logBytes[i], _ = rlp.EncodeToBytes(l)
	}
	for _, fni := range fnList {
		if fn, ok := fni.(func([]byte, core.Hash, [][]byte, *big.Int)); ok {
			fn(blockBytes, core.Hash(hash), logBytes, td)
		}
	}
}

func pluginNewHead(block *types.Block, hash common.Hash, logs []*types.Log, td *big.Int) {

	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting NewHead, but default PluginLoader has not been initialized")
		return
	}
	PluginNewHead(plugins.DefaultPluginLoader, block, hash, logs, td)
}

func PluginNewSideBlock(pl *plugins.PluginLoader, block *types.Block, hash common.Hash, logs []*types.Log) {
	start := time.Now()
	defer sideblockTimer.UpdateSince(start)
	fnList := pl.Lookup("NewSideBlock", func(item interface{}) bool {
		_, ok := item.(func([]byte, core.Hash, [][]byte))
		return ok
	})
	blockBytes, _ := rlp.EncodeToBytes(block)
	logBytes := make([][]byte, len(logs))
	for i, l := range logs {
		logBytes[i], _ = rlp.EncodeToBytes(l)
	}
	for _, fni := range fnList {
		if fn, ok := fni.(func([]byte, core.Hash, [][]byte)); ok {
			fn(blockBytes, core.Hash(hash), logBytes)
		}
	}
}
func pluginNewSideBlock(block *types.Block, hash common.Hash, logs []*types.Log) {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting NewSideBlock, but default PluginLoader has not been initialized")
		return
	}
	PluginNewSideBlock(plugins.DefaultPluginLoader, block, hash, logs)
}

func PluginReorg(pl *plugins.PluginLoader, commonBlock *types.Block, oldChain, newChain types.Blocks) {
	start := time.Now()
	defer reorgTimer.UpdateSince(start)
	fnList := pl.Lookup("Reorg", func(item interface{}) bool {
		_, ok := item.(func(core.Hash, []core.Hash, []core.Hash))
		return ok
	})
	oldChainHashes := make([]core.Hash, len(oldChain))
	for i, block := range oldChain {
		oldChainHashes[i] = core.Hash(block.Hash())
	}
	newChainHashes := make([]core.Hash, len(newChain))
	for i, block := range newChain {
		newChainHashes[i] = core.Hash(block.Hash())
	}
	for _, fni := range fnList {
		if fn, ok := fni.(func(core.Hash, []core.Hash, []core.Hash)); ok {
			fn(core.Hash(commonBlock.Hash()), oldChainHashes, newChainHashes)
		}
	}
}
func pluginReorg(commonBlock *types.Block, oldChain, newChain types.Blocks) {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting Reorg, but default PluginLoader has not been initialized")
		return
	}
	PluginReorg(plugins.DefaultPluginLoader, commonBlock, oldChain, newChain)
}

func PluginSetTrieFlushIntervalClone(pl *plugins.PluginLoader, flushInterval time.Duration) time.Duration {
	fnList := pl.Lookup("SetTrieFlushIntervalClone", func(item interface{}) bool{
		_, ok := item.(func(time.Duration) time.Duration)
		return ok
	})
	var snc sync.Once
	if len(fnList) > 1 {
		snc.Do(func() {log.Warn("The blockChain flushInterval value is being accessed by multiple plugins")})
	}
	for _, fni := range fnList {
		if fn, ok := fni.(func(time.Duration) time.Duration); ok {
			flushInterval = fn(flushInterval) 
		}
	}
	return flushInterval
}

func pluginSetTrieFlushIntervalClone(flushInterval time.Duration) time.Duration {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting setTreiFlushIntervalClone, but default PluginLoader has not been initialized")
		return flushInterval
	}
	return PluginSetTrieFlushIntervalClone(plugins.DefaultPluginLoader, flushInterval)
}
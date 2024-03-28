package state

import (
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/plugins"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/core/state/snapshot"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/openrelayxyz/plugeth-utils/core"
	"time"
)

var (
	acctCheckTimer = metrics.NewRegisteredTimer("plugeth/statedb/accounts/checks", nil)
)

type pluginSnapshot struct {
	root common.Hash
}

func (s *pluginSnapshot) Root() common.Hash {
	return s.root
}

func (s *pluginSnapshot) Account(hash common.Hash) (*types.SlimAccount, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *pluginSnapshot) AccountRLP(hash common.Hash) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *pluginSnapshot) Storage(accountHash, storageHash common.Hash) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

func PluginStateUpdate(pl *plugins.PluginLoader, blockRoot, parentRoot common.Hash, snap snapshot.Snapshot, trie Trie, destructs map[common.Hash]struct{}, accounts map[common.Hash][]byte, storage map[common.Hash]map[common.Hash][]byte, codeUpdates map[common.Hash][]byte) {
	checker := &acctChecker{snap, trie}
	fnList := pl.Lookup("StateUpdate", func(item interface{}) bool {
		_, ok := item.(func(core.Hash, core.Hash, map[core.Hash]struct{}, map[core.Hash][]byte, map[core.Hash]map[core.Hash][]byte, map[core.Hash][]byte))
		return ok
	})
	coreDestructs := make(map[core.Hash]struct{})
	for k, v := range destructs {
		coreDestructs[core.Hash(k)] = v
	}
	coreAccounts := make(map[core.Hash][]byte)
	start := time.Now()
	for k, v := range accounts {
		if checker.updated(k, v) {
			coreAccounts[core.Hash(k)] = v
		}
	}
	acctCheckTimer.UpdateSince(start)
	coreStorage := make(map[core.Hash]map[core.Hash][]byte)
	for k, v := range storage {
		coreStorage[core.Hash(k)] = make(map[core.Hash][]byte)
		for h, d := range v {
			coreStorage[core.Hash(k)][core.Hash(h)] = d
		}
	}
	coreCode := make(map[core.Hash][]byte)
	for k, v := range codeUpdates {
		coreCode[core.Hash(k)] = v
	}

	for _, fni := range fnList {
		if fn, ok := fni.(func(core.Hash, core.Hash, map[core.Hash]struct{}, map[core.Hash][]byte, map[core.Hash]map[core.Hash][]byte, map[core.Hash][]byte)); ok {
			fn(core.Hash(blockRoot), core.Hash(parentRoot), coreDestructs, coreAccounts, coreStorage, coreCode)
		}
	}
}

func pluginStateUpdate(blockRoot, parentRoot common.Hash, snap snapshot.Snapshot, trie Trie, destructs map[common.Hash]struct{}, accounts map[common.Hash][]byte, storage map[common.Hash]map[common.Hash][]byte, codeUpdates map[common.Hash][]byte) {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting StateUpdate, but default PluginLoader has not been initialized")
		return
	}
	PluginStateUpdate(plugins.DefaultPluginLoader, blockRoot, parentRoot, snap, trie, destructs, accounts, storage, codeUpdates)
}


type acctChecker struct {
	snap snapshot.Snapshot
	trie Trie
}

func (ac *acctChecker) updated(k common.Hash, v []byte) bool {
	if updated, ok := ac.snapUpdated(k, v); ok {
		return updated
	}
	return ac.trieUpdated(k, v)
}

func (ac *acctChecker) snapUpdated(k common.Hash, v []byte) (bool, bool) {
	acct, err := ac.snap.AccountRLP(k)
	if err != nil {
		return false, false
	}
	if len(acct) == 0 {
		return false, false
	}
	return !bytes.Equal(acct, v), true
}

type acctByHasher interface {
	GetAccountByHash(common.Hash) (*types.StateAccount, error)
}

func (ac *acctChecker) trieUpdated(k common.Hash, v []byte) bool {
	trie, ok := ac.trie.(acctByHasher)
	if !ok {
		log.Warn("Couldn't check trie updates, wrong trie type")
		return true
	}
	oAcct, err := trie.GetAccountByHash(k)
	if err != nil {
		return true
	}
	return !bytes.Equal(v, types.SlimAccountRLP(*oAcct))
}

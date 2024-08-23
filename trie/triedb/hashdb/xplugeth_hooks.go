package hashdb

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/xplugeth"
)

type preTrieCommit interface {
	PreTrieCommit(node core.Hash)
}

type postTrieCommit interface {
	PostTrieCommit(node core.Hash)
}

func init() {
	xplugeth.RegisterHook[preTrieCommit]()
	xplugeth.RegisterHook[postTrieCommit]()
}

func PluginPreTrieCommit(node common.Hash) {
	for _, m := range xplugeth.GetModules[preTrieCommit]() {
		m.PreTrieCommit(core.Hash(node))
	}
}

func PluginPostTrieCommit(node common.Hash) {
	for _, m := range xplugeth.GetModules[postTrieCommit]() {
		m.PostTrieCommit(core.Hash(node))
	}
}

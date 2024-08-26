package hashdb

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/openrelayxyz/xplugeth"
)

type preTrieCommit interface {
	PreTrieCommit(node common.Hash)
}

type postTrieCommit interface {
	PostTrieCommit(node common.Hash)
}

func init() {
	xplugeth.RegisterHook[preTrieCommit]()
	xplugeth.RegisterHook[postTrieCommit]()
}

func pluginPreTrieCommit(node common.Hash) {
	for _, m := range xplugeth.GetModules[preTrieCommit]() {
		m.PreTrieCommit(node)
	}
}

func pluginPostTrieCommit(node common.Hash) {
	for _, m := range xplugeth.GetModules[postTrieCommit]() {
		m.PostTrieCommit(node)
	}
}

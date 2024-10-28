package hashdb

import (
	"github.com/openrelayxyz/xplugeth"
	"github.com/openrelayxyz/xplugeth/hooks/triecommit"

	"github.com/ethereum/go-ethereum/common"
)

func pluginPreTrieCommit(node common.Hash) {
	for _, m := range xplugeth.GetModules[triecommit.PreTrieCommit]() {
		m.PreTrieCommit(node)
	}
}

func pluginPostTrieCommit(node common.Hash) {
	for _, m := range xplugeth.GetModules[triecommit.PostTrieCommit]() {
		m.PostTrieCommit(node)
	}
}

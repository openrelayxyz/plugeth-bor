package fetcher

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/openrelayxyz/xplugeth"
	"github.com/openrelayxyz/xplugeth/hooks/blockchain"
)

func pluginPeerEval(id string, headers []*types.Header) {
	for _, m := range xplugeth.GetModules[blockchain.PeerEvalPlugin]() {
		m.PeerEval(id, headers)
	}
}

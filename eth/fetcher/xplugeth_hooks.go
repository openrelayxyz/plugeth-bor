package fetcher

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/openrelayxyz/xplugeth"
	"github.com/openrelayxyz/xplugeth/hooks/blockchain"
)

func pluginPeerEval(id string, headers []*types.Header) {
	log.Error("peer", "peer id", id)
	for _, m := range xplugeth.GetModules[blockchain.PeerEvalPlugin]() {
		m.PeerEval(id, headers)
	}
}

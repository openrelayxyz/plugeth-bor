package eth

import (
	"github.com/ethereum/go-ethereum/core/types"
)

func (b *EthAPIBackend) InsertBlock(block *types.Block) error {
	_, err := b.eth.BlockChain().InsertChain(types.Blocks{block})
	return err
}
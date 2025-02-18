package core

import (
	"github.com/zenanet-network/go-zenanet/common"
	"github.com/zenanet-network/go-zenanet/core/rawdb"
	"github.com/zenanet-network/go-zenanet/core/types"
)

// GetEireneReceiptByHash retrieves the eirene block receipt in a given block.
func (bc *BlockChain) GetEireneReceiptByHash(hash common.Hash) *types.Receipt {
	if receipt, ok := bc.eireneReceiptsCache.Get(hash); ok {
		return receipt
	}

	// read header from hash
	number := rawdb.ReadHeaderNumber(bc.db, hash)
	if number == nil {
		return nil
	}

	// read eirene receipt by hash and number
	receipt := rawdb.ReadEireneReceipt(bc.db, hash, *number, bc.chainConfig)
	if receipt == nil {
		return nil
	}

	// add into eirene receipt cache
	bc.eireneReceiptsCache.Add(hash, receipt)

	return receipt
}

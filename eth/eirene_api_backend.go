package eth

import (
	"context"
	"errors"
	"fmt"

	"github.com/zenanet-network/go-zenanet"
	"github.com/zenanet-network/go-zenanet/common"
	"github.com/zenanet-network/go-zenanet/consensus/eirene"
	"github.com/zenanet-network/go-zenanet/core"
	"github.com/zenanet-network/go-zenanet/core/rawdb"
	"github.com/zenanet-network/go-zenanet/core/types"
	"github.com/zenanet-network/go-zenanet/event"
	"github.com/zenanet-network/go-zenanet/rpc"
)

var errEireneEngineNotAvailable error = errors.New("Only available in Eirene engine")

// GetRootHash returns root hash for given start and end block
func (b *EthAPIBackend) GetRootHash(ctx context.Context, starBlockNr uint64, endBlockNr uint64) (string, error) {
	var api *eirene.API

	for _, _api := range b.eth.Engine().APIs(b.eth.BlockChain()) {
		if _api.Namespace == "eirene" {
			api = _api.Service.(*eirene.API)
		}
	}

	if api == nil {
		return "", errEireneEngineNotAvailable
	}

	root, err := api.GetRootHash(starBlockNr, endBlockNr)
	if err != nil {
		return "", err
	}

	return root, nil
}

// GetVoteOnHash returns the vote on hash
func (b *EthAPIBackend) GetVoteOnHash(ctx context.Context, starBlockNr uint64, endBlockNr uint64, hash string, milestoneId string) (bool, error) {
	var api *eirene.API

	for _, _api := range b.eth.Engine().APIs(b.eth.BlockChain()) {
		if _api.Namespace == "eirene" {
			api = _api.Service.(*eirene.API)
		}
	}

	if api == nil {
		return false, errEireneEngineNotAvailable
	}

	// Confirmation of 16 blocks on the endblock
	tipConfirmationBlockNr := endBlockNr + uint64(16)

	// Check if tipConfirmation block exit
	_, err := b.BlockByNumber(ctx, rpc.BlockNumber(tipConfirmationBlockNr))
	if err != nil {
		return false, errTipConfirmationBlock
	}

	// Check if end block exist
	localEndBlock, err := b.BlockByNumber(ctx, rpc.BlockNumber(endBlockNr))
	if err != nil {
		return false, errEndBlock
	}

	localEndBlockHash := localEndBlock.Hash().String()

	downloader := b.eth.handler.downloader
	isLocked := downloader.LockMutex(endBlockNr)

	if !isLocked {
		downloader.UnlockMutex(false, "", endBlockNr, common.Hash{})
		return false, errors.New("whitelisted number or locked sprint number is more than the received end block number")
	}

	if localEndBlockHash != hash {
		downloader.UnlockMutex(false, "", endBlockNr, common.Hash{})
		return false, fmt.Errorf("hash mismatch: localChainHash %s, milestoneHash %s", localEndBlockHash, hash)
	}

	downloader.UnlockMutex(true, milestoneId, endBlockNr, localEndBlock.Hash())

	return true, nil
}

// GetEireneBlockReceipt returns eirene block receipt
func (b *EthAPIBackend) GetEireneBlockReceipt(ctx context.Context, hash common.Hash) (*types.Receipt, error) {
	receipt := b.eth.blockchain.GetEireneReceiptByHash(hash)
	if receipt == nil {
		return nil, zenanet.NotFound
	}

	return receipt, nil
}

// GetEireneBlockLogs returns eirene block logs
func (b *EthAPIBackend) GetEireneBlockLogs(ctx context.Context, hash common.Hash) ([]*types.Log, error) {
	receipt := b.eth.blockchain.GetEireneReceiptByHash(hash)
	if receipt == nil {
		return nil, nil
	}

	return receipt.Logs, nil
}

// GetEireneBlockTransaction returns eirene block tx
func (b *EthAPIBackend) GetEireneBlockTransaction(ctx context.Context, hash common.Hash) (*types.Transaction, common.Hash, uint64, uint64, error) {
	tx, blockHash, blockNumber, index := rawdb.ReadEireneTransaction(b.eth.ChainDb(), hash)
	return tx, blockHash, blockNumber, index, nil
}

func (b *EthAPIBackend) GetEireneBlockTransactionWithBlockHash(ctx context.Context, txHash common.Hash, blockHash common.Hash) (*types.Transaction, common.Hash, uint64, uint64, error) {
	tx, blockHash, blockNumber, index := rawdb.ReadEireneTransactionWithBlockHash(b.eth.ChainDb(), txHash, blockHash)
	return tx, blockHash, blockNumber, index, nil
}

// SubscribeStateSyncEvent subscribes to state sync event
func (b *EthAPIBackend) SubscribeStateSyncEvent(ch chan<- core.StateSyncEvent) event.Subscription {
	return b.eth.BlockChain().SubscribeStateSyncEvent(ch)
}

// SubscribeChain2HeadEvent subscribes to reorg/head/fork event
func (b *EthAPIBackend) SubscribeChain2HeadEvent(ch chan<- core.Chain2HeadEvent) event.Subscription {
	return b.eth.BlockChain().SubscribeChain2HeadEvent(ch)
}

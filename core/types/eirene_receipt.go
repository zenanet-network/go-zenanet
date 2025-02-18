package types

import (
	"encoding/binary"
	"math/big"
	"sort"

	"github.com/zenanet-network/go-zenanet/common"
	"github.com/zenanet-network/go-zenanet/crypto"
)

// TenToTheFive - To be used while sorting eirene logs
//
// Sorted using ( blockNumber * (10 ** 5) + logIndex )
const TenToTheFive uint64 = 100000

var (
	eireneReceiptPrefix = []byte("matic-eirene-receipt-") // eireneReceiptPrefix + number + block hash -> eirene block receipt

	// SystemAddress address for system sender
	SystemAddress = common.HexToAddress("0xffffFFFfFFffffffffffffffFfFFFfffFFFfFFfE")
)

// EireneReceiptKey = eireneReceiptPrefix + num (uint64 big endian) + hash
func EireneReceiptKey(number uint64, hash common.Hash) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, number)

	return append(append(eireneReceiptPrefix, enc...), hash.Bytes()...)
}

// GetDerivedEireneTxHash get derived tx hash from receipt key
func GetDerivedEireneTxHash(receiptKey []byte) common.Hash {
	return common.BytesToHash(crypto.Keccak256(receiptKey))
}

// NewEireneTransaction create new eirene transaction for eirene receipt
func NewEireneTransaction() *Transaction {
	return NewTransaction(0, common.Address{}, big.NewInt(0), 0, big.NewInt(0), make([]byte, 0))
}

// DeriveFieldsForEireneReceipt fills the receipts with their computed fields based on consensus
// data and contextual infos like containing block and transactions.
func DeriveFieldsForEireneReceipt(receipt *Receipt, hash common.Hash, number uint64, receipts Receipts) error {
	// get derived tx hash
	txHash := GetDerivedEireneTxHash(EireneReceiptKey(number, hash))
	txIndex := uint(len(receipts))

	// set tx hash and tx index
	receipt.TxHash = txHash
	receipt.TransactionIndex = txIndex
	receipt.BlockHash = hash
	receipt.BlockNumber = big.NewInt(0).SetUint64(number)

	logIndex := 0
	for i := 0; i < len(receipts); i++ {
		logIndex += len(receipts[i].Logs)
	}

	// The derived log fields can simply be set from the block and transaction
	for j := 0; j < len(receipt.Logs); j++ {
		receipt.Logs[j].BlockNumber = number
		receipt.Logs[j].BlockHash = hash
		receipt.Logs[j].TxHash = txHash
		receipt.Logs[j].TxIndex = txIndex
		receipt.Logs[j].Index = uint(logIndex)
		logIndex++
	}

	return nil
}

// DeriveFieldsForEireneLogs fills the receipts with their computed fields based on consensus
// data and contextual infos like containing block and transactions.
func DeriveFieldsForEireneLogs(logs []*Log, hash common.Hash, number uint64, txIndex uint, logIndex uint) {
	// get derived tx hash
	txHash := GetDerivedEireneTxHash(EireneReceiptKey(number, hash))

	// the derived log fields can simply be set from the block and transaction
	for j := 0; j < len(logs); j++ {
		logs[j].BlockNumber = number
		logs[j].BlockHash = hash
		logs[j].TxHash = txHash
		logs[j].TxIndex = txIndex
		logs[j].Index = logIndex
		logIndex++
	}
}

// MergeEireneLogs merges receipt logs and block receipt logs
func MergeEireneLogs(logs []*Log, eireneLogs []*Log) []*Log {
	result := append(logs, eireneLogs...)

	sort.SliceStable(result, func(i int, j int) bool {
		return (result[i].BlockNumber*TenToTheFive + uint64(result[i].Index)) < (result[j].BlockNumber*TenToTheFive + uint64(result[j].Index))
	})

	return result
}

package rawdb

import (
	"math/big"

	"github.com/zenanet-network/go-zenanet/common"
	"github.com/zenanet-network/go-zenanet/core/types"
	"github.com/zenanet-network/go-zenanet/ethdb"
	"github.com/zenanet-network/go-zenanet/log"
	"github.com/zenanet-network/go-zenanet/params"
	"github.com/zenanet-network/go-zenanet/rlp"
)

var (
	// eirene receipt key
	eireneReceiptKey = types.EireneReceiptKey

	// eireneTxLookupPrefix + hash -> transaction/receipt lookup metadata
	eireneTxLookupPrefix = []byte(eireneTxLookupPrefixStr)
)

const (
	eireneTxLookupPrefixStr = "zenanet-eirene-tx-lookup-"

	// freezerEireneReceiptTable indicates the name of the freezer eirene receipts table.
	freezerEireneReceiptTable = "zenanet-eirene-receipts"
)

// eireneTxLookupKey = eireneTxLookupPrefix + eirene tx hash
func eireneTxLookupKey(hash common.Hash) []byte {
	return append(eireneTxLookupPrefix, hash.Bytes()...)
}

func ReadEireneReceiptRLP(db ethdb.Reader, hash common.Hash, number uint64) rlp.RawValue {
	var data []byte

	err := db.ReadAncients(func(reader ethdb.AncientReaderOp) error {
		// Check if the data is in ancients
		if isCanon(reader, number, hash) {
			data, _ = reader.Ancient(freezerEireneReceiptTable, number)

			return nil
		}

		// If not, try reading from leveldb
		data, _ = db.Get(eireneReceiptKey(number, hash))

		return nil
	})

	if err != nil {
		log.Warn("during ReadEireneReceiptRLP", "number", number, "hash", hash, "err", err)
	}

	return data
}

// ReadRawEireneReceipt retrieves the block receipt belonging to a block.
// The receipt metadata fields are not guaranteed to be populated, so they
// should not be used. Use ReadEireneReceipt instead if the metadata is needed.
func ReadRawEireneReceipt(db ethdb.Reader, hash common.Hash, number uint64) *types.Receipt {
	// Retrieve the flattened receipt slice
	data := ReadEireneReceiptRLP(db, hash, number)
	if len(data) == 0 {
		return nil
	}

	// Convert the receipts from their storage form to their internal representation
	var storageReceipt types.ReceiptForStorage
	if err := rlp.DecodeBytes(data, &storageReceipt); err != nil {
		log.Error("Invalid eirene receipt RLP", "hash", hash, "err", err)
		return nil
	}

	return (*types.Receipt)(&storageReceipt)
}

// ReadEireneReceipt retrieves all the eirene block receipts belonging to a block, including
// its corresponding metadata fields. If it is unable to populate these metadata
// fields then nil is returned.
func ReadEireneReceipt(db ethdb.Reader, hash common.Hash, number uint64, config *params.ChainConfig) *types.Receipt {
	if config != nil && config.Eirene != nil && config.Eirene.Sprint != nil && !config.Eirene.IsSprintStart(number) {
		return nil
	}

	// We're deriving many fields from the block body, retrieve beside the receipt
	eireneReceipt := ReadRawEireneReceipt(db, hash, number)
	if eireneReceipt == nil {
		return nil
	}

	// We're deriving many fields from the block body, retrieve beside the receipt
	receipts := ReadRawReceipts(db, hash, number)
	if receipts == nil {
		return nil
	}

	body := ReadBody(db, hash, number)
	if body == nil {
		log.Error("Missing body but have eirene receipt", "hash", hash, "number", number)
		return nil
	}

	if err := types.DeriveFieldsForEireneReceipt(eireneReceipt, hash, number, receipts); err != nil {
		log.Error("Failed to derive eirene receipt fields", "hash", hash, "number", number, "err", err)
		return nil
	}

	return eireneReceipt
}

// WriteEireneReceipt stores all the eirene receipt belonging to a block.
func WriteEireneReceipt(db ethdb.KeyValueWriter, hash common.Hash, number uint64, eireneReceipt *types.ReceiptForStorage) {
	// Convert the eirene receipt into their storage form and serialize them
	bytes, err := rlp.EncodeToBytes(eireneReceipt)
	if err != nil {
		log.Crit("Failed to encode eirene receipt", "err", err)
	}

	// Store the flattened receipt slice
	if err := db.Put(eireneReceiptKey(number, hash), bytes); err != nil {
		log.Crit("Failed to store eirene receipt", "err", err)
	}
}

// DeleteEireneReceipt removes receipt data associated with a block hash.
func DeleteEireneReceipt(db ethdb.KeyValueWriter, hash common.Hash, number uint64) {
	key := eireneReceiptKey(number, hash)

	if err := db.Delete(key); err != nil {
		log.Crit("Failed to delete eirene receipt", "err", err)
	}
}

// ReadEireneTransactionWithBlockHash retrieves a specific eirene (fake) transaction by tx hash and block hash, along with
// its added positional metadata.
func ReadEireneTransactionWithBlockHash(db ethdb.Reader, txHash common.Hash, blockHash common.Hash) (*types.Transaction, common.Hash, uint64, uint64) {
	blockNumber := ReadEireneTxLookupEntry(db, txHash)
	if blockNumber == nil {
		return nil, common.Hash{}, 0, 0
	}

	body := ReadBody(db, blockHash, *blockNumber)
	if body == nil {
		log.Error("Transaction referenced missing", "number", blockNumber, "hash", blockHash)
		return nil, common.Hash{}, 0, 0
	}

	// fetch receipt and return it
	return types.NewEireneTransaction(), blockHash, *blockNumber, uint64(len(body.Transactions))
}

// ReadEireneTransaction retrieves a specific Eirene (fake) transaction by hash, along with
// its added positional metadata.
func ReadEireneTransaction(db ethdb.Reader, hash common.Hash) (*types.Transaction, common.Hash, uint64, uint64) {
	blockNumber := ReadEireneTxLookupEntry(db, hash)
	if blockNumber == nil {
		return nil, common.Hash{}, 0, 0
	}

	blockHash := ReadCanonicalHash(db, *blockNumber)
	if blockHash == (common.Hash{}) {
		return nil, common.Hash{}, 0, 0
	}

	body := ReadBody(db, blockHash, *blockNumber)
	if body == nil {
		log.Error("Transaction referenced missing", "number", blockNumber, "hash", blockHash)
		return nil, common.Hash{}, 0, 0
	}

	// fetch receipt and return it
	return types.NewEireneTransaction(), blockHash, *blockNumber, uint64(len(body.Transactions))
}

//
// Indexes for reverse lookup
//

// ReadEireneTxLookupEntry retrieves the positional metadata associated with a transaction
// hash to allow retrieving the eirene transaction or eirene receipt using tx hash.
func ReadEireneTxLookupEntry(db ethdb.Reader, txHash common.Hash) *uint64 {
	data, _ := db.Get(eireneTxLookupKey(txHash))
	if len(data) == 0 {
		return nil
	}

	number := new(big.Int).SetBytes(data).Uint64()

	return &number
}

// WriteEireneTxLookupEntry stores a positional metadata for eirene transaction using block hash and block number
func WriteEireneTxLookupEntry(db ethdb.KeyValueWriter, hash common.Hash, number uint64) {
	txHash := types.GetDerivedEireneTxHash(eireneReceiptKey(number, hash))
	if err := db.Put(eireneTxLookupKey(txHash), big.NewInt(0).SetUint64(number).Bytes()); err != nil {
		log.Crit("Failed to store eirene transaction lookup entry", "err", err)
	}
}

// DeleteEireneTxLookupEntry removes eirene transaction data associated with block hash and block number
func DeleteEireneTxLookupEntry(db ethdb.KeyValueWriter, hash common.Hash, number uint64) {
	txHash := types.GetDerivedEireneTxHash(eireneReceiptKey(number, hash))
	DeleteEireneTxLookupEntryByTxHash(db, txHash)
}

// DeleteEireneTxLookupEntryByTxHash removes eirene transaction data associated with a eirene tx hash.
func DeleteEireneTxLookupEntryByTxHash(db ethdb.KeyValueWriter, txHash common.Hash) {
	if err := db.Delete(eireneTxLookupKey(txHash)); err != nil {
		log.Crit("Failed to delete eirene transaction lookup entry", "err", err)
	}
}

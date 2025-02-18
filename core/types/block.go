// Copyright 2014 The go-zenanet Authors
// This file is part of the go-zenanet library.
//
// The go-zenanet library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-zenanet library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-zenanet library. If not, see <http://www.gnu.org/licenses/>.

// Package types contains data types related to Zenanet consensus.
package types

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
	"reflect"
	"slices"
	"sync/atomic"
	"time"

	"github.com/zenanet-network/go-zenanet/common"
	"github.com/zenanet-network/go-zenanet/common/hexutil"
	"github.com/zenanet-network/go-zenanet/log"
	"github.com/zenanet-network/go-zenanet/params"
	"github.com/zenanet-network/go-zenanet/rlp"
)

var (
	ExtraVanityLength = 32 // Fixed number of extra-data prefix bytes reserved for signer vanity
	ExtraSealLength   = 65 // Fixed number of extra-data suffix bytes reserved for signer seal
)

// A BlockNonce is a 64-bit hash which proves (combined with the
// mix-hash) that a sufficient amount of computation has been carried
// out on a block.
type BlockNonce [8]byte

// EncodeNonce converts the given integer to a block nonce.
func EncodeNonce(i uint64) BlockNonce {
	var n BlockNonce

	binary.BigEndian.PutUint64(n[:], i)

	return n
}

// Uint64 returns the integer value of a block nonce.
func (n BlockNonce) Uint64() uint64 {
	return binary.BigEndian.Uint64(n[:])
}

// MarshalText encodes n as a hex string with 0x prefix.
func (n BlockNonce) MarshalText() ([]byte, error) {
	return hexutil.Bytes(n[:]).MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (n *BlockNonce) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("BlockNonce", input, n[:])
}

//go:generate go run github.com/fjl/gencodec -type Header -field-override headerMarshaling -out gen_header_json.go
//go:generate go run ../../rlp/rlpgen -type Header -out gen_header_rlp.go

// Header represents a block header in the Zenanet blockchain.
type Header struct {
	ParentHash  common.Hash    `json:"parentHash"       gencodec:"required"`
	UncleHash   common.Hash    `json:"sha3Uncles"       gencodec:"required"`
	Coinbase    common.Address `json:"miner"`
	Root        common.Hash    `json:"stateRoot"        gencodec:"required"`
	TxHash      common.Hash    `json:"transactionsRoot" gencodec:"required"`
	ReceiptHash common.Hash    `json:"receiptsRoot"     gencodec:"required"`
	Bloom       Bloom          `json:"logsBloom"        gencodec:"required"`
	Difficulty  *big.Int       `json:"difficulty"       gencodec:"required"`
	Number      *big.Int       `json:"number"           gencodec:"required"`
	GasLimit    uint64         `json:"gasLimit"         gencodec:"required"`
	GasUsed     uint64         `json:"gasUsed"          gencodec:"required"`
	Time        uint64         `json:"timestamp"        gencodec:"required"`
	Extra       []byte         `json:"extraData"        gencodec:"required"`
	MixDigest   common.Hash    `json:"mixHash"`
	Nonce       BlockNonce     `json:"nonce"`

	// BaseFee was added by EIP-1559 and is ignored in legacy headers.
	BaseFee *big.Int `json:"baseFeePerGas" rlp:"optional"`

	// WithdrawalsHash was added by EIP-4895 and is ignored in legacy headers.
	WithdrawalsHash *common.Hash `json:"withdrawalsRoot" rlp:"optional"`

	// BlobGasUsed was added by EIP-4844 and is ignored in legacy headers.
	BlobGasUsed *uint64 `json:"blobGasUsed" rlp:"optional"`

	// ExcessBlobGas was added by EIP-4844 and is ignored in legacy headers.
	ExcessBlobGas *uint64 `json:"excessBlobGas" rlp:"optional"`

	// ParentBeaconRoot was added by EIP-4788 and is ignored in legacy headers.
	ParentBeaconRoot *common.Hash `json:"parentBeaconBlockRoot" rlp:"optional"`
}

// Used for Encoding and Decoding of the Extra Data Field
type BlockExtraData struct {
	ValidatorBytes []byte

	// length of TxDependency          ->   n (n = number of transactions in the block)
	// length of TxDependency[i]       ->   k (k = a whole number)
	// k elements in TxDependency[i]   ->   transaction indexes on which transaction i is dependent on
	TxDependency [][]uint64
}

// field type overrides for gencodec
type headerMarshaling struct {
	Difficulty    *hexutil.Big
	Number        *hexutil.Big
	GasLimit      hexutil.Uint64
	GasUsed       hexutil.Uint64
	Time          hexutil.Uint64
	Extra         hexutil.Bytes
	BaseFee       *hexutil.Big
	Hash          common.Hash `json:"hash"` // adds call to Hash() in MarshalJSON
	BlobGasUsed   *hexutil.Uint64
	ExcessBlobGas *hexutil.Uint64
}

// Hash returns the block hash of the header, which is simply the keccak256 hash of its
// RLP encoding.
func (h *Header) Hash() common.Hash {
	return rlpHash(h)
}

var headerSize = common.StorageSize(reflect.TypeOf(Header{}).Size())

// Size returns the approximate memory used by all internal contents. It is used
// to approximate and limit the memory consumption of various caches.
func (h *Header) Size() common.StorageSize {
	var baseFeeBits int
	if h.BaseFee != nil {
		baseFeeBits = h.BaseFee.BitLen()
	}

	return headerSize + common.StorageSize(len(h.Extra)+(h.Difficulty.BitLen()+h.Number.BitLen()+baseFeeBits)/8)
}

// SanityCheck checks a few basic things -- these checks are way beyond what
// any 'sane' production values should hold, and can mainly be used to prevent
// that the unbounded fields are stuffed with junk data to add processing
// overhead
func (h *Header) SanityCheck() error {
	if h.Number != nil && !h.Number.IsUint64() {
		return fmt.Errorf("too large block number: bitlen %d", h.Number.BitLen())
	}

	if h.Difficulty != nil {
		if diffLen := h.Difficulty.BitLen(); diffLen > 80 {
			return fmt.Errorf("too large block difficulty: bitlen %d", diffLen)
		}
	}

	if eLen := len(h.Extra); eLen > 100*1024 {
		return fmt.Errorf("too large block extradata: size %d", eLen)
	}

	if h.BaseFee != nil {
		if bfLen := h.BaseFee.BitLen(); bfLen > 256 {
			return fmt.Errorf("too large base fee: bitlen %d", bfLen)
		}
	}

	return nil
}

// EmptyBody returns true if there is no additional 'body' to complete the header
// that is: no transactions, no uncles and no withdrawals.
func (h *Header) EmptyBody() bool {
	if h.WithdrawalsHash != nil {
		return h.TxHash == EmptyTxsHash && *h.WithdrawalsHash == EmptyWithdrawalsHash
	}
	return h.TxHash == EmptyTxsHash && h.UncleHash == EmptyUncleHash
}

// EmptyReceipts returns true if there are no receipts for this header/block.
func (h *Header) EmptyReceipts() bool {
	return h.ReceiptHash == EmptyReceiptsHash
}

// ValidateBlockNumberOptionsPIP15 validates the block range passed as in the options parameter in the conditional transaction (PIP-15)
func (h *Header) ValidateBlockNumberOptionsPIP15(minBlockNumber *big.Int, maxBlockNumber *big.Int) error {
	currentBlockNumber := h.Number

	if minBlockNumber != nil {
		if currentBlockNumber.Cmp(minBlockNumber) == -1 {
			return fmt.Errorf("current block number %v is less than minimum block number: %v", currentBlockNumber, minBlockNumber)
		}
	}

	if maxBlockNumber != nil {
		if currentBlockNumber.Cmp(maxBlockNumber) == 1 {
			return fmt.Errorf("current block number %v is greater than maximum block number: %v", currentBlockNumber, maxBlockNumber)
		}
	}

	return nil
}

// ValidateBlockNumberOptionsPIP15 validates the timestamp range passed as in the options parameter in the conditional transaction (PIP-15)
func (h *Header) ValidateTimestampOptionsPIP15(minTimestamp *uint64, maxTimestamp *uint64) error {
	currentBlockTime := h.Time

	if minTimestamp != nil {
		if currentBlockTime < *minTimestamp {
			return fmt.Errorf("current block time %v is less than minimum timestamp: %v", currentBlockTime, minTimestamp)
		}
	}

	if maxTimestamp != nil {
		if currentBlockTime > *maxTimestamp {
			return fmt.Errorf("current block time %v is greater than maximum timestamp: %v", currentBlockTime, maxTimestamp)
		}
	}

	return nil
}

// Body is a simple (mutable, non-safe) data container for storing and moving
// a block's data contents (transactions and uncles) together.
type Body struct {
	Transactions []*Transaction
	Uncles       []*Header
	Withdrawals  []*Withdrawal `rlp:"optional"`
}

// Block represents an Zenanet block.
//
// Note the Block type tries to be 'immutable', and contains certain caches that rely
// on that. The rules around block immutability are as follows:
//
//   - We copy all data when the block is constructed. This makes references held inside
//     the block independent of whatever value was passed in.
//
//   - We copy all header data on access. This is because any change to the header would mess
//     up the cached hash and size values in the block. Calling code is expected to take
//     advantage of this to avoid over-allocating!
//
//   - When new body data is attached to the block, a shallow copy of the block is returned.
//     This ensures block modifications are race-free.
//
//   - We do not copy body data on access because it does not affect the caches, and also
//     because it would be too expensive.
type Block struct {
	header       *Header
	uncles       []*Header
	transactions Transactions
	withdrawals  Withdrawals

	// caches
	hash atomic.Pointer[common.Hash]
	size atomic.Uint64

	// These fields are used by package eth to track
	// inter-peer block relay.
	ReceivedAt   time.Time
	ReceivedFrom interface{}
	AnnouncedAt  *time.Time
}

// "external" block encoding. used for eth protocol, etc.
type extblock struct {
	Header      *Header
	Txs         []*Transaction
	Uncles      []*Header
	Withdrawals []*Withdrawal `rlp:"optional"`
}

// NewBlock creates a new block. The input data is copied, changes to header and to the
// field values will not affect the block.
//
// The body elements and the receipts are used to recompute and overwrite the
// relevant portions of the header.
func NewBlock(header *Header, body *Body, receipts []*Receipt, hasher TrieHasher) *Block {
	if body == nil {
		body = &Body{}
	}
	var (
		b           = NewBlockWithHeader(header)
		txs         = body.Transactions
		uncles      = body.Uncles
		withdrawals = body.Withdrawals
	)

	if len(txs) == 0 {
		b.header.TxHash = EmptyTxsHash
	} else {
		b.header.TxHash = DeriveSha(Transactions(txs), hasher)
		b.transactions = make(Transactions, len(txs))
		copy(b.transactions, txs)
	}

	if len(receipts) == 0 {
		b.header.ReceiptHash = EmptyReceiptsHash
	} else {
		b.header.ReceiptHash = DeriveSha(Receipts(receipts), hasher)
		b.header.Bloom = CreateBloom(receipts)
	}

	if len(uncles) == 0 {
		b.header.UncleHash = EmptyUncleHash
	} else {
		b.header.UncleHash = CalcUncleHash(uncles)
		b.uncles = make([]*Header, len(uncles))

		for i := range uncles {
			b.uncles[i] = CopyHeader(uncles[i])
		}
	}

	if withdrawals == nil {
		b.header.WithdrawalsHash = nil
	} else if len(withdrawals) == 0 {
		b.header.WithdrawalsHash = &EmptyWithdrawalsHash
		b.withdrawals = Withdrawals{}
	} else {
		hash := DeriveSha(Withdrawals(withdrawals), hasher)
		b.header.WithdrawalsHash = &hash
		b.withdrawals = slices.Clone(withdrawals)
	}

	return b
}

// CopyHeader creates a deep copy of a block header.
func CopyHeader(h *Header) *Header {
	cpy := *h
	if cpy.Difficulty = new(big.Int); h.Difficulty != nil {
		cpy.Difficulty.Set(h.Difficulty)
	}

	if cpy.Number = new(big.Int); h.Number != nil {
		cpy.Number.Set(h.Number)
	}

	if h.BaseFee != nil {
		cpy.BaseFee = new(big.Int).Set(h.BaseFee)
	}

	if len(h.Extra) > 0 {
		cpy.Extra = make([]byte, len(h.Extra))
		copy(cpy.Extra, h.Extra)
	}

	if h.WithdrawalsHash != nil {
		cpy.WithdrawalsHash = new(common.Hash)
		*cpy.WithdrawalsHash = *h.WithdrawalsHash
	}
	if h.ExcessBlobGas != nil {
		cpy.ExcessBlobGas = new(uint64)
		*cpy.ExcessBlobGas = *h.ExcessBlobGas
	}
	if h.BlobGasUsed != nil {
		cpy.BlobGasUsed = new(uint64)
		*cpy.BlobGasUsed = *h.BlobGasUsed
	}
	if h.ParentBeaconRoot != nil {
		cpy.ParentBeaconRoot = new(common.Hash)
		*cpy.ParentBeaconRoot = *h.ParentBeaconRoot
	}
	return &cpy
}

// DecodeRLP decodes a block from RLP.
func (b *Block) DecodeRLP(s *rlp.Stream) error {
	var eb extblock

	_, size, _ := s.Kind()

	if err := s.Decode(&eb); err != nil {
		return err
	}

	b.header, b.uncles, b.transactions, b.withdrawals = eb.Header, eb.Uncles, eb.Txs, eb.Withdrawals
	b.size.Store(rlp.ListSize(size))

	return nil
}

// EncodeRLP serializes a block as RLP.
func (b *Block) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &extblock{
		Header:      b.header,
		Txs:         b.transactions,
		Uncles:      b.uncles,
		Withdrawals: b.withdrawals,
	})
}

// Body returns the non-header content of the block.
// Note the returned data is not an independent copy.
func (b *Block) Body() *Body {
	return &Body{b.transactions, b.uncles, b.withdrawals}
}

// Accessors for body data. These do not return a copy because the content
// of the body slices does not affect the cached hash/size in block.

func (b *Block) Uncles() []*Header          { return b.uncles }
func (b *Block) Transactions() Transactions { return b.transactions }
func (b *Block) Withdrawals() Withdrawals   { return b.withdrawals }

func (b *Block) Transaction(hash common.Hash) *Transaction {
	for _, transaction := range b.transactions {
		if transaction.Hash() == hash {
			return transaction
		}
	}

	return nil
}

// Header returns the block header (as a copy).
func (b *Block) Header() *Header {
	return CopyHeader(b.header)
}

// Header value accessors. These do copy!

func (b *Block) Number() *big.Int     { return new(big.Int).Set(b.header.Number) }
func (b *Block) GasLimit() uint64     { return b.header.GasLimit }
func (b *Block) GasUsed() uint64      { return b.header.GasUsed }
func (b *Block) Difficulty() *big.Int { return new(big.Int).Set(b.header.Difficulty) }
func (b *Block) Time() uint64         { return b.header.Time }

func (b *Block) NumberU64() uint64        { return b.header.Number.Uint64() }
func (b *Block) MixDigest() common.Hash   { return b.header.MixDigest }
func (b *Block) Nonce() uint64            { return binary.BigEndian.Uint64(b.header.Nonce[:]) }
func (b *Block) Bloom() Bloom             { return b.header.Bloom }
func (b *Block) Coinbase() common.Address { return b.header.Coinbase }
func (b *Block) Root() common.Hash        { return b.header.Root }
func (b *Block) ParentHash() common.Hash  { return b.header.ParentHash }
func (b *Block) TxHash() common.Hash      { return b.header.TxHash }
func (b *Block) ReceiptHash() common.Hash { return b.header.ReceiptHash }
func (b *Block) UncleHash() common.Hash   { return b.header.UncleHash }
func (b *Block) Extra() []byte            { return common.CopyBytes(b.header.Extra) }

func (b *Block) GetTxDependency() [][]uint64 {
	if len(b.header.Extra) < ExtraVanityLength+ExtraSealLength {
		log.Error("length of extra less is than vanity and seal")
		return nil
	}

	var blockExtraData BlockExtraData
	if err := rlp.DecodeBytes(b.header.Extra[ExtraVanityLength:len(b.header.Extra)-ExtraSealLength], &blockExtraData); err != nil {
		log.Debug("error while decoding block extra data", "err", err)
		return nil
	}

	return blockExtraData.TxDependency
}

func (h *Header) GetValidatorBytes(chainConfig *params.ChainConfig) []byte {
	if !chainConfig.IsCancun(h.Number) {
		return h.Extra[ExtraVanityLength : len(h.Extra)-ExtraSealLength]
	}

	if len(h.Extra) < ExtraVanityLength+ExtraSealLength {
		log.Error("length of extra less is than vanity and seal")
		return nil
	}

	var blockExtraData BlockExtraData
	if err := rlp.DecodeBytes(h.Extra[ExtraVanityLength:len(h.Extra)-ExtraSealLength], &blockExtraData); err != nil {
		log.Debug("error while decoding block extra data", "err", err)
		return nil
	}

	return blockExtraData.ValidatorBytes
}

func (b *Block) BaseFee() *big.Int {
	if b.header.BaseFee == nil {
		return nil
	}

	return new(big.Int).Set(b.header.BaseFee)
}

func (b *Block) BeaconRoot() *common.Hash { return b.header.ParentBeaconRoot }

func (b *Block) ExcessBlobGas() *uint64 {
	var excessBlobGas *uint64
	if b.header.ExcessBlobGas != nil {
		excessBlobGas = new(uint64)
		*excessBlobGas = *b.header.ExcessBlobGas
	}
	return nil
}

func (b *Block) BlobGasUsed() *uint64 {
	var blobGasUsed *uint64
	if b.header.BlobGasUsed != nil {
		blobGasUsed = new(uint64)
		*blobGasUsed = *b.header.BlobGasUsed
	}
	return nil
}

// Size returns the true RLP encoded storage size of the block, either by encoding
// and returning it, or returning a previously cached value.
func (b *Block) Size() uint64 {
	if size := b.size.Load(); size > 0 {
		return size
	}

	c := writeCounter(0)
	rlp.Encode(&c, b)
	b.size.Store(uint64(c))

	return uint64(c)
}

// SanityCheck can be used to prevent that unbounded fields are
// stuffed with junk data to add processing overhead
func (b *Block) SanityCheck() error {
	return b.header.SanityCheck()
}

type writeCounter uint64

func (c *writeCounter) Write(b []byte) (int, error) {
	*c += writeCounter(len(b))
	return len(b), nil
}

func CalcUncleHash(uncles []*Header) common.Hash {
	if len(uncles) == 0 {
		return EmptyUncleHash
	}

	return rlpHash(uncles)
}

// NewBlockWithHeader creates a block with the given header data. The
// header data is copied, changes to header and to the field values
// will not affect the block.
func NewBlockWithHeader(header *Header) *Block {
	return &Block{header: CopyHeader(header)}
}

// WithSeal returns a new block with the data from b but the header replaced with
// the sealed one.
func (b *Block) WithSeal(header *Header) *Block {
	return &Block{
		header:       CopyHeader(header),
		transactions: b.transactions,
		uncles:       b.uncles,
		withdrawals:  b.withdrawals,
	}
}

// WithBody returns a new block with the original header and a deep copy of the
// provided body.
func (b *Block) WithBody(body Body) *Block {
	block := &Block{
		header:       b.header,
		transactions: slices.Clone(body.Transactions),
		uncles:       make([]*Header, len(body.Uncles)),
		withdrawals:  slices.Clone(body.Withdrawals),
	}
	for i := range body.Uncles {
		block.uncles[i] = CopyHeader(body.Uncles[i])
	}
	return block
}

// Hash returns the keccak256 hash of b's header.
// The hash is computed on the first call and cached thereafter.
func (b *Block) Hash() common.Hash {
	if hash := b.hash.Load(); hash != nil {
		return *hash
	}
	h := b.header.Hash()
	b.hash.Store(&h)
	return h
}

type Blocks []*Block

// HeaderParentHashFromRLP returns the parentHash of an RLP-encoded
// header. If 'header' is invalid, the zero hash is returned.
func HeaderParentHashFromRLP(header []byte) common.Hash {
	// parentHash is the first list element.
	listContent, _, err := rlp.SplitList(header)
	if err != nil {
		return common.Hash{}
	}

	parentHash, _, err := rlp.SplitString(listContent)
	if err != nil {
		return common.Hash{}
	}

	if len(parentHash) != 32 {
		return common.Hash{}
	}

	return common.BytesToHash(parentHash)
}

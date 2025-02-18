package miner

import (
	"errors"
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/zenanet-network/go-zenanet/common"
	"github.com/zenanet-network/go-zenanet/consensus"
	"github.com/zenanet-network/go-zenanet/consensus/eirene"
	"github.com/zenanet-network/go-zenanet/consensus/eirene/api"
	"github.com/zenanet-network/go-zenanet/consensus/eirene/validset"
	"github.com/zenanet-network/go-zenanet/core"
	"github.com/zenanet-network/go-zenanet/core/rawdb"
	"github.com/zenanet-network/go-zenanet/core/state"
	"github.com/zenanet-network/go-zenanet/core/txpool"
	"github.com/zenanet-network/go-zenanet/core/txpool/legacypool"
	"github.com/zenanet-network/go-zenanet/core/types"
	"github.com/zenanet-network/go-zenanet/core/vm"
	"github.com/zenanet-network/go-zenanet/ethdb"
	"github.com/zenanet-network/go-zenanet/ethdb/memorydb"
	"github.com/zenanet-network/go-zenanet/event"
	"github.com/zenanet-network/go-zenanet/params"
	"github.com/zenanet-network/go-zenanet/tests/eirene/mocks"
	"github.com/zenanet-network/go-zenanet/triedb"
)

type DefaultEireneMiner struct {
	Miner   *Miner
	Mux     *event.TypeMux //nolint:staticcheck
	Cleanup func(skipMiner bool)

	Ctrl               *gomock.Controller
	EthAPIMock         api.Caller
	HarmoniaClientMock eirene.IHarmoniaClient
	ContractMock       eirene.GenesisContract
}

func NewEireneDefaultMiner(t *testing.T) *DefaultEireneMiner {
	t.Helper()

	ctrl := gomock.NewController(t)

	ethAPI := api.NewMockCaller(ctrl)
	ethAPI.EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	spanner := eirene.NewMockSpanner(ctrl)
	spanner.EXPECT().GetCurrentValidatorsByHash(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*validset.Validator{
		{
			ID:               0,
			Address:          common.Address{0x1},
			VotingPower:      100,
			ProposerPriority: 0,
		},
	}, nil).AnyTimes()

	harmoniaClient := mocks.NewMockIHarmoniaClient(ctrl)
	harmoniaClient.EXPECT().Close().Times(1)

	genesisContracts := eirene.NewMockGenesisContract(ctrl)

	miner, mux, cleanup := createEireneMiner(t, ethAPI, spanner, harmoniaClient, genesisContracts)

	return &DefaultEireneMiner{
		Miner:              miner,
		Mux:                mux,
		Cleanup:            cleanup,
		Ctrl:               ctrl,
		EthAPIMock:         ethAPI,
		HarmoniaClientMock: harmoniaClient,
		ContractMock:       genesisContracts,
	}
}

// //nolint:staticcheck
func createEireneMiner(t *testing.T, ethAPIMock api.Caller, spanner eirene.Spanner, HarmoniaClientMock eirene.IHarmoniaClient, contractMock eirene.GenesisContract) (*Miner, *event.TypeMux, func(skipMiner bool)) {
	t.Helper()

	// Create Ethash config
	chainDB, genspec, chainConfig := NewDBForFakes(t)

	engine := NewFakeEirene(t, chainDB, chainConfig, ethAPIMock, spanner, HarmoniaClientMock, contractMock)

	// Create Zenanet backend
	bc, err := core.NewBlockChain(chainDB, nil, genspec, nil, engine, vm.Config{}, nil, nil, nil)
	if err != nil {
		t.Fatalf("can't create new chain %v", err)
	}

	statedb, _ := state.New(common.Hash{}, state.NewDatabase(chainDB), nil)
	blockchain := &testBlockChainEirene{chainConfig, statedb, 10000000, new(event.Feed)}

	pool := legacypool.New(testTxPoolConfigEirene, blockchain)
	txpool, _ := txpool.New(testTxPoolConfigEirene.PriceLimit, blockchain, []txpool.SubPool{pool})

	backend := NewMockBackendEirene(bc, txpool)

	// Create event Mux
	mux := new(event.TypeMux)

	config := Config{
		Zenbase: common.HexToAddress("123456789"),
	}

	// Create Miner
	miner := New(backend, &config, chainConfig, mux, engine, nil)

	cleanup := func(skipMiner bool) {
		bc.Stop()
		engine.Close()

		if !skipMiner {
			miner.Close()
		}
	}

	return miner, mux, cleanup
}

type TensingObject interface {
	Helper()
	Fatalf(format string, args ...any)
}

func NewDBForFakes(t TensingObject) (ethdb.Database, *core.Genesis, *params.ChainConfig) {
	t.Helper()

	memdb := memorydb.New()
	chainDB := rawdb.NewDatabase(memdb)
	addr := common.HexToAddress("12345")
	genesis := core.DeveloperGenesisBlock(11_500_000, &addr)

	chainConfig, _, err := core.SetupGenesisBlock(chainDB, triedb.NewDatabase(chainDB, triedb.HashDefaults), genesis)
	if err != nil {
		t.Fatalf("can't create new chain config: %v", err)
	}

	chainConfig.Eirene.Period = map[string]uint64{
		"0": 1,
	}
	chainConfig.Eirene.Sprint = map[string]uint64{
		"0": 64,
	}

	return chainDB, genesis, chainConfig
}

func NewFakeEirene(t TensingObject, chainDB ethdb.Database, chainConfig *params.ChainConfig, ethAPIMock api.Caller, spanner eirene.Spanner, harmoniaClientMock eirene.IHarmoniaClient, contractMock eirene.GenesisContract) consensus.Engine {
	t.Helper()

	if chainConfig.Eirene == nil {
		chainConfig.Eirene = params.EireneUnittestChainConfig.Eirene
	}

	return eirene.New(chainConfig, chainDB, ethAPIMock, spanner, harmoniaClientMock, contractMock, false)
}

var (
	// Test chain configurations
	testTxPoolConfigEirene legacypool.Config
)

// TODO - Arpit, Duplicate Functions
type mockBackendEirene struct {
	bc     *core.BlockChain
	txPool *txpool.TxPool
}

func NewMockBackendEirene(bc *core.BlockChain, txPool *txpool.TxPool) *mockBackendEirene {
	return &mockBackendEirene{
		bc:     bc,
		txPool: txPool,
	}
}

func (m *mockBackendEirene) BlockChain() *core.BlockChain {
	return m.bc
}

// PeerCount implements Backend.
func (*mockBackendEirene) PeerCount() int {
	panic("unimplemented")
}

func (m *mockBackendEirene) TxPool() *txpool.TxPool {
	return m.txPool
}

func (m *mockBackendEirene) StateAtBlock(block *types.Block, reexec uint64, base *state.StateDB, checkLive bool, preferDisk bool) (statedb *state.StateDB, err error) {
	return nil, errors.New("not supported")
}

// TODO - Arpit, Duplicate Functions
type testBlockChainEirene struct {
	config        *params.ChainConfig
	statedb       *state.StateDB
	gasLimit      uint64
	chainHeadFeed *event.Feed
}

func (bc *testBlockChainEirene) Config() *params.ChainConfig {
	return bc.config
}

func (bc *testBlockChainEirene) CurrentBlock() *types.Header {
	return &types.Header{
		Number:   new(big.Int),
		GasLimit: bc.gasLimit,
	}
}

func (bc *testBlockChainEirene) GetBlock(hash common.Hash, number uint64) *types.Block {
	return types.NewBlock(bc.CurrentBlock(), nil, nil, nil)
}

func (bc *testBlockChainEirene) StateAt(common.Hash) (*state.StateDB, error) {
	return bc.statedb, nil
}

func (bc *testBlockChainEirene) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return bc.chainHeadFeed.Subscribe(ch)
}

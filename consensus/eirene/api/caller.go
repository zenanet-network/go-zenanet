package api

import (
	"context"

	"github.com/zenanet-network/go-zenanet/common/hexutil"
	"github.com/zenanet-network/go-zenanet/core/state"
	"github.com/zenanet-network/go-zenanet/internal/ethapi"
	"github.com/zenanet-network/go-zenanet/rpc"
)

//go:generate mockgen -destination=./caller_mock.go -package=api . Caller
type Caller interface {
	Call(ctx context.Context, args ethapi.TransactionArgs, blockNrOrHash *rpc.BlockNumberOrHash, overrides *ethapi.StateOverride, blockOverrides *ethapi.BlockOverrides) (hexutil.Bytes, error)
	CallWithState(ctx context.Context, args ethapi.TransactionArgs, blockNrOrHash *rpc.BlockNumberOrHash, state *state.StateDB, overrides *ethapi.StateOverride, blockOverrides *ethapi.BlockOverrides) (hexutil.Bytes, error)
}

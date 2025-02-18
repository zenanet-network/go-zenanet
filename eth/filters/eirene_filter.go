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

package filters

import (
	"context"

	"github.com/zenanet-network/go-zenanet/common"
	"github.com/zenanet-network/go-zenanet/core/types"
	"github.com/zenanet-network/go-zenanet/ethdb"
	"github.com/zenanet-network/go-zenanet/params"
	"github.com/zenanet-network/go-zenanet/rpc"
)

// EireneBlockLogsFilter can be used to retrieve and filter logs.
type EireneBlockLogsFilter struct {
	backend   Backend
	eireneConfig *params.EireneConfig

	db        ethdb.Database
	addresses []common.Address
	topics    [][]common.Hash

	block      common.Hash // Block hash if filtering a single block
	begin, end int64       // Range interval if filtering multiple blocks
}

// NewEireneBlockLogsRangeFilter creates a new filter which uses a bloom filter on blocks to
// figure out whether a particular block is interesting or not.
func NewEireneBlockLogsRangeFilter(backend Backend, eireneConfig *params.EireneConfig, begin, end int64, addresses []common.Address, topics [][]common.Hash) *EireneBlockLogsFilter {
	// Create a generic filter and convert it into a range filter
	filter := newEireneBlockLogsFilter(backend, eireneConfig, addresses, topics)
	filter.begin = begin
	filter.end = end

	return filter
}

// NewEireneBlockLogsFilter creates a new filter which directly inspects the contents of
// a block to figure out whether it is interesting or not.
func NewEireneBlockLogsFilter(backend Backend, eireneConfig *params.EireneConfig, block common.Hash, addresses []common.Address, topics [][]common.Hash) *EireneBlockLogsFilter {
	// Create a generic filter and convert it into a block filter
	filter := newEireneBlockLogsFilter(backend, eireneConfig, addresses, topics)
	filter.block = block

	return filter
}

// newEireneBlockLogsFilter creates a generic filter that can either filter based on a block hash,
// or based on range queries. The search criteria needs to be explicitly set.
func newEireneBlockLogsFilter(backend Backend, eireneConfig *params.EireneConfig, addresses []common.Address, topics [][]common.Hash) *EireneBlockLogsFilter {
	return &EireneBlockLogsFilter{
		backend:   backend,
		eireneConfig: eireneConfig,
		addresses: addresses,
		topics:    topics,
		db:        backend.ChainDb(),
	}
}

// Logs searches the blockchain for matching log entries, returning all from the
// first block that contains matches, updating the start of the filter accordingly.
func (f *EireneBlockLogsFilter) Logs(ctx context.Context) ([]*types.Log, error) {
	// If we're doing singleton block filtering, execute and return
	if f.block != (common.Hash{}) {
		receipt, _ := f.backend.GetEireneBlockReceipt(ctx, f.block)
		if receipt == nil {
			return nil, nil
		}

		return f.eireneBlockLogs(ctx, receipt)
	}

	// Figure out the limits of the filter range
	header, _ := f.backend.HeaderByNumber(ctx, rpc.LatestBlockNumber)
	if header == nil {
		return nil, nil
	}

	head := header.Number.Uint64()

	if f.begin == -1 {
		f.begin = int64(head)
	}

	// adjust begin for sprint
	f.begin = currentSprintEnd(f.eireneConfig.CalculateSprint(uint64(f.begin)), f.begin)

	end := f.end
	if f.end == -1 {
		end = int64(head)
	}

	// Gather all indexed logs, and finish with non indexed ones
	return f.unindexedLogs(ctx, uint64(end))
}

// unindexedLogs returns the logs matching the filter criteria based on raw block
// iteration and bloom matching.
func (f *EireneBlockLogsFilter) unindexedLogs(ctx context.Context, end uint64) ([]*types.Log, error) {
	var logs []*types.Log

	sprintLength := f.eireneConfig.CalculateSprint(uint64(f.begin))

	for ; f.begin <= int64(end); f.begin = f.begin + int64(sprintLength) {
		header, err := f.backend.HeaderByNumber(ctx, rpc.BlockNumber(f.begin))
		if header == nil || err != nil {
			return logs, err
		}

		// get eirene block receipt
		receipt, err := f.backend.GetEireneBlockReceipt(ctx, header.Hash())
		if receipt == nil || err != nil {
			continue
		}

		// filter eirene block logs
		found, err := f.eireneBlockLogs(ctx, receipt)
		if err != nil {
			return logs, err
		}

		logs = append(logs, found...)
		sprintLength = f.eireneConfig.CalculateSprint(uint64(f.begin))
	}

	return logs, nil
}

// eireneBlockLogs returns the logs matching the filter criteria within a single block.
func (f *EireneBlockLogsFilter) eireneBlockLogs(ctx context.Context, receipt *types.Receipt) (logs []*types.Log, err error) {
	if bloomFilter(receipt.Bloom, f.addresses, f.topics) {
		logs = filterLogs(receipt.Logs, nil, nil, f.addresses, f.topics)
	}

	return logs, nil
}

func currentSprintEnd(sprint uint64, n int64) int64 {
	m := n % int64(sprint)
	if m == 0 {
		return n
	}

	return n + int64(sprint) - m
}

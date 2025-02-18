// Copyright 2021 The go-zenanet Authors
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

package tracetest

import (
	"context"
	"encoding/json"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zenanet-network/go-zenanet/common"
	"github.com/zenanet-network/go-zenanet/core"
	"github.com/zenanet-network/go-zenanet/core/rawdb"
	"github.com/zenanet-network/go-zenanet/core/types"
	"github.com/zenanet-network/go-zenanet/core/vm"
	"github.com/zenanet-network/go-zenanet/eth/tracers"
	"github.com/zenanet-network/go-zenanet/tests"
)

// prestateTrace is the result of a prestateTrace run.
type prestateTrace = map[common.Address]*account

type account struct {
	Balance string                      `json:"balance"`
	Code    string                      `json:"code"`
	Nonce   uint64                      `json:"nonce"`
	Storage map[common.Hash]common.Hash `json:"storage"`
}

// testcase defines a single test to check the stateDiff tracer against.
type testcase struct {
	Genesis      *core.Genesis   `json:"genesis"`
	Context      *callContext    `json:"context"`
	Input        string          `json:"input"`
	TracerConfig json.RawMessage `json:"tracerConfig"`
	Result       interface{}     `json:"result"`
}

func TestPrestateTracerLegacy(t *testing.T) {
	t.Parallel()
	testPrestateDiffTracer(t, "prestateTracerLegacy", "prestate_tracer_legacy")
}

func TestPrestateTracer(t *testing.T) {
	t.Parallel()
	testPrestateDiffTracer(t, "prestateTracer", "prestate_tracer")
}

func TestPrestateWithDiffModeTracer(t *testing.T) {
	t.Parallel()
	testPrestateDiffTracer(t, "prestateTracer", "prestate_tracer_with_diff_mode")
}

// nolint:gocognit
func testPrestateDiffTracer(t *testing.T, tracerName string, dirPath string) {
	t.Helper()

	files, err := os.ReadDir(filepath.Join("testdata", dirPath))
	if err != nil {
		t.Fatalf("failed to retrieve tracer test suite: %v", err)
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		file := file // capture range variable
		t.Run(camel(strings.TrimSuffix(file.Name(), ".json")), func(t *testing.T) {
			t.Parallel()

			var (
				test = new(testcase)
				tx   = new(types.Transaction)
			)
			// Call tracer test found, read if from disk
			if blob, err := os.ReadFile(filepath.Join("testdata", dirPath, file.Name())); err != nil {
				t.Fatalf("failed to read testcase: %v", err)
			} else if err := json.Unmarshal(blob, test); err != nil {
				t.Fatalf("failed to parse testcase: %v", err)
			}

			if err := tx.UnmarshalBinary(common.FromHex(test.Input)); err != nil {
				t.Fatalf("failed to parse testcase input: %v", err)
			}
			// Configure a blockchain with the given prestate
			var (
				signer       = types.MakeSigner(test.Genesis.Config, new(big.Int).SetUint64(uint64(test.Context.Number)), uint64(test.Context.Time))
				blockContext = test.Context.toBlockContext(test.Genesis)
				state        = tests.MakePreState(rawdb.NewMemoryDatabase(), test.Genesis.Alloc, false, rawdb.HashScheme)
			)
			defer state.Close()

			tracer, err := tracers.DefaultDirectory.New(tracerName, new(tracers.Context), test.TracerConfig)
			if err != nil {
				t.Fatalf("failed to create call tracer: %v", err)
			}

			state.StateDB.SetLogger(tracer.Hooks)
			msg, err := core.TransactionToMessage(tx, signer, blockContext.BaseFee)
			if err != nil {
				t.Fatalf("failed to prepare transaction for tracing: %v", err)
			}
			evm := vm.NewEVM(blockContext, core.NewEVMTxContext(msg), state.StateDB, test.Genesis.Config, vm.Config{Tracer: tracer.Hooks})
			tracer.OnTxStart(evm.GetVMContext(), tx, msg.From)
			vmRet, err := core.ApplyMessage(evm, msg, new(core.GasPool).AddGas(tx.Gas()), context.Background())
			if err != nil {
				t.Fatalf("failed to execute transaction: %v", err)
			}
			tracer.OnTxEnd(&types.Receipt{GasUsed: vmRet.UsedGas}, nil)
			// Retrieve the trace result and compare against the expected
			res, err := tracer.GetResult()
			if err != nil {
				t.Fatalf("failed to retrieve trace result: %v", err)
			}
			// The legacy javascript calltracer marshals json in js, which
			// is not deterministic (as opposed to the golang json encoder).
			if strings.HasSuffix(dirPath, "_legacy") {
				// This is a tweak to make it deterministic. Can be removed when
				// we remove the legacy tracer.
				var x prestateTrace

				_ = json.Unmarshal(res, &x)
				res, _ = json.Marshal(x)
			}

			want, err := json.Marshal(test.Result)
			if err != nil {
				t.Fatalf("failed to marshal test: %v", err)
			}

			if string(want) != string(res) {
				t.Fatalf("trace mismatch\n have: %v\n want: %v\n", string(res), string(want))
			}
		})
	}
}

// Copyright 2016 The go-zenanet Authors
// This file is part of go-zenanet.
//
// go-zenanet is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-zenanet is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-zenanet. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/zenanet-network/go-zenanet/internal/cmdtest"
	"github.com/zenanet-network/go-zenanet/internal/reexec"
	"github.com/zenanet-network/go-zenanet/rpc"
)

type testgze struct {
	*cmdtest.TestCmd

	// template variables for expect
	Datadir   string
	Zenabase string
}

func init() {
	// Run the app if we've been exec'd as "gze-test" in runGze.
	reexec.Register("gze-test", func() {
		if err := app.Run(os.Args); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Exit(0)
	})
}

func TestMain(m *testing.M) {
	// check if we have been reexec'd
	if reexec.Init() {
		return
	}
	os.Exit(m.Run())
}

func initGze(t *testing.T) string {
	args := []string{"--networkid=42", "init", "./testdata/clique.json"}
	t.Logf("Initializing gze: %v ", args)
	g := runGze(t, args...)
	datadir := g.Datadir
	g.WaitExit()
	return datadir
}

// spawns gze with the given command line args. If the args don't set --datadir, the
// child g gets a temporary data directory.
func runGze(t *testing.T, args ...string) *testgze {
	tt := &testgze{}
	tt.TestCmd = cmdtest.NewTestCmd(t, tt)
	for i, arg := range args {
		switch arg {
		case "--datadir":
			if i < len(args)-1 {
				tt.Datadir = args[i+1]
			}
		case "--miner.zenabase":
			if i < len(args)-1 {
				tt.Zenabase = args[i+1]
			}
		}
	}
	if tt.Datadir == "" {
		// The temporary datadir will be removed automatically if something fails below.
		tt.Datadir = t.TempDir()
		args = append([]string{"--datadir", tt.Datadir}, args...)
	}

	// Boot "gze". This actually runs the test binary but the TestMain
	// function will prevent any tests from running.
	tt.Run("gze-test", args...)

	return tt
}

// waitForEndpoint attempts to connect to an RPC endpoint until it succeeds.
func waitForEndpoint(t *testing.T, endpoint string, timeout time.Duration) {
	probe := func() bool {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		c, err := rpc.DialContext(ctx, endpoint)
		if c != nil {
			_, err = c.SupportedModules()
			c.Close()
		}
		return err == nil
	}

	start := time.Now()
	for {
		if probe() {
			return
		}
		if time.Since(start) > timeout {
			t.Fatal("endpoint", endpoint, "did not open within", timeout)
		}
		time.Sleep(200 * time.Millisecond)
	}
}
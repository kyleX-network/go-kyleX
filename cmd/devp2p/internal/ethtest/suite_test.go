// Copyright 2021 The go-kyleX Authors
// This file is part of go-kyleX.
//
// go-kyleX is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-kyleX is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-kyleX. If not, see <http://www.gnu.org/licenses/>.

package ethtest

import (
	"os"
	"testing"
	"time"

	"github.com/kyleX-network/go-kyleX/eth"
	"github.com/kyleX-network/go-kyleX/eth/ethconfig"
	"github.com/kyleX-network/go-kyleX/internal/utesting"
	"github.com/kyleX-network/go-kyleX/node"
	"github.com/kyleX-network/go-kyleX/p2p"
)

var (
	genesisFile   = "./testdata/genesis.json"
	halfchainFile = "./testdata/halfchain.rlp"
	fullchainFile = "./testdata/chain.rlp"
)

func TestEthSuite(t *testing.T) {
	gkyle, err := runGkyle()
	if err != nil {
		t.Fatalf("could not run gkyle: %v", err)
	}
	defer gkyle.Close()

	suite, err := NewSuite(gkyle.Server().Self(), fullchainFile, genesisFile)
	if err != nil {
		t.Fatalf("could not create new test suite: %v", err)
	}
	for _, test := range suite.EthTests() {
		t.Run(test.Name, func(t *testing.T) {
			result := utesting.RunTAP([]utesting.Test{{Name: test.Name, Fn: test.Fn}}, os.Stdout)
			if result[0].Failed {
				t.Fatal()
			}
		})
	}
}

func TestSnapSuite(t *testing.T) {
	gkyle, err := runGkyle()
	if err != nil {
		t.Fatalf("could not run gkyle: %v", err)
	}
	defer gkyle.Close()

	suite, err := NewSuite(gkyle.Server().Self(), fullchainFile, genesisFile)
	if err != nil {
		t.Fatalf("could not create new test suite: %v", err)
	}
	for _, test := range suite.SnapTests() {
		t.Run(test.Name, func(t *testing.T) {
			result := utesting.RunTAP([]utesting.Test{{Name: test.Name, Fn: test.Fn}}, os.Stdout)
			if result[0].Failed {
				t.Fatal()
			}
		})
	}
}

// runGkyle creates and starts a gkyle node
func runGkyle() (*node.Node, error) {
	stack, err := node.New(&node.Config{
		P2P: p2p.Config{
			ListenAddr:  "127.0.0.1:0",
			NoDiscovery: true,
			MaxPeers:    10, // in case a test requires multiple connections, can be changed in the future
			NoDial:      true,
		},
	})
	if err != nil {
		return nil, err
	}

	err = setupGkyle(stack)
	if err != nil {
		stack.Close()
		return nil, err
	}
	if err = stack.Start(); err != nil {
		stack.Close()
		return nil, err
	}
	return stack, nil
}

func setupGkyle(stack *node.Node) error {
	chain, err := loadChain(halfchainFile, genesisFile)
	if err != nil {
		return err
	}

	backend, err := eth.New(stack, &ethconfig.Config{
		Genesis:                 &chain.genesis,
		NetworkId:               chain.genesis.Config.ChainID.Uint64(), // 19763
		DatabaseCache:           10,
		TrieCleanCache:          10,
		TrieCleanCacheJournal:   "",
		TrieCleanCacheRejournal: 60 * time.Minute,
		TrieDirtyCache:          16,
		TrieTimeout:             60 * time.Minute,
		SnapshotCache:           10,
	})
	if err != nil {
		return err
	}

	_, err = backend.BlockChain().InsertChain(chain.blocks[1:])
	return err
}

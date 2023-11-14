package main

import (
	"context"
	"fmt"

	"github.com/zeta-chain/zetacore/common"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func (sm *SmokeTest) TestBlockHeaders() {
	LoudPrintf("TESTING BLOCK HEADERS\n")
	// test ethereum block headers; should have a chain
	checkBlock := func(chainID int64) {
		bhs, err := sm.observerClient.GetBlockHeaderStateByChain(context.TODO(), &observertypes.QueryGetBlockHeaderStateRequest{
			ChainId: chainID,
		})
		if err != nil {
			panic(err)
		}
		if bhs == nil || bhs.BlockHeaderState == nil {
			panic("no block header state")
		}
		earliestBlock := bhs.BlockHeaderState.EarliestHeight
		latestBlock := bhs.BlockHeaderState.LatestHeight
		if earliestBlock == 0 || latestBlock == earliestBlock {
			panic("no blocks")
		}
		latestBlockHash := bhs.BlockHeaderState.LatestBlockHash
		fmt.Printf("CHAIN %d: starting tracing back blocks; latest block %d\n", chainID, latestBlock)
		bn := latestBlock
		currentHash := latestBlockHash
		for {
			bhres, err := sm.observerClient.GetBlockHeaderByHash(context.TODO(), &observertypes.QueryGetBlockHeaderByHashRequest{
				BlockHash: currentHash,
			})
			if err != nil {
				fmt.Printf("cannot getting block header; tracing stops: %v\n", err)
				break
			}
			bn = bhres.BlockHeader.Height - 1
			currentHash = bhres.BlockHeader.ParentHash
			//fmt.Printf("found block header %d\n", bhres.BlockHeader.Height)
		}
		if bn > earliestBlock {
			panic(fmt.Sprintf("block header tracing failed; expected at most %d, got %d", earliestBlock, bn))
		}
		fmt.Printf("block header tracing succeeded; expected at most %d, got %d\n", earliestBlock, bn)
	}
	checkBlock(common.GoerliChain().ChainId)
	checkBlock(common.BtcRegtestChain().ChainId)
}

package smoketests

import (
	"context"
	"fmt"

	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestBlockHeaders(sm *runner.SmokeTestRunner) {
	// test ethereum block headers; should have a chain
	checkBlock := func(chainID int64) {
		bhs, err := sm.ObserverClient.GetBlockHeaderStateByChain(context.TODO(), &observertypes.QueryGetBlockHeaderStateRequest{
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
		sm.Logger.Info("CHAIN %d: starting tracing back blocks; latest block %d", chainID, latestBlock)
		bn := latestBlock
		currentHash := latestBlockHash
		for {
			bhres, err := sm.ObserverClient.GetBlockHeaderByHash(context.TODO(), &observertypes.QueryGetBlockHeaderByHashRequest{
				BlockHash: currentHash,
			})
			if err != nil {
				sm.Logger.Info("cannot getting block header; tracing stops: %v", err)
				break
			}
			bn = bhres.BlockHeader.Height - 1
			currentHash = bhres.BlockHeader.ParentHash
		}
		if bn > earliestBlock {
			panic(fmt.Sprintf("block header tracing failed; expected at most %d, got %d", earliestBlock, bn))
		}
		sm.Logger.Info("block header tracing succeeded; expected at most %d, got %d", earliestBlock, bn)
	}
	checkBlock(common.GoerliLocalnetChain().ChainId)
	checkBlock(common.BtcRegtestChain().ChainId)
}

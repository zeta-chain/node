package observer

import (
	"context"
	"time"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/zetaclient/chains/solana/rpc"
	"github.com/zeta-chain/node/zetaclient/common"
)

// watchRPCStatus watches the RPC status of the Solana chain
func (ob *Observer) watchRPCStatus(ctx context.Context) error {
	ob.Logger().Chain.Info().Msgf("watchRPCStatus started for chain %d", ob.Chain().ChainId)

	ticker := time.NewTicker(common.RPCStatusCheckInterval)
	for {
		select {
		case <-ticker.C:
			if !ob.GetChainParams().IsSupported {
				continue
			}

			ob.checkRPCStatus(ctx)
		case <-ob.StopChannel():
			return nil
		}
	}
}

// checkRPCStatus checks the RPC status of the Solana chain
func (ob *Observer) checkRPCStatus(ctx context.Context) {
	// Solana privnet doesn't have RPC 'GetHealth', need to differentiate
	privnet := ob.Chain().NetworkType == chains.NetworkType_privnet

	// check the RPC status
	blockTime, err := rpc.CheckRPCStatus(ctx, ob.solClient, privnet)
	if err != nil {
		ob.Logger().Chain.Error().Err(err).Msg("CheckRPCStatus failed")
		return
	}

	// alert if RPC latency is too high
	ob.AlertOnRPCLatency(blockTime, rpc.RPCAlertLatency)
}

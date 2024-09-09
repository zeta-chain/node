// Package observer implements the EVM chain observer
package observer

import (
	"context"
	"time"

	"github.com/zeta-chain/node/zetaclient/chains/evm/rpc"
	"github.com/zeta-chain/node/zetaclient/common"
)

// watchRPCStatus watches the RPC status of the EVM chain
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

// checkRPCStatus checks the RPC status of the EVM chain
func (ob *Observer) checkRPCStatus(ctx context.Context) {
	blockTime, err := rpc.CheckRPCStatus(ctx, ob.evmClient)
	if err != nil {
		ob.Logger().Chain.Error().Err(err).Msg("CheckRPCStatus failed")
		return
	}

	// alert if RPC latency is too high
	ob.AlertOnRPCLatency(blockTime, rpc.RPCAlertLatency)
}

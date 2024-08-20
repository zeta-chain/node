// Package observer implements the EVM chain observer
package observer

import (
	"context"
	"time"

	"github.com/zeta-chain/zetacore/zetaclient/chains/evm/rpc"
	"github.com/zeta-chain/zetacore/zetaclient/common"
)

// WatchRPCStatus watches the RPC status of the evm chain
func (ob *Observer) WatchRPCStatus(ctx context.Context) error {
	ob.Logger().Chain.Info().Msgf("WatchRPCStatus started for chain %d", ob.Chain().ChainId)

	ticker := time.NewTicker(common.RPCStatusCheckInterval)
	for {
		select {
		case <-ticker.C:
			if !ob.GetChainParams().IsSupported {
				continue
			}

			alertLatency := ob.RPCAlertLatency()
			err := rpc.CheckRPCStatus(ctx, ob.evmClient, alertLatency, ob.Logger().Chain)
			if err != nil {
				ob.Logger().Chain.Error().Err(err).Msg("RPC Status error")
			}
		case <-ob.StopChannel():
			return nil
		}
	}
}

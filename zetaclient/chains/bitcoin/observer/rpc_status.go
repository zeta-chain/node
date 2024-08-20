package observer

import (
	"context"
	"time"

	"github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin/rpc"
	"github.com/zeta-chain/zetacore/zetaclient/common"
)

// WatchRPCStatus watches the RPC status of the Bitcoin chain
func (ob *Observer) WatchRPCStatus(_ context.Context) error {
	ob.Logger().Chain.Info().Msgf("WatchRPCStatus started for chain %d", ob.Chain().ChainId)

	ticker := time.NewTicker(common.RPCStatusCheckInterval)
	for {
		select {
		case <-ticker.C:
			if !ob.GetChainParams().IsSupported {
				continue
			}

			alertLatency := ob.RPCAlertLatency()
			tssAddress := ob.TSS().BTCAddressWitnessPubkeyHash()
			err := rpc.CheckRPCStatus(ob.btcClient, alertLatency, tssAddress, ob.Logger().Chain)
			if err != nil {
				ob.Logger().Chain.Error().Err(err).Msg("RPC Status error")
			}

		case <-ob.StopChannel():
			return nil
		}
	}
}

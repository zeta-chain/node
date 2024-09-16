package observer

import (
	"context"
	"time"

	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/rpc"
	"github.com/zeta-chain/node/zetaclient/common"
)

// watchRPCStatus watches the RPC status of the Bitcoin chain
func (ob *Observer) watchRPCStatus(_ context.Context) error {
	ob.Logger().Chain.Info().Msgf("WatchRPCStatus started for chain %d", ob.Chain().ChainId)

	ticker := time.NewTicker(common.RPCStatusCheckInterval)
	for {
		select {
		case <-ticker.C:
			if !ob.GetChainParams().IsSupported {
				continue
			}

			ob.checkRPCStatus()
		case <-ob.StopChannel():
			return nil
		}
	}
}

// checkRPCStatus checks the RPC status of the Bitcoin chain
func (ob *Observer) checkRPCStatus() {
	tssAddress := ob.TSS().BTCAddressWitnessPubkeyHash()
	blockTime, err := rpc.CheckRPCStatus(ob.btcClient, tssAddress)
	if err != nil {
		ob.Logger().Chain.Error().Err(err).Msg("CheckRPCStatus failed")
		return
	}

	// alert if RPC latency is too high
	ob.AlertOnRPCLatency(blockTime, rpc.RPCAlertLatency)
}

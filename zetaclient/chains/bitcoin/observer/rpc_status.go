package observer

import (
	"context"

	"github.com/pkg/errors"

	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/rpc"
)

// CheckRPCStatus checks the RPC status of the Bitcoin chain
func (ob *Observer) CheckRPCStatus(_ context.Context) error {
	tssAddress, err := ob.TSS().PubKey().AddressBTC(ob.Chain().ChainId)
	if err != nil {
		return errors.Wrap(err, "unable to get TSS BTC address")
	}

	blockTime, err := rpc.CheckRPCStatus(ob.btcClient, tssAddress)
	if err != nil {
		return errors.Wrap(err, "unable to check RPC status")
	}

	ob.ReportBlockLatency(blockTime)

	return nil
}

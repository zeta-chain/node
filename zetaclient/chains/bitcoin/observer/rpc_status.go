package observer

import (
	"context"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/zetaclient/metrics"
)

// CheckRPCStatus checks the RPC status of the Bitcoin chain
func (ob *Observer) CheckRPCStatus(ctx context.Context) error {
	if !ob.isNodeEnabled() {
		return nil
	}

	// 1. Query last block timestamp
	blockTime, err := ob.rpc.Healthcheck(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to check rpc health")
	}

	metrics.ReportBlockLatency(ob.Chain().Name, blockTime)

	// 2. Query utxos owned by TSS address
	// This is to ensure that the Bitcoin node is configured to watch TSS address
	tssAddress, err := ob.TSS().PubKey().AddressBTC(ob.Chain().ChainId)
	if err != nil {
		return errors.Wrap(err, "unable to get TSS BTC address")
	}

	res, err := ob.rpc.ListUnspentMinMaxAddresses(ctx, 0, 1000000, []btcutil.Address{tssAddress})
	switch {
	case err != nil:
		return errors.Wrap(err, "unable to list TSS UTXOs")
	case len(res) == 0:
		return errors.New("no UTXOs found for TSS")
	}

	return nil
}

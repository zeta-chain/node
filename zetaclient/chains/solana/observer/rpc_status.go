package observer

import (
	"context"

	"cosmossdk.io/errors"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/zetaclient/chains/solana/rpc"
)

// checkRPCStatus checks the RPC status of the Solana chain
func (ob *Observer) CheckRPCStatus(ctx context.Context) error {
	// Solana privnet doesn't have RPC 'GetHealth', need to differentiate
	privnet := ob.Chain().NetworkType == chains.NetworkType_privnet

	// check the RPC status
	blockTime, err := rpc.CheckRPCStatus(ctx, ob.solClient, privnet)
	if err != nil {
		return errors.Wrap(err, "unable to check rpc status")
	}

	ob.ReportBlockLatency(blockTime)

	return nil
}

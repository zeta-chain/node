package observer

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/client"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
)

// PostGasPrice posts gas price to zetacore
func (ob *Observer) PostGasPrice(ctx context.Context) error {
	var (
		err              error
		feeRateEstimated int64
	)

	// special handle regnet and testnet gas rate
	if ob.Chain().NetworkType != chains.NetworkType_mainnet {
		feeRateEstimated, err = ob.GetFeeRateForRegnetAndTestnet(ctx)
		if err != nil {
			return errors.Wrap(err, "unable to execute specialHandleFeeRate")
		}
	} else {
		feeRateEstimated, err = ob.rpc.GetEstimatedFeeRate(ctx, 1, false)
		if err != nil {
			return errors.Wrap(err, "unable to get estimated fee rate")
		}
	}

	// query the current block number
	blockNumber, err := ob.rpc.GetBlockCount(ctx)
	if err != nil {
		return errors.Wrap(err, "GetBlockCount error")
	}

	// Bitcoin has no concept of priority fee (like eth)
	const priorityFee = 0

	// #nosec G115 always positive
	_, err = ob.ZetacoreClient().
		PostVoteGasPrice(ctx, ob.Chain(), uint64(feeRateEstimated), priorityFee, uint64(blockNumber))
	if err != nil {
		return errors.Wrap(err, "PostVoteGasPrice error")
	}

	return nil
}

// GetFeeRateForRegnetAndTestnet handles the fee rate for regnet and testnet
// regnet:  RPC 'EstimateSmartFee' is not available
// testnet: RPC 'EstimateSmartFee' can return unreasonable high fee rate
func (ob *Observer) GetFeeRateForRegnetAndTestnet(ctx context.Context) (int64, error) {
	switch ob.Chain().NetworkType {
	case chains.NetworkType_privnet:
		return client.FeeRateRegnet, nil
	case chains.NetworkType_testnet:
		feeRateEstimated, err := common.GetRecentFeeRate(ctx, ob.rpc, ob.netParams)
		if err != nil {
			return 0, errors.Wrapf(err, "error GetRecentFeeRate")
		}
		return feeRateEstimated, nil
	default:
		return 0, fmt.Errorf(" unsupported bitcoin network type %d", ob.Chain().NetworkType)
	}
}

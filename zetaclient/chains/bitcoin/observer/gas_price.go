package observer

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/client"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
	clienttypes "github.com/zeta-chain/node/zetaclient/types"
)

// WatchGasPrice watches Bitcoin chain for gas rate and post to zetacore
func (ob *Observer) WatchGasPrice(ctx context.Context) error {
	// report gas price right away as the ticker takes time to kick in
	err := ob.PostGasPrice(ctx)
	if err != nil {
		ob.logger.GasPrice.Error().Err(err).Msgf("PostGasPrice error for chain %d", ob.Chain().ChainId)
	}

	// start gas price ticker
	ticker, err := clienttypes.NewDynamicTicker("Bitcoin_WatchGasPrice", ob.ChainParams().GasPriceTicker)
	if err != nil {
		return errors.Wrapf(err, "NewDynamicTicker error")
	}
	ob.logger.GasPrice.Info().Msgf("WatchGasPrice started for chain %d with interval %d",
		ob.Chain().ChainId, ob.ChainParams().GasPriceTicker)

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C():
			if !ob.ChainParams().IsSupported {
				continue
			}
			err := ob.PostGasPrice(ctx)
			if err != nil {
				ob.logger.GasPrice.Error().Err(err).Msgf("PostGasPrice error for chain %d", ob.Chain().ChainId)
			}
			ticker.UpdateInterval(ob.ChainParams().GasPriceTicker, ob.logger.GasPrice)
		case <-ob.StopChannel():
			ob.logger.GasPrice.Info().Msgf("WatchGasPrice stopped for chain %d", ob.Chain().ChainId)
			return nil
		}
	}
}

// PostGasPrice posts gas price to zetacore
func (ob *Observer) PostGasPrice(ctx context.Context) error {
	var (
		err              error
		feeRateEstimated int64
	)

	// special handle regnet and testnet gas rate
	// regnet:  RPC 'EstimateSmartFee' is not available
	// testnet: RPC 'EstimateSmartFee' returns unreasonable high gas rate
	if ob.Chain().NetworkType != chains.NetworkType_mainnet {
		feeRateEstimated, err = ob.specialHandleFeeRate(ctx)
		if err != nil {
			return errors.Wrap(err, "unable to execute specialHandleFeeRate")
		}
	} else {
		isRegnet := chains.IsBitcoinRegnet(ob.Chain().ChainId)
		feeRateEstimated, err = ob.rpc.GetEstimatedFeeRate(ctx, 1, isRegnet)
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

// specialHandleFeeRate handles the fee rate for regnet and testnet
func (ob *Observer) specialHandleFeeRate(ctx context.Context) (int64, error) {
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

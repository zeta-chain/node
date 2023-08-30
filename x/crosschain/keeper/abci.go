package keeper

import (
	"fmt"
	"strconv"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

const (
	// EpochLength is the number of blocks in an epoch before triggering a gas price increase
	EpochLength = 100

	// GasPriceIncreasePercent is the percentage of median gas price by which to increase the gas price during an increment
	// 100 means the gas price is increased by the median gas price
	GasPriceIncreasePercent = 100
)

// IterateAndUpdateOutboundTxGasPrice iterates through all cctx and updates the gas price if pending for too long
func (k Keeper) IterateAndUpdateOutboundTxGasPrice(ctx sdk.Context, chainID int64) error {
	if chainID < 0 {
		return sdkerrors.Wrap(types.ErrInvalidChainID, "chain id cannot be negative")
	}

	// skip if haven't reached epoch end
	if ctx.BlockHeight()%EpochLength != 0 {
		return nil
	}

	unwrappedCtx := sdk.UnwrapSDKContext(ctx)

	// get all pending cctx
	res, err := k.CctxAllPending(unwrappedCtx, &types.QueryAllCctxPendingRequest{
		ChainId: uint64(chainID),
	})
	if err != nil {
		return err
	}

	// iterate through all pending cctx
	for _, pendingCctx := range res.CrossChainTx {
		if pendingCctx != nil {
			if pendingCctx.GetCurrentOutTxParam().OutboundTxGasPrice == "" {
				continue
			}

			// TODO: add block number in outbound tx

			// Compute gas price increase
			medianGasPrice, isFound := k.GetMedianGasPriceInUint(ctx, chainID)
			if !isFound {
				return sdkerrors.Wrap(
					types.ErrUnableToGetGasPrice,
					fmt.Sprintf("cannot get gas price for chain %d", chainID),
				)
			}
			gasPriceIncrease := medianGasPrice.MulUint64(GasPriceIncreasePercent).QuoUint64(100)

			// TODO: pay gas from gas stability pool

			// Increate the cctx value
			err := k.IncreaseCctxGasPrice(unwrappedCtx, *pendingCctx, gasPriceIncrease)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// IncreaseCctxGasPrice increases the gas price associated with a CCTX and updates the it in the store
func (k Keeper) IncreaseCctxGasPrice(ctx sdk.Context, cctx types.CrossChainTx, gasPriceIncrease math.Uint) error {
	currentGasPrice, err := strconv.ParseUint(cctx.GetCurrentOutTxParam().OutboundTxGasPrice, 10, 64)
	if err != nil {
		return fmt.Errorf("unable to parse cctx gas price %s: %s", cctx.GetCurrentOutTxParam().OutboundTxGasPrice, err.Error())
	}

	cctx.GetCurrentOutTxParam().OutboundTxGasPrice = math.NewUint(currentGasPrice).Add(gasPriceIncrease).String()
	k.SetCrossChainTx(ctx, cctx)
	return nil
}

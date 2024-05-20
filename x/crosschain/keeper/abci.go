package keeper

import (
	"fmt"
	"time"

	zetachains "github.com/zeta-chain/zetacore/pkg/chains"

	cosmoserrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

const (
	// RemainingFeesToStabilityPoolPercent is the percentage of remaining fees used to fund the gas stability pool
	RemainingFeesToStabilityPoolPercent = 95
)

// CheckAndUpdateCctxGasPriceFunc is a function type for checking and updating the gas price of a cctx
type CheckAndUpdateCctxGasPriceFunc func(
	ctx sdk.Context,
	k Keeper,
	cctx types.CrossChainTx,
	flags observertypes.GasPriceIncreaseFlags,
) (math.Uint, math.Uint, error)

// IterateAndUpdateCctxGasPrice iterates through all cctx and updates the gas price if pending for too long
// The function returns the number of cctxs updated and the gas price increase flags used
func (k Keeper) IterateAndUpdateCctxGasPrice(
	ctx sdk.Context,
	chains []*zetachains.Chain,
	updateFunc CheckAndUpdateCctxGasPriceFunc,
) (int, observertypes.GasPriceIncreaseFlags) {
	// fetch the gas price increase flags or use default
	gasPriceIncreaseFlags := observertypes.DefaultGasPriceIncreaseFlags
	crosschainFlags, found := k.zetaObserverKeeper.GetCrosschainFlags(ctx)
	if found && crosschainFlags.GasPriceIncreaseFlags != nil {
		gasPriceIncreaseFlags = *crosschainFlags.GasPriceIncreaseFlags
	}

	// skip if haven't reached epoch end
	if ctx.BlockHeight()%gasPriceIncreaseFlags.EpochLength != 0 {
		return 0, gasPriceIncreaseFlags
	}

	cctxCount := 0

IterateChains:
	for _, chain := range chains {
		// support only external evm chains
		if zetachains.IsEVMChain(chain.ChainId) && !zetachains.IsZetaChain(chain.ChainId) {
			res, err := k.ListPendingCctx(sdk.UnwrapSDKContext(ctx), &types.QueryListPendingCctxRequest{
				ChainId: chain.ChainId,
				Limit:   gasPriceIncreaseFlags.MaxPendingCctxs,
			})
			if err != nil {
				ctx.Logger().Info("GasStabilityPool: fetching pending cctx failed",
					"chainID", chain.ChainId,
					"err", err.Error(),
				)
				continue IterateChains
			}

			// iterate through all pending cctx
			for _, pendingCctx := range res.CrossChainTx {
				if pendingCctx != nil {
					gasPriceIncrease, additionalFees, err := updateFunc(ctx, k, *pendingCctx, gasPriceIncreaseFlags)
					if err != nil {
						ctx.Logger().Info("GasStabilityPool: updating gas price for pending cctx failed",
							"cctxIndex", pendingCctx.Index,
							"err", err.Error(),
						)
						continue IterateChains
					}
					if !gasPriceIncrease.IsNil() && !gasPriceIncrease.IsZero() {
						// Emit typed event for gas price increase
						if err := ctx.EventManager().EmitTypedEvent(
							&types.EventCCTXGasPriceIncreased{
								CctxIndex:        pendingCctx.Index,
								GasPriceIncrease: gasPriceIncrease.String(),
								AdditionalFees:   additionalFees.String(),
							}); err != nil {
							ctx.Logger().Error(
								"GasStabilityPool: failed to emit EventCCTXGasPriceIncreased",
								"err", err.Error(),
							)
						}
						cctxCount++
					}
				}
			}
		}
	}

	return cctxCount, gasPriceIncreaseFlags
}

// CheckAndUpdateCctxGasPrice checks if the retry interval is reached and updates the gas price if so
// The function returns the gas price increase and the additional fees paid from the gas stability pool
func CheckAndUpdateCctxGasPrice(
	ctx sdk.Context,
	k Keeper,
	cctx types.CrossChainTx,
	flags observertypes.GasPriceIncreaseFlags,
) (math.Uint, math.Uint, error) {
	// skip if gas price or gas limit is not set
	if cctx.GetCurrentOutboundParam().GasPrice == "" || cctx.GetCurrentOutboundParam().GasLimit == 0 {
		return math.ZeroUint(), math.ZeroUint(), nil
	}

	// skip if retry interval is not reached
	lastUpdated := time.Unix(cctx.CctxStatus.LastUpdateTimestamp, 0)
	if ctx.BlockTime().Before(lastUpdated.Add(flags.RetryInterval)) {
		return math.ZeroUint(), math.ZeroUint(), nil
	}

	// compute gas price increase
	chainID := cctx.GetCurrentOutboundParam().ReceiverChainId
	medianGasPrice, isFound := k.GetMedianGasPriceInUint(ctx, chainID)
	if !isFound {
		return math.ZeroUint(), math.ZeroUint(), cosmoserrors.Wrap(
			types.ErrUnableToGetGasPrice,
			fmt.Sprintf("cannot get gas price for chain %d", chainID),
		)
	}
	gasPriceIncrease := medianGasPrice.MulUint64(uint64(flags.GasPriceIncreasePercent)).QuoUint64(100)

	// compute new gas price
	currentGasPrice, err := cctx.GetCurrentOutboundParam().GetGasPriceUInt64()
	if err != nil {
		return math.ZeroUint(), math.ZeroUint(), err
	}
	newGasPrice := math.NewUint(currentGasPrice).Add(gasPriceIncrease)

	// check limit -- use default limit if not set
	gasPriceIncreaseMax := flags.GasPriceIncreaseMax
	if gasPriceIncreaseMax == 0 {
		gasPriceIncreaseMax = observertypes.DefaultGasPriceIncreaseFlags.GasPriceIncreaseMax
	}
	limit := medianGasPrice.MulUint64(uint64(gasPriceIncreaseMax)).QuoUint64(100)
	if newGasPrice.GT(limit) {
		return math.ZeroUint(), math.ZeroUint(), nil
	}

	// withdraw additional fees from the gas stability pool
	gasLimit := math.NewUint(cctx.GetCurrentOutboundParam().GasLimit)
	additionalFees := gasLimit.Mul(gasPriceIncrease)
	if err := k.fungibleKeeper.WithdrawFromGasStabilityPool(ctx, chainID, additionalFees.BigInt()); err != nil {
		return math.ZeroUint(), math.ZeroUint(), cosmoserrors.Wrap(
			types.ErrNotEnoughFunds,
			fmt.Sprintf("cannot withdraw %s from gas stability pool, error: %s", additionalFees.String(), err.Error()),
		)
	}

	// set new gas price and last update timestamp
	cctx.GetCurrentOutboundParam().GasPrice = newGasPrice.String()
	cctx.CctxStatus.LastUpdateTimestamp = ctx.BlockHeader().Time.Unix()
	k.SetCrossChainTx(ctx, cctx)

	return gasPriceIncrease, additionalFees, nil
}

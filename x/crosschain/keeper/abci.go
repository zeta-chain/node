package keeper

import (
	"fmt"
	"time"

	cosmoserrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/pkg/chains"
	mathpkg "github.com/zeta-chain/node/pkg/math"
	"github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

const (
	// RemainingFeesToStabilityPoolPercent is the percentage of remaining fees used to fund the gas stability pool
	RemainingFeesToStabilityPoolPercent = 95
)

// CheckAndUpdateCCTXGasPriceFunc is a function type for checking and updating the gas price of a cctx
type CheckAndUpdateCCTXGasPriceFunc func(
	ctx sdk.Context,
	k Keeper,
	cctx types.CrossChainTx,
	flags observertypes.GasPriceIncreaseFlags,
) (math.Uint, math.Uint, error)

// IterateAndUpdateCCTXGasPrice iterates through all cctx and updates the gas price if pending for too long
// The function returns the number of cctxs updated and the gas price increase flags used
func (k Keeper) IterateAndUpdateCCTXGasPrice(
	ctx sdk.Context,
	chains []chains.Chain,
	updateFunc CheckAndUpdateCCTXGasPriceFunc,
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

	additionalChains := k.GetAuthorityKeeper().GetAdditionalChainList(ctx)

	cctxCount := 0

IterateChains:
	for _, chain := range chains {
		if !IsCCTXGasPriceUpdateSupported(chain.ChainId, additionalChains) {
			continue
		}

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
			if pendingCctx == nil {
				continue
			}

			gasPriceIncrease, additionalFees, err := updateFunc(ctx, k, *pendingCctx, gasPriceIncreaseFlags)
			if err != nil {
				ctx.Logger().Info("GasStabilityPool: updating gas price for pending cctx failed",
					"cctxIndex", pendingCctx.Index,
					"err", err.Error(),
				)
				continue IterateChains
			}
			if gasPriceIncrease.IsNil() || gasPriceIncrease.IsZero() {
				continue
			}

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

	return cctxCount, gasPriceIncreaseFlags
}

// CheckAndUpdateCCTXGasPrice checks if the retry interval is reached and updates the gas price if so
// The function returns the gas price increase and the additional fees paid from the gas stability pool
func CheckAndUpdateCCTXGasPrice(
	ctx sdk.Context,
	k Keeper,
	cctx types.CrossChainTx,
	flags observertypes.GasPriceIncreaseFlags,
) (math.Uint, math.Uint, error) {
	// skip if gas price or gas limit is not set
	if cctx.GetCurrentOutboundParam().GasPrice == "" || cctx.GetCurrentOutboundParam().CallOptions.GasLimit == 0 {
		return math.ZeroUint(), math.ZeroUint(), nil
	}

	// get latest median gas price and priority fee
	chainID := cctx.GetCurrentOutboundParam().ReceiverChainId
	medianGasPrice, medianPriorityFee, isFound := k.GetMedianGasValues(ctx, chainID)
	if !isFound {
		return math.ZeroUint(), math.ZeroUint(), cosmoserrors.Wrap(
			types.ErrUnableToGetGasPrice,
			fmt.Sprintf("cannot get gas price for chain %d", chainID),
		)
	}

	// dispatch to chain-specific gas price update function
	additionalChains := k.GetAuthorityKeeper().GetAdditionalChainList(ctx)
	switch {
	case chains.IsEVMChain(chainID, additionalChains):
		return CheckAndUpdateCCTXGasPriceEVM(ctx, k, medianGasPrice, medianPriorityFee, cctx, flags)
	case chains.IsBitcoinChain(chainID, additionalChains):
		return CheckAndUpdateCCTXGasPriceBTC(ctx, k, medianGasPrice, cctx, flags)
	default:
		return math.ZeroUint(), math.ZeroUint(), nil
	}
}

// CheckAndUpdateCCTXGasPriceEVM updates the gas price for the given EVM chain CCTX
func CheckAndUpdateCCTXGasPriceEVM(
	ctx sdk.Context,
	k Keeper,
	medianGasPrice math.Uint,
	medianPriorityFee math.Uint,
	cctx types.CrossChainTx,
	flags observertypes.GasPriceIncreaseFlags,
) (gasPriceIncrease math.Uint, additionalFees math.Uint, err error) {
	// skip if retry interval is not reached
	lastUpdated := time.Unix(cctx.CctxStatus.LastUpdateTimestamp, 0)
	if ctx.BlockTime().Before(lastUpdated.Add(flags.RetryInterval)) {
		return math.ZeroUint(), math.ZeroUint(), nil
	}

	// compute gas price increase
	gasPriceIncrease = medianGasPrice.MulUint64(uint64(flags.GasPriceIncreasePercent)).QuoUint64(100)

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

	newPriorityFee, _ := mathpkg.IncreaseUintByPercent(medianPriorityFee, uint64(flags.GasPriceIncreasePercent))

	// should not happen
	if newPriorityFee.GT(newGasPrice) {
		return math.ZeroUint(), math.ZeroUint(), fmt.Errorf(
			"priorityFee %s is greater than new gasPrice %s",
			newPriorityFee.String(),
			newGasPrice.String(),
		)
	}

	// withdraw additional fees from the gas stability pool
	chainID := cctx.GetCurrentOutboundParam().ReceiverChainId
	gasLimit := math.NewUint(cctx.GetCurrentOutboundParam().CallOptions.GasLimit)
	additionalFees = gasLimit.Mul(gasPriceIncrease)
	if err := k.fungibleKeeper.WithdrawFromGasStabilityPool(ctx, chainID, additionalFees.BigInt()); err != nil {
		return math.ZeroUint(), math.ZeroUint(), cosmoserrors.Wrap(
			types.ErrNotEnoughFunds,
			fmt.Sprintf("cannot withdraw %s from gas stability pool, error: %s", additionalFees.String(), err.Error()),
		)
	}

	// set new gas price and last update timestamp
	cctx.GetCurrentOutboundParam().GasPrice = newGasPrice.String()
	cctx.GetCurrentOutboundParam().GasPriorityFee = newPriorityFee.String()
	k.SetCrossChainTx(ctx, cctx)

	return gasPriceIncrease, additionalFees, nil
}

// CheckAndUpdateCCTXGasPriceBTC updates the gas price for the given Bitcoin chain CCTX
func CheckAndUpdateCCTXGasPriceBTC(
	ctx sdk.Context,
	k Keeper,
	medianGasPrice math.Uint,
	cctx types.CrossChainTx,
	flags observertypes.GasPriceIncreaseFlags,
) (gasPriceIncrease math.Uint, additionalFees math.Uint, err error) {
	// check retry interval -- use default interval if not set
	retryIntervalBTC := flags.RetryIntervalBTC
	if retryIntervalBTC <= 0 {
		retryIntervalBTC = observertypes.DefaultGasPriceIncreaseFlags.RetryIntervalBTC
	}

	// skip if Bitcoin gas price retry interval is not reached
	lastUpdated := time.Unix(cctx.CctxStatus.LastUpdateTimestamp, 0)
	if ctx.BlockTime().Before(lastUpdated.Add(retryIntervalBTC)) {
		return math.ZeroUint(), math.ZeroUint(), nil
	}

	// get current gas price
	currentGasPrice, err := cctx.GetCurrentOutboundParam().GetGasPriceUInt64()
	if err != nil {
		return math.ZeroUint(), math.ZeroUint(), err
	}

	// use latest median gas price as new gas price, the reasons are:
	// 1. the goal is to increase the average gas price of all the stuck txs to market level
	// 2. zetaclient can't replace stuck tx individually, it only gives more funds to the last stuck tx (child tx)
	// 3. updating all pending CCTXs to the same 'mediaGasPrice' number simplifies the calculation in zetaclient
	newGasPrice := medianGasPrice

	// set priority fee to signal that fee bumping is allowed for this CCTX
	//
	// there is no priority fee for Bitcoin chain, the reason for setting 'GasPriorityFee' is:
	// zetaclient can't figure out whether the gas price has been updated or not, so setting the
	// 'GasPriorityFee' is to signal the gas price is updated, and the RBF (replace by fee) is allowed.
	cctx.GetCurrentOutboundParam().GasPriorityFee = newGasPrice.String()

	// early return if there is no need to withdraw additional fees
	if math.NewUint(currentGasPrice).GTE(newGasPrice) {
		k.SetCrossChainTx(ctx, cctx)
		return math.ZeroUint(), math.ZeroUint(), nil
	}

	// compute gas price increase
	gasPriceIncrease = newGasPrice.SubUint64(currentGasPrice)

	// withdraw additional fees from the gas stability pool
	chainID := cctx.GetCurrentOutboundParam().ReceiverChainId
	gasLimit := math.NewUint(cctx.GetCurrentOutboundParam().CallOptions.GasLimit)
	additionalFees = gasLimit.Mul(gasPriceIncrease)
	if err := k.fungibleKeeper.WithdrawFromGasStabilityPool(ctx, chainID, additionalFees.BigInt()); err != nil {
		return math.ZeroUint(), math.ZeroUint(), cosmoserrors.Wrap(
			types.ErrNotEnoughFunds,
			fmt.Sprintf(
				"cannot withdraw %s satoshis from gas stability pool, error: %s",
				additionalFees.String(),
				err.Error(),
			),
		)
	}

	// set new gas price and last update timestamp
	cctx.GetCurrentOutboundParam().GasPrice = newGasPrice.String()
	k.SetCrossChainTx(ctx, cctx)

	return gasPriceIncrease, additionalFees, nil
}

// IsCCTXGasPriceUpdateSupported checks if the given chain supports gas price update
func IsCCTXGasPriceUpdateSupported(chainID int64, additionalChains []chains.Chain) bool {
	switch {
	case chains.IsZetaChain(chainID, additionalChains):
		return false
	case chains.IsEVMChain(chainID, additionalChains),
		chains.IsBitcoinChain(chainID, additionalChains):
		return true
	default:
		return false
	}
}

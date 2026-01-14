package v12

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/x/observer/types"
)

type observerKeeper interface {
	GetChainParamsList(ctx sdk.Context) (val types.ChainParamsList, found bool)
	SetChainParamsList(ctx sdk.Context, chainParams types.ChainParamsList)
	GetCrosschainFlags(ctx sdk.Context) (val types.CrosschainFlags, found bool)
	SetCrosschainFlags(ctx sdk.Context, crosschainFlags types.CrosschainFlags)
}

// MigrateStore migrates the x/observer module state from the consensus version 11 to version 12.
// The migration updates the 'StabilityPoolPercentage' field of all chain params to 100
// The migration sets default gas price multipliers for all external chains.
// The migration disables V2 ZETA flows.
func MigrateStore(ctx sdk.Context, observerKeeper observerKeeper) error {
	if err := UpdateChainParams(ctx, observerKeeper); err != nil {
		return err
	}

	UpdateCrosschainFlags(ctx, observerKeeper)

	return nil
}

// UpdateChainParams updates the chain params with gas price multiplier and stability pool percentage.
// It also removes chain params for chains that no longer exist in the chain list.
func UpdateChainParams(ctx sdk.Context, observerKeeper observerKeeper) error {
	allChainParams, found := observerKeeper.GetChainParamsList(ctx)
	if !found {
		return errorsmod.Wrap(types.ErrChainParamsNotFound, "failed to get chain params")
	}

	// Filter out removed chains and update valid ones
	var updatedChainParams []*types.ChainParams
	for _, chainParams := range allChainParams.ChainParams {
		if chainParams != nil {
			chain, foundChain := chains.GetChainFromChainID(chainParams.ChainId, []chains.Chain{})
			if !foundChain {
				ctx.Logger().Warn("removing orphaned chain params for removed chain",
					"chain_id", chainParams.ChainId)
				continue
			}
			chainParams.GasPriceMultiplier = GetGasPriceMultiplierForChain(chain)
			chainParams.StabilityPoolPercentage = 100
			updatedChainParams = append(updatedChainParams, chainParams)
		}
	}

	allChainParams.ChainParams = updatedChainParams

	if err := allChainParams.Validate(); err != nil {
		return errorsmod.Wrap(types.ErrInvalidChainParams, err.Error())
	}

	observerKeeper.SetChainParamsList(ctx, allChainParams)

	return nil
}

// UpdateCrosschainFlags disables V2 ZETA flows in the crosschain flags.
func UpdateCrosschainFlags(ctx sdk.Context, observerKeeper observerKeeper) {
	flags, found := observerKeeper.GetCrosschainFlags(ctx)
	if !found {
		flags = *types.DefaultCrosschainFlags()
	}

	flags.IsV2ZetaEnabled = false
	observerKeeper.SetCrosschainFlags(ctx, flags)
}

// GetGasPriceMultiplierForChain returns the gas price multiplier for the given chain
func GetGasPriceMultiplierForChain(chain chains.Chain) sdkmath.LegacyDec {
	switch chain.Consensus {
	case chains.Consensus_ethereum:
		return types.DefaultEVMOutboundGasPriceMultiplier
	case chains.Consensus_bitcoin:
		return types.DefaultBTCOutboundGasPriceMultiplier
	default:
		return types.DefaultGasPriceMultiplier
	}
}

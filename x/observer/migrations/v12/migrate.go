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
}

// MigrateStore migrates the x/observer module state from the consensus version 11 to version 12.
// The migration sets default gas price multipliers for all external chains.
func MigrateStore(ctx sdk.Context, observerKeeper observerKeeper) error {
	allChainParams, found := observerKeeper.GetChainParamsList(ctx)
	if !found {
		return errorsmod.Wrap(types.ErrChainParamsNotFound, "failed to get chain params")
	}

	// set new fields to the same value as 'confirmation_count'
	for _, chainParams := range allChainParams.ChainParams {
		if chainParams != nil {
			chain, found := chains.GetChainFromChainID(chainParams.ChainId, []chains.Chain{})
			if !found {
				return errorsmod.Wrapf(types.ErrSupportedChains, "chain %d not found", chainParams.ChainId)
			}
			chainParams.GasPriceMultiplier = GetGasPriceMultiplierForChain(chain)
		}
	}

	// validate the updated chain params list
	if err := allChainParams.Validate(); err != nil {
		return errorsmod.Wrap(types.ErrInvalidChainParams, err.Error())
	}

	// set the updated chain params list
	observerKeeper.SetChainParamsList(ctx, allChainParams)

	return nil
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

package v10

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/observer/types"
)

type observerKeeper interface {
	GetChainParamsList(ctx sdk.Context) (val types.ChainParamsList, found bool)
	SetChainParamsList(ctx sdk.Context, chainParams types.ChainParamsList)
}

// MigrateStore migrates the x/observer module state from the consensus version 9 to version 10.
// The migration sets existing 'confirmation_count' as default value for newly added fields:
//   - 'safe_inbound_count'
//   - 'fast_inbound_count'
//   - 'safe_outbound_count'
//   - 'fast_outbound_count'
func MigrateStore(ctx sdk.Context, observerKeeper observerKeeper) error {
	allChainParams, found := observerKeeper.GetChainParamsList(ctx)
	if !found {
		return errorsmod.Wrap(types.ErrChainParamsNotFound, "failed to get chain params")
	}

	// set new fields to the same value as 'confirmation_count'
	for _, chainParams := range allChainParams.ChainParams {
		if chainParams != nil {
			chainParams.ConfirmationParams.SafeInboundCount = chainParams.ConfirmationCount
			chainParams.ConfirmationParams.FastInboundCount = chainParams.ConfirmationCount
			chainParams.ConfirmationParams.SafeOutboundCount = chainParams.ConfirmationCount
			chainParams.ConfirmationParams.FastOutboundCount = chainParams.ConfirmationCount
		}
	}

	// set the updated chain params list
	observerKeeper.SetChainParamsList(ctx, allChainParams)

	return nil
}

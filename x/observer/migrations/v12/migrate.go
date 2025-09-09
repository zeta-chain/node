package v12

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/observer/types"
)

type observerKeeper interface {
	GetChainParamsList(ctx sdk.Context) (val types.ChainParamsList, found bool)
	SetChainParamsList(ctx sdk.Context, chainParams types.ChainParamsList)
}

// MigrateStore migrates the x/observer module state from the consensus version 11 to version 12.
// The migration updates the 'StabilityPoolPercentage' field of all chain params to 60.
func MigrateStore(ctx sdk.Context, observerKeeper observerKeeper) error {
	allChainParams, found := observerKeeper.GetChainParamsList(ctx)
	if !found {
		return errorsmod.Wrap(types.ErrChainParamsNotFound, "failed to get chain params")
	}

	for _, chainParams := range allChainParams.ChainParams {
		if chainParams != nil {
			chainParams.StabilityPoolPercentage = 60
		}
	}

	if err := allChainParams.Validate(); err != nil {
		return errorsmod.Wrap(types.ErrInvalidChainParams, err.Error())
	}

	observerKeeper.SetChainParamsList(ctx, allChainParams)

	return nil
}

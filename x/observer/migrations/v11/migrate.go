package v11

import (
	"slices"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/x/observer/types"
)

type observerKeeper interface {
	GetChainParamsList(ctx sdk.Context) (val types.ChainParamsList, found bool)
	SetChainParamsList(ctx sdk.Context, chainParams types.ChainParamsList)
}

var disableBlockScanChainIDs = []int64{
	// localnet
	chains.GoerliLocalnet.ChainId,

	// testnets
	chains.AvalancheTestnet.ChainId,
	chains.ArbitrumSepolia.ChainId,
	chains.BscTestnet.ChainId,
	chains.BaseSepolia.ChainId,
	chains.Amoy.ChainId,

	// mainnets
	// the other mainnets will be deprecated later
	chains.AvalancheMainnet.ChainId,
	chains.ArbitrumMainnet.ChainId,
}

// MigrateStore migrates the x/observer module state from the consensus version 10 to version 11.
// The migration sets the skip_block_scan parameter correctly based on the chain ID
func MigrateStore(ctx sdk.Context, observerKeeper observerKeeper) error {
	allChainParams, found := observerKeeper.GetChainParamsList(ctx)
	if !found {
		return errorsmod.Wrap(types.ErrChainParamsNotFound, "failed to get chain params")
	}

	// set new fields to the same value as 'confirmation_count'
	for _, chainParams := range allChainParams.ChainParams {
		if chainParams == nil {
			continue
		}
		if slices.Contains(disableBlockScanChainIDs, chainParams.ChainId) {
			chainParams.DisableBlockScan = true
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

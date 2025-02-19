package v5

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/emissions/types"
)

type EmissionsKeeper interface {
	SetParams(ctx sdk.Context, params types.Params) error
	GetParams(ctx sdk.Context) (types.Params, bool)
}

// MigrateStore migrates the store from v4 to v5
// The v5 params include a new parameter, PendingBallotsBufferBlocks, which is set to the default value
func MigrateStore(
	ctx sdk.Context,
	emissionsKeeper EmissionsKeeper,
) error {
	updatedParams := types.DefaultParams()
	// If params are found, update all fields except the new one `PendingBallotsBufferBlocks`
	params, found := emissionsKeeper.GetParams(ctx)
	if found {
		if params.ValidatorEmissionPercentage != "" {
			updatedParams.ValidatorEmissionPercentage = params.ValidatorEmissionPercentage
		}

		if params.ObserverEmissionPercentage != "" {
			updatedParams.ObserverEmissionPercentage = params.ObserverEmissionPercentage
		}

		if params.TssSignerEmissionPercentage != "" {
			updatedParams.TssSignerEmissionPercentage = params.TssSignerEmissionPercentage
		}

		if params.BlockRewardAmount.GT(sdkmath.LegacyZeroDec()) {
			updatedParams.BlockRewardAmount = params.BlockRewardAmount
		}

		if params.ObserverSlashAmount.GTE(sdkmath.ZeroInt()) {
			updatedParams.ObserverSlashAmount = params.ObserverSlashAmount
		}
		if params.BallotMaturityBlocks > 0 {
			updatedParams.BallotMaturityBlocks = params.BallotMaturityBlocks
		}
	}

	return emissionsKeeper.SetParams(ctx, updatedParams)
}

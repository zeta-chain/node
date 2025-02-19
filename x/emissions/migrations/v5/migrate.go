package v5

import (
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
	params, found := emissionsKeeper.GetParams(ctx)
	if found {
		updatedParams.BlockRewardAmount = params.BlockRewardAmount
		updatedParams.ValidatorEmissionPercentage = params.ValidatorEmissionPercentage
		updatedParams.ObserverEmissionPercentage = params.ObserverEmissionPercentage
		updatedParams.TssSignerEmissionPercentage = params.TssSignerEmissionPercentage
		updatedParams.ObserverSlashAmount = params.ObserverSlashAmount
		updatedParams.BallotMaturityBlocks = params.BallotMaturityBlocks
	}
	return emissionsKeeper.SetParams(ctx, updatedParams)
}

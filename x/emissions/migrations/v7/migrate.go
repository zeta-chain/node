package v7

import (
	"errors"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/emissions/types"
)

type EmissionsKeeper interface {
	SetParams(ctx sdk.Context, params types.Params) error
	GetParams(ctx sdk.Context) (types.Params, bool)
}

// MigrateStore migrates the store from v6 to v7
// It just updates values for BlockRewardAmount and BallotMaturityBlocks
func MigrateStore(
	ctx sdk.Context,
	emissionsKeeper EmissionsKeeper,
) error {
	// If params are found, update fields
	params, found := emissionsKeeper.GetParams(ctx)
	if found {
		params.BallotMaturityBlocks = 300                  // increase from 150
		params.PendingBallotsDeletionBufferBlocks = 432000 // increase from 216000
		// decrease from 6751543209876543209.876543209876543210
		params.BlockRewardAmount = sdkmath.LegacyMustNewDecFromStr("3375771604938271604.938271604938271605")

		return emissionsKeeper.SetParams(ctx, params)
	}

	// should not happen, in previous migrations store was set properly
	return errors.New("emission params not found")
}

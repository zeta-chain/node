package v6

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

// MigrateStore migrates the store from v5 to v6
// It just updates values for BlockRewardAmount and BallotMaturityBlocks
func MigrateStore(
	ctx sdk.Context,
	emissionsKeeper EmissionsKeeper,
) error {
	// If params are found, update fields
	params, found := emissionsKeeper.GetParams(ctx)
	if found {
		params.BallotMaturityBlocks = 150                  // increase from 100
		params.PendingBallotsDeletionBufferBlocks = 216000 // increase from 144000
		// decrease from 9620949074074074074.074070733466756687
		params.BlockRewardAmount = sdkmath.LegacyMustNewDecFromStr("6751543209876543209.876543209876543210")

		return emissionsKeeper.SetParams(ctx, params)
	}

	// should not happen, in previous migrations store was set properly
	return errors.New("emission params not found")
}

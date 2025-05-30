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
		params.BallotMaturityBlocks = 133
		params.BlockRewardAmount = sdkmath.LegacyMustNewDecFromStr("7595486111111111680.000000000000000000")

		return emissionsKeeper.SetParams(ctx, params)
	}

	// should not happen, in previous migrations store was set properly
	return errors.New("emission params not found")
}

package v3

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/emissions/exported"
	"github.com/zeta-chain/zetacore/x/emissions/types"
)

type EmissionsKeeper interface {
	GetParams(ctx sdk.Context) (types.Params, bool)
	SetParams(ctx sdk.Context, params types.Params) error
}

// ObserverSlashAmount is the amount of tokens to be slashed from observer in case of incorrect vote
// by default it is set to 0.1 ZETA
var observerSlashAmountDefaultValue = sdkmath.NewInt(100000000000000000)

// BallotMaturityBlocks is amount of blocks needed for ballot to mature
// by default is set to 100
var ballotMaturityBlocksDefaultValue = 100

// Migrate migrates the x/emissions module state from the consensus version 2 to
// version 3. Specifically, it takes the parameters that are currently stored
// and managed by the x/params modules and stores them directly into the x/emissions
// module state.
func MigrateStore(
	ctx sdk.Context,
	emissionsKeeper EmissionsKeeper,
	legacySubspace exported.Subspace,
) error {
	var currParams types.Params
	legacySubspace.GetParamSet(ctx, &currParams)

	currParams.ObserverSlashAmount = observerSlashAmountDefaultValue
	currParams.BallotMaturityBlocks = int64(ballotMaturityBlocksDefaultValue)
	err := currParams.Validate()
	if err != nil {
		return err
	}

	return emissionsKeeper.SetParams(ctx, currParams)
}

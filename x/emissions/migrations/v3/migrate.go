package v3

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/emissions/exported"
	"github.com/zeta-chain/node/x/emissions/types"
)

type EmissionsKeeper interface {
	GetParams(ctx sdk.Context) (types.Params, bool)
	SetParams(ctx sdk.Context, params types.Params) error
}

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

	defaultParams := types.NewParams()

	if currParams.ValidatorEmissionPercentage == "" {
		currParams.ValidatorEmissionPercentage = defaultParams.ValidatorEmissionPercentage
	}
	if currParams.ObserverEmissionPercentage == "" {
		currParams.ObserverEmissionPercentage = defaultParams.ObserverEmissionPercentage
	}
	if currParams.TssSignerEmissionPercentage == "" {
		currParams.TssSignerEmissionPercentage = defaultParams.TssSignerEmissionPercentage
	}

	currParams.ObserverSlashAmount = types.ObserverSlashAmount
	currParams.BallotMaturityBlocks = int64(types.BallotMaturityBlocks)

	// BlockRewardAmount is a new parameter introduced in version 4
	currParams.BlockRewardAmount = types.BlockReward
	err := currParams.Validate()
	if err != nil {
		return err
	}

	return emissionsKeeper.SetParams(ctx, currParams)
}

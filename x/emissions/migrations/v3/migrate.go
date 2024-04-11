package v3

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/emissions/exported"
	"github.com/zeta-chain/zetacore/x/emissions/types"
)

type EmissionsKeeper interface {
	GetParams(ctx sdk.Context) types.Params
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

	currParams.ObserverSlashAmount = sdkmath.NewInt(100000000000000000)
	err := currParams.Validate()
	if err != nil {
		return err
	}

	return emissionsKeeper.SetParams(ctx, currParams)
}

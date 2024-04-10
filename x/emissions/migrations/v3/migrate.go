package v3

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/emissions/exported"
	"github.com/zeta-chain/zetacore/x/emissions/types"
)

type EmissionsKeeper interface {
	GetParamsIfExists(ctx sdk.Context) types.Params
	SetParams(ctx sdk.Context, params types.Params) error
}

// Migrate migrates the x/emissions module state from the consensus version 1 to
// version 2. Specifically, it takes the parameters that are currently stored
// and managed by the x/params modules and stores them directly into the x/emissions
// module state.
func Migrate(
	ctx sdk.Context,
	emissionsKeeper EmissionsKeeper,
	legacySubspace exported.Subspace,
) error {
	var currParams types.Params
	legacySubspace.GetParamSet(ctx, &currParams)

	if err := currParams.Validate(); err != nil {
		return err
	}
	emissionsKeeper.SetParams(ctx, currParams)

	return nil
}

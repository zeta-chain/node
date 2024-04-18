package v8

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/exported"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

type ObserverKeeper interface {
	GetParams(ctx sdk.Context) (types.Params, bool)
	SetParams(ctx sdk.Context, params types.Params) error
}

// Migrate migrates the x/observer module state from the consensus version 7 to
// version 8. Specifically, it takes the parameters that are currently stored
// and managed by the x/params modules and stores them directly into the x/observer
// module state.
func MigrateStore(
	ctx sdk.Context,
	observerKeeper ObserverKeeper,
	legacySubspace exported.Subspace,
) error {
	var currParams types.Params
	legacySubspace.GetParamSet(ctx, &currParams)

	err := currParams.Validate()
	if err != nil {
		return err
	}

	return observerKeeper.SetParams(ctx, currParams)
}

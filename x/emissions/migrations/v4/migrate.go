package v4

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/emissions/types"
)

type EmissionsKeeper interface {
	SetParams(ctx sdk.Context, params types.Params) error
}

// Migrate migrates the x/emissions module state from the consensus version 2 to
// version 3. Specifically, it takes the parameters that are currently stored
// and managed by the x/params modules and stores them directly into the x/emissions
// module state.
func MigrateStore(
	ctx sdk.Context,
	storeKey storetypes.StoreKey,
	cdc codec.Codec,
	emissionsKeeper EmissionsKeeper,
) error {
	currentParams, found := GetParamsLegacy(ctx, storeKey, cdc)
	if !found {
		err := fmt.Errorf("failed to get legacy params")
		ctx.Logger().Error("error :", err.Error())
		return err
	}

	defaultParams := types.NewParams()
	if currentParams.ValidatorEmissionPercentage != "" {
		defaultParams.ValidatorEmissionPercentage = currentParams.ValidatorEmissionPercentage
	}
	if currentParams.ObserverEmissionPercentage != "" {
		defaultParams.ObserverEmissionPercentage = currentParams.ObserverEmissionPercentage
	}
	if currentParams.TssSignerEmissionPercentage != "" {
		defaultParams.TssSignerEmissionPercentage = currentParams.TssSignerEmissionPercentage
	}
	defaultParams.ObserverSlashAmount = currentParams.ObserverSlashAmount
	defaultParams.BallotMaturityBlocks = currentParams.BallotMaturityBlocks

	err := emissionsKeeper.SetParams(ctx, defaultParams)
	if err != nil {
		return err
	}
	return nil
}

func GetParamsLegacy(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.Codec) (params types.LegacyParams, found bool) {
	store := ctx.KVStore(storeKey)
	bz := store.Get(types.KeyPrefix(types.ParamsKey))
	if bz == nil {
		return types.LegacyParams{}, false
	}
	err := cdc.Unmarshal(bz, &params)
	if err != nil {
		return types.LegacyParams{}, false
	}

	return params, true
}

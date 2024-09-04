package v4

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/emissions/types"
)

type EmissionsKeeper interface {
	SetParams(ctx sdk.Context, params types.Params) error
	GetCodec() codec.BinaryCodec
	GetStoreKey() storetypes.StoreKey
}

// Migrate migrates the x/emissions module state from the consensus version 2 to
// version 3. Specifically, it takes the parameters that are currently stored
// and managed by the x/params modules and stores them directly into the x/emissions
// module state.
func MigrateStore(
	ctx sdk.Context,
	emissionsKeeper EmissionsKeeper,
) error {
	currentParams, found := GetParamsLegacy(ctx, emissionsKeeper.GetStoreKey(), emissionsKeeper.GetCodec())
	if !found {
		return errorsmod.Wrap(types.ErrMigrationFailed, "failed to get legacy params")
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
		return errorsmod.Wrap(types.ErrMigrationFailed, err.Error())
	}
	return nil
}

func GetParamsLegacy(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) (params types.LegacyParams, found bool) {
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

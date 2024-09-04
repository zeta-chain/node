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

// MigrateStore migrates the store from v3 to v4
// The v3 params are copied to the v4 params , and the v4 params are set in the store
// v4 introduces a new parameter, BlockRewardAmount, which is set to the default value
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

func GetParamsLegacy(
	ctx sdk.Context,
	storeKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
) (params types.LegacyParams, found bool) {
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

package v4

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
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
// The v3 params are copied to the v4 params, and the v4 params are set in the store
// v4 params removes unused parameters from v3; these values are discarded.
// v4 introduces a new parameter, BlockRewardAmount, which is set to the default value
func MigrateStore(
	ctx sdk.Context,
	emissionsKeeper EmissionsKeeper,
) error {
	v3Params, found := GetParamsLegacy(ctx, emissionsKeeper.GetStoreKey(), emissionsKeeper.GetCodec())
	if !found {
		return errorsmod.Wrap(types.ErrMigrationFailed, "failed to get legacy params")
	}

	// New params initializes v4 params with default values
	v4Params := types.NewParams()
	if v3Params.ValidatorEmissionPercentage != "" {
		v4Params.ValidatorEmissionPercentage = v3Params.ValidatorEmissionPercentage
	}
	if v3Params.ObserverEmissionPercentage != "" {
		v4Params.ObserverEmissionPercentage = v3Params.ObserverEmissionPercentage
	}
	if v3Params.TssSignerEmissionPercentage != "" {
		v4Params.TssSignerEmissionPercentage = v3Params.TssSignerEmissionPercentage
	}
	if v3Params.ObserverSlashAmount.GTE(sdkmath.ZeroInt()) {
		v4Params.ObserverSlashAmount = v3Params.ObserverSlashAmount
	}
	if v3Params.BallotMaturityBlocks > 0 {
		v4Params.BallotMaturityBlocks = v3Params.BallotMaturityBlocks
	}

	err := emissionsKeeper.SetParams(ctx, v4Params)
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

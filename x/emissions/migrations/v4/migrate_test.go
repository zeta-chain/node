package v4_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/x/emissions/types"
)

func TestMigrateStore(t *testing.T) {
	t.Run("should successfully migrate to new params", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)

		cdc := k.GetCodec()
		emissionsStoreKey := sdk.NewKVStoreKey(types.StoreKey)

		err := SetLegacyParams(ctx, emissionsStoreKey, cdc, LegacyMainnetParams())
		require.NoError(t, err)

	})
}

func SetLegacyParams(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec, params types.LegacyParams) error {
	store := ctx.KVStore(storeKey)
	bz, err := cdc.Marshal(&params)
	if err != nil {
		return err
	}

	store.Set(types.KeyPrefix(types.ParamsKey), bz)
	return nil
}

// https://zetachain-api.lavenderfive.com/zeta-chain/emissions/params
func LegacyMainnetParams() types.LegacyParams {
	return types.LegacyParams{
		MaxBondFactor:               "1.25",
		MinBondFactor:               "0.75",
		AvgBlockTime:                "6.00",
		TargetBondRatio:             "0.67",
		ObserverEmissionPercentage:  "0.125",
		ValidatorEmissionPercentage: "0.75",
		TssSignerEmissionPercentage: "0.125",
		DurationFactorConstant:      "0.001877876953694702",
		ObserverSlashAmount:         sdkmath.NewIntFromUint64(100000000000000000),
		BallotMaturityBlocks:        100,
	}
}

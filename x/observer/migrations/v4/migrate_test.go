package v4_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	v4 "github.com/zeta-chain/zetacore/x/observer/migrations/v4"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMigrateCrosschainFlags(t *testing.T) {
	k, ctx := keepertest.ObserverKeeper(t)
	store := prefix.NewStore(ctx.KVStore(k.StoreKey()), types.KeyPrefix(types.CrosschainFlagsKey))
	legacyFlags := types.LegacyCrosschainFlags{
		IsInboundEnabled:      false,
		IsOutboundEnabled:     false,
		GasPriceIncreaseFlags: &types.DefaultGasPriceIncreaseFlags,
	}
	val := k.Codec().MustMarshal(&legacyFlags)
	store.Set([]byte{0}, val)

	err := v4.MigrateCrosschainFlags(ctx, k.StoreKey(), k.Codec())
	assert.NoError(t, err)

	flags, found := k.GetCrosschainFlags(ctx)
	assert.True(t, found)
	assert.True(t, flags.BlockHeaderVerificationFlags.IsBtcTypeChainEnabled)
	assert.True(t, flags.BlockHeaderVerificationFlags.IsEthTypeChainEnabled)
}

func TestMigrateObserverParams(t *testing.T) {
	k, ctx := keepertest.ObserverKeeper(t)

	// set core params
	previousCoreParamsList := types.CoreParamsList{
		CoreParams: []*types.CoreParams{
			sample.CoreParams(1),
			sample.CoreParams(2),
			sample.CoreParams(3),
			sample.CoreParams(4),
		},
	}
	k.SetCoreParamsList(ctx, previousCoreParamsList)

	// set observer params
	dec42, err := sdk.NewDecFromStr("0.42")
	require.NoError(t, err)
	dec43, err := sdk.NewDecFromStr("0.43")
	require.NoError(t, err)
	dec1000, err := sdk.NewDecFromStr("1000.0")
	require.NoError(t, err)
	dec1001, err := sdk.NewDecFromStr("1001.0")
	require.NoError(t, err)
	params := types.Params{
		ObserverParams: []*types.ObserverParams{
			{
				Chain:                 &common.Chain{ChainId: 2},
				BallotThreshold:       dec42,
				MinObserverDelegation: dec1000,
				IsSupported:           true,
			},
			{
				Chain:                 &common.Chain{ChainId: 3},
				BallotThreshold:       dec43,
				MinObserverDelegation: dec1001,
				IsSupported:           true,
			},
		},
	}
	k.SetParams(ctx, params)

	// perform migration
	err = v4.MigrateObserverParams(ctx, *k)
	require.NoError(t, err)

	// check core params
	newCoreParamsList, found := k.GetCoreParamsList(ctx)
	require.True(t, found)

	// unchanged values
	require.EqualValues(t, previousCoreParamsList.CoreParams[0], newCoreParamsList.CoreParams[0])
	require.EqualValues(t, previousCoreParamsList.CoreParams[3], newCoreParamsList.CoreParams[3])

	// changed values
	require.EqualValues(t, dec42, newCoreParamsList.CoreParams[1].BallotThreshold)
	require.EqualValues(t, dec1000, newCoreParamsList.CoreParams[1].MinObserverDelegation)
	require.EqualValues(t, dec43, newCoreParamsList.CoreParams[2].BallotThreshold)
	require.EqualValues(t, dec1001, newCoreParamsList.CoreParams[2].MinObserverDelegation)
	require.True(t, newCoreParamsList.CoreParams[1].IsSupported)
	require.True(t, newCoreParamsList.CoreParams[2].IsSupported)

	// check remaining values are unchanged
	previousCoreParamsList.CoreParams[1].BallotThreshold = dec42
	previousCoreParamsList.CoreParams[2].BallotThreshold = dec43
	previousCoreParamsList.CoreParams[1].MinObserverDelegation = dec1000
	previousCoreParamsList.CoreParams[2].MinObserverDelegation = dec1001
	previousCoreParamsList.CoreParams[1].IsSupported = true
	previousCoreParamsList.CoreParams[2].IsSupported = true
	require.EqualValues(t, previousCoreParamsList.CoreParams[1], newCoreParamsList.CoreParams[1])
	require.EqualValues(t, previousCoreParamsList.CoreParams[2], newCoreParamsList.CoreParams[2])
}

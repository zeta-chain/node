package v5_test

import (
	"testing"

	cosmossdk_io_math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/x/emissions/keeper"
	v5 "github.com/zeta-chain/node/x/emissions/migrations/v5"
	"github.com/zeta-chain/node/x/emissions/types"
)

type LegacyParams struct {
	ValidatorEmissionPercentage string                      `protobuf:"bytes,5,opt,name=validator_emission_percentage,json=validatorEmissionPercentage,proto3"                     json:"validator_emission_percentage,omitempty"`
	ObserverEmissionPercentage  string                      `protobuf:"bytes,6,opt,name=observer_emission_percentage,json=observerEmissionPercentage,proto3"                       json:"observer_emission_percentage,omitempty"`
	TssSignerEmissionPercentage string                      `protobuf:"bytes,7,opt,name=tss_signer_emission_percentage,json=tssSignerEmissionPercentage,proto3"                    json:"tss_signer_emission_percentage,omitempty"`
	ObserverSlashAmount         cosmossdk_io_math.Int       `protobuf:"bytes,9,opt,name=observer_slash_amount,json=observerSlashAmount,proto3,customtype=cosmossdk.io/math.Int"    json:"observer_slash_amount"`
	BallotMaturityBlocks        int64                       `protobuf:"varint,10,opt,name=ballot_maturity_blocks,json=ballotMaturityBlocks,proto3"                                 json:"ballot_maturity_blocks,omitempty"`
	BlockRewardAmount           cosmossdk_io_math.LegacyDec `protobuf:"bytes,11,opt,name=block_reward_amount,json=blockRewardAmount,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"block_reward_amount"`
}

func (l *LegacyParams) Reset() {
	*l = LegacyParams{}
}
func (*LegacyParams) ProtoMessage() {}

func (l *LegacyParams) String() string {
	out, err := yaml.Marshal(l)
	if err != nil {
		return ""
	}
	return string(out)
}

func SetLegacyParams(t *testing.T, k *keeper.Keeper, ctx sdk.Context, params LegacyParams) {
	store := ctx.KVStore(k.GetStoreKey())
	bz, err := k.GetCodec().Marshal(&params)
	require.NoError(t, err)
	store.Set(types.KeyPrefix(types.ParamsKey), bz)
}

func TestMigrateStore(t *testing.T) {
	t.Run("successfully migrate store", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)

		mainnetParams := LegacyParams{
			ValidatorEmissionPercentage: "0.75",
			ObserverEmissionPercentage:  "0.125",
			TssSignerEmissionPercentage: "0.125",
			ObserverSlashAmount:         cosmossdk_io_math.NewInt(100000000000000000),
			BallotMaturityBlocks:        100,
			BlockRewardAmount: cosmossdk_io_math.LegacyMustNewDecFromStr(
				"9620949074074074074.074070733466756687",
			),
		}
		SetLegacyParams(t, k, ctx, mainnetParams)

		err := v5.MigrateStore(ctx, k)

		// Assert
		require.NoError(t, err)
		updatedParams, found := k.GetParams(ctx)
		require.True(t, found)
		require.Equal(t, mainnetParams.ValidatorEmissionPercentage, updatedParams.ValidatorEmissionPercentage)
		require.Equal(t, mainnetParams.ObserverEmissionPercentage, updatedParams.ObserverEmissionPercentage)
		require.Equal(t, mainnetParams.TssSignerEmissionPercentage, updatedParams.TssSignerEmissionPercentage)
		require.Equal(t, mainnetParams.ObserverSlashAmount, updatedParams.ObserverSlashAmount)
		require.Equal(t, mainnetParams.BallotMaturityBlocks, updatedParams.BallotMaturityBlocks)
		require.Equal(t, mainnetParams.BlockRewardAmount, updatedParams.BlockRewardAmount)
		require.Equal(
			t,
			types.DefaultParams().PendingBallotsDeletionBufferBlocks,
			updatedParams.PendingBallotsDeletionBufferBlocks,
		)
	})

	t.Run("successfully migrate store even if tssSignerEmissionPercentage is missing", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)

		mainnetParams := LegacyParams{
			ValidatorEmissionPercentage: "0.75",
			ObserverEmissionPercentage:  "0.125",
			ObserverSlashAmount:         cosmossdk_io_math.NewInt(100000000000000000),
			BallotMaturityBlocks:        100,
			BlockRewardAmount: cosmossdk_io_math.LegacyMustNewDecFromStr(
				"9620949074074074074.074070733466756687",
			),
		}
		SetLegacyParams(t, k, ctx, mainnetParams)

		err := v5.MigrateStore(ctx, k)

		// Assert
		require.NoError(t, err)
		updatedParams, found := k.GetParams(ctx)
		require.True(t, found)
		require.Equal(t, mainnetParams.ValidatorEmissionPercentage, updatedParams.ValidatorEmissionPercentage)
		require.Equal(t, mainnetParams.ObserverEmissionPercentage, updatedParams.ObserverEmissionPercentage)
		require.Equal(t, types.DefaultParams().TssSignerEmissionPercentage, updatedParams.TssSignerEmissionPercentage)
		require.Equal(t, mainnetParams.ObserverSlashAmount, updatedParams.ObserverSlashAmount)
		require.Equal(t, mainnetParams.BallotMaturityBlocks, updatedParams.BallotMaturityBlocks)
		require.Equal(t, mainnetParams.BlockRewardAmount, updatedParams.BlockRewardAmount)
		require.Equal(
			t,
			types.DefaultParams().PendingBallotsDeletionBufferBlocks,
			updatedParams.PendingBallotsDeletionBufferBlocks,
		)
	})

	t.Run("successfully migrate store even if block reward is missing", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)

		mainnetParams := LegacyParams{
			ValidatorEmissionPercentage: "0.75",
			ObserverEmissionPercentage:  "0.125",
			TssSignerEmissionPercentage: "0.125",
			ObserverSlashAmount:         cosmossdk_io_math.NewInt(100000000000000000),
			BallotMaturityBlocks:        100,
		}
		SetLegacyParams(t, k, ctx, mainnetParams)

		err := v5.MigrateStore(ctx, k)

		// Assert
		require.NoError(t, err)
		updatedParams, found := k.GetParams(ctx)
		require.True(t, found)
		require.Equal(t, mainnetParams.ValidatorEmissionPercentage, updatedParams.ValidatorEmissionPercentage)
		require.Equal(t, mainnetParams.ObserverEmissionPercentage, updatedParams.ObserverEmissionPercentage)
		require.Equal(t, mainnetParams.TssSignerEmissionPercentage, updatedParams.TssSignerEmissionPercentage)
		require.Equal(t, mainnetParams.ObserverSlashAmount, updatedParams.ObserverSlashAmount)
		require.Equal(t, mainnetParams.BallotMaturityBlocks, updatedParams.BallotMaturityBlocks)
		require.Equal(t, types.DefaultParams().BlockRewardAmount, updatedParams.BlockRewardAmount)
		require.Equal(
			t,
			types.DefaultParams().PendingBallotsDeletionBufferBlocks,
			updatedParams.PendingBallotsDeletionBufferBlocks,
		)
	})

	t.Run("migrate store even if existing params are not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)

		err := v5.MigrateStore(ctx, k)

		// Assert
		require.NoError(t, err)
		params, found := k.GetParams(ctx)
		require.True(t, found)
		require.Equal(t, types.DefaultParams(), params)
	})
}

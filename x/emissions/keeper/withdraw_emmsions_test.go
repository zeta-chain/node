package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	emissionstypes "github.com/zeta-chain/zetacore/x/emissions/types"
)

func TestKeeper_WithdrawEmissions(t *testing.T) {
	t.Run("set withdraw emission", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		emission := sample.WithdrawEmission(t)
		k.SetWithdrawEmissions(ctx, emission)
		em, found := k.GetWithdrawEmissions(ctx, emission.Address)
		require.True(t, found)
		require.Equal(t, emission, em)
	})

	t.Run("replace withdraw emission", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		emission := sample.WithdrawEmission(t)
		k.SetWithdrawEmissions(ctx, emission)
		oldAmount := emission.Amount
		emission.Amount = sample.IntInRange(1, 100)
		k.SetWithdrawEmissions(ctx, emission)
		em, found := k.GetWithdrawEmissions(ctx, emission.Address)
		require.True(t, found)
		require.Equal(t, emission, em)
		require.NotEqual(t, oldAmount, em.Amount)
	})
	t.Run("unable to get withdraw emission which doesnt exist", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		_, found := k.GetWithdrawEmissions(ctx, sample.AccAddress())
		require.False(t, found)
	})
	t.Run("delete withdraw emission", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		emission := sample.WithdrawEmission(t)
		k.SetWithdrawEmissions(ctx, emission)
		k.DeleteWithdrawEmissions(ctx, emission.Address)
		_, found := k.GetWithdrawEmissions(ctx, emission.Address)
		require.False(t, found)
	})
	t.Run("get all withdraw emissions", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		emissions := make([]emissionstypes.WithdrawEmission, 10)
		for i := 0; i < 10; i++ {
			emission := sample.WithdrawEmission(t)
			k.SetWithdrawEmissions(ctx, emission)
			emissions[i] = emission
		}
		allEmissions := k.GetAllWithdrawEmissions(ctx)
		require.ElementsMatch(t, emissions, allEmissions)
	})
}

func TestKeeper_CreateWithdrawEmissions(t *testing.T) {
	t.Run("create withdraw emission", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		withdrawableEmission := sample.WithdrawableEmissions(t)
		k.SetWithdrawableEmission(ctx, withdrawableEmission)
		err := k.CreateWithdrawEmissions(ctx, withdrawableEmission.Address, withdrawableEmission.Amount)
		require.NoError(t, err)
		em, found := k.GetWithdrawEmissions(ctx, withdrawableEmission.Address)
		require.True(t, found)
		require.Equal(t, withdrawableEmission.Amount, em.Amount)
	})

	t.Run("create withdraw for max available withdrawable emission", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		withdrawableEmission := sample.WithdrawableEmissions(t)
		k.SetWithdrawableEmission(ctx, withdrawableEmission)
		err := k.CreateWithdrawEmissions(ctx, withdrawableEmission.Address, withdrawableEmission.Amount.Add(sample.IntInRange(1, 100)))
		require.NoError(t, err)
		em, found := k.GetWithdrawEmissions(ctx, withdrawableEmission.Address)
		require.True(t, found)
		require.Equal(t, withdrawableEmission.Amount, em.Amount)
	})

	t.Run("unable to create withdraw for zero amount", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		withdrawableEmission := sample.WithdrawableEmissions(t)
		withdrawableEmission.Amount = sdkmath.ZeroInt()
		k.SetWithdrawableEmission(ctx, withdrawableEmission)
		err := k.CreateWithdrawEmissions(ctx, withdrawableEmission.Address, sdkmath.ZeroInt())
		require.ErrorIs(t, err, emissionstypes.ErrNotEnoughEmissionsAvailable)
	})

	t.Run("unable to create withdraw for negative amount", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		withdrawableEmission := sample.WithdrawableEmissions(t)
		withdrawableEmission.Amount = sdkmath.NewInt(-1)
		k.SetWithdrawableEmission(ctx, withdrawableEmission)
		err := k.CreateWithdrawEmissions(ctx, withdrawableEmission.Address, sdkmath.NewInt(-1))
		require.ErrorIs(t, err, emissionstypes.ErrNotEnoughEmissionsAvailable)
	})

	t.Run("unable to create withdraw for non existing withdrawable emission", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		err := k.CreateWithdrawEmissions(ctx, sample.AccAddress(), sdkmath.NewInt(1))
		require.ErrorIs(t, err, emissionstypes.ErrEmissionsNotFound)
	})
}

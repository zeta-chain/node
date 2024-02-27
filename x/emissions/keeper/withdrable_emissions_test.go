package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	emissionstypes "github.com/zeta-chain/zetacore/x/emissions/types"
)

func Test_WithdrawableEmissions(t *testing.T) {
	t.Run("set valid withdrawable emission", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		we := sample.WithdrawableEmissions(t)
		k.SetWithdrawableEmission(ctx, we)
		we2, found := k.GetWithdrawableEmission(ctx, we.Address)
		require.True(t, found)
		require.Equal(t, we, we2)
	})

	t.Run("get all withdrawable emissions", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		withdrawableEmissionslist := make([]emissionstypes.WithdrawableEmissions, 10)
		for i := 0; i < 10; i++ {
			we := sample.WithdrawableEmissions(t)
			k.SetWithdrawableEmission(ctx, we)
			withdrawableEmissionslist[i] = we
		}
		allWithdrawableEmissions := k.GetAllWithdrawableEmission(ctx)
		require.ElementsMatch(t, withdrawableEmissionslist, allWithdrawableEmissions)
	})

	t.Run("add observer emission to an existing value", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		we := sample.WithdrawableEmissions(t)
		k.SetWithdrawableEmission(ctx, we)
		amount := sample.IntInRange(1, 100)
		k.AddObserverEmission(ctx, we.Address, amount)
		we2, found := k.GetWithdrawableEmission(ctx, we.Address)
		require.True(t, found)
		require.Equal(t, we.Amount.Add(amount), we2.Amount)
	})

	t.Run("add observer emission to a non-existing value", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		we := sample.WithdrawableEmissions(t)
		amount := sample.IntInRange(1, 100)
		k.AddObserverEmission(ctx, we.Address, amount)
		we2, found := k.GetWithdrawableEmission(ctx, we.Address)
		require.True(t, found)
		require.Equal(t, amount, we2.Amount)
	})

	t.Run("remove observer emission", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		we := sample.WithdrawableEmissions(t)
		k.AddObserverEmission(ctx, we.Address, we.Amount)

		we2, found := k.GetWithdrawableEmission(ctx, we.Address)
		require.True(t, found)
		require.Equal(t, we.Amount, we2.Amount)
	})

	t.Run("remove observer emission with not enough emissions available", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		we := sample.WithdrawableEmissions(t)
		k.SetWithdrawableEmission(ctx, we)

		err := k.RemoveObserverEmission(ctx, we.Address, we.Amount.Add(sdkmath.OneInt()))
		require.ErrorIs(t, err, emissionstypes.ErrNotEnoughEmissionsAvailable)
	})

	t.Run("SlashObserverEmission", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		we := sample.WithdrawableEmissions(t)
		k.SetWithdrawableEmission(ctx, we)
		k.SlashObserverEmission(ctx, we.Address, we.Amount.Sub(sdkmath.OneInt()))
		we2, found := k.GetWithdrawableEmission(ctx, we.Address)
		require.True(t, found)
		require.Equal(t, sdkmath.OneInt(), we2.Amount)
	})

	t.Run("SlashObserverEmission to zero if not enough emissions available", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		we := sample.WithdrawableEmissions(t)
		k.SetWithdrawableEmission(ctx, we)
		k.SlashObserverEmission(ctx, we.Address, we.Amount.Add(sdkmath.OneInt()))
		we2, found := k.GetWithdrawableEmission(ctx, we.Address)
		require.True(t, found)
		require.Equal(t, sdkmath.ZeroInt(), we2.Amount)
	})

	t.Run("try slashing non existing observer emission", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		address := sample.AccAddress()
		k.SlashObserverEmission(ctx, address, sdkmath.OneInt())
		we, found := k.GetWithdrawableEmission(ctx, address)
		require.True(t, found)
		require.Equal(t, sdkmath.ZeroInt(), we.Amount)
	})

}

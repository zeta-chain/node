package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	emissionstypes "github.com/zeta-chain/node/x/emissions/types"
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

	t.Run("slash observer emission", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		we := sample.WithdrawableEmissions(t)
		k.SetWithdrawableEmission(ctx, we)
		k.SlashObserverEmission(ctx, we.Address, we.Amount.Sub(sdkmath.OneInt()))
		we2, found := k.GetWithdrawableEmission(ctx, we.Address)
		require.True(t, found)
		require.Equal(t, sdkmath.OneInt(), we2.Amount)
	})

}
func TestKeeper_AddObserverEmission(t *testing.T) {
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
}

func TestKeeper_SlashWithdrawableEmission(t *testing.T) {
	t.Run("successfully slash withdrawable emission", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		we := sample.WithdrawableEmissions(t)
		k.SetWithdrawableEmission(ctx, we)
		k.SlashObserverEmission(ctx, we.Address, we.Amount.Sub(sdkmath.OneInt()))
		we2, found := k.GetWithdrawableEmission(ctx, we.Address)
		require.True(t, found)
		require.Equal(t, sdkmath.OneInt(), we2.Amount)
	})

	t.Run("slash observer emission to zero if not enough emissions available", func(t *testing.T) {
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

func TestKeeper_RemoveObserverEmission(t *testing.T) {
	t.Run("remove all observer emission successfully", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		we := sample.WithdrawableEmissions(t)
		k.SetWithdrawableEmission(ctx, we)
		err := k.RemoveWithdrawableEmission(ctx, we.Address, we.Amount)
		require.NoError(t, err)
		we2, found := k.GetWithdrawableEmission(ctx, we.Address)
		require.True(t, found)
		require.Equal(t, sdkmath.ZeroInt(), we2.Amount)
	})

	t.Run("unable to remove observer emission if requested amount is higher than available", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		we := sample.WithdrawableEmissions(t)
		k.SetWithdrawableEmission(ctx, we)
		err := k.RemoveWithdrawableEmission(ctx, we.Address, we.Amount.Add(sdkmath.OneInt()))
		require.ErrorIs(t, err, emissionstypes.ErrInvalidAmount)
		require.ErrorContains(t, err, "amount to be removed is greater than the available withdrawable emission")
	})

	t.Run("unable to remove non-existent emission ", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		address := sample.AccAddress()
		err := k.RemoveWithdrawableEmission(ctx, address, sdkmath.ZeroInt())
		require.Error(t, err, emissionstypes.ErrEmissionsNotFound)
	})

	t.Run("remove all portion of observer emission successfully", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		we := sample.WithdrawableEmissions(t)
		k.SetWithdrawableEmission(ctx, we)
		withdrawAmount := we.Amount.Quo(sdkmath.NewInt(2))
		err := k.RemoveWithdrawableEmission(ctx, we.Address, withdrawAmount)
		require.NoError(t, err)
		we2, found := k.GetWithdrawableEmission(ctx, we.Address)
		require.True(t, found)
		require.Equal(t, we.Amount.Sub(withdrawAmount).String(), we2.Amount.String())
	})

	t.Run("unable to withdraw negative amount", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		we := sample.WithdrawableEmissions(t)
		k.SetWithdrawableEmission(ctx, we)
		err := k.RemoveWithdrawableEmission(ctx, we.Address, sdkmath.NewInt(-1))
		require.ErrorIs(t, err, emissionstypes.ErrInvalidAmount)
	})

	t.Run("unable to withdraw zero amount", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		we := sample.WithdrawableEmissions(t)
		k.SetWithdrawableEmission(ctx, we)
		err := k.RemoveWithdrawableEmission(ctx, we.Address, sdkmath.ZeroInt())
		require.ErrorIs(t, err, emissionstypes.ErrInvalidAmount)
	})
}

package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	"testing"

	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
)

func TestKeeper_UpdateZRC20WithdrawFee(t *testing.T) {
	t.Run("can update the withdraw fee", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.FungibleKeeper(t)
		chainID := getValidChainID(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		// set coin admin
		admin := sample.AccAddress()
		setAdminDeployFungibleCoin(ctx, zk, admin)

		// deploy the system contract and a ZRC20 contract
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "alpha", "alpha")

		// initial protocol fee is zero
		fee, err := k.QueryProtocolFlatFee(ctx, zrc20)
		require.NoError(t, err)
		require.Zero(t, fee.Uint64())

		// can update the fee
		_, err = k.UpdateZRC20WithdrawFee(ctx, types.NewMsgUpdateZRC20WithdrawFee(
			admin,
			zrc20.String(),
			math.NewUint(42),
		))
		require.NoError(t, err)

		// can query the updated fee
		fee, err = k.QueryProtocolFlatFee(ctx, zrc20)
		require.NoError(t, err)
		require.Equal(t, uint64(42), fee.Uint64())
	})
}

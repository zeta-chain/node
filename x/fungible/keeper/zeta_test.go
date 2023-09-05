package keeper_test

import (
	"errors"
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	testkeeper "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestKeeper_MintZetaToEVMAccount(t *testing.T) {
	t.Run("should mint the token in the specified balance", func(t *testing.T) {
		k, ctx, sdkk := testkeeper.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		acc := sample.Bech32AccAddress()
		bal := sdkk.BankKeeper.GetBalance(ctx, acc, config.BaseDenom)
		require.True(t, bal.IsZero())

		err := k.MintZetaToEVMAccount(ctx, acc, big.NewInt(42))
		require.NoError(t, err)
		bal = sdkk.BankKeeper.GetBalance(ctx, acc, config.BaseDenom)
		require.True(t, bal.Amount.Equal(sdk.NewInt(42)))
	})

	coins42 := sdk.NewCoins(sdk.NewCoin(config.BaseDenom, sdk.NewInt(42)))

	t.Run("should fail if minting fail", func(t *testing.T) {
		k, ctx := testkeeper.FungibleKeeperAllMocks(t)

		mockBankKeeper := testkeeper.GetFungibleBankMock(t, k)

		mockBankKeeper.On(
			"MintCoins",
			ctx,
			types.ModuleName,
			coins42,
		).Return(errors.New("error"))

		err := k.MintZetaToEVMAccount(ctx, sample.Bech32AccAddress(), big.NewInt(42))
		require.Error(t, err)

		mockBankKeeper.AssertExpectations(t)
	})

	t.Run("should fail if sending coins fail", func(t *testing.T) {
		k, ctx := testkeeper.FungibleKeeperAllMocks(t)
		acc := sample.Bech32AccAddress()

		mockBankKeeper := testkeeper.GetFungibleBankMock(t, k)

		mockBankKeeper.On(
			"MintCoins",
			ctx,
			types.ModuleName,
			coins42,
		).Return(nil)

		mockBankKeeper.On(
			"SendCoinsFromModuleToAccount",
			ctx,
			types.ModuleName,
			acc,
			coins42,
		).Return(errors.New("error"))

		err := k.MintZetaToEVMAccount(ctx, acc, big.NewInt(42))
		require.Error(t, err)

		mockBankKeeper.AssertExpectations(t)
	})
}

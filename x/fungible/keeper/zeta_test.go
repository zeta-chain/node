package keeper_test

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/zeta-chain/zetacore/x/fungible/keeper"
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
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		acc := sample.Bech32AccAddress()
		bal := sdkk.BankKeeper.GetBalance(ctx, acc, config.BaseDenom)
		require.True(t, bal.IsZero())

		err := k.MintZetaToEVMAccount(ctx, acc, big.NewInt(42))
		require.NoError(t, err)
		bal = sdkk.BankKeeper.GetBalance(ctx, acc, config.BaseDenom)
		require.True(t, bal.Amount.Equal(sdk.NewInt(42)))
	})

	t.Run("mint the token to reach max supply", func(t *testing.T) {
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		acc := sample.Bech32AccAddress()
		bal := sdkk.BankKeeper.GetBalance(ctx, acc, config.BaseDenom)
		require.True(t, bal.IsZero())

		zetaMaxSupply, ok := sdk.NewIntFromString(keeper.ZETAMaxSupplyStr)
		require.True(t, ok)

		supply := sdkk.BankKeeper.GetSupply(ctx, config.BaseDenom).Amount

		newAmount := zetaMaxSupply.Sub(supply)

		err := k.MintZetaToEVMAccount(ctx, acc, newAmount.BigInt())
		require.NoError(t, err)
		bal = sdkk.BankKeeper.GetBalance(ctx, acc, config.BaseDenom)
		require.True(t, bal.Amount.Equal(newAmount))
	})

	t.Run("can't mint more than max supply", func(t *testing.T) {
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		acc := sample.Bech32AccAddress()
		bal := sdkk.BankKeeper.GetBalance(ctx, acc, config.BaseDenom)
		require.True(t, bal.IsZero())

		zetaMaxSupply, ok := sdk.NewIntFromString(keeper.ZETAMaxSupplyStr)
		require.True(t, ok)

		supply := sdkk.BankKeeper.GetSupply(ctx, config.BaseDenom).Amount

		newAmount := zetaMaxSupply.Sub(supply).Add(sdk.NewInt(1))

		err := k.MintZetaToEVMAccount(ctx, acc, newAmount.BigInt())
		require.ErrorIs(t, err, types.ErrMaxSupplyReached)
	})

	coins42 := sdk.NewCoins(sdk.NewCoin(config.BaseDenom, sdk.NewInt(42)))

	t.Run("should fail if minting fail", func(t *testing.T) {
		k, ctx := testkeeper.FungibleKeeperAllMocks(t)

		mockBankKeeper := testkeeper.GetFungibleBankMock(t, k)

		mockBankKeeper.On("GetSupply", ctx, mock.Anything, mock.Anything).
			Return(sdk.NewCoin(config.BaseDenom, sdk.NewInt(0))).
			Once()
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

		mockBankKeeper.On("GetSupply", ctx, mock.Anything, mock.Anything).
			Return(sdk.NewCoin(config.BaseDenom, sdk.NewInt(0))).
			Once()
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

func TestKeeper_MintZetaToFungibleModule(t *testing.T) {
	t.Run("should mint the token in the specified balance", func(t *testing.T) {
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		acc := k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName).GetAddress()

		bal := sdkk.BankKeeper.GetBalance(ctx, acc, config.BaseDenom)
		require.True(t, bal.IsZero())

		err := k.MintZetaToEVMAccount(ctx, acc, big.NewInt(42))
		require.NoError(t, err)
		bal = sdkk.BankKeeper.GetBalance(ctx, acc, config.BaseDenom)
		require.True(t, bal.Amount.Equal(sdk.NewInt(42)))
	})

	t.Run("can't mint more than max supply", func(t *testing.T) {
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		zetaMaxSupply, ok := sdk.NewIntFromString(keeper.ZETAMaxSupplyStr)
		require.True(t, ok)

		supply := sdkk.BankKeeper.GetSupply(ctx, config.BaseDenom).Amount

		newAmount := zetaMaxSupply.Sub(supply).Add(sdk.NewInt(1))

		err := k.MintZetaToFungibleModule(ctx, newAmount.BigInt())
		require.ErrorIs(t, err, types.ErrMaxSupplyReached)
	})
}

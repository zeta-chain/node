package keeper_test

import (
	"errors"
	"math/big"
	"testing"

	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/stretchr/testify/mock"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/cmd/zetacored/config"
	testkeeper "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/fungible/keeper"
	"github.com/zeta-chain/node/x/fungible/types"
)

func TestKeeper_ExecuteWithMintedZeta(t *testing.T) {
	t.Run("should execute the operation with minted ZETA", func(t *testing.T) {
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		executedOperation := false
		operationNoErr := func(sdk.Context) (*evmtypes.MsgEthereumTxResponse, bool, error) {
			executedOperation = true
			return nil, true, nil
		}
		amount := int64(42)

		_, ok, err := k.ExecuteWithMintedZeta(ctx, big.NewInt(amount), operationNoErr)
		require.NoError(t, err)
		require.True(t, ok)
		require.True(t, executedOperation)

		require.Equal(t, amount, sdkk.BankKeeper.GetBalance(ctx, types.ModuleAddress, config.BaseDenom).Amount.Int64())
	})

	t.Run("should not mint zeta if operation fails", func(t *testing.T) {
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		executedOperation := false
		operationErr := func(sdk.Context) (*evmtypes.MsgEthereumTxResponse, bool, error) {
			executedOperation = true
			return nil, false, errors.New("operation failed")
		}
		amount := int64(42)

		_, ok, err := k.ExecuteWithMintedZeta(ctx, big.NewInt(amount), operationErr)
		require.Error(t, err)
		require.False(t, ok)
		require.True(t, executedOperation)

		require.Equal(t, sdkmath.ZeroInt(), sdkk.BankKeeper.GetBalance(ctx, types.ModuleAddress, config.BaseDenom).Amount)
	})
}

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
		require.True(t, bal.Amount.Equal(sdkmath.NewInt(42)))
	})

	t.Run("mint the token to reach max supply", func(t *testing.T) {
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		acc := sample.Bech32AccAddress()
		bal := sdkk.BankKeeper.GetBalance(ctx, acc, config.BaseDenom)
		require.True(t, bal.IsZero())

		zetaMaxSupply, ok := sdkmath.NewIntFromString(keeper.ZETAMaxSupplyStr)
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

		zetaMaxSupply, ok := sdkmath.NewIntFromString(keeper.ZETAMaxSupplyStr)
		require.True(t, ok)

		supply := sdkk.BankKeeper.GetSupply(ctx, config.BaseDenom).Amount

		newAmount := zetaMaxSupply.Sub(supply).Add(sdkmath.NewInt(1))

		err := k.MintZetaToEVMAccount(ctx, acc, newAmount.BigInt())
		require.ErrorIs(t, err, types.ErrMaxSupplyReached)
	})

	coins42 := sdk.NewCoins(sdk.NewCoin(config.BaseDenom, sdkmath.NewInt(42)))

	t.Run("should fail if minting fail", func(t *testing.T) {
		k, ctx := testkeeper.FungibleKeeperAllMocks(t)

		mockBankKeeper := testkeeper.GetFungibleBankMock(t, k)

		mockBankKeeper.On("GetSupply", ctx, mock.Anything, mock.Anything).
			Return(sdk.NewCoin(config.BaseDenom, sdkmath.NewInt(0))).
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
			Return(sdk.NewCoin(config.BaseDenom, sdkmath.NewInt(0))).
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
		require.True(t, bal.Amount.Equal(sdkmath.NewInt(42)))
	})

	t.Run("can't mint more than max supply", func(t *testing.T) {
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		zetaMaxSupply, ok := sdkmath.NewIntFromString(keeper.ZETAMaxSupplyStr)
		require.True(t, ok)

		supply := sdkk.BankKeeper.GetSupply(ctx, config.BaseDenom).Amount

		newAmount := zetaMaxSupply.Sub(supply).Add(sdkmath.NewInt(1))

		err := k.MintZetaToFungibleModule(ctx, newAmount.BigInt())
		require.ErrorIs(t, err, types.ErrMaxSupplyReached)
	})
}

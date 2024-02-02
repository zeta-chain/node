package keeper_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestKeeper_GetSystemContract(t *testing.T) {
	k, ctx, _, _ := keepertest.FungibleKeeper(t)
	k.SetSystemContract(ctx, types.SystemContract{SystemContract: "test"})
	val, found := k.GetSystemContract(ctx)
	assert.True(t, found)
	assert.Equal(t, types.SystemContract{SystemContract: "test"}, val)

	// can remove contract
	k.RemoveSystemContract(ctx)
	_, found = k.GetSystemContract(ctx)
	assert.False(t, found)
}

func TestKeeper_GetSystemContractAddress(t *testing.T) {
	k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

	_, err := k.GetSystemContractAddress(ctx)
	assert.Error(t, err)
	assert.ErrorIs(t, err, types.ErrStateVariableNotFound)

	_, _, _, _, systemContract := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
	found, err := k.GetSystemContractAddress(ctx)
	assert.NoError(t, err)
	assert.Equal(t, systemContract, found)
}

func TestKeeper_GetWZetaContractAddress(t *testing.T) {
	k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

	_, err := k.GetWZetaContractAddress(ctx)
	assert.Error(t, err)
	assert.ErrorIs(t, err, types.ErrStateVariableNotFound)

	wzeta, _, _, _, _ := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
	found, err := k.GetWZetaContractAddress(ctx)
	assert.NoError(t, err)
	assert.Equal(t, wzeta, found)
}

func TestKeeper_GetUniswapV2FactoryAddress(t *testing.T) {
	k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

	_, err := k.GetUniswapV2FactoryAddress(ctx)
	assert.Error(t, err)
	assert.ErrorIs(t, err, types.ErrStateVariableNotFound)

	_, factory, _, _, _ := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
	found, err := k.GetUniswapV2FactoryAddress(ctx)
	assert.NoError(t, err)
	assert.Equal(t, factory, found)
}

func TestKeeper_GetUniswapV2Router02Address(t *testing.T) {
	k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

	_, err := k.GetUniswapV2Router02Address(ctx)
	assert.Error(t, err)
	assert.ErrorIs(t, err, types.ErrStateVariableNotFound)

	_, _, router, _, _ := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
	found, err := k.GetUniswapV2Router02Address(ctx)
	assert.NoError(t, err)
	assert.Equal(t, router, found)
}

func TestKeeper_CallWZetaDeposit(t *testing.T) {
	k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

	// mint tokens
	addr := sample.Bech32AccAddress()
	ethAddr := common.BytesToAddress(addr.Bytes())
	coins := sample.Coins()
	err := sdkk.BankKeeper.MintCoins(ctx, types.ModuleName, sample.Coins())
	assert.NoError(t, err)
	err = sdkk.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, coins)
	assert.NoError(t, err)

	// fail if no system contract
	err = k.CallWZetaDeposit(ctx, ethAddr, big.NewInt(42))
	assert.Error(t, err)

	deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

	// deposit
	err = k.CallWZetaDeposit(ctx, ethAddr, big.NewInt(42))
	assert.NoError(t, err)

	balance, err := k.QueryWZetaBalanceOf(ctx, ethAddr)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(42), balance)
}

func TestKeeper_QuerySystemContractGasCoinZRC20(t *testing.T) {
	k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
	chainID := getValidChainID(t)

	_, err := k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chainID))
	assert.Error(t, err)
	assert.ErrorIs(t, err, types.ErrStateVariableNotFound)

	deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
	zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

	found, err := k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chainID))
	assert.NoError(t, err)
	assert.Equal(t, zrc20, found)
}

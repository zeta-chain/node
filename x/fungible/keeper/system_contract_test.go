package keeper_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestKeeper_GetSystemContract(t *testing.T) {
	k, ctx, _, _ := keepertest.FungibleKeeper(t)
	k.SetSystemContract(ctx, types.SystemContract{SystemContract: "test"})
	val, found := k.GetSystemContract(ctx)
	require.True(t, found)
	require.Equal(t, types.SystemContract{SystemContract: "test"}, val)

	// can remove contract
	k.RemoveSystemContract(ctx)
	_, found = k.GetSystemContract(ctx)
	require.False(t, found)
}

func TestKeeper_GetSystemContractAddress(t *testing.T) {
	k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

	_, err := k.GetSystemContractAddress(ctx)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrStateVariableNotFound)

	_, _, _, _, systemContract := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
	found, err := k.GetSystemContractAddress(ctx)
	require.NoError(t, err)
	require.Equal(t, systemContract, found)
}

func TestKeeper_GetWZetaContractAddress(t *testing.T) {
	k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

	_, err := k.GetWZetaContractAddress(ctx)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrStateVariableNotFound)

	wzeta, _, _, _, _ := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
	found, err := k.GetWZetaContractAddress(ctx)
	require.NoError(t, err)
	require.Equal(t, wzeta, found)
}

func TestKeeper_GetWZetaContractAddressFails(t *testing.T) {
	k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
		UseEVMMock: true,
	})
	mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

	_, err := k.GetWZetaContractAddress(ctx)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrStateVariableNotFound)

	mockEVMKeeper.SetupMockEVMKeeperForSystemContractDeployment()
	deploySystemContracts(t, ctx, k, mockEVMKeeper)

	mockEVMKeeper.MockEVMFailCallOnce()
	_, err = k.GetWZetaContractAddress(ctx)
	require.ErrorIs(t, err, types.ErrContractCall)
}

func TestKeeper_GetWZetaContractAddressFailsToUnpack(t *testing.T) {
	k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
		UseEVMMock: true,
	})
	mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

	_, err := k.GetWZetaContractAddress(ctx)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrStateVariableNotFound)

	mockEVMKeeper.SetupMockEVMKeeperForSystemContractDeployment()
	deploySystemContracts(t, ctx, k, mockEVMKeeper)

	mockEVMKeeper.MockEVMSuccessCallOnce()
	_, err = k.GetWZetaContractAddress(ctx)
	require.ErrorIs(t, err, types.ErrABIUnpack)
}

func TestKeeper_GetUniswapV2FactoryAddress(t *testing.T) {
	k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

	_, err := k.GetUniswapV2FactoryAddress(ctx)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrStateVariableNotFound)

	_, factory, _, _, _ := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
	found, err := k.GetUniswapV2FactoryAddress(ctx)
	require.NoError(t, err)
	require.Equal(t, factory, found)
}

func TestKeeper_GetUniswapV2FactoryAddressFails(t *testing.T) {
	k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
		UseEVMMock: true,
	})
	mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

	_, err := k.GetUniswapV2FactoryAddress(ctx)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrStateVariableNotFound)

	mockEVMKeeper.SetupMockEVMKeeperForSystemContractDeployment()
	deploySystemContracts(t, ctx, k, mockEVMKeeper)

	mockEVMKeeper.MockEVMFailCallOnce()
	_, err = k.GetUniswapV2FactoryAddress(ctx)
	require.ErrorIs(t, err, types.ErrContractCall)
}

func TestKeeper_GetUniswapV2FactoryAddressFailsToUnpack(t *testing.T) {
	k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
		UseEVMMock: true,
	})
	mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

	_, err := k.GetUniswapV2FactoryAddress(ctx)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrStateVariableNotFound)

	mockEVMKeeper.SetupMockEVMKeeperForSystemContractDeployment()
	deploySystemContracts(t, ctx, k, mockEVMKeeper)

	mockEVMKeeper.MockEVMSuccessCallOnce()
	_, err = k.GetUniswapV2FactoryAddress(ctx)
	require.ErrorIs(t, err, types.ErrABIUnpack)
}

func TestKeeper_GetUniswapV2Router02Address(t *testing.T) {
	k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

	_, err := k.GetUniswapV2Router02Address(ctx)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrStateVariableNotFound)

	_, _, router, _, _ := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
	found, err := k.GetUniswapV2Router02Address(ctx)
	require.NoError(t, err)
	require.Equal(t, router, found)
}

func TestKeeper_GetUniswapV2Router02AddressFails(t *testing.T) {
	k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
		UseEVMMock: true,
	})
	mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

	_, err := k.GetUniswapV2Router02Address(ctx)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrStateVariableNotFound)

	mockEVMKeeper.SetupMockEVMKeeperForSystemContractDeployment()
	deploySystemContracts(t, ctx, k, mockEVMKeeper)

	mockEVMKeeper.MockEVMFailCallOnce()
	_, err = k.GetUniswapV2Router02Address(ctx)
	require.ErrorIs(t, err, types.ErrContractCall)
}

func TestKeeper_GetUniswapV2Router02AddressFailsToUnpack(t *testing.T) {
	k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
		UseEVMMock: true,
	})
	mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

	_, err := k.GetUniswapV2Router02Address(ctx)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrStateVariableNotFound)

	mockEVMKeeper.SetupMockEVMKeeperForSystemContractDeployment()
	deploySystemContracts(t, ctx, k, mockEVMKeeper)

	mockEVMKeeper.MockEVMSuccessCallOnce()
	_, err = k.GetUniswapV2Router02Address(ctx)
	require.ErrorIs(t, err, types.ErrABIUnpack)
}

func TestKeeper_CallWZetaDeposit(t *testing.T) {
	k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

	// mint tokens
	addr := sample.Bech32AccAddress()
	ethAddr := common.BytesToAddress(addr.Bytes())
	coins := sample.Coins()
	err := sdkk.BankKeeper.MintCoins(ctx, types.ModuleName, sample.Coins())
	require.NoError(t, err)
	err = sdkk.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, coins)
	require.NoError(t, err)

	// fail if no system contract
	err = k.CallWZetaDeposit(ctx, ethAddr, big.NewInt(42))
	require.Error(t, err)

	deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

	// deposit
	err = k.CallWZetaDeposit(ctx, ethAddr, big.NewInt(42))
	require.NoError(t, err)

	balance, err := k.QueryWZetaBalanceOf(ctx, ethAddr)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(42), balance)
}

func TestKeeper_CallWZetaDepositFails(t *testing.T) {
	k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
		UseEVMMock: true,
	})
	mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

	// mint tokens
	addr := sample.Bech32AccAddress()
	ethAddr := common.BytesToAddress(addr.Bytes())
	coins := sample.Coins()
	err := sdkk.BankKeeper.MintCoins(ctx, types.ModuleName, sample.Coins())
	require.NoError(t, err)
	err = sdkk.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, coins)
	require.NoError(t, err)

	// fail if no system contract
	err = k.CallWZetaDeposit(ctx, ethAddr, big.NewInt(42))
	require.Error(t, err)

	mockEVMKeeper.SetupMockEVMKeeperForSystemContractDeployment()
	deploySystemContracts(t, ctx, k, mockEVMKeeper)

	// deposit
	mockEVMKeeper.MockEVMFailCallOnce()
	err = k.CallWZetaDeposit(ctx, ethAddr, big.NewInt(42))
	require.ErrorIs(t, err, types.ErrContractCall)
}

func TestKeeper_QueryWZetaBalanceOfFails(t *testing.T) {
	k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
		UseEVMMock: true,
	})
	mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

	// mint tokens
	addr := sample.Bech32AccAddress()
	ethAddr := common.BytesToAddress(addr.Bytes())
	coins := sample.Coins()
	err := sdkk.BankKeeper.MintCoins(ctx, types.ModuleName, sample.Coins())
	require.NoError(t, err)
	err = sdkk.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, coins)
	require.NoError(t, err)

	// fail if no system contract
	err = k.CallWZetaDeposit(ctx, ethAddr, big.NewInt(42))
	require.Error(t, err)

	mockEVMKeeper.SetupMockEVMKeeperForSystemContractDeployment()
	deploySystemContracts(t, ctx, k, mockEVMKeeper)

	// deposit
	mockEVMKeeper.MockEVMFailCallOnce()
	_, err = k.QueryWZetaBalanceOf(ctx, ethAddr)
	require.ErrorIs(t, err, types.ErrContractCall)
}

func TestKeeper_QueryWZetaBalanceOfFailsToUnpack(t *testing.T) {
	k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
		UseEVMMock: true,
	})
	mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

	// mint tokens
	addr := sample.Bech32AccAddress()
	ethAddr := common.BytesToAddress(addr.Bytes())
	coins := sample.Coins()
	err := sdkk.BankKeeper.MintCoins(ctx, types.ModuleName, sample.Coins())
	require.NoError(t, err)
	err = sdkk.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, coins)
	require.NoError(t, err)

	// fail if no system contract
	err = k.CallWZetaDeposit(ctx, ethAddr, big.NewInt(42))
	require.Error(t, err)

	mockEVMKeeper.SetupMockEVMKeeperForSystemContractDeployment()
	deploySystemContracts(t, ctx, k, mockEVMKeeper)

	// deposit
	mockEVMKeeper.MockEVMSuccessCallOnce()
	_, err = k.QueryWZetaBalanceOf(ctx, ethAddr)
	require.ErrorIs(t, err, types.ErrABIUnpack)
}

func TestKeeper_QuerySystemContractGasCoinZRC20(t *testing.T) {
	k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
	chainID := getValidChainID(t)

	_, err := k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chainID))
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrStateVariableNotFound)

	deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
	zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

	found, err := k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chainID))
	require.NoError(t, err)
	require.Equal(t, zrc20, found)
}

func TestKeeper_QuerySystemContractGasCoinZRC20Fails(t *testing.T) {
	k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
		UseEVMMock: true,
	})
	mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

	chainID := getValidChainID(t)

	_, err := k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chainID))
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrStateVariableNotFound)

	mockEVMKeeper.SetupMockEVMKeeperForSystemContractDeployment()
	deploySystemContracts(t, ctx, k, mockEVMKeeper)

	mockEVMKeeper.MockEVMFailCallOnce()
	_, err = k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chainID))
	require.ErrorIs(t, err, types.ErrContractCall)
}

func TestKeeper_QuerySystemContractGasCoinZRC20FailsToUnpack(t *testing.T) {
	k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
		UseEVMMock: true,
	})
	mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

	chainID := getValidChainID(t)

	_, err := k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chainID))
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrStateVariableNotFound)

	mockEVMKeeper.SetupMockEVMKeeperForSystemContractDeployment()
	deploySystemContracts(t, ctx, k, mockEVMKeeper)

	mockEVMKeeper.MockEVMSuccessCallOnce()
	_, err = k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chainID))
	require.ErrorIs(t, err, types.ErrABIUnpack)
}

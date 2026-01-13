package keeper_test

import (
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/fungible/types"
)

func TestKeeper_GetSystemContract(t *testing.T) {
	t.Run("should get and remove system contract", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		system := *sample.SystemContract()
		k.SetSystemContract(ctx, system)
		val, found := k.GetSystemContract(ctx)
		require.True(t, found)
		require.Equal(t, system, val)

		// can remove contract
		k.RemoveSystemContract(ctx)
		_, found = k.GetSystemContract(ctx)
		require.False(t, found)
	})
}

func TestKeeper_GetSystemContractAddress(t *testing.T) {
	t.Run("should fail to get system contract address if system contracts are not deployed", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		_, err := k.GetSystemContractAddress(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrStateVariableNotFound)
	})

	t.Run("should get system contract address if system contracts are deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		_, _, _, _, systemContract := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		found, err := k.GetSystemContractAddress(ctx)
		require.NoError(t, err)
		require.Equal(t, systemContract, found)
	})
}

func TestKeeper_GetWZetaContractAddress(t *testing.T) {
	t.Run("should fail to get wzeta contract address if system contracts are not deployed", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		_, err := k.GetWZetaContractAddress(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrStateVariableNotFound)
	})

	t.Run("should get wzeta contract address if system contracts are deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		wzeta, _, _, _, _ := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		found, err := k.GetWZetaContractAddress(ctx)
		require.NoError(t, err)
		require.Equal(t, wzeta, found)
	})

	t.Run("should fail if wzeta not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsConfigurable(t, ctx, k, sdkk.EvmKeeper, &SystemContractDeployConfig{
			DeployWZeta:            false,
			DeployUniswapV2Router:  true,
			DeployUniswapV2Factory: true,
		})

		_, err := k.GetWZetaContractAddress(ctx)
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})

	t.Run("should fail if evm call fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsWithMockEvmKeeper(t, ctx, k, mockEVMKeeper)

		mockEVMKeeper.MockEVMFailCallOnce()

		_, err := k.GetWZetaContractAddress(ctx)
		require.ErrorIs(t, err, types.ErrContractCall)
	})

	t.Run("should fail if abi unpack fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsWithMockEvmKeeper(t, ctx, k, mockEVMKeeper)

		mockEVMKeeper.MockEVMSuccessCallOnce()

		_, err := k.GetWZetaContractAddress(ctx)
		require.ErrorIs(t, err, types.ErrABIUnpack)
	})
}

func TestKeeper_GetUniswapV2FactoryAddress(t *testing.T) {
	t.Run(
		"should fail to get uniswapfactory contract address if system contracts are not deployed",
		func(t *testing.T) {
			k, ctx, _, _ := keepertest.FungibleKeeper(t)
			k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

			_, err := k.GetUniswapV2FactoryAddress(ctx)
			require.Error(t, err)
			require.ErrorIs(t, err, types.ErrStateVariableNotFound)
		},
	)

	t.Run("should get uniswapfactory contract address if system contracts are deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		_, factory, _, _, _ := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		found, err := k.GetUniswapV2FactoryAddress(ctx)
		require.NoError(t, err)
		require.Equal(t, factory, found)
	})

	t.Run("should fail in factory not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsConfigurable(t, ctx, k, sdkk.EvmKeeper, &SystemContractDeployConfig{
			DeployWZeta:            true,
			DeployUniswapV2Router:  true,
			DeployUniswapV2Factory: false,
		})

		_, err := k.GetUniswapV2FactoryAddress(ctx)
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})

	t.Run("should fail if evm call fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsWithMockEvmKeeper(t, ctx, k, mockEVMKeeper)

		mockEVMKeeper.MockEVMFailCallOnce()

		_, err := k.GetUniswapV2FactoryAddress(ctx)
		require.ErrorIs(t, err, types.ErrContractCall)
	})

	t.Run("should fail if abi unpack fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsWithMockEvmKeeper(t, ctx, k, mockEVMKeeper)

		mockEVMKeeper.MockEVMSuccessCallOnce()

		_, err := k.GetUniswapV2FactoryAddress(ctx)
		require.ErrorIs(t, err, types.ErrABIUnpack)
	})
}

func TestKeeper_GetUniswapV2Router02Address(t *testing.T) {
	t.Run("should fail to get uniswaprouter contract address if system contracts are not deployed", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		_, err := k.GetUniswapV2Router02Address(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrStateVariableNotFound)
	})

	t.Run("should get uniswaprouter contract address if system contracts are deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		_, _, router, _, _ := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		found, err := k.GetUniswapV2Router02Address(ctx)
		require.NoError(t, err)
		require.Equal(t, router, found)
	})

	t.Run("should fail in router not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsConfigurable(t, ctx, k, sdkk.EvmKeeper, &SystemContractDeployConfig{
			DeployWZeta:            true,
			DeployUniswapV2Router:  false,
			DeployUniswapV2Factory: true,
		})

		_, err := k.GetUniswapV2Router02Address(ctx)
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})

	t.Run("should fail if evm call fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsWithMockEvmKeeper(t, ctx, k, mockEVMKeeper)

		mockEVMKeeper.MockEVMFailCallOnce()

		_, err := k.GetUniswapV2Router02Address(ctx)
		require.ErrorIs(t, err, types.ErrContractCall)
	})

	t.Run("should fail if abi unpack fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsWithMockEvmKeeper(t, ctx, k, mockEVMKeeper)

		mockEVMKeeper.MockEVMSuccessCallOnce()

		_, err := k.GetUniswapV2Router02Address(ctx)
		require.ErrorIs(t, err, types.ErrABIUnpack)
	})
}

func TestKeeper_CallWZetaDeposit(t *testing.T) {
	t.Run("should fail to deposit if system contracts are not deployed", func(t *testing.T) {
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
	})

	t.Run("should deposit if system contracts are deployed", func(t *testing.T) {
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

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		// deposit
		err = k.CallWZetaDeposit(ctx, ethAddr, big.NewInt(42))
		require.NoError(t, err)

		balance, err := k.QueryWZetaBalanceOf(ctx, ethAddr)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(42), balance)
	})
}

func TestKeeper_QueryWZetaBalanceOf(t *testing.T) {
	t.Run("should fail if system contracts are not deployed", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		// fail if no system contract
		_, err := k.QueryWZetaBalanceOf(ctx, sample.EthAddress())
		require.ErrorIs(t, err, types.ErrStateVariableNotFound)
	})
}

func TestKeeper_QuerySystemContractGasCoinZRC20(t *testing.T) {
	t.Run("should fail if system contracts are not deployed", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		chainID := getValidChainID(t)

		_, err := k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chainID))
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrStateVariableNotFound)
	})

	t.Run("should query if system contracts are deployed", func(t *testing.T) {
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
	})

	t.Run("should fail if gas coin not setup", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		chainID := getValidChainID(t)

		_, err := k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chainID))
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrStateVariableNotFound)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		_, err = k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chainID))
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})

	t.Run("should fail if evm call fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsWithMockEvmKeeper(t, ctx, k, mockEVMKeeper)

		mockEVMKeeper.MockEVMFailCallOnce()

		_, err := k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(1))
		require.ErrorIs(t, err, types.ErrContractCall)
	})

	t.Run("should fail if abi unpack fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsWithMockEvmKeeper(t, ctx, k, mockEVMKeeper)

		mockEVMKeeper.MockEVMSuccessCallOnce()

		_, err := k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(1))
		require.ErrorIs(t, err, types.ErrABIUnpack)
	})
}

func TestKeeper_CallUniswapV2RouterSwapExactETHForToken(t *testing.T) {
	t.Run("should fail if system contracts are not deployed", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		// fail if no system contract
		_, err := k.CallUniswapV2RouterSwapExactETHForToken(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			big.NewInt(1),
			sample.EthAddress(),
			true,
		)
		require.ErrorIs(t, err, types.ErrStateVariableNotFound)
	})

	t.Run("should swap if system contracts are deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		chainID := getValidChainID(t)

		// deploy system contracts and swap exact eth for 1 token
		tokenAmount := big.NewInt(1)
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		amountToSwap, err := k.QueryUniswapV2RouterGetZetaAmountsIn(ctx, tokenAmount, zrc20)
		require.NoError(t, err)
		err = sdkk.BankKeeper.MintCoins(
			ctx,
			types.ModuleName,
			sdk.NewCoins(sdk.NewCoin("azeta", sdkmath.NewIntFromBigInt(amountToSwap))),
		)
		require.NoError(t, err)

		amounts, err := k.CallUniswapV2RouterSwapExactETHForToken(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			amountToSwap,
			zrc20,
			true,
		)
		require.NoError(t, err)

		require.Equal(t, 2, len(amounts))
		require.Equal(t, tokenAmount, amounts[1])
	})

	t.Run("should fail if missing zeta balance", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		chainID := getValidChainID(t)

		// deploy system contracts and swap 1 token fails because of missing wrapped balance
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		amountToSwap, err := k.QueryUniswapV2RouterGetZetaAmountsIn(ctx, big.NewInt(1), zrc20)
		require.NoError(t, err)

		_, err = k.CallUniswapV2RouterSwapExactETHForToken(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			amountToSwap,
			zrc20,
			true,
		)
		require.ErrorIs(t, err, types.ErrContractCall)
	})

	t.Run("should fail if wzeta not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsConfigurable(t, ctx, k, sdkk.EvmKeeper, &SystemContractDeployConfig{
			DeployWZeta:            false,
			DeployUniswapV2Factory: true,
			DeployUniswapV2Router:  true,
		})

		_, err := k.CallUniswapV2RouterSwapExactETHForToken(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			big.NewInt(1),
			sample.EthAddress(),
			true,
		)
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})

	t.Run("should fail if router not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsConfigurable(t, ctx, k, sdkk.EvmKeeper, &SystemContractDeployConfig{
			DeployWZeta:            true,
			DeployUniswapV2Factory: true,
			DeployUniswapV2Router:  false,
		})

		_, err := k.CallUniswapV2RouterSwapExactETHForToken(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			big.NewInt(1),
			sample.EthAddress(),
			true,
		)
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})
}

func TestKeeper_CallUniswapV2RouterSwapEthForExactToken(t *testing.T) {
	t.Run("should fail if system contracts are not deployed", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		// fail if no system contract
		_, err := k.CallUniswapV2RouterSwapExactETHForToken(
			ctx, types.ModuleAddressEVM, types.ModuleAddressEVM, big.NewInt(1), sample.EthAddress(), true)
		require.ErrorIs(t, err, types.ErrStateVariableNotFound)
	})

	t.Run("should swap if system contracts are deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		chainID := getValidChainID(t)

		// deploy system contracts and swap exact 1 token
		tokenAmount := big.NewInt(1)
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		amountToSwap, err := k.QueryUniswapV2RouterGetZetaAmountsIn(ctx, tokenAmount, zrc20)
		require.NoError(t, err)
		err = sdkk.BankKeeper.MintCoins(
			ctx,
			types.ModuleName,
			sdk.NewCoins(sdk.NewCoin("azeta", sdkmath.NewIntFromBigInt(amountToSwap))),
		)
		require.NoError(t, err)

		amounts, err := k.CallUniswapV2RouterSwapEthForExactToken(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			amountToSwap,
			tokenAmount,
			zrc20,
		)
		require.NoError(t, err)

		require.Equal(t, 2, len(amounts))
		require.Equal(t, big.NewInt(1), amounts[1])
	})

	t.Run("should fail if missing zeta balance", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		chainID := getValidChainID(t)

		// deploy system contracts and swap 1 token fails because of missing wrapped balance
		tokenAmount := big.NewInt(1)
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		amountToSwap, err := k.QueryUniswapV2RouterGetZetaAmountsIn(ctx, tokenAmount, zrc20)
		require.NoError(t, err)

		_, err = k.CallUniswapV2RouterSwapEthForExactToken(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			amountToSwap,
			tokenAmount,
			zrc20,
		)
		require.ErrorIs(t, err, types.ErrContractCall)
	})

	t.Run("should fail if wzeta not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsConfigurable(t, ctx, k, sdkk.EvmKeeper, &SystemContractDeployConfig{
			DeployWZeta:            false,
			DeployUniswapV2Factory: true,
			DeployUniswapV2Router:  true,
		})

		_, err := k.CallUniswapV2RouterSwapEthForExactToken(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			big.NewInt(1),
			big.NewInt(1),
			sample.EthAddress(),
		)
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})

	t.Run("should fail if router not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsConfigurable(t, ctx, k, sdkk.EvmKeeper, &SystemContractDeployConfig{
			DeployWZeta:            true,
			DeployUniswapV2Factory: true,
			DeployUniswapV2Router:  false,
		})

		_, err := k.CallUniswapV2RouterSwapEthForExactToken(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			big.NewInt(1),
			big.NewInt(1),
			sample.EthAddress(),
		)
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})
}

func TestKeeper_CallUniswapV2RouterSwapExactTokensForETH(t *testing.T) {
	t.Run("should fail if system contracts are not deployed", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		// fail if no system contract
		_, err := k.CallUniswapV2RouterSwapExactTokensForETH(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			big.NewInt(1),
			sample.EthAddress(),
			true,
		)
		require.ErrorIs(t, err, types.ErrStateVariableNotFound)
	})

	t.Run("should swap if system contracts are deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		chainID := getValidChainID(t)

		// fail if no system contract
		_, err := k.CallUniswapV2RouterSwapExactTokensForETH(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			big.NewInt(1),
			sample.EthAddress(),
			true,
		)
		require.Error(t, err)

		// deploy system contracts and swap exact eth for 1 token
		ethAmount := big.NewInt(1)
		_, _, router, _, _ := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		amountToSwap, err := k.QueryUniswapV2RouterGetZRC4AmountsIn(ctx, ethAmount, zrc20)
		require.NoError(t, err)

		_, err = k.DepositZRC20(ctx, zrc20, types.ModuleAddressEVM, amountToSwap)
		require.NoError(t, err)
		k.CallZRC20Approve(
			ctx,
			types.ModuleAddressEVM,
			zrc20,
			router,
			amountToSwap,
			false,
		)

		amounts, err := k.CallUniswapV2RouterSwapExactTokensForETH(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			amountToSwap,
			zrc20,
			true,
		)
		require.NoError(t, err)

		require.Equal(t, 2, len(amounts))
		require.Equal(t, ethAmount, amounts[0])
	})

	t.Run("should fail if missing tokens balance", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		chainID := getValidChainID(t)

		// fail if no system contract
		_, err := k.CallUniswapV2RouterSwapExactTokensForETH(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			big.NewInt(1),
			sample.EthAddress(),
			true,
		)
		require.Error(t, err)

		// deploy system contracts and swap fails because of missing balance
		ethAmount := big.NewInt(1)
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		amountToSwap, err := k.QueryUniswapV2RouterGetZRC4AmountsIn(ctx, ethAmount, zrc20)
		require.NoError(t, err)

		_, err = k.CallUniswapV2RouterSwapExactTokensForETH(
			ctx, types.ModuleAddressEVM, types.ModuleAddressEVM, amountToSwap, zrc20, true)
		require.ErrorIs(t, err, types.ErrContractCall)
	})

	t.Run("should fail if wzeta not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsConfigurable(t, ctx, k, sdkk.EvmKeeper, &SystemContractDeployConfig{
			DeployWZeta:            false,
			DeployUniswapV2Factory: true,
			DeployUniswapV2Router:  true,
		})
		_, err := k.CallUniswapV2RouterSwapExactTokensForETH(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			big.NewInt(1),
			sample.EthAddress(),
			true,
		)
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})

	t.Run("should fail if router not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsConfigurable(t, ctx, k, sdkk.EvmKeeper, &SystemContractDeployConfig{
			DeployWZeta:            true,
			DeployUniswapV2Factory: true,
			DeployUniswapV2Router:  false,
		})
		_, err := k.CallUniswapV2RouterSwapExactTokensForETH(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			big.NewInt(1),
			sample.EthAddress(),
			true,
		)
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})
}

func TestKeeper_CallUniswapV2RouterSwapExactTokensForTokens(t *testing.T) {
	t.Run("should fail if system contracts are not deployed", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		// fail if no system contract
		_, err := k.CallUniswapV2RouterSwapExactTokensForTokens(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			big.NewInt(1),
			sample.EthAddress(),
			sample.EthAddress(),
			true,
		)
		require.ErrorIs(t, err, types.ErrStateVariableNotFound)
	})

	t.Run("should swap if system contracts are deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		chainID := getValidChainID(t)

		// fail if no system contract
		_, err := k.CallUniswapV2RouterSwapExactTokensForTokens(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			big.NewInt(1),
			sample.EthAddress(),
			sample.EthAddress(),
			true,
		)
		require.Error(t, err)

		// deploy system contracts and swap exact token for 1 token
		tokenAmount := big.NewInt(1)
		_, _, router, _, _ := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		inzrc20 := deployZRC20(t, ctx, k, sdkk.EvmKeeper, chainID, "foo", sample.EthAddress().String(), "foo")
		outzrc20 := deployZRC20(t, ctx, k, sdkk.EvmKeeper, chainID, "bar", sample.EthAddress().String(), "bar")
		setupZRC20Pool(t, ctx, k, sdkk.BankKeeper, inzrc20)
		setupZRC20Pool(t, ctx, k, sdkk.BankKeeper, outzrc20)

		amountToSwap, err := k.QueryUniswapV2RouterGetZRC4ToZRC4AmountsIn(ctx, tokenAmount, inzrc20, outzrc20)
		require.NoError(t, err)

		_, err = k.DepositZRC20(ctx, inzrc20, types.ModuleAddressEVM, amountToSwap)
		require.NoError(t, err)
		k.CallZRC20Approve(
			ctx,
			types.ModuleAddressEVM,
			inzrc20,
			router,
			amountToSwap,
			false,
		)

		amounts, err := k.CallUniswapV2RouterSwapExactTokensForTokens(
			ctx, types.ModuleAddressEVM, types.ModuleAddressEVM, amountToSwap, inzrc20, outzrc20, true)
		require.NoError(t, err)
		require.Equal(t, 3, len(amounts))
		require.Equal(t, amounts[2], tokenAmount)
	})

	t.Run("should fail if missing tokens balance", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		chainID := getValidChainID(t)

		// deploy system contracts and swap fails because of missing balance
		tokenAmount := big.NewInt(1)
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		inzrc20 := deployZRC20(t, ctx, k, sdkk.EvmKeeper, chainID, "foo", sample.EthAddress().String(), "foo")
		outzrc20 := deployZRC20(t, ctx, k, sdkk.EvmKeeper, chainID, "bar", sample.EthAddress().String(), "bar")
		setupZRC20Pool(t, ctx, k, sdkk.BankKeeper, inzrc20)
		setupZRC20Pool(t, ctx, k, sdkk.BankKeeper, outzrc20)

		amountToSwap, err := k.QueryUniswapV2RouterGetZRC4ToZRC4AmountsIn(ctx, tokenAmount, inzrc20, outzrc20)
		require.NoError(t, err)

		_, err = k.CallUniswapV2RouterSwapExactTokensForTokens(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			amountToSwap,
			inzrc20,
			outzrc20,
			true,
		)
		require.ErrorIs(t, err, types.ErrContractCall)
	})

	t.Run("should fail if wzeta not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		// fail if no system contract
		_, err := k.CallUniswapV2RouterSwapExactTokensForTokens(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			big.NewInt(1),
			sample.EthAddress(),
			sample.EthAddress(),
			true,
		)
		require.Error(t, err)

		// deploy system contracts except router
		deploySystemContractsConfigurable(t, ctx, k, sdkk.EvmKeeper, &SystemContractDeployConfig{
			DeployUniswapV2Router:  true,
			DeployWZeta:            false,
			DeployUniswapV2Factory: true,
		})

		_, err = k.CallUniswapV2RouterSwapExactTokensForTokens(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			big.NewInt(1),
			sample.EthAddress(),
			sample.EthAddress(),
			true,
		)
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})

	t.Run("should fail if router not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		// fail if no system contract
		_, err := k.CallUniswapV2RouterSwapExactTokensForTokens(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			big.NewInt(1),
			sample.EthAddress(),
			sample.EthAddress(),
			true,
		)
		require.Error(t, err)

		// deploy system contracts except router
		deploySystemContractsConfigurable(t, ctx, k, sdkk.EvmKeeper, &SystemContractDeployConfig{
			DeployUniswapV2Router:  false,
			DeployWZeta:            true,
			DeployUniswapV2Factory: true,
		})

		_, err = k.CallUniswapV2RouterSwapExactTokensForTokens(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			big.NewInt(1),
			sample.EthAddress(),
			sample.EthAddress(),
			true,
		)
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})

}

func TestKeeper_QueryUniswapV2RouterGetZRC4AmountsIn(t *testing.T) {
	t.Run("should fail if no amounts out", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		_, err := k.QueryUniswapV2RouterGetZRC4AmountsIn(ctx, big.NewInt(1), sample.EthAddress())
		require.ErrorIs(t, err, types.ErrContractCall)
	})

	t.Run("should fail if wzeta not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsConfigurable(t, ctx, k, sdkk.EvmKeeper, &SystemContractDeployConfig{
			DeployWZeta:            false,
			DeployUniswapV2Factory: true,
			DeployUniswapV2Router:  true,
		})

		_, err := k.QueryUniswapV2RouterGetZRC4AmountsIn(ctx, big.NewInt(1), sample.EthAddress())
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})

	t.Run("should fail if router not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsConfigurable(t, ctx, k, sdkk.EvmKeeper, &SystemContractDeployConfig{
			DeployWZeta:            true,
			DeployUniswapV2Factory: true,
			DeployUniswapV2Router:  false,
		})

		_, err := k.QueryUniswapV2RouterGetZRC4AmountsIn(ctx, big.NewInt(1), sample.EthAddress())
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})
}

func TestKeeper_QueryUniswapV2RouterGetZetaAmountsIn(t *testing.T) {
	t.Run("should fail if no amounts out", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		_, err := k.QueryUniswapV2RouterGetZetaAmountsIn(ctx, big.NewInt(1), sample.EthAddress())
		require.ErrorIs(t, err, types.ErrContractCall)
	})

	t.Run("should fail if wzeta not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsConfigurable(t, ctx, k, sdkk.EvmKeeper, &SystemContractDeployConfig{
			DeployWZeta:            false,
			DeployUniswapV2Factory: true,
			DeployUniswapV2Router:  true,
		})

		_, err := k.QueryUniswapV2RouterGetZetaAmountsIn(ctx, big.NewInt(1), sample.EthAddress())
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})

	t.Run("should fail if router not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsConfigurable(t, ctx, k, sdkk.EvmKeeper, &SystemContractDeployConfig{
			DeployWZeta:            true,
			DeployUniswapV2Factory: true,
			DeployUniswapV2Router:  false,
		})

		_, err := k.QueryUniswapV2RouterGetZetaAmountsIn(ctx, big.NewInt(1), sample.EthAddress())
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})
}

func TestKeeper_QueryUniswapV2RouterGetZRC4ToZRC4AmountsIn(t *testing.T) {
	t.Run("should fail if no amounts out", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		_, err := k.QueryUniswapV2RouterGetZRC4ToZRC4AmountsIn(
			ctx,
			big.NewInt(1),
			sample.EthAddress(),
			sample.EthAddress(),
		)
		require.ErrorIs(t, err, types.ErrContractCall)
	})

	t.Run("should fail if wzeta not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsConfigurable(t, ctx, k, sdkk.EvmKeeper, &SystemContractDeployConfig{
			DeployWZeta:            false,
			DeployUniswapV2Factory: true,
			DeployUniswapV2Router:  true,
		})

		_, err := k.QueryUniswapV2RouterGetZRC4ToZRC4AmountsIn(
			ctx,
			big.NewInt(1),
			sample.EthAddress(),
			sample.EthAddress(),
		)
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})

	t.Run("should fail if router not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsConfigurable(t, ctx, k, sdkk.EvmKeeper, &SystemContractDeployConfig{
			DeployWZeta:            true,
			DeployUniswapV2Factory: true,
			DeployUniswapV2Router:  false,
		})

		_, err := k.QueryUniswapV2RouterGetZRC4ToZRC4AmountsIn(
			ctx,
			big.NewInt(1),
			sample.EthAddress(),
			sample.EthAddress(),
		)
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})
}

func TestKeeper_CallZRC20Burn(t *testing.T) {
	t.Run("should fail if evm call fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsWithMockEvmKeeper(t, ctx, k, mockEVMKeeper)

		mockEVMKeeper.MockEVMFailCallOnce()
		err := k.CallZRC20Burn(ctx, types.ModuleAddressEVM, sample.EthAddress(), big.NewInt(1), false)
		require.ErrorIs(t, err, types.ErrContractCall)
	})
}

func TestKeeper_CallZRC20Approve(t *testing.T) {
	t.Run("should fail if evm call fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsWithMockEvmKeeper(t, ctx, k, mockEVMKeeper)

		mockEVMKeeper.MockEVMFailCallOnce()
		err := k.CallZRC20Approve(
			ctx,
			types.ModuleAddressEVM,
			sample.EthAddress(),
			types.ModuleAddressEVM,
			big.NewInt(1),
			false,
		)
		require.ErrorIs(t, err, types.ErrContractCall)
	})
}

func TestKeeper_CallZRC20Deposit(t *testing.T) {
	t.Run("should fail if evm call fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsWithMockEvmKeeper(t, ctx, k, mockEVMKeeper)

		mockEVMKeeper.MockEVMFailCallOnce()
		err := k.CallZRC20Deposit(
			ctx,
			types.ModuleAddressEVM,
			sample.EthAddress(),
			types.ModuleAddressEVM,
			big.NewInt(1),
			false,
		)
		require.ErrorIs(t, err, types.ErrContractCall)
	})
}

func TestKeeper_GetGatewayGasLimitSafe(t *testing.T) {
	t.Run("get default if not set", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)

		limit := k.GetGatewayGasLimitSafe(ctx)
		require.Equal(t, types.DefaultGatewayGasLimit, limit.Uint64())
	})

	t.Run("get value from state if set", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		k.SetSystemContract(ctx, *sample.SystemContract())
		newGasLimit := uint64(10)

		k.SetGatewayGasLimit(ctx, newGasLimit)

		limit := k.GetGatewayGasLimitSafe(ctx)
		require.Equal(t, newGasLimit, limit.Uint64())
	})
}

func TestKeeper_SetGatewayGasLimit(t *testing.T) {
	t.Run("set gas limit if system contract is found", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		defaultSystemContract := *sample.SystemContract()
		k.SetSystemContract(ctx, defaultSystemContract)
		newGasLimit := uint64(10)

		// Act
		k.SetGatewayGasLimit(ctx, newGasLimit)

		// Assert
		updatedSystemContract, found := k.GetSystemContract(ctx)
		require.True(t, found)
		require.Equal(t, newGasLimit, updatedSystemContract.GatewayGasLimit)
		require.Equal(t, defaultSystemContract.SystemContract, updatedSystemContract.SystemContract)
		require.Equal(t, defaultSystemContract.Gateway, updatedSystemContract.Gateway)
		require.Equal(t, defaultSystemContract.ConnectorZevm, updatedSystemContract.ConnectorZevm)

	})

	t.Run("set gas limit if system contract is not found", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		newGasLimit := uint64(10)

		// Act
		k.SetGatewayGasLimit(ctx, newGasLimit)

		// Assert
		updatedSystemContract, found := k.GetSystemContract(ctx)
		require.True(t, found)
		require.Equal(t, newGasLimit, updatedSystemContract.GatewayGasLimit)
		require.Equal(t, "", updatedSystemContract.SystemContract)
		require.Equal(t, "", updatedSystemContract.Gateway)
		require.Equal(t, "", updatedSystemContract.ConnectorZevm)
	})
}

func TestGetGatewayGasLimit(t *testing.T) {
	t.Run("unable to get gas limit if system contract is not set", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.FungibleKeeper(t)

		// Act
		_, err := k.GetGatewayGasLimit(ctx)
		require.ErrorIs(t, err, types.ErrSystemContractNotFound)
	})

	t.Run("unable to get gas limit if gas limit is not set", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.FungibleKeeper(t)

		k.SetSystemContract(ctx, types.SystemContract{})

		// Act
		_, err := k.GetGatewayGasLimit(ctx)

		// Assert
		require.ErrorIs(t, err, types.ErrGasLimitNotSet)
	})

	t.Run("get gas limit if system contract is set", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		defaultSystemContract := *sample.SystemContract()
		k.SetSystemContract(ctx, defaultSystemContract)

		// Act
		gasLimit, err := k.GetGatewayGasLimit(ctx)

		// Assert
		require.NoError(t, err)
		require.Equal(t, defaultSystemContract.GatewayGasLimit, gasLimit)
	})
}

package keeper_test

import (
	"errors"
	"math/big"
	"testing"

	"github.com/zeta-chain/node/pkg/contracts/uniswap/v2-periphery/contracts/uniswapv2router02.sol"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	evmkeeper "github.com/cosmos/evm/x/vm/keeper"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/systemcontract.sol"

	"github.com/zeta-chain/node/cmd/zetacored/config"
	"github.com/zeta-chain/node/pkg/ptr"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	fungiblekeeper "github.com/zeta-chain/node/x/fungible/keeper"
	"github.com/zeta-chain/node/x/fungible/types"
)

// setupGasCoin is a helper function to setup the gas coin for testing
func setupGasCoin(
	t *testing.T,
	ctx sdk.Context,
	k *fungiblekeeper.Keeper,
	evmk types.EVMKeeper,
	chainID int64,
	assetName string,
	symbol string,
) (zrc20 common.Address) {
	addr, err := k.SetupChainGasCoinAndPool(
		ctx,
		chainID,
		assetName,
		symbol,
		8,
		nil,
		ptr.Ptr(sdkmath.NewUint(1000)),
	)
	require.NoError(t, err)
	assertContractDeployment(t, evmk, ctx, addr)

	// increase the default liquidity cap
	foreignCoin, found := k.GetForeignCoins(ctx, addr.Hex())
	require.True(t, found)

	foreignCoin.LiquidityCap = sdkmath.NewUint(1e18).MulUint64(1e12)
	k.SetForeignCoins(ctx, foreignCoin)

	return addr
}

func deployZRC20(
	t *testing.T,
	ctx sdk.Context,
	k *fungiblekeeper.Keeper,
	evmk *evmkeeper.Keeper,
	chainID int64,
	assetName string,
	assetAddress string,
	symbol string,
) (zrc20 common.Address) {
	addr, err := k.DeployZRC20Contract(
		ctx,
		assetName,
		symbol,
		8,
		chainID,
		0,
		assetAddress,
		big.NewInt(21_000),
		ptr.Ptr(sdkmath.NewUint(1000)),
	)
	require.NoError(t, err)
	assertContractDeployment(t, evmk, ctx, addr)

	// increase the default liquidity cap
	foreignCoin, found := k.GetForeignCoins(ctx, addr.Hex())
	require.True(t, found)
	foreignCoin.LiquidityCap = sdkmath.NewUint(1e18).MulUint64(1e12)
	k.SetForeignCoins(ctx, foreignCoin)

	return addr
}

// setupZRC20Pool setup a Uniswap pool with liquidity for the pair zeta/asset
func setupZRC20Pool(
	t *testing.T,
	ctx sdk.Context,
	k *fungiblekeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
	zrc20Addr common.Address,
) {
	routerAddress, err := k.GetUniswapV2Router02Address(ctx)
	require.NoError(t, err)
	routerABI, err := uniswapv2router02.UniswapV2Router02MetaData.GetAbi()
	require.NoError(t, err)

	// enough for the small numbers used in test
	liquidityAmount := big.NewInt(1e17)

	// mint some zrc20 and zeta
	_, err = k.DepositZRC20(ctx, zrc20Addr, types.ModuleAddressEVM, liquidityAmount)
	require.NoError(t, err)
	err = bankKeeper.MintCoins(
		ctx,
		types.ModuleName,
		sdk.NewCoins(sdk.NewCoin(config.BaseDenom, sdkmath.NewIntFromBigInt(liquidityAmount))),
	)
	require.NoError(t, err)

	// approve the router to spend the zrc20
	err = k.CallZRC20Approve(
		ctx,
		types.ModuleAddressEVM,
		zrc20Addr,
		routerAddress,
		liquidityAmount,
		false,
	)
	require.NoError(t, err)

	// k2 := liquidityAmount.Sub(liquidityAmount, big.NewInt(1000))
	// add the liquidity
	//function addLiquidityETH(
	//	address token,
	//	uint amountTokenDesired,
	//	uint amountTokenMin,
	//	uint amountETHMin,
	//	address to,
	//	uint deadline
	//)
	_, err = k.CallEVM(
		ctx,
		*routerABI,
		types.ModuleAddressEVM,
		routerAddress,
		liquidityAmount,
		big.NewInt(5_000_000),
		true,
		false,
		"addLiquidityETH",
		zrc20Addr,
		liquidityAmount,
		fungiblekeeper.BigIntZero,
		fungiblekeeper.BigIntZero,
		types.ModuleAddressEVM,
		liquidityAmount,
	)
	require.NoError(t, err)
}

func TestKeeper_SetupChainGasCoinAndPool(t *testing.T) {
	t.Run("can setup a new chain gas coin", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := getValidChainID(t)

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		// can retrieve the gas coin
		found, err := k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chainID))
		require.NoError(t, err)
		require.Equal(t, zrc20, found)
	})

	t.Run("should error if system contracts not deployed", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := getValidChainID(t)

		addr, err := k.SetupChainGasCoinAndPool(
			ctx,
			chainID,
			"test",
			"test",
			8,
			nil,
			ptr.Ptr(sdkmath.NewUint(1000)),
		)
		require.Error(t, err)
		require.Empty(t, addr)
	})

	t.Run("should error if mint coins fails", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseBankMock: true,
		})
		bankMock := keepertest.GetFungibleBankMock(t, k)
		bankMock.On("MintCoins", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("err"))
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := getValidChainID(t)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		addr, err := k.SetupChainGasCoinAndPool(
			ctx,
			chainID,
			"test",
			"test",
			8,
			nil,
			ptr.Ptr(sdkmath.NewUint(1000)),
		)
		require.Error(t, err)
		require.Empty(t, addr)
	})

	t.Run("should error if set gas coin fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := getValidChainID(t)
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		deploySystemContractsWithMockEvmKeeper(t, ctx, k, mockEVMKeeper)

		// deployZrc20 success
		mockSuccessfulContractDeployment(ctx, t, k)

		// setGasCoin fail
		mockEVMKeeper.MockEVMFailCallOnce()

		addr, err := k.SetupChainGasCoinAndPool(
			ctx,
			chainID,
			"test",
			"test",
			8,
			nil,
			ptr.Ptr(sdkmath.NewUint(1000)),
		)
		require.Error(t, err)
		require.Empty(t, addr)
	})

	t.Run("should error if deposit zrc20 fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := getValidChainID(t)
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		deploySystemContractsWithMockEvmKeeper(t, ctx, k, mockEVMKeeper)

		// deployZrc20 success
		mockSuccessfulContractDeployment(ctx, t, k)

		// setGasCoin success
		mockEVMKeeper.MockEVMSuccessCallOnce()

		// depositZrc20 fails
		mockEVMKeeper.MockEVMFailCallOnce()

		addr, err := k.SetupChainGasCoinAndPool(
			ctx,
			chainID,
			"test",
			"test",
			8,
			nil,
			ptr.Ptr(sdkmath.NewUint(1000)),
		)
		require.Error(t, err)
		require.Empty(t, addr)
	})

	t.Run("should error if set gas pool call fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := getValidChainID(t)
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		deploySystemContractsWithMockEvmKeeper(t, ctx, k, mockEVMKeeper)

		// deployZrc20 success
		mockSuccessfulContractDeployment(ctx, t, k)

		// setGasCoin success
		// depositZrc20 success
		mockEVMKeeper.MockEVMSuccessCallTimes(2)

		// set gas pool call fail
		mockEVMKeeper.MockEVMFailCallOnce()

		addr, err := k.SetupChainGasCoinAndPool(
			ctx,
			chainID,
			"test",
			"test",
			8,
			nil,
			ptr.Ptr(sdkmath.NewUint(1000)),
		)
		require.Error(t, err)
		require.Empty(t, addr)
	})

	t.Run("should error if get uniswap router fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := getValidChainID(t)
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		deploySystemContractsWithMockEvmKeeper(t, ctx, k, mockEVMKeeper)

		// deployZrc20 success
		mockSuccessfulContractDeployment(ctx, t, k)

		// setGasCoin success
		// depositZrc20 success
		// set gas pool call success
		mockEVMKeeper.MockEVMSuccessCallTimes(3)

		// get uniswap router fails
		mockEVMKeeper.MockEVMFailCallOnce()

		addr, err := k.SetupChainGasCoinAndPool(
			ctx,
			chainID,
			"test",
			"test",
			8,
			nil,
			ptr.Ptr(sdkmath.NewUint(1000)),
		)
		require.Error(t, err)
		require.Empty(t, addr)
	})

	t.Run("should error if approve fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := getValidChainID(t)
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		deploySystemContractsWithMockEvmKeeper(t, ctx, k, mockEVMKeeper)

		// deployZrc20 success
		mockSuccessfulContractDeployment(ctx, t, k)

		// setGasCoin success
		// depositZrc20 success
		// set gas pool call success
		mockEVMKeeper.MockEVMSuccessCallTimes(3)

		// get uniswap router success
		sysABI, err := systemcontract.SystemContractMetaData.GetAbi()
		require.NoError(t, err)
		routerAddr, err := sysABI.Methods["uniswapv2Router02Address"].Outputs.Pack(sample.EthAddress())
		require.NoError(t, err)
		mockEVMKeeper.MockEVMSuccessCallOnceWithReturn(&evmtypes.MsgEthereumTxResponse{Ret: routerAddr})

		// get approve fails
		mockEVMKeeper.MockEVMFailCallOnce()

		addr, err := k.SetupChainGasCoinAndPool(
			ctx,
			chainID,
			"test",
			"test",
			8,
			nil,
			ptr.Ptr(sdkmath.NewUint(1000)),
		)
		require.Error(t, err)
		require.Empty(t, addr)
	})

	t.Run("should error if add liquidity fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := getValidChainID(t)
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		deploySystemContractsWithMockEvmKeeper(t, ctx, k, mockEVMKeeper)

		// deployZrc20 success
		mockSuccessfulContractDeployment(ctx, t, k)

		// setGasCoin success
		// depositZrc20 success
		// set gas pool call success
		mockEVMKeeper.MockEVMSuccessCallTimes(3)

		// get uniswap router success
		sysABI, err := systemcontract.SystemContractMetaData.GetAbi()
		require.NoError(t, err)
		routerAddr, err := sysABI.Methods["uniswapv2Router02Address"].Outputs.Pack(sample.EthAddress())
		require.NoError(t, err)
		mockEVMKeeper.MockEVMSuccessCallOnceWithReturn(&evmtypes.MsgEthereumTxResponse{Ret: routerAddr})

		// get approve success
		mockEVMKeeper.MockEVMSuccessCallOnce()

		// add liquidity fails
		mockEVMKeeper.MockEVMFailCallOnce()

		addr, err := k.SetupChainGasCoinAndPool(
			ctx,
			chainID,
			"test",
			"test",
			8,
			nil,
			ptr.Ptr(sdkmath.NewUint(1000)),
		)
		require.Error(t, err)
		require.Empty(t, addr)
	})

	t.Run("should error if add liquidity fails to unpack", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := getValidChainID(t)
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		deploySystemContractsWithMockEvmKeeper(t, ctx, k, mockEVMKeeper)

		// deployZrc20 success
		mockSuccessfulContractDeployment(ctx, t, k)

		// setGasCoin success
		// depositZrc20 success
		// set gas pool call success
		mockEVMKeeper.MockEVMSuccessCallTimes(3)

		// get uniswap router success
		sysABI, err := systemcontract.SystemContractMetaData.GetAbi()
		require.NoError(t, err)
		routerAddr, err := sysABI.Methods["uniswapv2Router02Address"].Outputs.Pack(sample.EthAddress())
		require.NoError(t, err)
		mockEVMKeeper.MockEVMSuccessCallOnceWithReturn(&evmtypes.MsgEthereumTxResponse{Ret: routerAddr})

		// get approve success
		mockEVMKeeper.MockEVMSuccessCallOnce()

		// add liquidity success
		mockEVMKeeper.MockEVMSuccessCallOnce()

		addr, err := k.SetupChainGasCoinAndPool(
			ctx,
			chainID,
			"test",
			"test",
			8,
			nil,
			ptr.Ptr(sdkmath.NewUint(1000)),
		)
		require.Error(t, err)
		require.Empty(t, addr)
	})
}

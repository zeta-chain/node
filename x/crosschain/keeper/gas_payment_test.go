package keeper_test

import (
	"math/big"
	"testing"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/ethereum/go-ethereum/common"
	evmkeeper "github.com/evmos/ethermint/x/evm/keeper"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/uniswap/v2-periphery/contracts/uniswapv2router02.sol"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	zetacommon "github.com/zeta-chain/zetacore/common"
	testkeeper "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungiblekeeper "github.com/zeta-chain/zetacore/x/fungible/keeper"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// get a valid chain id independently of the build flag
func getValidEthChainID(t *testing.T) int64 {
	list := zetacommon.DefaultChainsList()
	require.True(t, len(list) > 1)
	require.NotNil(t, list[1])
	require.False(t, zetacommon.IsBitcoinChain(list[1].ChainId))

	return list[1].ChainId
}

// assert that a contract has been deployed by checking stored code is non-empty.
func assertContractDeployment(t *testing.T, k *evmkeeper.Keeper, ctx sdk.Context, contractAddress common.Address) {
	acc := k.GetAccount(ctx, contractAddress)
	require.NotNil(t, acc)

	code := k.GetCode(ctx, common.BytesToHash(acc.CodeHash))
	require.NotEmpty(t, code)
}

// deploySystemContracts deploys the system contracts and returns their addresses.
func deploySystemContracts(
	t *testing.T,
	ctx sdk.Context,
	k *fungiblekeeper.Keeper,
	evmk *evmkeeper.Keeper,
) (wzeta, uniswapV2Factory, uniswapV2Router, connector, systemContract common.Address) {
	var err error

	wzeta, err = k.DeployWZETA(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, wzeta)
	assertContractDeployment(t, evmk, ctx, wzeta)

	uniswapV2Factory, err = k.DeployUniswapV2Factory(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, uniswapV2Factory)
	assertContractDeployment(t, evmk, ctx, uniswapV2Factory)

	uniswapV2Router, err = k.DeployUniswapV2Router02(ctx, uniswapV2Factory, wzeta)
	require.NoError(t, err)
	require.NotEmpty(t, uniswapV2Router)
	assertContractDeployment(t, evmk, ctx, uniswapV2Router)

	connector, err = k.DeployConnectorZEVM(ctx, wzeta)
	require.NoError(t, err)
	require.NotEmpty(t, connector)
	assertContractDeployment(t, evmk, ctx, connector)

	systemContract, err = k.DeploySystemContract(ctx, wzeta, uniswapV2Factory, uniswapV2Router)
	require.NoError(t, err)
	require.NotEmpty(t, systemContract)
	assertContractDeployment(t, evmk, ctx, systemContract)

	return
}

// setupGasCoin is a helper function to setup the gas coin for testing
func setupGasCoin(
	t *testing.T,
	ctx sdk.Context,
	k *fungiblekeeper.Keeper,
	evmk *evmkeeper.Keeper,
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
	)
	require.NoError(t, err)
	assertContractDeployment(t, evmk, ctx, addr)
	return addr
}

// deployZRC20 deploys a ZRC20 contract and returns its address
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
	)
	require.NoError(t, err)
	assertContractDeployment(t, evmk, ctx, addr)
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
		sdk.NewCoins(sdk.NewCoin(config.BaseDenom, sdk.NewIntFromBigInt(liquidityAmount))),
	)
	require.NoError(t, err)

	// approve the router to spend the zeta
	err = k.CallZRC20Approve(
		ctx,
		types.ModuleAddressEVM,
		zrc20Addr,
		routerAddress,
		liquidityAmount,
		false,
	)
	require.NoError(t, err)

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

func setAdminPolicies(ctx sdk.Context, zk testkeeper.ZetaKeepers, admin string) {
	params := zk.ObserverKeeper.GetParams(ctx)
	params.AdminPolicy = []*observertypes.Admin_Policy{
		{
			PolicyType: observertypes.Policy_Type_group1,
			Address:    admin,
		},
		{
			PolicyType: observertypes.Policy_Type_group2,
			Address:    admin,
		},
	}
	zk.ObserverKeeper.SetParams(ctx, params)
}

var (
	// gasLimit = big.NewInt(21_000) - value used in SetupChainGasCoinAndPool for gas limit initialization
	withdrawFee uint64 = 1000
	gasPrice    uint64 = 2
	inputAmount uint64 = 100000
)

func TestKeeper_PayGasNativeAndUpdateCctx(t *testing.T) {
	t.Run("can pay gas in native gas", func(t *testing.T) {
		k, ctx, sdkk, zk := testkeeper.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)

		// deploy gas coin and set fee params
		chainID := getValidEthChainID(t)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "foobar", "foobar")
		_, err := zk.FungibleKeeper.UpdateZRC20WithdrawFee(
			sdk.UnwrapSDKContext(ctx),
			fungibletypes.NewMsgUpdateZRC20WithdrawFee(admin, zrc20.String(), sdk.NewUint(withdrawFee)),
		)
		require.NoError(t, err)
		k.SetGasPrice(ctx, types.GasPrice{
			ChainId:     chainID,
			MedianIndex: 0,
			Prices:      []uint64{gasPrice},
		})

		// create a cctx reverted from zeta
		cctx := types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType: zetacommon.CoinType_Gas,
			},
			OutboundTxParams: []*types.OutboundTxParams{
				{
					ReceiverChainId: zetacommon.ZetaChain().ChainId,
				},
				{
					ReceiverChainId: chainID,
				},
			},
		}

		// total fees must be 21000*2+1000=43000
		// if the input amount of the cctx is 100000, the output amount must be 100000-43000=57000
		err = k.PayGasNativeAndUpdateCctx(ctx, chainID, &cctx, math.NewUint(inputAmount))
		require.NoError(t, err)
		require.Equal(t, uint64(57000), cctx.GetCurrentOutTxParam().Amount.Uint64())
		require.Equal(t, uint64(21_000), cctx.GetCurrentOutTxParam().OutboundTxGasLimit)
		require.Equal(t, "2", cctx.GetCurrentOutTxParam().OutboundTxGasPrice)
	})

	t.Run("should fail if not coin type gas", func(t *testing.T) {
		k, ctx, _, _ := testkeeper.CrosschainKeeper(t)
		chainID := getValidEthChainID(t)
		cctx := types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType: zetacommon.CoinType_Zeta,
			},
		}
		err := k.PayGasNativeAndUpdateCctx(ctx, chainID, &cctx, math.NewUint(inputAmount))
		require.ErrorIs(t, err, types.ErrInvalidCoinType)
	})

	t.Run("should fail if chain is not supported", func(t *testing.T) {
		k, ctx, _, _ := testkeeper.CrosschainKeeper(t)
		cctx := types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType: zetacommon.CoinType_Gas,
			},
		}
		err := k.PayGasNativeAndUpdateCctx(ctx, 999999, &cctx, math.NewUint(inputAmount))
		require.ErrorIs(t, err, observertypes.ErrSupportedChains)
	})

	t.Run("should fail if can't query the gas price", func(t *testing.T) {
		k, ctx, sdkk, zk := testkeeper.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)

		// deploy gas coin and set fee params
		chainID := getValidEthChainID(t)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		// create a cctx reverted from zeta
		cctx := types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType: zetacommon.CoinType_Gas,
			},
			OutboundTxParams: []*types.OutboundTxParams{
				{
					ReceiverChainId: zetacommon.ZetaChain().ChainId,
				},
				{
					ReceiverChainId: chainID,
				},
			},
		}

		err := k.PayGasNativeAndUpdateCctx(ctx, chainID, &cctx, math.NewUint(inputAmount))
		require.ErrorIs(t, err, types.ErrCannotFindGasParams)
	})

	t.Run("should fail if not enough amount for the fee", func(t *testing.T) {
		k, ctx, sdkk, zk := testkeeper.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)

		// deploy gas coin and set fee params
		chainID := getValidEthChainID(t)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "foobar", "foobar")
		_, err := zk.FungibleKeeper.UpdateZRC20WithdrawFee(
			sdk.UnwrapSDKContext(ctx),
			fungibletypes.NewMsgUpdateZRC20WithdrawFee(admin, zrc20.String(), sdk.NewUint(withdrawFee)),
		)
		require.NoError(t, err)
		k.SetGasPrice(ctx, types.GasPrice{
			ChainId:     chainID,
			MedianIndex: 0,
			Prices:      []uint64{gasPrice},
		})

		cctx := types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType: zetacommon.CoinType_Gas,
			},
			OutboundTxParams: []*types.OutboundTxParams{
				{
					ReceiverChainId: zetacommon.ZetaChain().ChainId,
				},
				{
					ReceiverChainId: chainID,
				},
			},
		}

		// 42999 < 43000
		err = k.PayGasNativeAndUpdateCctx(ctx, chainID, &cctx, math.NewUint(42999))
		require.ErrorIs(t, err, types.ErrNotEnoughGas)
	})
}

func TestKeeper_PayGasInERC20AndUpdateCctx(t *testing.T) {
	t.Run("can pay gas in erc20", func(t *testing.T) {
		k, ctx, sdkk, zk := testkeeper.CrosschainKeeper(t)

		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)

		// deploy gas coin, erc20 and set fee params
		chainID := getValidEthChainID(t)
		assetAddress := sample.EthAddress().String()
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		gasZRC20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "foo", "foo")
		zrc20Addr := deployZRC20(
			t,
			ctx,
			zk.FungibleKeeper,
			sdkk.EvmKeeper,
			chainID,
			"bar",
			assetAddress,
			"bar",
		)
		_, err := zk.FungibleKeeper.UpdateZRC20WithdrawFee(
			sdk.UnwrapSDKContext(ctx),
			fungibletypes.NewMsgUpdateZRC20WithdrawFee(admin, gasZRC20.String(), sdk.NewUint(withdrawFee)),
		)
		require.NoError(t, err)
		k.SetGasPrice(ctx, types.GasPrice{
			ChainId:     chainID,
			MedianIndex: 0,
			Prices:      []uint64{gasPrice},
		})

		setupZRC20Pool(
			t,
			ctx,
			zk.FungibleKeeper,
			sdkk.BankKeeper,
			zrc20Addr,
		)

		// create a cctx reverted from zeta
		cctx := types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType: zetacommon.CoinType_ERC20,
				Asset:    assetAddress,
			},
			OutboundTxParams: []*types.OutboundTxParams{
				{
					ReceiverChainId: zetacommon.ZetaChain().ChainId,
				},
				{
					ReceiverChainId: chainID,
				},
			},
		}

		// total fees in gas must be 21000*2+1000=43000
		// we calculate what it represents in erc20
		expectedInZeta, err := zk.FungibleKeeper.QueryUniswapV2RouterGetZetaAmountsIn(ctx, big.NewInt(43000), gasZRC20)
		require.NoError(t, err)
		expectedInZRC20, err := zk.FungibleKeeper.QueryUniswapV2RouterGetZRC4AmountsIn(ctx, expectedInZeta, zrc20Addr)
		require.NoError(t, err)

		err = k.PayGasInERC20AndUpdateCctx(ctx, chainID, &cctx, math.NewUint(inputAmount), false)
		require.NoError(t, err)
		require.Equal(t, inputAmount-expectedInZRC20.Uint64(), cctx.GetCurrentOutTxParam().Amount.Uint64())
		require.Equal(t, uint64(21_000), cctx.GetCurrentOutTxParam().OutboundTxGasLimit)
		require.Equal(t, "2", cctx.GetCurrentOutTxParam().OutboundTxGasPrice)
	})

	t.Run("should fail if not coin type erc20", func(t *testing.T) {
		k, ctx, _, _ := testkeeper.CrosschainKeeper(t)
		chainID := getValidEthChainID(t)
		cctx := types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType: zetacommon.CoinType_Gas,
			},
		}
		err := k.PayGasInERC20AndUpdateCctx(ctx, chainID, &cctx, math.NewUint(inputAmount), false)
		require.ErrorIs(t, err, types.ErrInvalidCoinType)
	})

	t.Run("should fail if chain is not supported", func(t *testing.T) {
		k, ctx, _, _ := testkeeper.CrosschainKeeper(t)
		cctx := types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType: zetacommon.CoinType_ERC20,
			},
		}
		err := k.PayGasInERC20AndUpdateCctx(ctx, 999999, &cctx, math.NewUint(inputAmount), false)
		require.ErrorIs(t, err, observertypes.ErrSupportedChains)
	})

	t.Run("should fail if can't query the gas price", func(t *testing.T) {
		k, ctx, sdkk, zk := testkeeper.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)

		// deploy gas coin and set fee params
		chainID := getValidEthChainID(t)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		// create a cctx reverted from zeta
		cctx := types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType: zetacommon.CoinType_ERC20,
			},
			OutboundTxParams: []*types.OutboundTxParams{
				{
					ReceiverChainId: zetacommon.ZetaChain().ChainId,
				},
				{
					ReceiverChainId: chainID,
				},
			},
		}

		err := k.PayGasInERC20AndUpdateCctx(ctx, chainID, &cctx, math.NewUint(inputAmount), false)
		require.ErrorIs(t, err, types.ErrCannotFindGasParams)
	})

	t.Run("should fail if can't find the ZRC20", func(t *testing.T) {
		k, ctx, sdkk, zk := testkeeper.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)

		// deploy gas coin, erc20 and set fee params
		chainID := getValidEthChainID(t)
		assetAddress := sample.EthAddress().String()
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		gasZRC20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "foo", "foo")
		_, err := zk.FungibleKeeper.UpdateZRC20WithdrawFee(
			sdk.UnwrapSDKContext(ctx),
			fungibletypes.NewMsgUpdateZRC20WithdrawFee(admin, gasZRC20.String(), sdk.NewUint(withdrawFee)),
		)
		require.NoError(t, err)
		k.SetGasPrice(ctx, types.GasPrice{
			ChainId:     chainID,
			MedianIndex: 0,
			Prices:      []uint64{gasPrice},
		})

		// zrc20 not deployed

		// create a cctx reverted from zeta
		cctx := types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType: zetacommon.CoinType_ERC20,
				Asset:    assetAddress,
			},
			OutboundTxParams: []*types.OutboundTxParams{
				{
					ReceiverChainId: zetacommon.ZetaChain().ChainId,
				},
				{
					ReceiverChainId: chainID,
				},
			},
		}

		err = k.PayGasInERC20AndUpdateCctx(ctx, chainID, &cctx, math.NewUint(inputAmount), false)
		require.ErrorIs(t, err, types.ErrForeignCoinNotFound)
	})

	t.Run("should fail if liquidity pool not setup", func(t *testing.T) {
		k, ctx, sdkk, zk := testkeeper.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)

		// deploy gas coin, erc20 and set fee params
		chainID := getValidEthChainID(t)
		assetAddress := sample.EthAddress().String()
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		gasZRC20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "foo", "foo")
		deployZRC20(
			t,
			ctx,
			zk.FungibleKeeper,
			sdkk.EvmKeeper,
			chainID,
			"bar",
			assetAddress,
			"bar",
		)
		_, err := zk.FungibleKeeper.UpdateZRC20WithdrawFee(
			sdk.UnwrapSDKContext(ctx),
			fungibletypes.NewMsgUpdateZRC20WithdrawFee(admin, gasZRC20.String(), sdk.NewUint(withdrawFee)),
		)
		require.NoError(t, err)
		k.SetGasPrice(ctx, types.GasPrice{
			ChainId:     chainID,
			MedianIndex: 0,
			Prices:      []uint64{gasPrice},
		})

		// liquidity pool not set

		// create a cctx reverted from zeta
		cctx := types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType: zetacommon.CoinType_ERC20,
				Asset:    assetAddress,
			},
			OutboundTxParams: []*types.OutboundTxParams{
				{
					ReceiverChainId: zetacommon.ZetaChain().ChainId,
				},
				{
					ReceiverChainId: chainID,
				},
			},
		}

		err = k.PayGasInERC20AndUpdateCctx(ctx, chainID, &cctx, math.NewUint(inputAmount), false)
		require.ErrorIs(t, err, types.ErrNoLiquidityPool)
	})

	t.Run("should fail if not enough amount for the fee", func(t *testing.T) {
		k, ctx, sdkk, zk := testkeeper.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)

		// deploy gas coin, erc20 and set fee params
		chainID := getValidEthChainID(t)
		assetAddress := sample.EthAddress().String()
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		gasZRC20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "foo", "foo")
		zrc20Addr := deployZRC20(
			t,
			ctx,
			zk.FungibleKeeper,
			sdkk.EvmKeeper,
			chainID,
			"bar",
			assetAddress,
			"bar",
		)
		_, err := zk.FungibleKeeper.UpdateZRC20WithdrawFee(
			sdk.UnwrapSDKContext(ctx),
			fungibletypes.NewMsgUpdateZRC20WithdrawFee(admin, gasZRC20.String(), sdk.NewUint(withdrawFee)),
		)
		require.NoError(t, err)
		k.SetGasPrice(ctx, types.GasPrice{
			ChainId:     chainID,
			MedianIndex: 0,
			Prices:      []uint64{gasPrice},
		})

		setupZRC20Pool(
			t,
			ctx,
			zk.FungibleKeeper,
			sdkk.BankKeeper,
			zrc20Addr,
		)

		// create a cctx reverted from zeta
		cctx := types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType: zetacommon.CoinType_ERC20,
				Asset:    assetAddress,
			},
			OutboundTxParams: []*types.OutboundTxParams{
				{
					ReceiverChainId: zetacommon.ZetaChain().ChainId,
				},
				{
					ReceiverChainId: chainID,
				},
			},
		}

		// total fees in gas must be 21000*2+1000=43000
		// we calculate what it represents in erc20
		expectedInZeta, err := zk.FungibleKeeper.QueryUniswapV2RouterGetZetaAmountsIn(ctx, big.NewInt(43000), gasZRC20)
		require.NoError(t, err)
		expectedInZRC20, err := zk.FungibleKeeper.QueryUniswapV2RouterGetZRC4AmountsIn(ctx, expectedInZeta, zrc20Addr)
		require.NoError(t, err)

		// Provide expected value minus 1
		err = k.PayGasInERC20AndUpdateCctx(ctx, chainID, &cctx, math.NewUintFromBigInt(expectedInZRC20).SubUint64(1), false)
		require.ErrorIs(t, err, types.ErrNotEnoughGas)
	})
}

func TestKeeper_PayGasInZetaAndUpdateCctx(t *testing.T) {
	t.Run("can pay gas in zeta", func(t *testing.T) {
		k, ctx, sdkk, zk := testkeeper.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)

		// deploy gas coin and set fee params
		chainID := getValidEthChainID(t)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "foobar", "foobar")
		k.SetGasPrice(ctx, types.GasPrice{
			ChainId:     chainID,
			MedianIndex: 0,
			Prices:      []uint64{gasPrice},
		})

		// create a cctx reverted from zeta
		cctx := types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType: zetacommon.CoinType_Zeta,
			},
			OutboundTxParams: []*types.OutboundTxParams{
				{
					ReceiverChainId:    chainID,
					OutboundTxGasLimit: 1000,
				},
			},
			ZetaFees: math.NewUint(100),
		}
		// gasLimit * gasPrice * 2 = 1000 * 2 * 2 = 4000
		expectedOutTxGasFeeInZeta, err := zk.FungibleKeeper.QueryUniswapV2RouterGetZetaAmountsIn(ctx, big.NewInt(4000), zrc20)
		require.NoError(t, err)

		// the output amount must be input amount - (out tx fee in zeta + protocol flat fee)
		expectedFeeInZeta := types.GetProtocolFee().Add(math.NewUintFromBigInt(expectedOutTxGasFeeInZeta))
		inputAmount := expectedFeeInZeta.Add(math.NewUint(100000))
		err = k.PayGasInZetaAndUpdateCctx(ctx, chainID, &cctx, inputAmount, false)
		require.NoError(t, err)
		require.Equal(t, "100000", cctx.GetCurrentOutTxParam().Amount.String())
		require.Equal(t, "4", cctx.GetCurrentOutTxParam().OutboundTxGasPrice) // gas price is doubled
		require.True(t, cctx.ZetaFees.Equal(expectedFeeInZeta.Add(math.NewUint(100))), "expected %s, got %s", expectedFeeInZeta.String(), cctx.ZetaFees.String())

		// can call with undefined zeta fees
		cctx = types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType: zetacommon.CoinType_Zeta,
			},
			OutboundTxParams: []*types.OutboundTxParams{
				{
					ReceiverChainId:    chainID,
					OutboundTxGasLimit: 1000,
				},
			},
		}
		expectedOutTxGasFeeInZeta, err = zk.FungibleKeeper.QueryUniswapV2RouterGetZetaAmountsIn(ctx, big.NewInt(4000), zrc20)
		require.NoError(t, err)
		expectedFeeInZeta = types.GetProtocolFee().Add(math.NewUintFromBigInt(expectedOutTxGasFeeInZeta))
		inputAmount = expectedFeeInZeta.Add(math.NewUint(100000))
		err = k.PayGasInZetaAndUpdateCctx(ctx, chainID, &cctx, inputAmount, false)
		require.NoError(t, err)
		require.Equal(t, "100000", cctx.GetCurrentOutTxParam().Amount.String())
		require.Equal(t, "4", cctx.GetCurrentOutTxParam().OutboundTxGasPrice) // gas price is doubled
		require.True(t, cctx.ZetaFees.Equal(expectedFeeInZeta), "expected %s, got %s", expectedFeeInZeta.String(), cctx.ZetaFees.String())
	})

	t.Run("should fail if pay gas in zeta with coin type other than zeta", func(t *testing.T) {
		k, ctx, _, _ := testkeeper.CrosschainKeeper(t)
		chainID := getValidEthChainID(t)
		cctx := types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType: zetacommon.CoinType_Gas,
			},
		}
		err := k.PayGasInZetaAndUpdateCctx(ctx, chainID, &cctx, math.NewUint(100000), false)
		require.ErrorIs(t, err, types.ErrInvalidCoinType)
	})

	t.Run("should fail if chain is not supported", func(t *testing.T) {
		k, ctx, _, _ := testkeeper.CrosschainKeeper(t)
		cctx := types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType: zetacommon.CoinType_Zeta,
			},
		}
		err := k.PayGasInZetaAndUpdateCctx(ctx, 999999, &cctx, math.NewUint(100000), false)
		require.ErrorIs(t, err, observertypes.ErrSupportedChains)
	})

	t.Run("should fail if can't query the gas price", func(t *testing.T) {
		k, ctx, sdkk, zk := testkeeper.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)

		// deploy gas coin and set fee params
		chainID := getValidEthChainID(t)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		// gas price not set

		// create a cctx reverted from zeta
		cctx := types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType: zetacommon.CoinType_Zeta,
			},
			OutboundTxParams: []*types.OutboundTxParams{
				{
					ReceiverChainId:    chainID,
					OutboundTxGasLimit: 1000,
				},
			},
		}

		err := k.PayGasInZetaAndUpdateCctx(ctx, chainID, &cctx, math.NewUint(100000), false)
		require.ErrorIs(t, err, types.ErrUnableToGetGasPrice)
	})

	t.Run("should fail if not enough amount for the fee", func(t *testing.T) {
		k, ctx, sdkk, zk := testkeeper.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)

		// deploy gas coin and set fee params
		chainID := getValidEthChainID(t)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "foobar", "foobar")
		k.SetGasPrice(ctx, types.GasPrice{
			ChainId:     chainID,
			MedianIndex: 0,
			Prices:      []uint64{gasPrice},
		})

		// create a cctx reverted from zeta
		cctx := types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType: zetacommon.CoinType_Zeta,
			},
			OutboundTxParams: []*types.OutboundTxParams{
				{
					ReceiverChainId:    chainID,
					OutboundTxGasLimit: 1000,
				},
			},
			ZetaFees: math.NewUint(100),
		}
		expectedOutTxGasFeeInZeta, err := zk.FungibleKeeper.QueryUniswapV2RouterGetZetaAmountsIn(ctx, big.NewInt(4000), zrc20)
		require.NoError(t, err)
		expectedFeeInZeta := types.GetProtocolFee().Add(math.NewUintFromBigInt(expectedOutTxGasFeeInZeta))

		// set input amount lower than total zeta fee
		inputAmount := expectedFeeInZeta.Sub(math.NewUint(1))
		err = k.PayGasInZetaAndUpdateCctx(ctx, chainID, &cctx, inputAmount, false)
		require.ErrorIs(t, err, types.ErrNotEnoughZetaBurnt)
	})
}

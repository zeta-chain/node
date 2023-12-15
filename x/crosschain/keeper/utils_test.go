// This file contains helper functions for testing the crosschain module
package keeper_test

import (
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/ethereum/go-ethereum/common"
	evmkeeper "github.com/evmos/ethermint/x/evm/keeper"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/uniswap/v2-periphery/contracts/uniswapv2router02.sol"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	zetacommon "github.com/zeta-chain/zetacore/common"
	testkeeper "github.com/zeta-chain/zetacore/testutil/keeper"
	fungiblekeeper "github.com/zeta-chain/zetacore/x/fungible/keeper"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// getValidEthChainID get a valid eth chain id
func getValidEthChainID(t *testing.T) int64 {
	return getValidEthChain(t).ChainId
}

// getValidEthChain get a valid eth chain
func getValidEthChain(_ *testing.T) *zetacommon.Chain {
	goerli := zetacommon.GoerliLocalnetChain()
	return &goerli
}

// getValidEthChainIDWithIndex get a valid eth chain id with index
func getValidEthChainIDWithIndex(t *testing.T, index int) int64 {
	switch index {
	case 0:
		return zetacommon.GoerliLocalnetChain().ChainId
	case 1:
		return zetacommon.GoerliChain().ChainId
	default:
		require.Fail(t, "invalid index")
	}
	return 0
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
		nil,
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

// setAdminPolicies sets the admin policies for the observer module with group 1 and 2
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

// setSupportedChain sets the supported chains for the observer module
func setSupportedChain(ctx sdk.Context, zk testkeeper.ZetaKeepers, chainIDs ...int64) {
	coreParamsList := make([]*observertypes.CoreParams, len(chainIDs))
	for i, chainID := range chainIDs {
		coreParams := sample.CoreParams(chainID)
		coreParams.IsSupported = true
		coreParamsList[i] = coreParams
	}
	zk.ObserverKeeper.SetCoreParamsList(ctx, observertypes.CoreParamsList{
		CoreParams: coreParamsList,
	})
}

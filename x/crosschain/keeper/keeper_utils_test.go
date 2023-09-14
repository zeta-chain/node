package keeper_test

import (
	"testing"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	evmkeeper "github.com/evmos/ethermint/x/evm/keeper"
	"github.com/stretchr/testify/require"
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

func setAdminDeployFungibleCoin(ctx sdk.Context, zk testkeeper.ZetaKeepers, admin string) {
	params := zk.ObserverKeeper.GetParams(ctx)
	params.AdminPolicy = []*observertypes.Admin_Policy{
		{
			PolicyType: observertypes.Policy_Type_deploy_fungible_coin,
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
		setAdminDeployFungibleCoin(ctx, zk, admin)

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
		require.Equal(t, uint64(57000), cctx.GetCurrentOutTxParam().Amount.Uint64(), "output amount must be 57000 but is %d", cctx.GetCurrentOutTxParam().Amount.Uint64())
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
		require.ErrorIs(t, err, observertypes.ErrInvalidCoinType)
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
		setAdminDeployFungibleCoin(ctx, zk, admin)

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
		require.ErrorIs(t, err, types.ErrUnableToGetGasPrice)
	})

	t.Run("should fail if not enough amount for the fee", func(t *testing.T) {
		k, ctx, sdkk, zk := testkeeper.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		admin := sample.AccAddress()
		setAdminDeployFungibleCoin(ctx, zk, admin)

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

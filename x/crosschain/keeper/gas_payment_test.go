package keeper_test

import (
	"math/big"
	"testing"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/testutil/sample"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"

	"cosmossdk.io/math"

	"github.com/stretchr/testify/require"
	testkeeper "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
)

var (
	// gasLimit = big.NewInt(21_000) - value used in SetupChainGasCoinAndPool for gas limit initialization
	withdrawFee int64  = 1000
	gasPrice    uint64 = 2
	inputAmount uint64 = 1e16
)

func TestKeeper_PayGasNativeAndUpdateCctx(t *testing.T) {
	t.Run("can pay gas in native gas", func(t *testing.T) {
		k, ctx, sdkk, zk := testkeeper.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		// deploy gas coin and set fee params
		chainID := getValidEthChainID()
		setSupportedChain(ctx, zk, chainID)

		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		_, err := zk.FungibleKeeper.UpdateZRC20ProtocolFlatFee(ctx, zrc20, big.NewInt(withdrawFee))
		require.NoError(t, err)

		k.SetGasPrice(ctx, types.GasPrice{
			ChainId:     chainID,
			MedianIndex: 0,
			Prices:      []uint64{gasPrice},
		})

		// create a cctx reverted from zeta
		cctx := types.CrossChainTx{
			InboundParams: &types.InboundParams{
				CoinType: coin.CoinType_Gas,
			},
			OutboundParams: []*types.OutboundParams{
				{
					ReceiverChainId: chains.ZetaPrivnetChain.ChainId,
					CoinType:        coin.CoinType_Gas,
				},
				{
					ReceiverChainId: chainID,
				},
			},
		}

		// total fees must be 21000*2+1000=43000
		// if the input amount of the cctx is 1e16, the output amount must be 1e16-43000=9999999999957000
		err = k.PayGasNativeAndUpdateCctx(ctx, chainID, &cctx, math.NewUint(inputAmount))
		require.NoError(t, err)
		require.Equal(t, uint64(9999999999957000), cctx.GetCurrentOutboundParam().Amount.Uint64())
		require.Equal(t, uint64(21_000), cctx.GetCurrentOutboundParam().GasLimit)
		require.Equal(t, "2", cctx.GetCurrentOutboundParam().GasPrice)
	})

	t.Run("should fail if not coin type gas", func(t *testing.T) {
		k, ctx, _, _ := testkeeper.CrosschainKeeper(t)
		chainID := getValidEthChainID()
		cctx := types.CrossChainTx{
			InboundParams: &types.InboundParams{
				CoinType: coin.CoinType_Zeta,
			},
		}
		err := k.PayGasNativeAndUpdateCctx(ctx, chainID, &cctx, math.NewUint(inputAmount))
		require.ErrorIs(t, err, types.ErrInvalidCoinType)
	})

	t.Run("should fail if chain is not supported", func(t *testing.T) {
		k, ctx, _, _ := testkeeper.CrosschainKeeper(t)
		cctx := types.CrossChainTx{
			InboundParams: &types.InboundParams{
				CoinType: coin.CoinType_Gas,
			},
		}
		err := k.PayGasNativeAndUpdateCctx(ctx, 999999, &cctx, math.NewUint(inputAmount))
		require.ErrorIs(t, err, observertypes.ErrSupportedChains)
	})

	t.Run("should fail if can't query the gas price", func(t *testing.T) {
		k, ctx, sdkk, zk := testkeeper.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		// deploy gas coin and set fee params
		chainID := getValidEthChainID()
		setSupportedChain(ctx, zk, chainID)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		// create a cctx reverted from zeta
		cctx := types.CrossChainTx{
			InboundParams: &types.InboundParams{
				CoinType: coin.CoinType_Gas,
			},
			OutboundParams: []*types.OutboundParams{
				{
					ReceiverChainId: chains.ZetaPrivnetChain.ChainId,
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

		// deploy gas coin and set fee params
		chainID := getValidEthChainID()
		setSupportedChain(ctx, zk, chainID)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		_, err := zk.FungibleKeeper.UpdateZRC20ProtocolFlatFee(ctx, zrc20, big.NewInt(withdrawFee))
		require.NoError(t, err)

		k.SetGasPrice(ctx, types.GasPrice{
			ChainId:     chainID,
			MedianIndex: 0,
			Prices:      []uint64{gasPrice},
		})

		cctx := types.CrossChainTx{
			InboundParams: &types.InboundParams{
				SenderChainId: chainID,
				Sender:        sample.EthAddress().String(),
				CoinType:      coin.CoinType_Gas,
			},
			OutboundParams: []*types.OutboundParams{
				{
					ReceiverChainId: chains.ZetaPrivnetChain.ChainId,
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

		// deploy gas coin, erc20 and set fee params
		chainID := getValidEthChainID()
		setSupportedChain(ctx, zk, chainID)
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

		_, err := zk.FungibleKeeper.UpdateZRC20ProtocolFlatFee(ctx, gasZRC20, big.NewInt(withdrawFee))
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

			InboundParams: &types.InboundParams{
				Asset:    assetAddress,
				CoinType: coin.CoinType_ERC20,
			},
			OutboundParams: []*types.OutboundParams{
				{
					ReceiverChainId: chains.ZetaPrivnetChain.ChainId,
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
		require.Equal(t, inputAmount-expectedInZRC20.Uint64(), cctx.GetCurrentOutboundParam().Amount.Uint64())
		require.Equal(t, uint64(21_000), cctx.GetCurrentOutboundParam().GasLimit)
		require.Equal(t, "2", cctx.GetCurrentOutboundParam().GasPrice)
	})

	t.Run("should fail if not coin type erc20", func(t *testing.T) {
		k, ctx, _, _ := testkeeper.CrosschainKeeper(t)
		chainID := getValidEthChainID()
		cctx := types.CrossChainTx{
			InboundParams: &types.InboundParams{
				CoinType: coin.CoinType_Gas,
			},
		}
		err := k.PayGasInERC20AndUpdateCctx(ctx, chainID, &cctx, math.NewUint(inputAmount), false)
		require.ErrorIs(t, err, types.ErrInvalidCoinType)
	})

	t.Run("should fail if chain is not supported", func(t *testing.T) {
		k, ctx, _, _ := testkeeper.CrosschainKeeper(t)
		cctx := types.CrossChainTx{
			InboundParams: &types.InboundParams{
				CoinType: coin.CoinType_ERC20,
			},
		}
		err := k.PayGasInERC20AndUpdateCctx(ctx, 999999, &cctx, math.NewUint(inputAmount), false)
		require.ErrorIs(t, err, observertypes.ErrSupportedChains)
	})

	t.Run("should fail if can't query the gas price", func(t *testing.T) {
		k, ctx, sdkk, zk := testkeeper.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		// deploy gas coin and set fee params
		chainID := getValidEthChainID()
		setSupportedChain(ctx, zk, chainID)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		// create a cctx reverted from zeta
		cctx := types.CrossChainTx{
			InboundParams: &types.InboundParams{
				CoinType: coin.CoinType_ERC20,
			},
			OutboundParams: []*types.OutboundParams{
				{
					ReceiverChainId: chains.ZetaPrivnetChain.ChainId,
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

		// deploy gas coin, erc20 and set fee params
		chainID := getValidEthChainID()
		setSupportedChain(ctx, zk, chainID)
		assetAddress := sample.EthAddress().String()
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		gasZRC20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "foo", "foo")

		_, err := zk.FungibleKeeper.UpdateZRC20ProtocolFlatFee(ctx, gasZRC20, big.NewInt(withdrawFee))
		require.NoError(t, err)

		k.SetGasPrice(ctx, types.GasPrice{
			ChainId:     chainID,
			MedianIndex: 0,
			Prices:      []uint64{gasPrice},
		})

		// zrc20 not deployed

		// create a cctx reverted from zeta
		cctx := types.CrossChainTx{
			InboundParams: &types.InboundParams{
				CoinType: coin.CoinType_ERC20,
				Asset:    assetAddress,
			},
			OutboundParams: []*types.OutboundParams{
				{
					ReceiverChainId: chains.ZetaPrivnetChain.ChainId,
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

		// deploy gas coin, erc20 and set fee params
		chainID := getValidEthChainID()
		setSupportedChain(ctx, zk, chainID)
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

		_, err := zk.FungibleKeeper.UpdateZRC20ProtocolFlatFee(ctx, gasZRC20, big.NewInt(withdrawFee))
		require.NoError(t, err)

		k.SetGasPrice(ctx, types.GasPrice{
			ChainId:     chainID,
			MedianIndex: 0,
			Prices:      []uint64{gasPrice},
		})

		// liquidity pool not set

		// create a cctx reverted from zeta
		cctx := types.CrossChainTx{
			InboundParams: &types.InboundParams{
				Asset:    assetAddress,
				CoinType: coin.CoinType_ERC20,
			},
			OutboundParams: []*types.OutboundParams{
				{
					ReceiverChainId: chains.ZetaPrivnetChain.ChainId,
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

		// deploy gas coin, erc20 and set fee params
		chainID := getValidEthChainID()
		setSupportedChain(ctx, zk, chainID)
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

		_, err := zk.FungibleKeeper.UpdateZRC20ProtocolFlatFee(ctx, gasZRC20, big.NewInt(withdrawFee))
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
			InboundParams: &types.InboundParams{
				CoinType: coin.CoinType_ERC20,
				Asset:    assetAddress,
			},
			OutboundParams: []*types.OutboundParams{
				{
					ReceiverChainId: chains.ZetaPrivnetChain.ChainId,
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

		// deploy gas coin and set fee params
		chainID := getValidEthChainID()
		setSupportedChain(ctx, zk, chainID)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "foobar", "foobar")
		k.SetGasPrice(ctx, types.GasPrice{
			ChainId:     chainID,
			MedianIndex: 0,
			Prices:      []uint64{gasPrice},
		})

		// create a cctx reverted from zeta
		cctx := types.CrossChainTx{
			InboundParams: &types.InboundParams{
				CoinType: coin.CoinType_Zeta,
			},
			OutboundParams: []*types.OutboundParams{
				{
					ReceiverChainId: chainID,
					GasLimit:        1000,
				},
			},
			ZetaFees: math.NewUint(100),
		}
		// gasLimit * gasPrice * 2 = 1000 * 2 * 2 = 4000
		expectedOutboundGasFeeInZeta, err := zk.FungibleKeeper.QueryUniswapV2RouterGetZetaAmountsIn(ctx, big.NewInt(4000), zrc20)
		require.NoError(t, err)

		// the output amount must be input amount - (out tx fee in zeta + protocol flat fee)
		expectedFeeInZeta := types.GetProtocolFee().Add(math.NewUintFromBigInt(expectedOutboundGasFeeInZeta))
		inputAmount := expectedFeeInZeta.Add(math.NewUint(100000))
		err = k.PayGasInZetaAndUpdateCctx(ctx, chainID, &cctx, inputAmount, false)
		require.NoError(t, err)
		require.Equal(t, "100000", cctx.GetCurrentOutboundParam().Amount.String())
		require.Equal(t, "4", cctx.GetCurrentOutboundParam().GasPrice) // gas price is doubled
		require.True(t, cctx.ZetaFees.Equal(expectedFeeInZeta.Add(math.NewUint(100))), "expected %s, got %s", expectedFeeInZeta.String(), cctx.ZetaFees.String())

		// can call with undefined zeta fees
		cctx = types.CrossChainTx{
			InboundParams: &types.InboundParams{
				CoinType: coin.CoinType_Zeta,
			},
			OutboundParams: []*types.OutboundParams{
				{
					ReceiverChainId: chainID,
					GasLimit:        1000,
				},
			},
		}
		expectedOutboundGasFeeInZeta, err = zk.FungibleKeeper.QueryUniswapV2RouterGetZetaAmountsIn(ctx, big.NewInt(4000), zrc20)
		require.NoError(t, err)
		expectedFeeInZeta = types.GetProtocolFee().Add(math.NewUintFromBigInt(expectedOutboundGasFeeInZeta))
		inputAmount = expectedFeeInZeta.Add(math.NewUint(100000))
		err = k.PayGasInZetaAndUpdateCctx(ctx, chainID, &cctx, inputAmount, false)
		require.NoError(t, err)
		require.Equal(t, "100000", cctx.GetCurrentOutboundParam().Amount.String())
		require.Equal(t, "4", cctx.GetCurrentOutboundParam().GasPrice) // gas price is doubled
		require.True(t, cctx.ZetaFees.Equal(expectedFeeInZeta), "expected %s, got %s", expectedFeeInZeta.String(), cctx.ZetaFees.String())
	})

	t.Run("should fail if pay gas in zeta with coin type other than zeta", func(t *testing.T) {
		k, ctx, _, _ := testkeeper.CrosschainKeeper(t)
		chainID := getValidEthChainID()
		cctx := types.CrossChainTx{
			InboundParams: &types.InboundParams{
				CoinType: coin.CoinType_Gas,
			},
		}
		err := k.PayGasInZetaAndUpdateCctx(ctx, chainID, &cctx, math.NewUint(100000), false)
		require.ErrorIs(t, err, types.ErrInvalidCoinType)
	})

	t.Run("should fail if chain is not supported", func(t *testing.T) {
		k, ctx, _, _ := testkeeper.CrosschainKeeper(t)
		cctx := types.CrossChainTx{
			InboundParams: &types.InboundParams{
				CoinType: coin.CoinType_Zeta,
			},
		}
		err := k.PayGasInZetaAndUpdateCctx(ctx, 999999, &cctx, math.NewUint(100000), false)
		require.ErrorIs(t, err, observertypes.ErrSupportedChains)
	})

	t.Run("should fail if can't query the gas price", func(t *testing.T) {
		k, ctx, sdkk, zk := testkeeper.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		// deploy gas coin and set fee params
		chainID := getValidEthChainID()
		setSupportedChain(ctx, zk, chainID)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		// gas price not set

		// create a cctx reverted from zeta
		cctx := types.CrossChainTx{
			InboundParams: &types.InboundParams{
				SenderChainId: chainID,
				Sender:        sample.EthAddress().String(),
				CoinType:      coin.CoinType_Zeta,
			},
			OutboundParams: []*types.OutboundParams{
				{
					ReceiverChainId: chainID,
					GasLimit:        1000,
				},
			},
		}

		err := k.PayGasInZetaAndUpdateCctx(ctx, chainID, &cctx, math.NewUint(100000), false)
		require.ErrorIs(t, err, types.ErrUnableToGetGasPrice)
	})

	t.Run("should fail if not enough amount for the fee", func(t *testing.T) {
		k, ctx, sdkk, zk := testkeeper.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		// deploy gas coin and set fee params
		chainID := getValidEthChainID()
		setSupportedChain(ctx, zk, chainID)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "foobar", "foobar")
		k.SetGasPrice(ctx, types.GasPrice{
			ChainId:     chainID,
			MedianIndex: 0,
			Prices:      []uint64{gasPrice},
		})

		// create a cctx reverted from zeta
		cctx := types.CrossChainTx{
			InboundParams: &types.InboundParams{
				SenderChainId: chainID,
				Sender:        sample.EthAddress().String(),
				CoinType:      coin.CoinType_Zeta,
			},
			OutboundParams: []*types.OutboundParams{
				{
					ReceiverChainId: chainID,
					GasLimit:        1000,
				},
			},
			ZetaFees: math.NewUint(100),
		}
		expectedOutboundGasFeeInZeta, err := zk.FungibleKeeper.QueryUniswapV2RouterGetZetaAmountsIn(ctx, big.NewInt(4000), zrc20)
		require.NoError(t, err)
		expectedFeeInZeta := types.GetProtocolFee().Add(math.NewUintFromBigInt(expectedOutboundGasFeeInZeta))

		// set input amount lower than total zeta fee
		inputAmount := expectedFeeInZeta.Sub(math.NewUint(1))
		err = k.PayGasInZetaAndUpdateCctx(ctx, chainID, &cctx, inputAmount, false)
		require.ErrorIs(t, err, types.ErrNotEnoughZetaBurnt)
	})
}

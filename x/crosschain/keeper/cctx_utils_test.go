package keeper_test

import (
	"math/big"
	"testing"

	"cosmossdk.io/math"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestKeeper_RefundAmountOnZetaChain(t *testing.T) {
	t.Run("should refund amount on zeta chain", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		asset := sample.EthAddress().String()
		sender := sample.EthAddress()
		chainID := getValidEthChainID(t)

		// deploy zrc20
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		zrc20Addr := deployZRC20(
			t,
			ctx,
			zk.FungibleKeeper,
			sdkk.EvmKeeper,
			chainID,
			"bar",
			asset,
			"bar",
		)

		err := k.RefundAmountOnZetaChain(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_ERC20,
				SenderChainId: chainID,
				Sender:        sender.String(),
				Asset:         asset,
			}},
			math.NewUint(42),
		)
		require.NoError(t, err)

		// check amount deposited in balance
		balance, err := zk.FungibleKeeper.BalanceOfZRC4(ctx, zrc20Addr, sender)
		require.NoError(t, err)
		require.Equal(t, uint64(42), balance.Uint64())

		// can refund again
		err = k.RefundAmountOnZetaChain(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_ERC20,
				SenderChainId: chainID,
				Sender:        sender.String(),
				Asset:         asset,
			}},
			math.NewUint(42),
		)
		require.NoError(t, err)
		balance, err = zk.FungibleKeeper.BalanceOfZRC4(ctx, zrc20Addr, sender)
		require.NoError(t, err)
		require.Equal(t, uint64(84), balance.Uint64())
	})

	t.Run("should fail with invalid cctx", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)

		err := k.RefundAmountOnZetaChain(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType: common.CoinType_Zeta,
			}},
			math.NewUint(42),
		)
		require.ErrorContains(t, err, "unsupported coin type")

		err = k.RefundAmountOnZetaChain(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType: common.CoinType_Gas,
			}},
			math.NewUint(42),
		)
		require.ErrorContains(t, err, "unsupported coin type")

		err = k.RefundAmountOnZetaChain(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_ERC20,
				SenderChainId: 999999,
			}},
			math.NewUint(42),
		)
		require.ErrorContains(t, err, "only EVM chains are supported")

		err = k.RefundAmountOnZetaChain(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_ERC20,
				SenderChainId: getValidEthChainID(t),
				Sender:        "invalid",
			}},
			math.NewUint(42),
		)
		require.ErrorContains(t, err, "invalid sender address")

		err = k.RefundAmountOnZetaChain(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_ERC20,
				SenderChainId: getValidEthChainID(t),
				Sender:        sample.EthAddress().String(),
			},
		},
			math.Uint{},
		)
		require.ErrorContains(t, err, "no amount to refund")

		err = k.RefundAmountOnZetaChain(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_ERC20,
				SenderChainId: getValidEthChainID(t),
				Sender:        sample.EthAddress().String(),
			}},
			math.ZeroUint(),
		)
		require.ErrorContains(t, err, "no amount to refund")

		// the foreign coin has not been set
		err = k.RefundAmountOnZetaChain(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_ERC20,
				SenderChainId: getValidEthChainID(t),
				Sender:        sample.EthAddress().String(),
				Asset:         sample.EthAddress().String(),
			}},
			math.NewUint(42),
		)
		require.ErrorContains(t, err, "zrc not found")
	})
}

func TestGetRevertGasLimit(t *testing.T) {
	t.Run("should return 0 if no inbound tx params", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)

		gasLimit, err := k.GetRevertGasLimit(ctx, types.CrossChainTx{})
		require.NoError(t, err)
		require.Equal(t, uint64(0), gasLimit)
	})

	t.Run("should return 0 if coin type is not gas or erc20", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)

		gasLimit, err := k.GetRevertGasLimit(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType: common.CoinType_Zeta,
			}})
		require.NoError(t, err)
		require.Equal(t, uint64(0), gasLimit)
	})

	t.Run("should return the gas limit of the gas token", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := getValidEthChainID(t)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		gas := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "foo", "FOO")

		_, err := zk.FungibleKeeper.UpdateZRC20GasLimit(ctx, gas, big.NewInt(42))
		require.NoError(t, err)

		gasLimit, err := k.GetRevertGasLimit(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_Gas,
				SenderChainId: chainID,
			}})
		require.NoError(t, err)
		require.Equal(t, uint64(42), gasLimit)
	})

	t.Run("should return the gas limit of the associated asset", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := getValidEthChainID(t)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		asset := sample.EthAddress().String()
		zrc20Addr := deployZRC20(
			t,
			ctx,
			zk.FungibleKeeper,
			sdkk.EvmKeeper,
			chainID,
			"bar",
			asset,
			"bar",
		)

		_, err := zk.FungibleKeeper.UpdateZRC20GasLimit(ctx, zrc20Addr, big.NewInt(42))
		require.NoError(t, err)

		gasLimit, err := k.GetRevertGasLimit(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_ERC20,
				SenderChainId: chainID,
				Asset:         asset,
			}})
		require.NoError(t, err)
		require.Equal(t, uint64(42), gasLimit)
	})

	t.Run("should fail if no gas coin found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)

		_, err := k.GetRevertGasLimit(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_Gas,
				SenderChainId: 999999,
			}})
		require.ErrorIs(t, err, types.ErrForeignCoinNotFound)
	})

	t.Run("should fail if query gas limit for gas coin fails", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := getValidEthChainID(t)

		zk.FungibleKeeper.SetForeignCoins(ctx, fungibletypes.ForeignCoin{
			Zrc20ContractAddress: sample.EthAddress().String(),
			ForeignChainId:       chainID,
			CoinType:             common.CoinType_Gas,
		})

		// no contract deployed therefore will fail
		_, err := k.GetRevertGasLimit(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_Gas,
				SenderChainId: chainID,
			}})
		require.ErrorIs(t, err, fungibletypes.ErrContractCall)
	})

	t.Run("should fail if no asset found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)

		_, err := k.GetRevertGasLimit(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_ERC20,
				SenderChainId: 999999,
			}})
		require.ErrorIs(t, err, types.ErrForeignCoinNotFound)
	})

	t.Run("should fail if query gas limit for asset fails", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := getValidEthChainID(t)
		asset := sample.EthAddress().String()

		zk.FungibleKeeper.SetForeignCoins(ctx, fungibletypes.ForeignCoin{
			Zrc20ContractAddress: sample.EthAddress().String(),
			ForeignChainId:       chainID,
			CoinType:             common.CoinType_ERC20,
			Asset:                asset,
		})

		// no contract deployed therefore will fail
		_, err := k.GetRevertGasLimit(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_ERC20,
				SenderChainId: chainID,
				Asset:         asset,
			}})
		require.ErrorIs(t, err, fungibletypes.ErrContractCall)
	})
}
